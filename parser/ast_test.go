package parser

import (
	"github.com/planklang/goplank/lexer"
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
	axis, ok := tree.Body[0].(*Axis)
	if !ok {
		t.Error("Expected Axis, got", tree.Body[0])
	}
	if axis.Target != "x" {
		t.Error("Expected x, got", axis.Target)
	}
	if axis.Label != "" {
		t.Error("Expected empty label, got", axis.Label)
	}
	//TODO: check range
	if len(axis.Modifiers) != 0 {
		t.Error("Expected zero modifiers, got", len(axis.Modifiers))
	}
}
