package resetor

import (
	"fmt"
	"go/token"
	"io"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/analyzer"
	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/generator"
	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/model"
	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/parser"
	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/writer"
)

func Reset(filePath string, src any, w io.Writer) error {
	fset := token.NewFileSet()

	pi, err := parser.ParseFile(fset, filePath, src)
	if err != nil {
		return fmt.Errorf("parse file error %s: %w", filePath, err)
	}

	if len(pi.Structs) == 0 {
		return nil
	}

	resets := make([]model.ResetFunc, len(pi.Structs))
	for name, strct := range pi.Structs {
		si := analyzer.AnalyzeStruct(name, strct)

		gen := generator.New(si.Name)
		commands := make(map[model.FieldName]model.Command, len(si.Fields))
		for _, fi := range si.Fields {
			if fi.IsSlice {
				commands[fi.Name] = gen.GenSliceResetCommand(fi.Name)
			} else if fi.IsMap {
				commands[fi.Name] = gen.GenMapResetCommand(fi.Name)
			} else if fi.IsStruct {
				commands[fi.Name] = gen.GenStructResetCommand(fi.Name, fi.ZeroValue)
			} else if fi.IsArray {
				if fi.IsPtr {
					commands[fi.Name] = gen.GenArrayPResetCommand(fi.Name, fi.ZeroValue)
				} else {
					commands[fi.Name] = gen.GenArrayResetCommand(fi.Name, fi.ZeroValue)
				}
			} else if fi.IsPtr {
				commands[fi.Name] = gen.GenPtrResetCommand(fi.Name, fi.ZeroValue)
			} else {
				commands[fi.Name] = gen.GenBasicResetCommand(fi.Name, fi.ZeroValue)
			}
		}

		f, err := gen.GenResetFunc(commands)
		if err != nil {
			return fmt.Errorf("generate func error %s: %w", filePath, err)
		}

		resets = append(resets, f)
	}

	return writer.WriteResetFile(pi.Name, resets, w)
}
