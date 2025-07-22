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
	KeywordType            LexType = "keyword"
	LiteralType            LexType = "literal"
	DelimiterType          LexType = "delimiter"
	PropertyDelimiterType  LexType = "property_delimiter"
	FigureDelimiterType    LexType = "figure_delimiter"
	StatementDelimiterType LexType = "statement_delimiter"
	IdentifierType         LexType = "identifier"
	VariableType           LexType = "variable"
	StringType             LexType = "string"
	NumberType             LexType = "number"

	ImplicitDelimiter = "implicit"
	FigureDelimiter   = "---"
)

var (
	keywords            = []string{"plot", "default", "overwrite", "ow", "axis"}
	propertyDelimiters  = []string{"|"}
	statementDelimiters = []string{";;"}

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
	for ln, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			word := words[i]
			isDelim, typ := isDelimiter(word)
			if !delimiterAdded && i == 0 && !isDelim { // implicit delimiter
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = false
				lexs = append(lexs, &Lexer{StatementDelimiterType, ImplicitDelimiter})
			}
			if isDelim {
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = true
				if typ == FigureDelimiterType {
					lexs = append(lexs, &Lexer{typ, FigureDelimiter})
				} else {
					lexs = append(lexs, &Lexer{typ, word})
					if typ == PropertyDelimiterType {
						inProperty = true
						inStatement = true
					}
				}
			} else if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
				inStatement = true
			} else if !inStatement {
				fmt.Println(genErrorMessage(ErrStatementExcepted, i, words, ln))
				return nil, ErrStatementExcepted
			} else if !identifierAdded && inProperty {
				lexs = append(lexs, &Lexer{IdentifierType, word})
				identifierAdded = true
			} else {
				ls, err := parseLiteral(word, &i, words)
				if err != nil {
					fmt.Println(genErrorMessage(err, i-1, words, ln)) // i-1 because every error here leads to i == len(words)
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

func parseLiteral(word string, i *int, words []string) ([]*Lexer, error) {
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
	case '(', ')', '[', ']':
		var lexs []*Lexer
		lexs = append(lexs, &Lexer{DelimiterType, string(f)})
		if len(word) == 1 {
			*i++
			return lexs, nil
		}
		ls, err := parseLiteral(word[1:], i, words)
		if err != nil {
			return nil, err
		}
		lexs = append(lexs, ls...)
		return lexs, nil
	}
	_, err := strconv.ParseFloat(word, 64)
	if err != nil {
		return []*Lexer{{LiteralType, word}}, nil
	}
	return []*Lexer{{NumberType, word}}, nil
}

func genErrorMessage(err error, i int, words []string, line int) string {
	s := ""
	for j := range i {
		s += words[j] + " "
	}
	l1 := len(s)
	s += words[i]
	for j := range len(words) - i - 1 {
		s += " " + words[j+i+1]
	}
	l2 := len(s) - 1
	s += "\n"
	if i == len(words)-1 {
		for range l2 {
			s += "-"
		}
		s += "^"
	} else {
		for range l1 {
			s += "-"
		}
		s += "^"
		for range l2 - l1 {
			s += "-"
		}
	}
	return fmt.Sprintf("%s (line %d)\n\n%s", s, line+1, err.Error())
}

func isDelimiter(word string) (bool, LexType) {
	if slices.Contains(statementDelimiters, word) {
		return true, StatementDelimiterType
	} else if slices.Contains(propertyDelimiters, word) {
		return true, PropertyDelimiterType
	}
	if len(word) >= 3 && word[:3] == "---" && strings.Count(word, "-") == len(word) {
		return true, FigureDelimiterType
	}
	return false, ""
}
