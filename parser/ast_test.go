package parser

import (
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"testing"
)

func TestParse(t *testing.T) {
	lex, err := lexer.Lex("axis x")
	if err != nil {
		t.Log(lex)
		t.Fatal(err)
	}
	tree, err := Parse(lex)
	if err != nil {
		t.Log(tree)
		t.Fatal(err)
	}
	if len(tree.Body) != 1 {
		t.Errorf("Excepted 1, got %d", len(tree.Body))
	}
	fig := tree.Body[0]
	if len(fig.Stmts) != 1 {
		t.Errorf("Excepted 1, got %d", len(fig.Stmts))
	}
	axis := fig.Stmts[0]
	if axis.Keyword != "axis" {
		t.Error("Expected axis, got", axis.Keyword)
	}
	vs := axis.Arguments.GetValues()
	if len(vs) != 1 {
		t.Errorf("Excepted 1, got %d", len(vs))
	}
	arg := vs[0]
	if arg.Type() != types.DefaultLiteralType {
		t.Errorf("Excepted %s, got %s", types.DefaultLiteralType, arg.Type())
	}
	p, ok := arg.Value().(string)
	if !ok {
		t.Errorf("Cannot convert %s to literal", arg.Value())
	}
	if p != "x" {
		t.Error("Expected x, got", p)
	}
	//TODO: check range
	if len(axis.Modifiers) != 0 {
		t.Error("Expected zero modifiers, got", len(axis.Modifiers))
	}
	if t.Failed() {
		println(tree)
		t.FailNow()
	}

	lex, err = lexer.Lex("axis x 'Label'")
	if err != nil {
		t.Log(lex)
		t.Fatal(err)
	}
	tree, err = Parse(lex)
	if err != nil {
		t.Log(tree)
		t.Fatal(err)
	}
	if len(tree.Body) != 1 {
		t.Errorf("Excepted 1, got %d", len(tree.Body))
	}
	if len(tree.Body) != 1 {
		t.Errorf("Excepted 1, got %d", len(tree.Body))
	}
	fig = tree.Body[0]
	axis = fig.Stmts[0]
	if axis.Keyword != "axis" {
		t.Error("Expected axis, got", axis.Keyword)
	}
	vs = axis.Arguments.GetValues()
	if len(vs) != 2 {
		t.Errorf("Excepted 2, got %d", len(vs))
	}
	arg1 := vs[0]
	if arg1.Type() != types.DefaultLiteralType {
		t.Errorf("Excepted %s, got %s", types.DefaultLiteralType, arg1.Type())
	}
	p, ok = arg1.Value().(string)
	if !ok {
		t.Errorf("Cannot convert %s to literal", arg1.Type())
	}
	if p != "x" {
		t.Error("Expected x, got", p)
	}
	arg2 := vs[1]
	if arg2.Type() != types.StringType {
		t.Errorf("Excepted %s, got %s", types.StringType, arg1.Type())
	}
	p2, ok := arg2.Value().(string)
	if !ok {
		t.Errorf("Cannot convert %s to string", arg2.Type())
	}
	if p2 != "Label" {
		t.Error("Expected Label, got", p2)
	}
	//TODO: check range
	if len(axis.Modifiers) != 0 {
		t.Error("Expected zero modifiers, got", len(axis.Modifiers))
	}
	if t.Failed() {
		println(tree.String())
		t.FailNow()
	}
}
