package parser

type Statement interface {
	Eval() error
	AddModifier(Modifier) error
	String() string
}
