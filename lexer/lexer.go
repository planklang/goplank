package lexer

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type LexType string

const (
	KeywordType    LexType = "keyword"
	LiteralType    LexType = "literal"
	DelimiterType  LexType = "delimiter"
	IdentifierType LexType = "identifier"
	VariableType   LexType = "variable"
	StringType     LexType = "string"
	NumberType     LexType = "number"

	ImplicitDelimiter = "implicit"
	FigureDelimiter   = "---"
)

var (
	keywords   = []string{"plot", "default", "overwrite", "ow", "axis"}
	delimiters = []string{";;", "|"}

	ErrStatementExcepted = fmt.Errorf("statement excepted")
	ErrInvalidExpression = fmt.Errorf("invalid expression")
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
				fmt.Println(genErrorMessage(ErrStatementExcepted, i, words))
				return nil, ErrStatementExcepted
			} else if !identifierAdded && inProperty {
				lexs = append(lexs, &Lexer{IdentifierType, word})
				identifierAdded = true
			} else {
				ls, err := parseLiteral(&i, words)
				if err != nil {
					fmt.Println(genErrorMessage(err, i-1, words)) // i-1 because every error here leads to i == len(words)
					return nil, err
				}
				lexs = append(lexs, ls...)
			}
			i++
		}
		delimiterAdded = false
		inProperty = false
		identifierAdded = false
	}
	return lexs, nil
}

func parseLiteral(i *int, words []string) ([]*Lexer, error) {
	word := words[*i]
	f := word[0]
	switch f {
	case '$':
		if len(word) == 1 {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("$ is reserved to call variables"))
		}
		return []*Lexer{{VariableType, word[1:]}}, nil
	case '"', '\'', '`':
		s := ""
		finished := false
		j := 1
		for *i < len(words) && !finished {
			c := words[*i]
			for j < len(c) && !finished {
				s += string(c[j])
				if j < len(c)-1 {
					finished = c[j+1] == f
				}
				j++
			}
			s += " "
			j = 0
			*i++
		}
		if !finished {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("string is not finished"))
		}
		return []*Lexer{{StringType, s[:len(s)-1]}}, nil
	case '(':
		var lexs []*Lexer
		lexs = append(lexs, &Lexer{DelimiterType, "("})
		if len(word) > 1 {
			//
		}
		finished := false
		for *i < len(words) && !finished {
			l, err := parseLiteral(i, words)
			if err != nil {
				return nil, err
			}
			c := words[*i]
			finished = c[len(c)-1] != ')'
			lexs = append(lexs, l...)
		}
		if !finished {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("tuple is not finished"))
		}
		lexs = append(lexs, &Lexer{DelimiterType, ")"})
		return lexs, nil
	case '[':
		var lexs []*Lexer
		lexs = append(lexs, &Lexer{DelimiterType, "["})
		if len(word) > 1 {
			//
		}
		finished := false
		for *i < len(words) && !finished {
			l, err := parseLiteral(i, words)
			if err != nil {
				return nil, err
			}
			c := words[*i]
			finished = c[len(c)-1] != ']'
			lexs = append(lexs, l...)
		}
		if !finished {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("list is not finished"))
		}
		lexs = append(lexs, &Lexer{DelimiterType, "]"})
		return lexs, nil
	}
	_, err := strconv.ParseFloat(word, 64)
	if err != nil {
		return []*Lexer{{LiteralType, word}}, nil
	}
	return []*Lexer{{NumberType, word}}, nil
}

func genErrorMessage(err error, i int, words []string) string {
	s := ""
	for j := range i {
		s += words[j] + " "
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
	return s + "\n" + err.Error()
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
