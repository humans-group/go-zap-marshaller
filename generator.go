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

type fieldGenerator struct {
	fieldDef    astparser.FieldDef
	fileContent *bytes.Buffer
}

type fieldGenContext struct {
	encoderMethod encMethod
	zapFieldName  string
	fieldName     string
	keyName       string
	secured       bool
	recursiveCall bool
	fieldType     astparser.Type
}

func (g *fieldGenerator) writeDef(ctx fieldGenContext) {
	if !ctx.recursiveCall {
		g.fileContent.WriteString(fmt.Sprintf("\nkeyName = \"%s\"\n", ctx.zapFieldName))
	}

	switch t := ctx.fieldType.(type) {
	case astparser.TypeInterfaceValue:
		ctx.fieldName = recursiveFieldName(ctx.recursiveCall, ctx.fieldName)
		g.writeInterfaceDef(ctx)
	case astparser.TypeArray:
		g.writeArrayDef(ctx)
	case astparser.TypeMap:
		g.writeMapDef(ctx)
	case astparser.TypePointer:
		ctx.fieldName = recursiveFieldName(ctx.recursiveCall, ctx.fieldName)
		g.writePointerDef(ctx)
	case astparser.TypeCustom:
		ctx.fieldName = recursiveFieldName(ctx.recursiveCall, ctx.fieldName)
		g.writeCustomDef(ctx)
	case astparser.TypeSimple:
		ctx.fieldName = recursiveFieldName(ctx.recursiveCall, ctx.fieldName)
		g.writeSimpleDef(ctx)
	default:
		panic(fmt.Sprintf("unknown field name %s type %T", ctx.fieldName, t))
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

func (g *fieldGenerator) writeMapDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	g.fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.%s {
	`, ctx.fieldName))

	t := ctx.fieldType.(astparser.TypeMap)

	switch tt := t.ValueType.(type) {
	case astparser.TypeInterfaceValue:
		ctx.fieldName = "value"
		ctx.keyName = "key"
		g.writeInterfaceDef(ctx)
	case astparser.TypeCustom:
		ctx.fieldName = "value"
		ctx.keyName = "key"
		ctx.fieldType = tt
		g.writeCustomDef(ctx)
	case astparser.TypePointer:
		ctx.fieldName = "value"
		ctx.keyName = "key"
		ctx.fieldType = tt
		g.writePointerDef(ctx)

	case astparser.TypeSimple:
		ctx.fieldName = "value"
		ctx.keyName = "key"
		ctx.fieldType = tt
		g.writeSimpleDef(ctx)

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, ctx.fieldName))
	}

	g.fileContent.WriteString("}\n")
	g.fileContent.WriteString("return nil\n")
	g.fileContent.WriteString("}))\n\n")
}

func (g *fieldGenerator) writeArrayDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	t := ctx.fieldType.(astparser.TypeArray)

	if s, ok := t.InnerType.(astparser.TypeSimple); ok && s.Name == "byte" {
		g.fileContent.WriteString(fmt.Sprintf("enc.AddByteString(keyName, m.%s)\n", ctx.fieldName))
		return
	}

	g.fileContent.WriteString(fmt.Sprintf(
		`_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.%s {
	`, ctx.fieldName))

	switch tt := t.InnerType.(type) {
	case astparser.TypeInterfaceValue:
		ctx.fieldName = "value"
		ctx.keyName = "keyName"
		ctx.encoderMethod = encMethodAppend
		g.writeInterfaceDef(ctx)
	case astparser.TypeCustom:
		ctx.fieldName = "value"
		ctx.keyName = "keyName"
		ctx.fieldType = tt
		ctx.encoderMethod = encMethodAppend
		g.writeCustomDef(ctx)
	case astparser.TypePointer:
		ctx.fieldName = "value"
		ctx.keyName = "key"
		ctx.fieldType = tt
		ctx.encoderMethod = encMethodAppend
		g.writePointerDef(ctx)

	case astparser.TypeSimple:
		ctx.fieldName = "value"
		ctx.keyName = "keyName"
		ctx.fieldType = tt
		ctx.encoderMethod = encMethodAppend
		g.writeSimpleDef(ctx)

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, ctx.fieldName))
	}

	g.fileContent.WriteString("}\n")
	g.fileContent.WriteString("return nil\n")
	g.fileContent.WriteString("}))\n\n")
}

func (g *fieldGenerator) writePointerDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	switch ctx.encoderMethod {
	case encMethodAdd:
		g.fileContent.WriteString(fmt.Sprintf("if %s != nil {\n", ctx.fieldName))
	case encMethodAppend:
		g.fileContent.WriteString(fmt.Sprintf(`if %s == nil {
				continue
			}
		`, ctx.fieldName))
	default:
		panic(fmt.Sprintf("unexpected enc method %s", ctx.encoderMethod))
	}

	t := ctx.fieldType.(astparser.TypePointer)

	switch ttt := t.InnerType.(type) {
	case astparser.TypeSimple:
		ctx.fieldName = "*" + ctx.fieldName
		ctx.fieldType = ttt
		g.writeSimpleDef(ctx)

	case astparser.TypeCustom:
		ctx.fieldType = ttt
		g.writeCustomDef(ctx)
	default:
		panic(fmt.Sprintf("unsupported array pointer innter type %T field name %s", ttt, ctx.fieldName))
	}

	if ctx.encoderMethod == encMethodAdd {
		g.fileContent.WriteString("}\n")
	}
}

func (g *fieldGenerator) writeCustomDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	t := ctx.fieldType.(astparser.TypeCustom)

	// we cant resolve alias type, so lets just stringify it
	if t.Alias {
		g.fileContent.WriteString(fmt.Sprintf("enc.AddString(keyName, %s)\n", ctx.fieldName))
		return
	}

	if t.AliasType != nil {
		ctx.recursiveCall = true
		ctx.fieldType = t.AliasType
		g.writeDef(ctx)
		return
	}

	switch ctx.encoderMethod {
	case encMethodAppend:
		g.fileContent.WriteString(
			`vv = value
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = aenc.AppendObject(marshaler)
			}
		`)

	case encMethodAdd:
		g.fileContent.WriteString(fmt.Sprintf(
			`vv = %s
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = enc.AddObject(%s, marshaler)
			}
		`, ctx.fieldName, ctx.keyName))

	default:
		panic(fmt.Sprintf("unexpected enc method %s", ctx.encoderMethod))
	}
}

func (g *fieldGenerator) writeInterfaceDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	switch ctx.encoderMethod {
	case encMethodAdd:
		g.fileContent.WriteString(fmt.Sprintf("_ = enc.AddReflected(%s, %s)\n", ctx.keyName, ctx.fieldName))
	case encMethodAppend:
		g.fileContent.WriteString(fmt.Sprintf("_ = aenc.AppendReflected(%s)\n", ctx.fieldName))
	default:
		panic(fmt.Sprintf("unknown method name %s", ctx.encoderMethod))
	}
}

func (g *fieldGenerator) writeSimpleDef(ctx fieldGenContext) {
	if ctx.secured {
		g.fileContent.WriteString("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	simpleType := ctx.fieldType.(astparser.TypeSimple)

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
		panic(fmt.Sprintf("unknown simple type %s fieldName %s\n", simpleType.Name, ctx.fieldName))
	}

	switch ctx.encoderMethod {
	case encMethodAdd:
		g.fileContent.WriteString(fmt.Sprintf(
			"enc.Add%s(%s, %s(%s))\n",
			zapMethodType,
			ctx.keyName,
			simpleType.Name,
			ctx.fieldName,
		))
	case encMethodAppend:
		g.fileContent.WriteString(fmt.Sprintf(
			"aenc.Append%s(%s(%s))\n",
			zapMethodType,
			simpleType.Name,
			ctx.fieldName,
		))
	default:
		panic(fmt.Sprintf("unknown method name %s", ctx.encoderMethod))
	}
}
