package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

	log.Println("agent starting")

	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("agent stoped")
}

func run() error {
	cfg := &agent.Config{}
	cfg.LoadConfig("AGENT_")

	cli := agent.NewClient(*cfg)
	if err := cli.WithCrypto(cfg.CryptoKey); err != nil {
		log.Printf("agent encrypt opt disabled")
	}

	log.Printf("agent starts on %s\n", cfg.Address)

	return cli.Start()
}

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
