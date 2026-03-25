package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	di := config.DI{}
	di.Init("SERVER_")

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)

	go func() {
		err := <-errCh
		log.Fatal(err)

		stopCh <- os.Kill
	}()

	di.Start(errCh)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stopCh
	log.Println("server stoping")

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	di.Stop(stopCtx)
	log.Println("server stoped")
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
