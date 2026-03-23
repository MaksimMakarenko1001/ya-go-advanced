package model

import "go/ast"

type PackageName = string
type StructName = string
type TypeName = string
type FieldName = string
type Command = string
type ResetFunc = string

type PackageInfo struct {
	Name    PackageName
	Structs map[StructName]*ast.StructType
}

// StructInfo represents a struct that needs reset functionality
type StructInfo struct {
	Name   StructName
	Fields []FieldInfo
}

type ResetCommands map[FieldName]Command

// FieldInfo represents a field within a struct
type FieldInfo struct {
	Name      FieldName
	Type      TypeName
	ZeroValue string
	IsPtr     bool
	IsArray   bool
	IsSlice   bool
	IsMap     bool
	IsStruct  bool
}
