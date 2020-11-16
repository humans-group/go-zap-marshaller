package zapmarshaller

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

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
				zapFieldName := fieldDef.FieldName
				if fieldDef.JsonName != "" {
					zapFieldName = fieldDef.JsonName
				}

				writeDef(fileContent, fieldDef.FieldType, fieldDef.FieldName, zapFieldName)
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

func writeDef(fileContent *bytes.Buffer, fieldType astparser.Type, fieldName, zapFieldName string) {
	fileContent.WriteString(fmt.Sprintf("\nkeyName = \"%s\"\n", zapFieldName))

	switch t := fieldType.(type) {
	case astparser.TypeArray:
		writeArrayDef(fileContent, fieldName, t)
	case astparser.TypeMap:
		writeMapDef(fileContent, fieldName, t)
	case astparser.TypePointer:
		writePointerDef(fileContent, fmt.Sprintf("m.%s", fieldName), "keyName", encMethodAdd, t)
	case astparser.TypeCustom:
		writeCustomDef(fileContent, encMethodAdd, fmt.Sprintf("m.%s", fieldName), "keyName")
	case astparser.TypeSimple:
		writeSimpleDef(fileContent, t, encMethodAdd, fmt.Sprintf("m.%s", fieldName), "keyName")
	default:
		panic(fmt.Sprintf("unknown field name %s type %T", fieldName, t))
	}
}

func writeMapDef(fileContent *bytes.Buffer, fieldName string, t astparser.TypeMap) {
	fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.%s {
	`, fieldName))

	switch tt := t.ValueType.(type) {
	case astparser.TypeCustom:
		writeCustomDef(fileContent, encMethodAdd, "value", "key")
	case astparser.TypePointer:
		writePointerDef(fileContent, "value", "key", encMethodAdd, tt)

	case astparser.TypeSimple:
		writeSimpleDef(fileContent, tt, encMethodAdd, "value", "key")

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, fieldName))
	}

	fileContent.WriteString("}\n")
	fileContent.WriteString("return nil\n")
	fileContent.WriteString("}))\n\n")
}

func writeArrayDef(fileContent *bytes.Buffer, fieldName string, t astparser.TypeArray) {
	if s, ok := t.InnerType.(astparser.TypeSimple); ok && s.Name == "byte" {
		fileContent.WriteString(fmt.Sprintf("enc.AddByteString(keyName, m.%s)\n", fieldName))
		return
	}

	fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.%s {
	`, fieldName))

	switch tt := t.InnerType.(type) {
	case astparser.TypeCustom:
		writeCustomDef(fileContent, encMethodAppend, "value", "keyName")
	case astparser.TypePointer:
		writePointerDef(fileContent, "value", "key", encMethodAppend, tt)

	case astparser.TypeSimple:
		writeSimpleDef(fileContent, tt, encMethodAppend, "value", "keyName")

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, fieldName))
	}

	fileContent.WriteString("}\n")
	fileContent.WriteString("return nil\n")
	fileContent.WriteString("}))\n\n")
}

func writePointerDef(fileContent *bytes.Buffer, fieldName, keyName string, method encMethod, t astparser.TypePointer) {
	switch method {
	case encMethodAdd:
		fileContent.WriteString(fmt.Sprintf("if %s != nil {\n", fieldName))
	case encMethodAppend:
		fileContent.WriteString(fmt.Sprintf(`if m.%s == nil {
				continue
			}
		`, fieldName))
	default:
		panic(fmt.Sprintf("unexpected enc method %s", method))
	}

	switch ttt := t.InnerType.(type) {
	case astparser.TypeSimple:
		writeSimpleDef(fileContent, ttt, method, "*"+fieldName, keyName)

	case astparser.TypeCustom:
		writeCustomDef(fileContent, method, "*"+fieldName, keyName)
	default:
		panic(fmt.Sprintf("unsupported array pointer innter type %T field name %s", ttt, fieldName))
	}

	if method == encMethodAdd {
		fileContent.WriteString("}\n")
	}
}

func writeCustomDef(fileContent *bytes.Buffer, method encMethod, fieldName, keyName string) {
	switch method {
	case encMethodAppend:
		fileContent.WriteString(
			`vv = value
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = aenc.AppendObject(marshaler)
			}
		`)

	case encMethodAdd:
		fileContent.WriteString(fmt.Sprintf(
			`vv = %s
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = enc.AddObject(%s, marshaler)
			}
		`, fieldName, keyName))

	default:
		panic(fmt.Sprintf("unexpected enc method %s", method))
	}
}

func writeSimpleDef(
	fileContent *bytes.Buffer,
	simpleType astparser.TypeSimple,
	methodName encMethod,
	fieldName, keyName string,
) {
	var zapMethodType string
	switch simpleType.Name {
	case "bool",
		"complex128",
		"complex64",
		"int",
		"int64",
		"int32",
		"int16",
		"int8",
		"float64",
		"float32",
		"string",
		"uint",
		"uint64",
		"uint32",
		"uint16",
		"uint8",
		"uintptr":
		zapMethodType = strings.ToTitle(simpleType.Name[:1]) + simpleType.Name[1:]
	default:
		panic(fmt.Sprintf("unknown simple type %s fieldName %s\n", simpleType.Name, fieldName))
	}

	switch methodName {
	case encMethodAdd:
		fileContent.WriteString(fmt.Sprintf("enc.Add%s(%s, %s)\n", zapMethodType, keyName, fieldName))
	case encMethodAppend:
		fileContent.WriteString(fmt.Sprintf("aenc.Append%s(%s)\n", zapMethodType, fieldName))
	default:
		panic(fmt.Sprintf("unknown method name %s", methodName))
	}
}
