package allarm

import (
	"log"
	"os"
)

// This file should trigger warnings for panic, log.Fatal, and os.Exit

func allarmPanic() {
	panic("should be reported") // want "panic call should not be used in production code"
}

func allarmLogFatal() {
	log.Fatal("should be reported") // want "log.Fatal call should only be used in main function of main package"
}

func allarmLogFatalf() {
	log.Fatalf("should be reported: %s", "error") // want "log.Fatal call should only be used in main function of main package"
}

func allarmLogFatalln() {
	log.Fatalln("should be reported") // want "log.Fatal call should only be used in main function of main package"
}

func allarmOsExit() {
	os.Exit(1) // want "os.Exit call should only be used in main function of main package"
}
