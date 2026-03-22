package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/resetor"
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "cmd" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, ".gen.go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		dir := filepath.Dir(path)

		fname := filepath.Join(dir, "reset.gen.go")
		file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return fmt.Errorf("open file error %s: %w", fname, err)
		}

		return resetor.Reset(path, nil, file)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "walk directory error: %v\n", err)
		os.Exit(1)
	}
}
