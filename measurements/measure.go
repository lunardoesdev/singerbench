package measurements

import (
	"context"
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
	port := getFreePort()

	proxy, err := singerbox.FromSharedLink(
		link,
		singerbox.ProxyConfig{
			ListenAddr: fmt.Sprintf("127.0.0.1:%d", port),
		},
	)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer proxy.Stop()

	proxyURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var (
		connectStart time.Time
		connectDone  time.Time
		gotFirstByte time.Time
	)

	start := time.Now()

	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			connectDone = time.Now()
		},
		GotFirstResponseByte: func() {
			gotFirstByte = time.Now()
		},
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		httptrace.WithClientTrace(ctx, trace),
		http.MethodGet,
		"http://cachefly.cachefly.net/1mb.test",
		nil,
	)
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

	end := time.Now()

	var pingDur time.Duration
	if !connectStart.IsZero() && !connectDone.IsZero() {
		pingDur = connectDone.Sub(connectStart)
	}

	var firstByteDur time.Duration
	if !gotFirstByte.IsZero() {
		firstByteDur = gotFirstByte.Sub(start)
	}

	lastByteDur := end.Sub(start)

	_ = pingDur
	return start.UnixMilli(), firstByteDur.Milliseconds(), lastByteDur.Milliseconds(), firstByteDur.Milliseconds(), nil
}
