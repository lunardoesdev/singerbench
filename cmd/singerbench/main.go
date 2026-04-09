package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lunardoesdev/singerbench/db2"
	"github.com/napalu/goopt/v2"
)

type Config struct {
	Help struct {
	} `goopt:"kind:command;name:help;desc:help"`
	AddSubscription struct {
	} `goopt:"kind:command;name:add-subscription;desc:add subscription"`
	ListSubscriptoins struct {
	} `goopt:"kind:command;name:list-subscriptions;desc:list subscriptions"`
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

	if parser.HasCommand("list-subscriptions") {
		for sub := range db2.IterateSubscriptions(24) {
			fmt.Println(sub.Link.String)
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
