package parser

import "github.com/planklang/goplank/parser/types"

type Modifier interface {
	Modify(Statement) error
	SetArgument(*types.Tuple) error
	ValidArgument(types.Type) bool
	ValidStatement(Statement) bool
	String() string
}
