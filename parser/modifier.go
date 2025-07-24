package parser

type Modifier interface {
	Modify(Statement) error
	SetArgument(*Tuple) error
	ValidArgument(Type) bool
	ValidStatement(Statement) bool
	String() string
}
