package parser

type Modifier interface {
	Modify(Statement) error
	SetArgument(*Argument) error
	ValidArgument(Type) bool
	ValidStatement(Statement) bool
	String() string
}
