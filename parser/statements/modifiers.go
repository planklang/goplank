package statements

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/parser/types"
	"image/color"
	"strings"
)

var (
	TypeError = errors.New("type error")
)

type Modifier struct {
	Name      string
	Arguments *types.Tuple
}

type Colorable interface {
	GetColor() color.Color
	SetColor(color.Color)
}

type Labelable interface {
	GetLabel() string
	SetLabel(string)
}

type Rangeable interface {
	GetRange() (float64, float64)
	SetRange(float64, float64)
}

func matchColorTuple(tup *types.Tuple) (color.Color, bool, string) {
	expected1 := types.NewTupleType(types.IntType, types.IntType, types.IntType)
	expected2 := types.NewTupleType(types.IntType, types.IntType, types.IntType, types.IntType)

	t := tup.Type()
	v := tup.Value().([]uint8)

	if t.Is(expected1) {
		return color.RGBA{
			R: v[0],
			G: v[1],
			B: v[2],
			A: 255,
		}, true, ""
	} else if t.Is(expected2) {
		return color.RGBA{
			R: v[0],
			G: v[1],
			B: v[2],
			A: v[3],
		}, true, ""
	}

	return nil, false, strings.Join([]string{expected1.String(), expected1.String()}, "\n")
}

func matchLabelTuple(tup *types.Tuple) (string, bool, string) {
	expected := types.StringType

	if tup.Type().Is(expected) {
		return tup.Value().([]string)[0], true, ""
	}
	return "", false, expected.String()
}

func matchRangeTuple(tup *types.Tuple) (float64, float64, bool, string) {
	expected := types.NewTupleType(types.FloatType, types.FloatType)

	if tup.Type().Is(expected) {
		return tup.GetValues()[0].Value().(float64), tup.GetValues()[1].Value().(float64), true, ""
	}

	return 0, 0, false, expected.String()
}

func (m *Modifier) Apply(s Statement) error {
	errPropIncompatible := errors.Join(TypeError, fmt.Errorf("%s statement is not compatible with %s property", s.Keyword(), m.Name))
	errTypeIncompatible := func(exp string) error {
		return errors.Join(TypeError, fmt.Errorf("%s property expects argument one of:\n%s\n; got %s", m.Name, exp, m.Arguments.Type().String()))
	}
	switch m.Name {
	case "color":
		colorable, ok := s.(Colorable)
		if !ok {
			return errPropIncompatible
		}

		c, ok, exp := matchColorTuple(m.Arguments)
		if !ok {
			return errTypeIncompatible(exp)
		}

		colorable.SetColor(c)
	case "label":
		labelable, ok := s.(Labelable)
		if !ok {
			return errPropIncompatible
		}

		t, ok, exp := matchLabelTuple(m.Arguments)
		if !ok {
			return errTypeIncompatible(exp)
		}

		labelable.SetLabel(t)
	case "range":
		rangeable, ok := s.(Rangeable)
		if !ok {
			return errPropIncompatible
		}

		rMin, rMax, ok, exp := matchRangeTuple(m.Arguments)
		if !ok {
			return errTypeIncompatible(exp)
		}

		rangeable.SetRange(rMin, rMax)
	default:
		panic(fmt.Sprintf("%s property is not implemented yet", m.Name))
	}

	return nil
}
