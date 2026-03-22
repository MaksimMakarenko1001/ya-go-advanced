package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/model"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
)

func ParseFile(fset *token.FileSet, filePath string, src any) (*model.PackageInfo, error) {
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file error %s: %w", filePath, err)
	}

	genDecls := make([]*ast.GenDecl, 0, len(node.Decls))
	for _, decl := range node.Decls {
		if item, ok := decl.(*ast.GenDecl); ok && item.Doc != nil {
			comments := pkg.SliceFilter(item.Doc.List, func(c *ast.Comment) bool {
				return strings.TrimSpace(c.Text) == "// generate:reset"
			})
			if len(comments) > 0 {
				genDecls = append(genDecls, item)
			}
		}
	}

	structs := make(map[model.StructName]*ast.StructType, len(genDecls))
	for _, genDecl := range genDecls {
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				structs[typeSpec.Name.Name] = structType
			}
		}
	}

	return &model.PackageInfo{
		Name:    node.Name.Name,
		Structs: structs,
	}, nil
}
