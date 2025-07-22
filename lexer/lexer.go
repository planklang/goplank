package lexer

import (
	"fmt"
	"slices"
	"strings"
)

type LexType string

const (
	KeywordType    LexType = "keyword"
	LiteralType    LexType = "literal"
	DelimiterType  LexType = "delimiter"
	IdentifierType LexType = "identifier"

	ImplicitDelimiter = "implicit"
	FigureDelimiter   = "---"
)

var (
	keywords   = []string{"plot", "default", "overwrite", "ow", "axis"}
	delimiters = []string{";;", "|"}

	ErrStatementExcepted = fmt.Errorf("statement excepted")
)

type Lexer struct {
	Type    LexType
	Literal string
}

func (lex *Lexer) String() string {
	return fmt.Sprintf("%s(%s)", lex.Type, lex.Literal)
}

func Lex(content string) ([]*Lexer, error) {
	var lexs []*Lexer
	lines := strings.Split(content, "\n")
	delimiterAdded := true
	inStatement := false
	inProperty := false
	identifierAdded := false
	for _, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			word := words[i]
			isDelim, isDelimFig := isDelimiter(word)
			if !delimiterAdded && i == 0 && !isDelim { // implicit delimiter
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = false
				lexs = append(lexs, &Lexer{DelimiterType, ImplicitDelimiter})
			}
			if isDelim {
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = true
				if isDelimFig {
					lexs = append(lexs, &Lexer{DelimiterType, FigureDelimiter})
				} else {
					lexs = append(lexs, &Lexer{DelimiterType, word})
					if word == "|" {
						inProperty = true
						inStatement = true
					}
				}
			} else if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
				inStatement = true
			} else if !inStatement {
				genErrorMessage(ErrStatementExcepted, i, words)
				return nil, ErrStatementExcepted
			} else if !inProperty {
				lexs = append(lexs, &Lexer{LiteralType, word})
			} else if !identifierAdded {
				lexs = append(lexs, &Lexer{IdentifierType, word})
				identifierAdded = true
			} else {
				lexs = append(lexs, &Lexer{LiteralType, word})
			}
			i++
		}
		delimiterAdded = false
		inProperty = false
		identifierAdded = false
	}
	return lexs, nil
}

func genErrorMessage(err error, i int, words []string) string {
	s := ""
	if i > 0 {
		s += words[i-1] + " "
	}
	s += words[i]
	l1 := len(s) - 1
	if i < len(words)-1 {
		s += " " + words[i+1]
	}
	l2 := len(s) - 1
	s += "\n"
	for range l1 {
		s += "-"
	}
	s += "^"
	for range l2 - l1 {
		s += "-"
	}
	return "\n" + err.Error()
}

func isDelimiter(word string) (bool, bool) {
	if slices.Contains(delimiters, word) {
		return true, false
	}
	if len(word) >= 3 && word[:3] == "---" && strings.Count(word, "-") == len(word) {
		return true, true
	}
	return false, false
}
