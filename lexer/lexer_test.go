package lexer

import "testing"

func TestLex(t *testing.T) {
	res := Lex("axis")
	if len(res) != 1 {
		t.Error("Expected 1, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	res = Lex("axis x")
	if len(res) != 2 {
		t.Error("Expected 2, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	if res[1].Type != LiteralType || res[1].Literal != "x" {
		t.Error("Expected literal(x), got", res[1])
	}
	res = Lex("axis # hello world")
	if len(res) != 1 {
		t.Error("Expected 1, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}

	res = Lex("axis\nplot")
	if len(res) != 3 {
		t.Error("Expected 3, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	if res[1].Type != DelimiterType || res[1].Literal != ImplicitDelimiter {
		t.Error("Expected delimiter(axis), got", res[1])
	}
	if res[2].Type != KeywordType || res[2].Literal != "plot" {
		t.Error("Expected keyword(plot), got", res[2])
	}
	res = Lex("axis\n| color")
	if len(res) != 3 {
		t.Error("Expected 3, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	if res[1].Type != DelimiterType || res[1].Literal != "|" {
		t.Error("Expected delimiter(|), got", res[1])
	}
	if res[2].Type != LiteralType || res[2].Literal != "color" {
		t.Error("Expected literal(color), got", res[2])
	}
	res = Lex("axis\n--- axis")
	if len(res) != 3 {
		t.Error("Expected 3, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	if res[1].Type != DelimiterType || res[1].Literal != FigureDelimiter {
		t.Error("Expected delimiter(---), got", res[1])
	}
	if res[2].Type != KeywordType || res[2].Literal != "axis" {
		t.Error("Expected keyword(color), got", res[2])
	}
}
