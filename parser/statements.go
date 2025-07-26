package parser

import "github.com/planklang/goplank/parser/types"

type Statement struct {
	Keyword   string
	Arguments *types.Tuple
	Modifiers []*Modifier
}

//type Axis struct {
//	Target    string
//	Label     string
//	Range     [2]float64
//	Modifiers []Modifier
//}
//
//func (a *Axis) Eval() error {
//	return nil
//}
//
//func (a *Axis) AddModifier(m Modifier) error {
//	if !m.ValidStatement(a) {
//		return errors.Join(ErrInvalidModifier, fmt.Errorf("cannot apply modifier %s to statement %s", m, a))
//	}
//	a.Modifiers = append(a.Modifiers, m)
//	return nil
//}
//
//func (a *Axis) SetArgument(arg *types.Tuple) error {
//	if !a.ValidArgument(arg.Type()) {
//		return errors.Join(ErrInvalidArgument, fmt.Errorf("cannot apply argument %v to statement %s", arg, a))
//	}
//	values := arg.GetValues()
//	a.Target = values[0].Value().(string) // inferred by ValidArgument
//	for _, v := range values {
//		switch v.Type() {
//		case types.StringType:
//			a.Label = v.Value().(string) // inferred by ValidArgument
//		case types.NewListType(types.FloatType):
//			vl := v.Value().([]float64) // inferred by ValidArgument
//			if len(vl) != 2 {
//				return errors.Join(ErrInvalidArgument, fmt.Errorf("invalid length for range %v", vl))
//			}
//			a.Range = [2]float64{vl[0], vl[1]}
//		}
//	}
//	return nil
//}
//
//func (a *Axis) ValidArgument(t types.Type) bool {
//	valids := []types.Type{
//		types.NewTupleType(types.DefaultLiteralType),
//		types.NewTupleType(types.DefaultLiteralType, types.NewListType(types.FloatType)),
//		types.NewTupleType(types.DefaultLiteralType, types.StringType),
//		types.NewTupleType(types.DefaultLiteralType, types.NewListType(types.FloatType), types.StringType),
//		types.NewTupleType(types.DefaultLiteralType, types.StringType, types.NewListType(types.FloatType)),
//	}
//	for _, v := range valids {
//		if t.Is(v) {
//			return true
//		}
//	}
//	return false
//}
//
//func (a *Axis) String() string {
//	return fmt.Sprintf("Axis{%s %s [%f %f]}", a.Target, a.Label, a.Range[0], a.Range[1])
//}
