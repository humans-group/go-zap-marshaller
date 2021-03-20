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

	secured := secured(fieldDef.AllTags)

	keyName := "keyName"

	writeDef(fileContent, encMethodAdd, fieldDef.FieldType, fieldDef.FieldName, zapFieldName, keyName, secured, false)
}

func writeDef(fileContent *bytes.Buffer, method encMethod, fieldType astparser.Type, fieldName, zapFieldName, keyName string, secured, recursive bool) {
	if !recursive {
		fileContent.WriteString(fmt.Sprintf("\nkeyName = \"%s\"\n", zapFieldName))
	}

	switch t := fieldType.(type) {
	case astparser.TypeInterfaceValue:
		writeInterfaceDef(fileContent, method, recursiveFieldName(recursive, fieldName), keyName, secured)
	case astparser.TypeArray:
		writeArrayDef(fileContent, fieldName, zapFieldName, t, secured)
	case astparser.TypeMap:
		writeMapDef(fileContent, fieldName, zapFieldName, t, secured)
	case astparser.TypePointer:
		writePointerDef(fileContent, recursiveFieldName(recursive, fieldName), keyName, zapFieldName, encMethodAdd, t, secured)
	case astparser.TypeCustom:
		writeCustomDef(t, fileContent, method, recursiveFieldName(recursive, fieldName), keyName, zapFieldName, secured)
	case astparser.TypeSimple:
		writeSimpleDef(fileContent, t, method, recursiveFieldName(recursive, fieldName), keyName, secured)
	default:
		panic(fmt.Sprintf("unknown field name %s type %T", fieldName, t))
	}
}

func recursiveFieldName(recursive bool, fieldName string) string {
	if recursive {
		return fieldName
	}

	return "m." + fieldName
}

const securedTagName = "secured"

func secured(tags map[string]string) bool {
	for k, _ := range tags {
		if k == securedTagName {
			return true
		}
	}
	return false
}

func writeMapDef(fileContent *bytes.Buffer, fieldName, zapFieldName string, t astparser.TypeMap, secured bool) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.%s {
	`, fieldName))

	switch tt := t.ValueType.(type) {
	case astparser.TypeInterfaceValue:
		writeInterfaceDef(fileContent, encMethodAdd, "value", "key", secured)
	case astparser.TypeCustom:
		writeCustomDef(tt, fileContent, encMethodAdd, "value", "key", zapFieldName, secured)
	case astparser.TypePointer:
		writePointerDef(fileContent, "value", "key", zapFieldName, encMethodAdd, tt, secured)

	case astparser.TypeSimple:
		writeSimpleDef(fileContent, tt, encMethodAdd, "value", "key", secured)

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, fieldName))
	}

	fileContent.WriteString("}\n")
	fileContent.WriteString("return nil\n")
	fileContent.WriteString("}))\n\n")
}

func writeArrayDef(fileContent *bytes.Buffer, fieldName, zapFieldName string, t astparser.TypeArray, secured bool) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	if s, ok := t.InnerType.(astparser.TypeSimple); ok && s.Name == "byte" {
		fileContent.WriteString(fmt.Sprintf("enc.AddByteString(keyName, m.%s)\n", fieldName))
		return
	}

	fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.%s {
	`, fieldName))

	switch tt := t.InnerType.(type) {
	case astparser.TypeInterfaceValue:
		writeInterfaceDef(fileContent, encMethodAppend, "value", "keyName", secured)
	case astparser.TypeCustom:
		writeCustomDef(tt, fileContent, encMethodAppend, "value", "keyName", zapFieldName, secured)
	case astparser.TypePointer:
		writePointerDef(fileContent, "value", "key", zapFieldName, encMethodAppend, tt, secured)

	case astparser.TypeSimple:
		writeSimpleDef(fileContent, tt, encMethodAppend, "value", "keyName", secured)

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, fieldName))
	}

	fileContent.WriteString("}\n")
	fileContent.WriteString("return nil\n")
	fileContent.WriteString("}))\n\n")
}

func writePointerDef(
	fileContent *bytes.Buffer,
	fieldName, keyName, zapFieldName string,
	method encMethod,
	t astparser.TypePointer,
	secured bool,
) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	switch method {
	case encMethodAdd:
		fileContent.WriteString(fmt.Sprintf("if %s != nil {\n", fieldName))
	case encMethodAppend:
		fileContent.WriteString(fmt.Sprintf(`if %s == nil {
				continue
			}
		`, fieldName))
	default:
		panic(fmt.Sprintf("unexpected enc method %s", method))
	}

	switch ttt := t.InnerType.(type) {
	case astparser.TypeSimple:
		writeSimpleDef(fileContent, ttt, method, "*"+fieldName, keyName, secured)

	case astparser.TypeCustom:
		writeCustomDef(ttt, fileContent, method, fieldName, keyName, zapFieldName, secured)
	default:
		panic(fmt.Sprintf("unsupported array pointer innter type %T field name %s", ttt, fieldName))
	}

	if method == encMethodAdd {
		fileContent.WriteString("}\n")
	}
}

func writeCustomDef(
	t astparser.TypeCustom,
	fileContent *bytes.Buffer,
	method encMethod,
	fieldName, keyName string,
	zapFieldName string,
	secured bool,
) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	// we cant resolve alias type, so lets just stringify it
	if t.Alias {
		fileContent.WriteString(fmt.Sprintf("enc.AddString(keyName, %s)\n", fieldName))
		return
	}

	if t.AliasType != nil {
		writeDef(fileContent, method, t.AliasType, fieldName, zapFieldName, keyName, secured, true)
		return
	}

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

func writeInterfaceDef(
	fileContent *bytes.Buffer,
	methodName encMethod,
	fieldName, keyName string,
	secured bool,
) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	switch methodName {
	case encMethodAdd:
		fileContent.WriteString(fmt.Sprintf("_ = enc.AddReflected(%s, %s)\n", keyName, fieldName))
	case encMethodAppend:
		fileContent.WriteString(fmt.Sprintf("_ = aenc.AppendReflected(%s)\n", fieldName))
	default:
		panic(fmt.Sprintf("unknown method name %s", methodName))
	}
}

func writeSimpleDef(
	fileContent *bytes.Buffer,
	simpleType astparser.TypeSimple,
	methodName encMethod,
	fieldName, keyName string,
	secured bool,
) {
	if secured {
		fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

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
		fileContent.WriteString(fmt.Sprintf("enc.Add%s(%s, %s(%s))\n", zapMethodType, keyName, simpleType.Name, fieldName))
	case encMethodAppend:
		fileContent.WriteString(fmt.Sprintf("aenc.Append%s(%s(%s))\n", zapMethodType, simpleType.Name, fieldName))
	default:
		panic(fmt.Sprintf("unknown method name %s", methodName))
	}
}
