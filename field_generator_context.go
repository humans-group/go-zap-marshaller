package zapmarshaller

import (
	"github.com/mkorolyov/astparser"
)

type fieldGenContext struct {
	encoderMethod encMethod
	zapFieldName  string
	fieldName     string
	keyName       string
	secured       bool
	recursiveCall bool
	fieldType     astparser.Type
}

func (c fieldGenContext) withFieldName(fn string) fieldGenContext {
	c.fieldName = fn
	return c
}

func (c fieldGenContext) withEncoderMethod(m encMethod) fieldGenContext {
	c.encoderMethod = m
	return c
}

func (c fieldGenContext) withKeyName(kn string) fieldGenContext {
	c.keyName = kn
	return c
}

func (c fieldGenContext) withFieldType(ft astparser.Type) fieldGenContext {
	c.fieldType = ft
	return c
}

func (c fieldGenContext) recursive() fieldGenContext {
	c.recursiveCall = true
	return c
}

func (c fieldGenContext) withNestedFieldName() fieldGenContext {
	if c.recursiveCall {
		return c
	}

	c.fieldName = "m." + c.fieldName
	return c
}
