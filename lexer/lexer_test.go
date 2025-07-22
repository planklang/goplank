package lexer

import (
	"errors"
	"testing"
)

func TestLex(t *testing.T) {
	res, err := Lex("axis")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Error("Expected 1, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}
	res, err = Lex("axis x")
	if err != nil {
		t.Fatal(err)
	}
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
	res, err = Lex("axis # hello world")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Error("Expected 1, got", len(res))
		t.Log(res)
	}
	if res[0].Type != KeywordType || res[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", res[0])
	}

	res, err = Lex("axis\nplot")
	if err != nil {
		t.Fatal(err)
	}
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
	res, err = Lex("axis\n| color")
	if err != nil {
		t.Fatal(err)
	}
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
	if res[2].Type != IdentifierType || res[2].Literal != "color" {
		t.Error("Expected identifier(color), got", res[2])
	}
	res, err = Lex("axis\n--- axis")
	if err != nil {
		t.Fatal(err)
	}
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

func TestLexLiteral(t *testing.T) {
	res, err := Lex("axis | color 'bonsoir je marche' ")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 4 {
		t.Error("Expected 4, got", len(res))
	}
	if res[3].Type != StringType || res[3].Literal != "bonsoir je marche" {
		t.Error("Expected string(bonsoir je marche), got", res[3])
	}
}

func TestLexError(t *testing.T) {
	res, err := Lex("1")
	if err == nil {
		t.Error("Expected error, got ", res)
	}
	if !errors.Is(err, ErrStatementExcepted) {
		t.Error("Expected ErrStatementExcepted, got", err)
	}

	res, err = Lex("axis\n;; 1 hello")
	if err == nil {
		t.Error("Expected error, got ", res)
	}
	if !errors.Is(err, ErrStatementExcepted) {
		t.Error("Expected ErrStatementExcepted, got", err)
	}

	res, err = Lex("axis | color 'bonsoir je marche pas ")
	if err == nil {
		t.Error("Expected error, got ", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}
}
