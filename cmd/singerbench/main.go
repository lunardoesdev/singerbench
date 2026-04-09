package main

import (
	"fmt"
	"log"
	"os"

	"github.com/napalu/goopt/v2"
)

type Config struct {
	Help struct {
	} `goopt:"kind:command;name:help;desc:help"`
}

func main() {
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

	if parser.HasCommand("add-subscription") {

	}

	_ = os.Stdout
}
