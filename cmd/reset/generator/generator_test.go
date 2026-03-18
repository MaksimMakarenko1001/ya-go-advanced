package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenResetFunc(t *testing.T) {
	tests := []struct {
		name     string
		resetter string
		fields   map[string]string
		want     string
	}{
		{
			name:     "no field",
			resetter: "Empty",
			fields:   map[string]string{},
			want:     "func (e *Empty) Reset() {\n\tif e == nil { return }\n}\n",
		},
		{
			name:     "one field",
			resetter: "One",
			fields: map[string]string{
				"First": "f = \"\"",
			},
			want: "func (o *One) Reset() {\n\tif o == nil { return }\n\n\t// resets First\n\tf = \"\"\n}\n",
		},
		{
			name:     "two field",
			resetter: "Two",
			fields: map[string]string{
				"First":  "f = \"\"",
				"Second": "s = \"\"",
			},
			want: "func (t *Two) Reset() {\n\tif t == nil { return }\n\n\t// resets First\n\tf = \"\"\n\n\t// resets Second\n\ts = \"\"\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result, err := gen.GenResetFunc(tt.fields)

			require.NoError(t, err)
			assert.Equal(t, result, tt.want)
		})
	}
}

func TestGenBasicResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		zeroValue string
		want      string
	}{
		{
			name:      "string zero value",
			resetter:  "User",
			fieldName: "Name",
			zeroValue: `""`,
			want:      "u.Name = \"\"",
		},
		{
			name:      "boolean zero value",
			resetter:  "Config",
			fieldName: "Enabled",
			zeroValue: "false",
			want:      "c.Enabled = false",
		},
		{
			name:      "float zero value",
			resetter:  "Metrics",
			fieldName: "Rate",
			zeroValue: "0.0",
			want:      "m.Rate = 0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenBasicResetCommand(tt.fieldName, tt.zeroValue)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenArrayResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		zeroValue string
		want      string
	}{
		{
			name:      "array",
			resetter:  "Array",
			fieldName: "vector",
			zeroValue: "0",
			want:      "for i := 0; i < len(a.vector); i++ {\n\ta.vector[i] = 0\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenArrayResetCommand(tt.fieldName, tt.zeroValue)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenArrayPResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		zeroValue string
		want      string
	}{
		{
			name:      "array",
			resetter:  "Array",
			fieldName: "vector",
			zeroValue: "0",
			want:      "for i := 0; i < len(a.vector); i++ {\n\t*a.vector[i] = 0\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenArrayPResetCommand(tt.fieldName, tt.zeroValue)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenSliceResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		want      string
	}{
		{
			name:      "slice",
			resetter:  "Slice",
			fieldName: "Items",
			want:      "if s.Items != nil {\n\ts.Items = s.Items[:0]\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenSliceResetCommand(tt.fieldName)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenMapResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		want      string
	}{
		{
			name:      "map",
			resetter:  "Map",
			fieldName: "Data",
			want:      "if m.Data != nil {\n\tclear(m.Data)\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenMapResetCommand(tt.fieldName)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenPtrResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		zeroValue string
		want      string
	}{
		{
			name:      "int zero value",
			resetter:  "User",
			fieldName: "Point",
			zeroValue: "0",
			want:      "if u.Point != nil {\n\t*u.Point = 0\n}",
		},
		{
			name:      "string zero value",
			resetter:  "PointerStruct",
			fieldName: "Data",
			zeroValue: `""`,
			want:      "if p.Data != nil {\n\t*p.Data = \"\"\n}",
		},
		{
			name:      "boolean zero value",
			resetter:  "Config",
			fieldName: "Enabled",
			zeroValue: "false",
			want:      "if c.Enabled != nil {\n\t*c.Enabled = false\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenPtrResetCommand(tt.fieldName, tt.zeroValue)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenStructResetCommand(t *testing.T) {
	tests := []struct {
		name      string
		resetter  string
		fieldName string
		zeroValue string
		want      string
	}{
		{
			name:      "struct",
			resetter:  "MyStruct",
			fieldName: "Nested",
			zeroValue: "Value{}",
			want:      "if r, ok := m.Nested.(interface{ Reset() }); ok {\n\tr.Reset()\n} else {\n\t*m.Nested = Value{}\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.resetter)
			result := gen.GenStructResetCommand(tt.fieldName, tt.zeroValue)

			assert.Equal(t, tt.want, result)
		})
	}
}
