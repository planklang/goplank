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
	resList := res.list
	if len(resList) != 1 {
		t.Error("Expected 1, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}
	res, err = Lex("axis x")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 2 {
		t.Error("Expected 2, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}
	if resList[1].Type != LiteralType || resList[1].Literal != "x" {
		t.Error("Expected literal(x), got", resList[1])
	}
	res, err = Lex("axis # hello world")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 1 {
		t.Error("Expected 1, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}

	res, err = Lex("axis\nplot")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 3 {
		t.Error("Expected 3, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}
	if resList[1].Type != StatementDelimiterType || resList[1].Literal != ImplicitDelimiter {
		t.Error("Expected statement_delimiter(axis), got", resList[1])
	}
	if resList[2].Type != KeywordType || resList[2].Literal != "plot" {
		t.Error("Expected keyword(plot), got", resList[2])
	}
	res, err = Lex("axis\n| color")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 3 {
		t.Error("Expected 3, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}
	if resList[1].Type != ModifierDelimiterType || resList[1].Literal != "|" {
		t.Error("Expected property_delimiter(|), got", resList[1])
	}
	if resList[2].Type != ModifierType || resList[2].Literal != "color" {
		t.Error("Expected modifier(color), got", resList[2])
	}
	res, err = Lex("axis\n--- axis")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 3 {
		t.Error("Expected 3, got", len(resList))
		t.Log(resList)
	}
	if resList[0].Type != KeywordType || resList[0].Literal != "axis" {
		t.Error("Expected keyword(axis), got", resList[0])
	}
	if resList[1].Type != FigureDelimiterType || resList[1].Literal != FigureDelimiter {
		t.Error("Expected figure_delimiter(---), got", resList[1])
	}
	if resList[2].Type != KeywordType || resList[2].Literal != "axis" {
		t.Error("Expected keyword(color), got", resList[2])
	}

	res, err = Lex("axis ;;")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 1 {
		t.Error("Expected 1, got", len(resList))
		t.Log(resList)
	}
}

func TestLexLiteral(t *testing.T) {
	res, err := Lex("axis | color 'bonsoir je marche' ")
	if err != nil {
		t.Fatal(err)
	}
	resList := res.list
	if len(resList) != 4 {
		t.Error("Expected 4, got", len(resList))
	}
	if resList[3].Type != StringType || resList[3].Literal != "bonsoir je marche" {
		t.Error("Expected string(bonsoir je marche), got", resList[3])
	}

	res, err = Lex("axis 1 0.2 .5")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 4 {
		t.Error("Expected 4, got", len(resList))
		t.Log(resList)
	}
	if resList[1].Type != IntType || resList[1].Literal != "1" {
		t.Error("Expected int(1), got", resList[1])
	}
	if resList[2].Type != FloatType || resList[2].Literal != "0.2" {
		t.Error("Expected float(0.2), got", resList[2])
	}
	if resList[3].Type != FloatType || resList[3].Literal != "0.5" {
		t.Error("Expected float(0.5), got", resList[3])
	}

	res, err = Lex("axis $hello")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 2 {
		t.Error("Expected 2, got", len(resList))
		t.Log(resList)
	}
	if resList[1].Type != VariableType || resList[1].Literal != "hello" {
		t.Error("Expected variable(hello), got", resList[1])
	}

	res, err = Lex("axis x")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 2 {
		t.Error("Expected 2, got", len(resList))
		t.Log(resList)
	}
	if resList[1].Type != LiteralType || resList[1].Literal != "x" {
		t.Error("Expected literal(x), got", resList[1])
	}

	res, err = Lex("axis (1 2 3 4)")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 7 {
		t.Error("Expected 7, got", len(resList))
		t.Log(resList)
	}
	if resList[1].Type != WeakDelimiterType || resList[1].Literal != "(" {
		t.Error("Expected delimiter((), got", resList[1])
	}
	if resList[2].Type != IntType || resList[2].Literal != "1" {
		t.Error("Expected int(1), got", resList[2])
	}
	if resList[3].Type != IntType || resList[3].Literal != "2" {
		t.Error("Expected int(2), got", resList[3])
	}
	if resList[4].Type != IntType || resList[4].Literal != "3" {
		t.Error("Expected int(3), got", resList[4])
	}
	if resList[5].Type != IntType || resList[5].Literal != "4" {
		t.Error("Expected int(4), got", resList[5])
	}
	if resList[6].Type != WeakDelimiterType || resList[6].Literal != ")" {
		t.Error("Expected delimiter()), got", resList[6])
	}

	res, err = Lex("axis [1 2 3 4]")
	if err != nil {
		t.Fatal(err)
	}
	resList = res.list
	if len(resList) != 7 {
		t.Error("Expected 7, got", len(resList))
		t.Log(resList)
	}
	if resList[1].Type != WeakDelimiterType || resList[1].Literal != "[" {
		t.Error("Expected delimiter([), got", resList[1])
	}
	if resList[6].Type != WeakDelimiterType || resList[6].Literal != "]" {
		t.Error("Expected delimiter(]), got", resList[1])
	}
}

func TestLexError(t *testing.T) {
	res, err := Lex("12")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrStatementExcepted) {
		t.Error("Expected ErrStatementExcepted, got", err)
	}

	res, err = Lex("axis\n;; 1 hello")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrStatementExcepted) {
		t.Error("Expected ErrStatementExcepted, got", err)
	}

	res, err = Lex("axis | color 'bonsoir je marche pas ")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}

	res, err = Lex("axis | color 1.23.3 ")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}

	res, err = Lex("axis (5 .2")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}

	res, err = Lex("axis 5 .2)")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}

	res, err = Lex("axis [5 .2")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}

	res, err = Lex("axis 5 .2]")
	if err == nil {
		t.Error("Expected error, got", res)
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Error("Expected ErrInvalidExpression, got", err)
	}
}
