package zapmarshaller

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/mkorolyov/astparser"
)

type encMethod string

const encMethodAdd encMethod = "Add"
const encMethodAppend encMethod = "Append"

type Generator struct {
	Cfg Config
}

type Config struct {
	OutPackage string
}

func (g *Generator) Generate(sources map[string]astparser.ParsedFile) map[string][]byte {
	result := make(map[string][]byte, len(sources))

	for fileName, file := range sources {
		fileContent := &bytes.Buffer{}

		filePackage := file.Package
		if g.Cfg.OutPackage != "" {
			filePackage = g.Cfg.OutPackage
		}

		fileContent.WriteString(fmt.Sprintf(
			`package %s

			import (
				"go.uber.org/zap/zapcore"
			)
			`, filePackage))

		for _, structDef := range file.Structs {
			fileContent.WriteString(fmt.Sprintf(
				"func (m *%s) MarshalLogObject(enc zapcore.ObjectEncoder) error {\n", structDef.Name))

			if len(structDef.Fields) > 0 {
				fileContent.WriteString(
					`var keyName string
					var vv interface{}
					_ = vv
				`)
			}

			for _, fieldDef := range structDef.Fields {
				if !fieldDef.CompositionField {
					writeFieldDef(fieldDef, fileContent)
					continue
				}

				switch v := fieldDef.FieldType.(type) {
				case astparser.TypeCustom:
					fieldDef.FieldName = v.Name
					writeFieldDef(fieldDef, fileContent)
				default:
					panic(fmt.Sprintf("unexpected composition for struct %s field %#v", structDef.Name, fieldDef))
				}
			}

			fileContent.WriteString("return nil\n")
			fileContent.WriteString("}\n\n")
		}

		formatted, err := format.Source(fileContent.Bytes())
		if err != nil {
			println(fmt.Sprintf("gofmt error %v\nresult:\n%s", err, string(fileContent.Bytes())))
		} else {
			result[fileName] = formatted
		}
	}

	return result
}

func writeFieldDef(fieldDef astparser.FieldDef, fileContent *bytes.Buffer) {
	zapFieldName := fieldDef.FieldName
	if fieldDef.JsonName != "" {
		zapFieldName = fieldDef.JsonName
	}

	fieldGenerator := fieldGenerator{
		fieldDef:    fieldDef,
		fileContent: fileContent,
	}

	fieldGenerator.writeDef(fieldGenContext{
		encoderMethod: encMethodAdd,
		zapFieldName:  zapFieldName,
		fieldName:     fieldDef.FieldName,
		keyName:       "keyName",
		secured:       secured(fieldDef.AllTags),
		recursiveCall: false,
		fieldType:     fieldDef.FieldType,
	})
}

