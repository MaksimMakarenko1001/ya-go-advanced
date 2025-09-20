package main

import (
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/config"
)

func main() {
	log.Println("application starting")

	if err := run(); err != nil {
		panic(err)
	}

	log.Println("application stoped")
}

func run() error {
	di := config.DI{}
	di.Init()

	return di.Start()
}
