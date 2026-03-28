package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	log.Printf("agent starts on %s\n", cfg.Address)

	if err := cli.Run(ctx); !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
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
