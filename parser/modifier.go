package parser

type Modifier interface {
	Modify(Statement) error
	AddArgument(*Argument) error
	ValidArg(Type) bool
	ValidStatement(Statement) bool
	String() string
}
