package statements

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/parser/types"
)

type Axis struct {
	Target string
	Label  string
	Min    float64
	Max    float64
}

func (a *Axis) UnpackArgs(args *types.Tuple) error {
	if !a.validArgument(args.Type()) {
		return errors.Join(ErrInvalidArgument, fmt.Errorf("cannot apply argument %v to statement %s", args, a))
	}
	values := args.GetValues()
	a.Target = values[0].Value().(string) // inferred by validArgument
	for _, v := range values {
		switch v.Type() {
		case types.StringType:
			a.SetLabel(v.Value().(string)) // inferred by validArgument
		case types.NewTupleType(types.FloatType, types.FloatType):
			vl := v.Value().([]float64) // inferred by validArgument
			a.SetRange(vl[0], vl[1])
		}
	}
	return nil
}

func (a *Axis) Keyword() string {
	return "axis"
}

func (a *Axis) String() string {
	return fmt.Sprintf("Axis{%s %s [%f %f]}", a.Target, a.Label, a.Min, a.Max)
}

func (a *Axis) GetLabel() string {
	return a.Label
}

func (a *Axis) SetLabel(label string) {
	a.Label = label
}

func (a *Axis) GetRange() (float64, float64) {
	return a.Min, a.Max
}

func (a *Axis) SetRange(rMin float64, rMax float64) {
	a.Min = rMin
	a.Max = rMax
}

func (a *Axis) validArgument(t types.Type) bool {
	valids := []types.Type{
		types.NewTupleType(types.DefaultLiteralType),
		types.NewTupleType(types.DefaultLiteralType, types.NewTupleType(types.FloatType, types.FloatType)),
		types.NewTupleType(types.DefaultLiteralType, types.StringType),
		types.NewTupleType(types.DefaultLiteralType, types.NewTupleType(types.FloatType, types.FloatType), types.StringType),
		types.NewTupleType(types.DefaultLiteralType, types.StringType, types.NewTupleType(types.FloatType, types.FloatType)),
	}
	for _, v := range valids {
		if t.Is(v) {
			return true
		}
	}
	return false
}
