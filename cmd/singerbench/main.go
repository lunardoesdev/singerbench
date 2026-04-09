package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lunardoesdev/singerbench/db2"
	"github.com/lunardoesdev/singerbench/measurements"
	"github.com/lunardoesdev/singerbench/mydb"
	"github.com/napalu/goopt/v2"
)

type Config struct {
	Help struct {
	} `goopt:"kind:command;name:help;desc:help"`
	AddSubscription struct {
	} `goopt:"kind:command;name:add-subscription;desc:add subscription"`
	ListSubscriptoins struct {
	} `goopt:"kind:command;name:list-subscriptions;desc:list subscriptions"`
	FetchSubscriptions struct {
	} `goopt:"kind:command;name:fetch-subscriptions;desc:fetch subscriptions"`
	RemoveSubscriptions struct {
	} `goopt:"kind:command;name:remove-subscriptions;desc:remove subscriptions"`
	ListProxies struct {
	} `goopt:"kind:command;name:list-proxies;desc:list proxies"`
	Measure struct {
	} `goopt:"kind:command;name:measure;desc:measure all proxies"`
	Gc struct {
	} `goopt:"kind:command;name:gc;desc:remove bad proxies"`
	Print struct {
	} `goopt:"kind:command;name:print;desc:print good proxies subscription-style"`
}

func spawnMeasureWorker(links chan string) {
	ctx := context.Background()
	var timer *time.Timer
myloop:
	for {
		timer = time.NewTimer(3 * time.Second)
		select {
		case link := <-links:
			when, fbyte, lbyte, ping, err := measurements.Measure(link)
			if err != nil {
				log.Printf("Warning: %v\n", err)
				break
			}

			prx, err := db2.Queries.GetProxyIdByLink(ctx, sql.NullString{
				String: link, Valid: true})
			if err != nil {
				log.Printf("Warning: %v\n", err)
				break
			}

			err = db2.Queries.SaveMeasurement(ctx, mydb.SaveMeasurementParams{
				Serverid: sql.NullInt64{
					Int64: prx.ID,
					Valid: true,
				},
				Datewhen: sql.NullInt64{
					Int64: when,
					Valid: true,
				},
				Ping: sql.NullInt64{
					Int64: ping,
					Valid: true,
				},
				Firstbyte: sql.NullInt64{
					Int64: fbyte,
					Valid: true,
				},
				Lastbyte: sql.NullInt64{
					Int64: lbyte,
					Valid: true,
				},
			})

			if err != nil {
				break
			}

			log.Printf("Saved proxy: %v", link)
		case <-timer.C:
			timer.Stop()
			break myloop
		}
		timer.Stop()
	}
}

func spawnGoMeasurer(threads int) chan string {
	channel := make(chan string, threads)
	for i := 0; i < threads; i++ {
		go spawnMeasureWorker(channel)
	}

	return channel
}

func run() error {
	ctx := context.Background()
	cfg := &Config{}
	parser, err := goopt.NewParserFromStruct(cfg)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// Parse returns false on failure or if --help was requested
	if !parser.Parse(os.Args) {
		// goopt handles printing errors and help text by default
		os.Exit(1)
	}

	if parser.HasCommand("help") {
		fmt.Println("Helpful help message")
	}

	args := parser.GetPositionalArgs()

	if parser.HasCommand("add-subscription") {
		if len(args) > 0 {
			link := args[0].Value
			fmt.Printf("adding subscription: %v\n", link)
			if err = db2.AddSubscription(ctx, link); err != nil {
				return err
			}
		}
	}

	if parser.HasCommand("remove-subscriptions") {
		for _, link := range args {
			fmt.Printf("removing %v\n", link.Value)
			err = db2.Queries.RemoveSubscription(ctx, sql.NullString{
				String: link.Value, Valid: true,
			})
			if err != nil {
				//probably doesn't matter; do nothing
			}
		}
	}

	if parser.HasCommand("measure") {
		links := spawnGoMeasurer(100)
		for sub := range db2.IterateProxies(24) {
			links <- sub.Link.String
		}
	}

	if parser.HasCommand("print") {
		proxies, err := db2.Queries.SelectBestProxies(ctx)
		if err != nil {
			return err
		}
		for _, proxy := range proxies {
			fmt.Println(proxy.Link.String)
		}
	}

	if parser.HasCommand("list-subscriptions") {
		for sub := range db2.IterateSubscriptions(24) {
			fmt.Println(sub.Link.String)
		}
	}

	if parser.HasCommand("list-proxies") {
		for sub := range db2.IterateProxies(24) {
			fmt.Println(sub.Link.String)
		}
	}

	if parser.HasCommand("fetch-subscriptions") {
		for sub := range db2.IterateSubscriptions(24) {
			resp, err := http.Get(sub.Link.String)
			if err != nil {
				log.Printf("Cant update %v; skipping to next ones...\n", sub.Link.String)
				continue //just skip and go to next sub
			}

			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				line = strings.Trim(line, " \t\r\n")
				if (len(line) > 0) && (line[0] != '#') {
					log.Printf("adding proxy: %v", line)
					db2.Addproxy(ctx, line)
				}
			}
		}
	}
	_ = os.Stdout
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
