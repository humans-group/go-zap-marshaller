package fixtures_test

import (
	"go.uber.org/zap/zapcore"
)

func (m *Dep1) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Str"
	enc.AddString(keyName, string(m.Str))
	return nil
}

func (m *Dep2) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Str"
	enc.AddString(keyName, string(m.Str))
	return nil
}

func (m *Dep3) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Str"
	enc.AddString(keyName, string(m.Str))
	return nil
}

func (m *Dep4) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Str"
	enc.AddString(keyName, string(m.Str))
	return nil
}

func (m *Dep5) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Str"
	enc.AddString(keyName, string(m.Str))
	return nil
}

func (m *Dep6) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Dep5"
	vv = m.Dep5
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}
	return nil
}

func (m *Dep) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Int"
	enc.AddInt(keyName, int(m.Int))

	keyName = "Dep1"
	vv = m.Dep1
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Dep2Opt"
	if m.Dep2Opt != nil {
		vv = m.Dep2Opt
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}

	keyName = "Dep3Array"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.Dep3Array {
			vv = value
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = aenc.AppendObject(marshaler)
			}
		}
		return nil
	}))

	keyName = "Dep4Map"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.Dep4Map {
			vv = value
			if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
				_ = enc.AddObject(key, marshaler)
			}
		}
		return nil
	}))

	keyName = "Dep4MapPtr"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.Dep4MapPtr {
			if value != nil {
				vv = value
				if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
					_ = enc.AddObject(key, marshaler)
				}
			}
		}
		return nil
	}))

	keyName = "DepWithDep"
	vv = m.DepWithDep
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "ConstStr"
	enc.AddString(keyName, string(m.ConstStr))

	keyName = "ConstInt"
	enc.AddInt(keyName, int(m.ConstInt))

	keyName = "ConstStrArray"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.ConstStrArray {
			aenc.AppendString(string(value))
		}
		return nil
	}))

	keyName = "ConstStrMap"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.ConstStrMap {
			enc.AddString(key, string(value))
		}
		return nil
	}))

	keyName = "ConstStrPtr"
	if m.ConstStrPtr != nil {
		enc.AddString(keyName, string(*m.ConstStrPtr))
	}

	keyName = "Error"
	if m.Error != nil {
		enc.AddString(keyName, m.Error.Error())
	}

	keyName = "Errors"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.Errors {
			if value != nil {
				aenc.AppendString(value.Error())
			}
		}
		return nil
	}))

	return nil
}

func (m *Optional) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Int"
	enc.AddInt(keyName, int(m.Int))
	return nil
}

func (m *StructV1) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Dep1"
	vv = m.Dep1
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Dep"
	vv = m.Dep
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Optional"
	if m.Optional != nil {
		vv = m.Optional
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}
	return nil
}
