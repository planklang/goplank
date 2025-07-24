package types

import "fmt"

type Value interface {
	Type() Type
	Value() any
}

type ValueContainer interface {
	AddValues(...Value)
	GetValues() []Value
	CanContain(Value) bool
}

type Tuple []Value

func (t *Tuple) Type() Type {
	typs := make([]Type, len(*t))
	for i, v := range *t {
		typs[i] = v.Type()
	}
	return NewTupleType(typs...)
}

func (t *Tuple) Value() any {
	return t.GetValues()
}

func (t *Tuple) AddValues(v ...Value) {
	*t = append(*t, v...)
}

func (t *Tuple) GetValues() []Value {
	return *t
}

func (t *Tuple) CanContain(Value) bool {
	return true
}

type List []Value

func (l *List) Type() Type {
	return NewListType((*l)[0].Type()) // if the list is empty, the program crashes
}

func (l *List) Value() any {
	return l.GetValues()
}

func (l *List) AddValues(v ...Value) {
	for _, value := range v {
		if len(*l) == 0 {
			*l = append(*l, value)
		} else {
			if !value.Type().Is((*l)[0].Type()) {
				panic(fmt.Errorf("invalid list: cannot create a non homogene list"))
			}
			*l = append(*l, value)
		}
	}
}

func (l *List) GetValues() []Value {
	return *l
}

func (l *List) CanContain(v Value) bool {
	if len(*l) == 0 {
		return true
	}
	return v.Type().Is((*l)[0].Type())
}

type Literal struct {
	string
	t Type
}

func NewLiteral(s string, t Type) Literal {
	return Literal{s, t}
}

func NewDefaultLiteral(s string) Literal {
	return NewLiteral(s, DefaultLiteralType)
}

func (v Literal) Type() Type {
	return v.t
}

func (v Literal) Value() any {
	return v.string
}

type String string

func (v String) Type() Type {
	return StringType
}

func (v String) Value() any {
	return string(v)
}

type Int int

func (v Int) Type() Type {
	return IntType
}

func (v Int) Value() any {
	return int(v)
}

type Float float64

func (v Float) Type() Type {
	return FloatType
}

func (v Float) Value() any {
	return float64(v)
}
