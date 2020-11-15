package fixtures_test

type Primitives struct {
	// comment here
	ByteSlice    []byte
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Float32      float32
	Float64      float64
	Bool         bool
	String       string
	Map          map[string]string
	Slice        []int
	MapWithNulls map[string]*string
	Ptr          *int
}
