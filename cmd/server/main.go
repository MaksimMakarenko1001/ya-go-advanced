package main

import (
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config"
)

func main() {
	log.Println("server starting")

	if err := run(); err != nil {
		panic(err)
	}

	log.Println("server stoped")
}

func run() error {
	di := config.DI{}
	di.Init("SERVER_")

	return di.Start()
}
