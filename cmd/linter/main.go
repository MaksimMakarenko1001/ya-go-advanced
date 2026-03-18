package main

import (
	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/linter/exitcontrol"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(exitcontrol.Analyzer)
}
