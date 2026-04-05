package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime"

	_ "embed"

	"github.com/lunardoesdev/singerbench/db2"
)

func run() error {
	ctx := context.Background()

	err := db2.Addproxy(ctx, "http://127.0.0.1:7777")
	if err != nil {
		return err
	}

	theproxy, err := db2.Queries.GetProxyIdByLink(ctx, sql.NullString{
		String: "http://127.0.0.1:7777", Valid: true,
	})
	if err != nil {
		//proxy doesnt exist
		log.Fatalln("Proxy doesn't exist but it should....")
	}
	log.Println("found " + theproxy.Link.String)

	return nil
}

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
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}
