package analyzer

import (
	"go/ast"
	"testing"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/model"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeStruct(t *testing.T) {
	tests := []struct {
		name       string
		structName model.StructName
		structType *ast.StructType
		want       *model.StructInfo
	}{
		{
			name:       "empty struct",
			structName: "Empty",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{},
				},
			},
			want: &model.StructInfo{
				Name:   "Empty",
				Fields: []model.FieldInfo{},
			},
		},
		{
			name:       "nil fields",
			structName: "NilFields",
			structType: &ast.StructType{
				Fields: nil,
			},
			want: nil,
		},
		{
			name:       "single basic field",
			structName: "SingleField",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Name"}},
							Type:  &ast.Ident{Name: "string"},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "SingleField",
				Fields: []model.FieldInfo{
					{
						Name:      "Name",
						Type:      "string",
						ZeroValue: `""`,
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "multiple basic fields",
			structName: "User",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Name"}},
							Type:  &ast.Ident{Name: "string"},
						},
						{
							Names: []*ast.Ident{{Name: "Age"}},
							Type:  &ast.Ident{Name: "int"},
						},
						{
							Names: []*ast.Ident{{Name: "Active"}},
							Type:  &ast.Ident{Name: "bool"},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "User",
				Fields: []model.FieldInfo{
					{
						Name:      "Name",
						Type:      "string",
						ZeroValue: `""`,
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
					{
						Name:      "Age",
						Type:      "int",
						ZeroValue: "0",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
					{
						Name:      "Active",
						Type:      "bool",
						ZeroValue: "false",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "pointer field",
			structName: "PointerStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Data"}},
							Type: &ast.StarExpr{
								X: &ast.Ident{Name: "string"},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "PointerStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Data",
						Type:      "string",
						ZeroValue: `""`,
						IsPtr:     true,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "slice field",
			structName: "SliceStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Items"}},
							Type: &ast.ArrayType{
								Elt: &ast.Ident{Name: "int"},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "SliceStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Items",
						Type:      "int",
						ZeroValue: "0",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   true,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "array field",
			structName: "ArrayStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Vector"}},
							Type: &ast.ArrayType{
								Len: &ast.BasicLit{Value: "10"},
								Elt: &ast.Ident{Name: "float64"},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "ArrayStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Vector",
						Type:      "float64",
						ZeroValue: "0",
						IsPtr:     false,
						IsArray:   true,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "map field",
			structName: "MapStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Data"}},
							Type: &ast.MapType{
								Key:   &ast.Ident{Name: "string"},
								Value: &ast.Ident{Name: "int"},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "MapStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Data",
						Type:      "",
						ZeroValue: "{}",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     true,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "struct field",
			structName: "NestedStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Nested"}},
							Type:  &ast.Ident{Name: "User"},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "NestedStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Nested",
						Type:      "User",
						ZeroValue: "User{}",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  true,
					},
				},
			},
		},
		{
			name:       "selector expression field",
			structName: "ExternalStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Time"}},
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "time"},
								Sel: &ast.Ident{Name: "Time"},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "ExternalStruct",
				Fields: []model.FieldInfo{
					{
						Name:      "Time",
						Type:      "time.Time",
						ZeroValue: "time.Time{}",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  true,
					},
				},
			},
		},
		{
			name:       "anonymous field",
			structName: "AnonymousStruct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: nil,
							Type:  &ast.Ident{Name: "string"},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name:   "AnonymousStruct",
				Fields: []model.FieldInfo{},
			},
		},
		{
			name:       "multiple named fields",
			structName: "MultiNamed",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "A"}, {Name: "B"}},
							Type:  &ast.Ident{Name: "int"},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "MultiNamed",
				Fields: []model.FieldInfo{
					{
						Name:      "A",
						Type:      "int",
						ZeroValue: "0",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
					{
						Name:      "B",
						Type:      "int",
						ZeroValue: "0",
						IsPtr:     false,
						IsArray:   false,
						IsSlice:   false,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
		{
			name:       "complex nested types",
			structName: "Complex",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "PtrSlice"}},
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: &ast.Ident{Name: "string"},
								},
							},
						},
						{
							Names: []*ast.Ident{{Name: "SlicePtr"}},
							Type: &ast.ArrayType{
								Elt: &ast.StarExpr{
									X: &ast.Ident{Name: "int"},
								},
							},
						},
					},
				},
			},
			want: &model.StructInfo{
				Name: "Complex",
				Fields: []model.FieldInfo{
					{
						Name:      "PtrSlice",
						Type:      "string",
						ZeroValue: `""`,
						IsPtr:     true,
						IsArray:   false,
						IsSlice:   true,
						IsMap:     false,
						IsStruct:  false,
					},
					{
						Name:      "SlicePtr",
						Type:      "int",
						ZeroValue: "0",
						IsPtr:     true,
						IsArray:   false,
						IsSlice:   true,
						IsMap:     false,
						IsStruct:  false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnalyzeStruct(tt.structName, tt.structType)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIsBasicType(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		{
			name:     "basic type int",
			typeName: "int",
			want:     true,
		},
		{
			name:     "basic type string",
			typeName: "string",
			want:     true,
		},
		{
			name:     "basic type bool",
			typeName: "bool",
			want:     true,
		},
		{
			name:     "basic type float64",
			typeName: "float64",
			want:     true,
		},
		{
			name:     "all basic types",
			typeName: "int8",
			want:     true,
		},
		{
			name:     "all basic types int16",
			typeName: "int16",
			want:     true,
		},
		{
			name:     "all basic types int32",
			typeName: "int32",
			want:     true,
		},
		{
			name:     "all basic types int64",
			typeName: "int64",
			want:     true,
		},
		{
			name:     "all basic types uint",
			typeName: "uint",
			want:     true,
		},
		{
			name:     "all basic types uint8",
			typeName: "uint8",
			want:     true,
		},
		{
			name:     "all basic types uint16",
			typeName: "uint16",
			want:     true,
		},
		{
			name:     "all basic types uint32",
			typeName: "uint32",
			want:     true,
		},
		{
			name:     "all basic types uint64",
			typeName: "uint64",
			want:     true,
		},
		{
			name:     "all basic types float32",
			typeName: "float32",
			want:     true,
		},
		{
			name:     "all basic types byte",
			typeName: "byte",
			want:     true,
		},
		{
			name:     "all basic types rune",
			typeName: "rune",
			want:     true,
		},
		{
			name:     "custom type",
			typeName: "User",
			want:     false,
		},
		{
			name:     "empty string",
			typeName: "",
			want:     false,
		},
		{
			name:     "basic type pointer",
			typeName: "*int",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBasicType(tt.typeName)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestAnalyzeFieldType(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want typeInfo
	}{
		{
			name: "basic type",
			expr: &ast.Ident{Name: "int"},
			want: typeInfo{
				name:     "int",
				isPtr:    false,
				isArray:  false,
				isSlice:  false,
				isMap:    false,
				isStruct: false,
			},
		},
		{
			name: "custom type",
			expr: &ast.Ident{Name: "User"},
			want: typeInfo{
				name:     "User",
				isPtr:    false,
				isArray:  false,
				isSlice:  false,
				isMap:    false,
				isStruct: true,
			},
		},
		{
			name: "pointer type",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "string"},
			},
			want: typeInfo{
				name:     "string",
				isPtr:    true,
				isArray:  false,
				isSlice:  false,
				isMap:    false,
				isStruct: false,
			},
		},
		{
			name: "slice type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "int"},
			},
			want: typeInfo{
				name:     "int",
				isPtr:    false,
				isArray:  false,
				isSlice:  true,
				isMap:    false,
				isStruct: false,
			},
		},
		{
			name: "array type",
			expr: &ast.ArrayType{
				Len: &ast.BasicLit{Value: "5"},
				Elt: &ast.Ident{Name: "float64"},
			},
			want: typeInfo{
				name:     "float64",
				isPtr:    false,
				isArray:  true,
				isSlice:  false,
				isMap:    false,
				isStruct: false,
			},
		},
		{
			name: "map type",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			want: typeInfo{
				name:     "",
				isPtr:    false,
				isArray:  false,
				isSlice:  false,
				isMap:    true,
				isStruct: false,
			},
		},
		{
			name: "selector expression",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			want: typeInfo{
				name:     "time.Time",
				isPtr:    false,
				isArray:  false,
				isSlice:  false,
				isMap:    false,
				isStruct: true,
			},
		},
		{
			name: "struct type",
			expr: &ast.StructType{},
			want: typeInfo{
				name:     "",
				isPtr:    false,
				isArray:  false,
				isSlice:  false,
				isMap:    false,
				isStruct: true,
			},
		},
		{
			name: "unknown type",
			expr: &ast.BasicLit{Value: "test"},
			want: typeInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzeFieldType(tt.expr)
			assert.Equal(t, tt.want, result)
		})
	}
}
