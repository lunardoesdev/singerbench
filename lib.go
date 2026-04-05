package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime"

	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/lunardoesdev/singerbench/mydb"
)

//go:embed schema.sql
var initsql string
var queries *mydb.Queries
var db sql.DB

func init() {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("%w\n", err)
	}

	if _, err := db.ExecContext(ctx, initsql); err != nil {
		log.Fatalf("%w\n", err)
	}

	queries = mydb.New(db)
}

func addproxy(ctx context.Context, proxy string) error {
	prox := sql.NullString{
		String: proxy,
		Valid:  true,
	}
	_, err := queries.GetProxyIdByLink(ctx, prox)
	if err != nil {
		//proxy doesnt exist
		if err := queries.AddProxy(ctx, prox); err != nil {
			return err
		}
	}

	return nil
}

func run() error {
	ctx := context.Background()

	err := addproxy(ctx, "http://127.0.0.1:7777")
	if err != nil {
		return err
	}

	theproxy, err := queries.GetProxyIdByLink(ctx, sql.NullString{
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
