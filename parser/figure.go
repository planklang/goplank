package parser

import "github.com/planklang/goplank/parser/statements"

type Figure struct {
	Stmts []statements.Statement
}

func (f *Figure) Eval() error { return nil }
