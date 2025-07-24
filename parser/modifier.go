package parser

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/parser/types"
	"image/color"
	"math"
)

type Modifier interface {
	Modify(Statement) error
	SetArgument(*types.Tuple) error
	ValidArgument(types.Type) bool
	ValidStatement(Statement) bool
	String() string
}

type ModifierColor struct {
	*color.RGBA
}

func (c *ModifierColor) Modify(s Statement) error {
	return nil
}

func (c *ModifierColor) SetArgument(t *types.Tuple) error {
	if !c.ValidArgument(t.Type()) {
		return errors.Join(ErrInvalidModifier, fmt.Errorf("cannot apply argument %v to modifier %s", t, c))
	}
	val := t.GetValues()
	c.RGBA = &color.RGBA{
		R: uint8(val[0].Value().(int)),                       // inferred by ValidArgument
		G: uint8(val[1].Value().(int)),                       // inferred by ValidArgument
		B: uint8(val[2].Value().(int)),                       // inferred by ValidArgument
		A: uint8(math.Floor(val[3].Value().(float64) * 255)), // inferred by ValidArgument
	}
	return nil
}

func (c *ModifierColor) ValidArgument(t types.Type) bool {
	return t.Is(types.NewTupleType(types.IntType, types.IntType, types.IntType, types.FloatType))
}

func (c *ModifierColor) ValidStatement(s Statement) bool {
	//_, ok := s.(StatementPlot)
	//return ok
	return false
}

func (c *ModifierColor) String() string {
	return fmt.Sprintf("color{%d %d %d %d}", c.R, c.G, c.B, c.A)
}
