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
	enc.AddInt(keyName, int(m.Int))

	keyName = "Int8"
	enc.AddInt8(keyName, int8(m.Int8))

	keyName = "Int16"
	enc.AddInt16(keyName, int16(m.Int16))

	keyName = "Int32"
	enc.AddInt32(keyName, int32(m.Int32))

	keyName = "Int64"
	enc.AddInt64(keyName, int64(m.Int64))

	keyName = "Uint"
	enc.AddUint(keyName, uint(m.Uint))

	keyName = "Uint8"
	enc.AddUint8(keyName, uint8(m.Uint8))

	keyName = "Uint16"
	enc.AddUint16(keyName, uint16(m.Uint16))

	keyName = "Uint32"
	enc.AddUint32(keyName, uint32(m.Uint32))

	keyName = "Uint64"
	enc.AddUint64(keyName, uint64(m.Uint64))

	keyName = "Float32"
	enc.AddFloat32(keyName, float32(m.Float32))

	keyName = "Float64"
	enc.AddFloat64(keyName, float64(m.Float64))

	keyName = "Bool"
	enc.AddBool(keyName, bool(m.Bool))

	keyName = "String"
	enc.AddString(keyName, string(m.String))

	keyName = "Map"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.Map {
			enc.AddString(key, string(value))
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
			aenc.AppendInt(int(value))
		}
		return nil
	}))

	keyName = "MapWithNulls"
	_ = enc.AddObject(keyName, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for key, value := range m.MapWithNulls {
			if value != nil {
				enc.AddString(key, string(*value))
			}
		}
		return nil
	}))

	keyName = "Ptr"
	if m.Ptr != nil {
		enc.AddInt(keyName, int(*m.Ptr))
	}

	keyName = "Secured"
	enc.AddString(keyName, "<secured>")
	return nil
}
