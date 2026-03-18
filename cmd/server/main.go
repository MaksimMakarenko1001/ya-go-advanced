package main

import (
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/config"
)

func main() {
	log.Println("server starting")

	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("server stoped")
}

func run() error {
	di := config.DI{}
	di.Init("SERVER_")

	return di.Start()
}
