package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/MaksimMakarenko1001/ya-go-advanced/cmd/reset/model"
)

const resetFunc = `func ({{.Receiver}} *{{.TypeName}}) Reset() {
	if {{.Receiver}} == nil { return }
{{range $fieldName, $command := .Commands}}
	// resets {{$fieldName}}
	{{$command}}
{{end}}}
`

// templateData holds the data for the reset function template
type templateData struct {
	TypeName string
	Receiver string
	Commands model.ResetCommands
}

type Generator struct {
	resetter model.StructName
}

func New(resetter model.StructName) Generator {
	return Generator{
		resetter: resetter,
	}
}

func (g Generator) receiver() string {
	if len(g.resetter) == 0 {
		return ""
	}

	return strings.ToLower(string(g.resetter[0]))
}

// GenerateResetFunction generates a Reset() method for a struct based on the provided reset rules
func (g Generator) GenResetFunc(commands model.ResetCommands) (model.ResetFunc, error) {
	if len(g.resetter) == 0 {
		return "", nil
	}

	data := templateData{
		TypeName: g.resetter,
		Receiver: g.receiver(),
		Commands: commands,
	}

	tmpl, err := template.New("reset").Parse(resetFunc)
	if err != nil {
		return "", fmt.Errorf("parse template error: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("execute template error: %w", err)
	}

	return buf.String(), nil
}

func (g Generator) GenBasicResetCommand(name model.FieldName, zeroValue string) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%[1]s.%[2]s = %[3]s",
		g.receiver(), name, zeroValue,
	)
}

func (g Generator) GenArrayResetCommand(name model.FieldName, zeroValue string) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"for i := 0; i < len(%[1]s.%[2]s); i++ {\n\t%[1]s.%[2]s[i] = %[3]s\n}",
		g.receiver(), name, zeroValue,
	)
}

func (g Generator) GenArrayPResetCommand(name model.FieldName, zeroValue string) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"for i := 0; i < len(%[1]s.%[2]s); i++ {\n\t*%[1]s.%[2]s[i] = %[3]s\n}",
		g.receiver(), name, zeroValue,
	)
}

func (g Generator) GenSliceResetCommand(name model.FieldName) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"if %[1]s.%[2]s != nil {\n\t%[1]s.%[2]s = %[1]s.%[2]s[:0]\n}",
		g.receiver(), name,
	)
}

func (g Generator) GenMapResetCommand(name model.FieldName) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"if %[1]s.%[2]s != nil {\n\tclear(%[1]s.%[2]s)\n}",
		g.receiver(), name,
	)
}

func (g Generator) GenPtrResetCommand(name model.FieldName, zeroValue string) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"if %[1]s.%[2]s != nil {\n\t*%[1]s.%[2]s = %[3]s\n}",
		g.receiver(), name, zeroValue,
	)
}

func (g Generator) GenStructResetCommand(name model.FieldName, zeroValue string) model.Command {
	if len(g.resetter) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"if r, ok := %[1]s.%[2]s.(interface{ Reset() }); ok {\n\tr.Reset()\n} else {\n\t*%[1]s.%[2]s = %[3]s\n}",
		g.receiver(), name, zeroValue,
	)
}
