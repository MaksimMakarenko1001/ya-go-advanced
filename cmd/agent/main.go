package main

import (
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/agent"
)

func main() {
	log.Println("agent starting")

	if err := run(); err != nil {
		panic(err)
	}

	log.Println("agent stoped")
}

func run() error {
	cfg := &agent.Config{}
	cfg.LoadConfig()

	cli := agent.NewClient(cfg.HTTP)

	return cli.Srart(cfg.PollInterval, cfg.ReportInterval)
}
