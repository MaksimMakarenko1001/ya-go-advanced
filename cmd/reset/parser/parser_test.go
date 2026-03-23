package parser

import (
	"testing"

	"go/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantLen int
	}{
		{
			name: "no generate reset comment",
			source: `package test

type User struct {
	Name string
	Age  int
}`,
			wantLen: 0,
		},
		{
			name: "single struct with generate reset comment",
			source: `package test

// generate:reset
type User struct {
	Name string
	Age  int
}`,
			wantLen: 1,
		},
		{
			name: "multiple structs with generate reset comment",
			source: `package test

// generate:reset
type User struct {
	Name string
	Age  int
}

// generate:reset
type Config struct {
	Enabled bool
	Value   float64
}`,
			wantLen: 2,
		},
		{
			name: "mixed structs with and without generate reset comment",
			source: `package test

// generate:reset
type User struct {
	Name string
	Age  int
}

type Config struct {
	Enabled bool
	Value   float64
}

// generate:reset
type Metrics struct {
	Count int
	Rate  float64
}`,
			wantLen: 2,
		},
		{
			name: "generate reset comment with extra spaces",
			source: `package test

// generate:reset  
type User struct {
	Name string
	Age  int
}`,
			wantLen: 1,
		},
		{
			name: "generate reset comment with tabs",
			source: `package test

// generate:reset	
type User struct {
	Name string
	Age  int
}`,
			wantLen: 1,
		},
		{
			name: "invalid generate comment format",
			source: `package test

// generate: reset
type User struct {
	Name string
	Age  int
}`,
			wantLen: 0,
		},
		{
			name: "generate reset comment on non-struct type",
			source: `package test

// generate:reset
type User string

// generate:reset
type Config struct {
	Enabled bool
}`,
			wantLen: 1,
		},
		{
			name:    "empty file",
			source:  `package test`,
			wantLen: 0,
		},
		{
			name: "multiple comments with generate reset",
			source: `package test

// Some comment
// generate:reset
// Another comment
type User struct {
	Name string
	Age  int
}`,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			result, err := ParseFile(fset, "", tt.source)

			require.NoError(t, err)
			assert.Equal(t, result.Name, "test")

			assert.Len(t, result.Structs, tt.wantLen)

			// Verify that all returned structs are actually struct types
			for _, structType := range result.Structs {
				assert.NotNil(t, structType)
				assert.NotNil(t, structType.Fields)
			}
		})
	}
}
