package measurements

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"github.com/lunardoesdev/singerbox"
)

func getFreePort() int {
	for {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			continue
		}
		defer ln.Close()

		return ln.Addr().(*net.TCPAddr).Port
	}
}

func Measure(link string) (datewhen int64, firstbyte int64, lastbyte int64, ping int64, err error) {
	// Create and start proxy from any share link (replace with your actual server)
	port := getFreePort()
	proxy, err := singerbox.FromSharedLink(
		link,
		singerbox.ProxyConfig{
			ListenAddr: fmt.Sprintf("127.0.0.1:%v", port),
			// LogLevel: "info",  // Uncomment to see connection logs
		},
	)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer proxy.Stop()

	var connectStart, connectDone, firstByte int64
	start := time.Now().UnixMilli()

	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			connectStart = time.Now().UnixMilli() - start
		},
		ConnectDone: func(network, addr string, err error) {
			connectDone = time.Now().UnixMilli() - start
		},
		GotFirstResponseByte: func() {
			firstByte = time.Now().UnixMilli() - start
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {},
	}

	proxyURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", port))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(context.Background(), trace),
		"GET", "http://cachefly.cachefly.net/1mb.test", nil)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	lastbyte = time.Now().UnixMilli() - start

	_ = connectStart
	_ = connectDone
	return start, firstByte, lastbyte, connectDone, nil
}
