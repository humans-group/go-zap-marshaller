package fixtures_test

type Dep1 struct {
	Str string
}
type Dep2 struct {
	Str string
}
type Dep3 struct {
	Str string
}
type Dep4 struct {
	Str string
}
type Dep5 struct {
	Str string
}
type Dep6 struct {
	Dep5 Dep5
}

type ConstStr string

const (
	ConstStr1 ConstStr = "str_val"
)

type ConstInt int

const (
	ConstInt1 ConstInt = 1
)

type Dep struct {
	Int           int
	Dep1          Dep1
	Dep2Opt       *Dep2
	Dep3Array     []Dep3
	Dep4Map       map[string]Dep4
	Dep4MapPtr    map[string]*Dep4
	DepWithDep    Dep6
	ConstStr      ConstStr
	ConstInt      ConstInt
	ConstStrArray []ConstStr
	ConstStrMap   map[string]ConstStr
}

type Optional struct {
	Int int
}

// skip for now in the generation
type OptionalSlice []Optional

type StructV1 struct {
	Dep1
	Dep      Dep
	Optional *Optional
}
