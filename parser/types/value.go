package types

import (
	"fmt"
	"strconv"
)

type Value interface {
	Type() Type
	Cast(Type) (Value, bool)
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

func (t *Tuple) Cast(target Type) (Value, bool) {
	if target.Is(t.Type()) {
		return t, true
	}

	if target.Is(NewTupleType(t.Type())) {
		var res Tuple = []Value{t}
		return &res, true
	}

	if len(t.GetValues()) == 1 {
		return t.GetValues()[0].Cast(target)
	}

	return nil, false
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

func (l *List) Cast(target Type) (Value, bool) {
	if target.Is(l.Type()) {
		return l, true
	}
	if target.Is(NewTupleType(l.Type())) {
		var res Tuple = []Value{l}
		return &res, true
	}

	return nil, false
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

func (v Literal) Cast(target Type) (Value, bool) {
	if target.Is(v.Type()) {
		return v, true
	}

	if target.Is(NewTupleType(v.t)) {
		t := (Tuple)([]Value{v})
		return &t, true
	}

	if target.Is(StringType) {
		return (String)(v.string), true
	}

	return nil, false
}

func (v Literal) Value() any {
	return v.string
}

type String string

func (v String) Type() Type {
	return StringType
}

func (v String) Cast(target Type) (Value, bool) {
	if target.Is(v.Type()) {
		return v, true
	}

	if target.Is(NewTupleType(StringType)) {
		t := (Tuple)([]Value{v})
		return &t, true
	}

	return nil, false
}

func (v String) Value() any {
	return string(v)
}

type Int int

func (v Int) Type() Type {
	return IntType
}

func (v Int) Cast(target Type) (Value, bool) {
	if target.Is(v.Type()) {
		return v, true
	}

	if target.Is(NewTupleType(StringType)) {
		t := (Tuple)([]Value{v})
		return &t, true
	}

	if target.Is(StringType) {
		return (String)(strconv.Itoa(int(v))), true
	}

	if target.Is(FloatType) {
		return (Float)(float64(v)), true
	}

	return nil, false
}

func (v Int) Value() any {
	return int(v)
}

type Float float64

func (v Float) Type() Type {
	return FloatType
}

func (v Float) Cast(target Type) (Value, bool) {
	if target.Is(v.Type()) {
		return v, true
	}

	if target.Is(NewTupleType(StringType)) {
		t := (Tuple)([]Value{v})
		return &t, true
	}

	if target.Is(StringType) {
		return (String)(strconv.FormatFloat(float64(v), 'f', 16, 64)), true
	}

	return nil, false
}

func (v Float) Value() any {
	return float64(v)
}
