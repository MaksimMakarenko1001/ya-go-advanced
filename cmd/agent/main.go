package main

import (
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/agent"
)

func main() {
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

	log.Printf("agent starts on %s\n", cfg.Address)

	return cli.Start()
}
