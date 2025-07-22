package lexer

import "testing"

func TestLex(t *testing.T) {
	res := Lex("axis")
	if len(res) != 1 {
		t.Error("Expected 1, got ", len(res))
		printLex(t, res)
	}
	if res[0].Type != keywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got ", res[0])
	}
	res = Lex("axis x")
	if len(res) != 2 {
		t.Error("Expected 2, got ", len(res))
	}
	if res[0].Type != keywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got ", res[0])
	}
	if res[1].Type != literalType || res[1].Literal != "x" {
		t.Error("Expected literal(x), got ", res[1])
	}
	res = Lex("axis # hello world")
	if len(res) != 1 {
		t.Error("Expected 1, got ", len(res))
		printLex(t, res)
	}
	if res[0].Type != keywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got ", res[0])
	}
}

func printLex(t *testing.T, lexs []*Lexer) {
	s := ""
	for _, l := range lexs {
		s += l.String() + " "
	}
	t.Log(s[:len(s)-1])
}
