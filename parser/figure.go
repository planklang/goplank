package parser

type Figure struct {
	Stmts []*Statement
}

func (f *Figure) Eval() error { return nil }
