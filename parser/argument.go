package parser

import "github.com/planklang/goplank/lexer"

type Argument struct {
	t      Type
	Values []*lexer.Lexer
}

func (a *Argument) Type() Type {
	return a.t
}
