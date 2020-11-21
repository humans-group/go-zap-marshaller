package fixtures_test

import (
	"go.uber.org/zap/zapcore"
)

func (m *Primitives) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "ByteSlice"
	enc.AddByteString(keyName, m.ByteSlice)

	keyName = "Int"
	enc.AddInt(keyName, m.Int)

	keyName = "Int8"
	enc.AddInt8(keyName, m.Int8)

	keyName = "Int16"
	enc.AddInt16(keyName, m.Int16)

	keyName = "Int32"
	enc.AddInt32(keyName, m.Int32)

	keyName = "Int64"
	enc.AddInt64(keyName, m.Int64)

	keyName = "Uint"
	enc.AddUint(keyName, m.Uint)

	keyName = "Uint8"
	enc.AddUint8(keyName, m.Uint8)

	keyName = "Uint16"
	enc.AddUint16(keyName, m.Uint16)

	keyName = "Uint32"
	enc.AddUint32(keyName, m.Uint32)

	keyName = "Uint64"
	enc.AddUint64(keyName, m.Uint64)

	keyName = "Float32"
	enc.AddFloat32(keyName, m.Float32)

	keyName = "Float64"
	enc.AddFloat64(keyName, m.Float64)

	keyName = "Bool"
	enc.AddBool(keyName, m.Bool)

	keyName = "String"
	enc.AddString(keyName, m.String)

	keyName = "Map"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.Map {
			enc.AddString(key, value)
		}
		return nil
	}))

	keyName = "MapInterface"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.MapInterface {
			_ = enc.AddReflected(key, value)
		}
		return nil
	}))

	keyName = "InterfaceSlice"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.InterfaceSlice {
			_ = aenc.AppendReflected(value)
		}
		return nil
	}))

	keyName = "Interface"
	_ = enc.AddReflected(keyName, m.Interface)

	keyName = "Slice"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.Slice {
			aenc.AppendInt(value)
		}
		return nil
	}))

	keyName = "MapWithNulls"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.MapWithNulls {
			if value != nil {
				enc.AddString(key, *value)
			}
		}
		return nil
	}))

	keyName = "Ptr"
	if m.Ptr != nil {
		enc.AddInt(keyName, *m.Ptr)
	}

	keyName = "Secured"
	enc.AddString(keyName, "<secured>")
	return nil
}
