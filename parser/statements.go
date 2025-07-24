package parser

import (
	"errors"
	"fmt"
)

type Statement interface {
	Eval() error
	AddModifier(Modifier) error
	SetArgument(*Tuple) error
	ValidArgument(Type) bool
	String() string
}

type Axis struct {
	Target    string
	Label     string
	Range     [2]float64
	Modifiers []Modifier
}

func (a *Axis) Eval() error {
	return nil
}

func (a *Axis) AddModifier(m Modifier) error {
	if !m.ValidStatement(a) {
		return errors.Join(ErrInvalidModifier, fmt.Errorf("cannot apply modifier %s to statement %s", m, a))
	}
	a.Modifiers = append(a.Modifiers, m)
	return nil
}

func (a *Axis) SetArgument(arg *Tuple) error {
	if !a.ValidArgument(arg.Type()) {
		return errors.Join(ErrInvalidArgument, fmt.Errorf("cannot apply argument %s to statement %s", arg, a))
	}
	values := arg.Value().([]Value)       // inferred by Tuple type
	a.Target = values[0].Value().(string) // inferred by ValidArgument
	for _, v := range values {
		switch v.Type() {
		case StringType:
			a.Label = v.Value().(string) // inferred by ValidArgument
		case NewListType(FloatType):
			vl := v.Value().([]float64) // inferred by ValidArgument
			if len(vl) != 2 {
				return errors.Join(ErrInvalidArgument, fmt.Errorf("invalid length for range %v", vl))
			}
			a.Range = [2]float64{vl[0], vl[1]}
		}
	}
	return nil
}

func (a *Axis) ValidArgument(t Type) bool {
	valids := []Type{
		NewTupleType(DefaultLiteralType),
		NewTupleType(DefaultLiteralType, NewListType(FloatType)),
		NewTupleType(DefaultLiteralType, StringType),
		NewTupleType(DefaultLiteralType, NewListType(FloatType), StringType),
		NewTupleType(DefaultLiteralType, StringType, NewListType(FloatType)),
	}
	for _, v := range valids {
		if t.Is(v) {
			return true
		}
	}
	return false
}

func (a *Axis) String() string {
	return fmt.Sprintf("Axis{%s %s [%f %f]}", a.Target, a.Label, a.Range[0], a.Range[1])
}
