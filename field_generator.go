package zapmarshaller

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mkorolyov/astparser"
)

type fieldGenerator struct {
	fieldDef    astparser.FieldDef
	fileContent *bytes.Buffer
}

func (g *fieldGenerator) write(str string) {
	g.fileContent.WriteString(str)
}

func (g *fieldGenerator) writef(str string, args ...interface{}) {
	g.write(fmt.Sprintf(str, args...))
}

func (g *fieldGenerator) writeDef(ctx fieldGenContext) {
	if !ctx.recursiveCall {
		g.writef("\nkeyName = \"%s\"\n", ctx.zapFieldName)
	}

	switch t := ctx.fieldType.(type) {
	case astparser.TypeArray:
		g.secure(ctx, g.writeArrayDef)
	case astparser.TypeMap:
		g.secure(ctx, g.writeMapDef)
	case astparser.TypeInterfaceValue:
		g.secure(ctx.withNestedFieldName(), g.writeInterfaceDef)
	case astparser.TypePointer:
		g.secure(ctx.withNestedFieldName(), g.writePointerDef)
	case astparser.TypeCustom:
		g.secure(ctx.withNestedFieldName(), g.writeCustomDef)
	case astparser.TypeSimple:
		g.secure(ctx.withNestedFieldName(), g.writeSimpleDef)
	default:
		panic(fmt.Sprintf("unknown field name %s type %T", ctx.fieldName, t))
	}
}

func (g *fieldGenerator) writeMapDef(ctx fieldGenContext) {
	g.writef(
		`_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.%s {
	`, ctx.fieldName)

	t := ctx.fieldType.(astparser.TypeMap)

	ctx = ctx.
		withFieldName("value").
		withKeyName("key")

	switch tt := t.ValueType.(type) {
	case astparser.TypeInterfaceValue:
		g.writeInterfaceDef(ctx)
	case astparser.TypeCustom:
		g.writeCustomDef(ctx.withFieldType(tt))
	case astparser.TypePointer:
		g.writePointerDef(ctx.withFieldType(tt))
	case astparser.TypeSimple:
		g.writeSimpleDef(ctx.withFieldType(tt))

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, ctx.fieldName))
	}

	g.write("}\n")
	g.write("return nil\n")
	g.write("}))\n\n")
}

func (g *fieldGenerator) writeArrayDef(ctx fieldGenContext) {
	t := ctx.fieldType.(astparser.TypeArray)

	if s, ok := t.InnerType.(astparser.TypeSimple); ok && s.Name == "byte" {
		g.writef("enc.AddByteString(keyName, m.%s)\n", ctx.fieldName)
		return
	}

	g.writef(
		`_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.%s {
	`, ctx.fieldName)

	ctx = ctx.
		withEncoderMethod(encMethodAppend).
		withFieldName("value")

	switch tt := t.InnerType.(type) {
	case astparser.TypeInterfaceValue:
		g.writeInterfaceDef(ctx)
	case astparser.TypeCustom:
		g.writeCustomDef(ctx.withFieldType(tt))
	case astparser.TypePointer:
		g.writePointerDef(ctx.withFieldType(tt).withKeyName("key"))

	case astparser.TypeSimple:
		g.writeSimpleDef(ctx.withFieldType(tt))

	default:
		panic(fmt.Sprintf("unsupported array inner type %T in field %s", tt, ctx.fieldName))
	}

	g.write("}\n")
	g.write("return nil\n")
	g.write("}))\n\n")
}

func (g *fieldGenerator) writePointerDef(ctx fieldGenContext) {
	switch ctx.encoderMethod {
	case encMethodAdd:
		g.writef("if %s != nil {\n", ctx.fieldName)
	case encMethodAppend:
		g.writef(`if %s == nil {
				continue
			}
		`, ctx.fieldName)
	default:
		panic(fmt.Sprintf("unexpected enc method %s", ctx.encoderMethod))
	}

	t := ctx.fieldType.(astparser.TypePointer)
	ctx = ctx.withFieldType(t.InnerType)

	switch tt := t.InnerType.(type) {
	case astparser.TypeSimple:
		g.writeSimpleDef(ctx.withFieldName("*" + ctx.fieldName))
	case astparser.TypeCustom:
		if tt.AliasType != nil {
			ctx = ctx.withFieldName("*" + ctx.fieldName)
		}
		g.writeCustomDef(ctx)
	default:
		panic(fmt.Sprintf("unsupported array pointer innter type %T field name %s", tt, ctx.fieldName))
	}

	if ctx.encoderMethod == encMethodAdd {
		g.write("}\n")
	}
}

func (g *fieldGenerator) writeCustomDef(ctx fieldGenContext) {
	t := ctx.fieldType.(astparser.TypeCustom)

	if t.Name == "error" {
		switch ctx.encoderMethod {
		case encMethodAppend:
			g.write(
				`if value != nil {
				aenc.AppendString(value.Error())
			}`,
			)

		case encMethodAdd:
			g.writef(
				`if %s != nil {
				enc.AddString(%s, %[1]s.Error())
			}
`, ctx.fieldName, ctx.keyName,
			)

		default:
			panic(fmt.Sprintf("unexpected enc method %s", ctx.encoderMethod))
		}

		return
	}

	// we cant resolve alias type, so lets just stringify it
	if t.Alias {
		g.writef("enc.AddString(keyName, %s)\n", ctx.fieldName)
		return
	}

	if t.AliasType != nil {
		g.writeDef(ctx.withFieldType(t.AliasType).recursive())
		return
	}

	switch ctx.encoderMethod {
	case encMethodAppend:
		g.write(
			`vv = value
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = aenc.AppendObject(marshaler)
			}
		`)

	case encMethodAdd:
		g.writef(
			`vv = %s
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = enc.AddObject(%s, marshaler)
			}
		`, ctx.fieldName, ctx.keyName)

	default:
		panic(fmt.Sprintf("unexpected enc method %s", ctx.encoderMethod))
	}
}

func (g *fieldGenerator) writeInterfaceDef(ctx fieldGenContext) {
	switch ctx.encoderMethod {
	case encMethodAdd:
		g.writef("_ = enc.AddReflected(%s, %s)\n", ctx.keyName, ctx.fieldName)
	case encMethodAppend:
		g.writef("_ = aenc.AppendReflected(%s)\n", ctx.fieldName)
	default:
		panic(fmt.Sprintf("unknown method name %s", ctx.encoderMethod))
	}
}

func (g *fieldGenerator) writeSimpleDef(ctx fieldGenContext) {
	t := ctx.fieldType.(astparser.TypeSimple)

	var zapMethodType string
	switch t.Name {
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
		zapMethodType = strings.ToTitle(t.Name[:1]) + t.Name[1:]
	default:
		panic(fmt.Sprintf("unknown simple type %s fieldName %s\n", t.Name, ctx.fieldName))
	}

	switch ctx.encoderMethod {
	case encMethodAdd:
		g.writef(
			"enc.Add%s(%s, %s(%s))\n",
			zapMethodType,
			ctx.keyName,
			t.Name,
			ctx.fieldName,
		)
	case encMethodAppend:
		g.writef(
			"aenc.Append%s(%s(%s))\n",
			zapMethodType,
			t.Name,
			ctx.fieldName,
		)
	default:
		panic(fmt.Sprintf("unknown method name %s", ctx.encoderMethod))
	}
}

func (g *fieldGenerator) secure(ctx fieldGenContext, plain func(ctx fieldGenContext)) {
	if ctx.secured {
		g.write("enc.AddString(keyName, \"<secured>\")\n")
		return
	}

	plain(ctx)
}
