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
	if res[1].Type != StatementDelimiterType || res[1].Literal != ImplicitDelimiter {
		t.Error("Expected statement_delimiter(axis), got", res[1])
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
	if res[1].Type != PropertyDelimiterType || res[1].Literal != "|" {
		t.Error("Expected property_delimiter(|), got", res[1])
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
	if res[1].Type != FigureDelimiterType || res[1].Literal != FigureDelimiter {
		t.Error("Expected figure_delimiter(---), got", res[1])
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

	res, err = Lex("axis (1 2 3 4)")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 7 {
		t.Error("Expected 7, got", len(res))
		t.Log(res)
	}
	if res[1].Type != WeakDelimiterType || res[1].Literal != "(" {
		t.Error("Expected delimiter((), got", res[1])
	}
	if res[2].Type != NumberType || res[2].Literal != "1" {
		t.Error("Expected number(1), got", res[2])
	}
	if res[3].Type != NumberType || res[3].Literal != "2" {
		t.Error("Expected number(2), got", res[3])
	}
	if res[4].Type != NumberType || res[4].Literal != "3" {
		t.Error("Expected number(3), got", res[4])
	}
	if res[5].Type != NumberType || res[5].Literal != "4" {
		t.Error("Expected number(4), got", res[5])
	}
	if res[6].Type != WeakDelimiterType || res[6].Literal != ")" {
		t.Error("Expected delimiter()), got", res[6])
	}
}

func TestLexError(t *testing.T) {
	res, err := Lex("12")
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

	res, err = Lex("axis | color 1.23.3 ")
	if err == nil {
		t.Error("Expected error, got ", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}
}
