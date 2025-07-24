package parser

type statement interface {
	Eval() error
}
