package types

type Value interface {
	Type() Type
	Value() interface{}
}

type ValueContainer interface {
	AddValues(...Value)
	GetValues() []Value
}

type Tuple []Value

func (t *Tuple) Type() Type {
	typs := make([]Type, len(*t))
	for i, v := range *t {
		typs[i] = v.Type()
	}
	return NewTupleType(typs...)
}

func (t *Tuple) Value() interface{} {
	return t.GetValues()
}

func (t *Tuple) AddValues(v ...Value) {
	*t = append(*t, v...)
}

func (t *Tuple) GetValues() []Value {
	return *t
}

type Literal struct {
	string
	t Type
}

func NewLiteral(s string, t Type) Literal {
	return Literal{s, t}
}

func NewDefaultLiteral(s string) Literal {
	return NewLiteral(s, DefaultLiteralType)
}

func (v Literal) Type() Type {
	return v.t
}

func (v Literal) Value() interface{} {
	return v.string
}

type String string

func (v String) Type() Type {
	return StringType
}

func (v String) Value() interface{} {
	return string(v)
}

type Int int

func (v Int) Type() Type {
	return IntType
}

func (v Int) Value() interface{} {
	return int(v)
}

type Float float64

func (v Float) Type() Type {
	return FloatType
}

func (v Float) Value() interface{} {
	return float64(v)
}
