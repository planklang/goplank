package parser

import "strings"

type Type interface {
	IsEqual(other Type) bool
	String() string
}

const (
	stringLiteral = "string"
	intLiteral    = "int"
	floatLiteral  = "float"
)

type LiteralType struct {
	t string
}

func (t *LiteralType) IsEqual(other Type) bool {
	otherLiteral, ok := other.(*LiteralType)
	return ok && t.t == otherLiteral.t
}

func (t *LiteralType) String() string {
	return t.t
}

var (
	StringType = &LiteralType{stringLiteral}
	IntType    = &LiteralType{intLiteral}
	FloatType  = &LiteralType{floatLiteral}
)

type TupleType struct {
	types []Type
}

func (t *TupleType) IsEqual(other Type) bool {
	otherTuple, ok := other.(*TupleType)
	if !ok || (len(t.types) != len(otherTuple.types)) {
		return false
	}

	for i := range t.types {
		if !t.types[i].IsEqual(otherTuple.types[i]) {
			return false
		}
	}

	return true
}

func (t *TupleType) String() string {
	t_names := []string{}
	for _, x := range t.types {
		t_names = append(t_names, x.String())
	}

	return "(" + strings.Join(t_names, ", ") + ")"
}

type ListType struct {
	t Type
}

func (t *ListType) IsEqual(other Type) bool {
	otherList, ok := other.(*ListType)

	return ok && t.t.IsEqual(otherList.t)
}

func (t *ListType) String() string {
	return "[" + t.t.String() + "]"
}
