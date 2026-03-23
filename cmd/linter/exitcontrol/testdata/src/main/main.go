package main

import (
	"log"
	"os"
)

// This file should NOT trigger warnings for log.Fatal and os.Exit in main function

func main() {
	log.Fatal("allowed in main")
	log.Fatalf("allowed in main: %s", "error")
	log.Fatalln("allowed in main")
	os.Exit(1)
}
