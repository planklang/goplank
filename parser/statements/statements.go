package statements

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/parser/types"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

type Statement interface {
	UnpackArgs(args *types.Tuple) error
	Keyword() string
	String() string
}

func NewStatement(keyword string) Statement {
	var stmt Statement

	switch keyword {
	case "axis":
		stmt = &Axis{}
	default:
		panic(fmt.Sprintf("statement %s not implemented yet", keyword))
	}

	return stmt
}
