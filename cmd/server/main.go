package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	certFile, keyFile, err := generateTLS()
	if err != nil {
		log.Printf("runtime error: %v", err)
		os.Exit(1)
	}

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

	di.Start(errCh, certFile, keyFile)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stopCh
	log.Println("server stoping")

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	di.Stop(stopCtx)
	log.Println("server stoped")
}

func generateTLS() (certPath string, keyPath string, err error) {
	tmp, err := os.MkdirTemp("", "tmp")
	if err != nil {
		return "", "", fmt.Errorf("create temporary directory error: %w", err)
	}

	cert, key := filepath.Join(tmp, "cert.pem"), filepath.Join(tmp, "key.pem")

	if err := exec.Command("go", "run", "./cmd/tls/main.go", "-cert", cert, "-private", key).Run(); err != nil {
		return "", "", fmt.Errorf("tls gen error: %w", err)
	}

	return cert, key, nil
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
