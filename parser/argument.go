package parser

type Value interface {
	Type() Type
	Value() interface{}
}

type Argument struct {
	t      Type
	Values []Value
}

func (a *Argument) Type() Type {
	return a.t
}

func (a *Argument) String() string {
	return a.t.String()
}

func NewArgument(values ...Value) *Argument {
	types := make([]Type, len(values))
	for i, v := range values {
		types[i] = v.Type()
	}
	return &Argument{
		t:      NewTupleType(types...),
		Values: values,
	}
}

type Tuple []Value

func (t *Tuple) Type() Type {
	typs := make([]Type, len(*t))
	for i, v := range *t {
		typs[i] = v.Type()
	}
	return NewTupleType(typs...)
}

func (t *Tuple) Value() interface{} {
	return *t
}

func (t *Tuple) AddValue(v Value) {
	*t = append(*t, v)
}

func (t *Tuple) Argument() *Argument {
	return &Argument{t: t.Type(), Values: *t}
}
