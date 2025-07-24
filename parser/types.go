package parser

import "strings"

type Type interface {
	Is(Type) bool
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

func (t *LiteralType) Is(other Type) bool {
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

func (t *TupleType) Is(other Type) bool {
	otherTuple, ok := other.(*TupleType)
	if !ok || (len(t.types) != len(otherTuple.types)) {
		return false
	}

	for i := range t.types {
		if !t.types[i].Is(otherTuple.types[i]) {
			return false
		}
	}

	return true
}

func (t *TupleType) String() string {
	var tNames []string
	for _, x := range t.types {
		tNames = append(tNames, x.String())
	}

	return "(" + strings.Join(tNames, ", ") + ")"
}

type ListType struct {
	t Type
}

func (t *ListType) Is(other Type) bool {
	otherList, ok := other.(*ListType)

	return ok && t.t.Is(otherList.t)
}

func (t *ListType) String() string {
	return "[" + t.t.String() + "]"
}
