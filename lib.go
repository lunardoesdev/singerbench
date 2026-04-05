package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/lunardoesdev/singerbox"
)

func printMemStats(label string) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: Alloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
		label,
		m.HeapAlloc,
		m.HeapSys,
		m.NumGC,
	)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {
	printMemStats("Zero")
	port := 8989
	prox, err := singerbox.FromSharedLink(fmt.Sprintf("http://127.0.0.1:%v", port),
		singerbox.ProxyConfig{})
	port = port + 1

	if err != nil {
		log.Fatalf("Error creating proxy")
	}

	prox.Stop()

	printMemStats("One")
	for i := 0; i < 100; i++ {
		//prox.Stop()
		log.Println(fmt.Sprintf("http://127.0.0.1:%v", port))
		prox, err = singerbox.FromSharedLink(fmt.Sprintf("http://127.0.0.1:%v", port),
			singerbox.ProxyConfig{
				ListenAddr: fmt.Sprintf("http://127.0.0.1:%v", port),
			})
		port = port + 1
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Error creating proxy")
		}
		printMemStats(fmt.Sprintf("%v", i))
	}
}
