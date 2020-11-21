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

type Dep struct {
	Int        int
	Dep1       Dep1
	Dep2Opt    *Dep2
	Dep3Array  []Dep3
	Dep4Map    map[string]Dep4
	DepWithDep Dep6
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
