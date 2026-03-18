package analyzer

import (
	"fmt"
	"go/ast"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/model"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
)

type typeInfo struct {
	name     string
	isPtr    bool
	isArray  bool
	isSlice  bool
	isMap    bool
	isStruct bool
}

func AnalyzeStruct(name model.StructName, structType *ast.StructType) *model.StructInfo {
	if structType.Fields == nil {
		return nil
	}

	list := pkg.SliceFilter(structType.Fields.List, func(x *ast.Field) bool {
		return len(x.Names) > 0
	})

	fields := make([]model.FieldInfo, 0, len(list))
	for _, field := range list {
		typeInfo := analyzeFieldType(field.Type)
		for _, name := range field.Names {
			fields = append(fields, model.FieldInfo{
				Name:      name.Name,
				Type:      typeInfo.name,
				ZeroValue: getZeroValue(typeInfo.name),
				IsPtr:     typeInfo.isPtr,
				IsArray:   typeInfo.isArray,
				IsSlice:   typeInfo.isSlice,
				IsMap:     typeInfo.isMap,
				IsStruct:  typeInfo.isStruct,
			})
		}
	}

	return &model.StructInfo{
		Name:   name,
		Fields: fields,
	}
}

func analyzeFieldType(expr ast.Expr) typeInfo {
	switch t := expr.(type) {
	case *ast.Ident:
		ti := typeInfo{name: t.Name}
		if !isBasicType(t.Name) {
			ti.isStruct = true
		}
		return ti

	case *ast.StarExpr:
		ti := analyzeFieldType(t.X)
		ti.isPtr = true
		return ti

	case *ast.ArrayType:
		ti := analyzeFieldType(t.Elt)
		ti.isArray = t.Len != nil
		ti.isSlice = t.Len == nil
		return ti

	case *ast.MapType:
		return typeInfo{isMap: true}

	case *ast.SelectorExpr:
		return typeInfo{
			name:     fmt.Sprintf("%s.%s", t.X, t.Sel),
			isStruct: true,
		}

	case *ast.StructType:
		return typeInfo{isStruct: true}
	}
	return typeInfo{}
}

func isBasicType(typeName string) bool {
	basicTypes := map[string]bool{
		"int":     true,
		"int8":    true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint":    true,
		"uint8":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
		"float32": true,
		"float64": true,
		"bool":    true,
		"string":  true,
		"byte":    true,
		"rune":    true,
	}
	return basicTypes[typeName]
}

func getZeroValue(name model.TypeName) string {
	switch name {
	case "int", "int8", "int16", "int32", "int64", "rune":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		return "0"
	case "float32", "float64":
		return "0"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return name + "{}"
	}
}
