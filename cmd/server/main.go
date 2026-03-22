package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

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
