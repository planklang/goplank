package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
)

type AstType uint

const (
	AstTypeDefault AstType = 0
)

var (
	ErrUnknownValue    = errors.New("unknown value")
	ErrInternal        = errors.New("internal error")
	ErrInvalidModifier = errors.New("invalid modifier")
	ErrInvalidArgument = errors.New("invalid argument")
)

type Ast struct {
	Type AstType
	Body []Statement
}

func (a *Ast) Eval() error {
	for _, s := range a.Body {
		if err := s.Eval(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Ast) String() string {
	m, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}
	return string(m)
}

func Parse(lex []*lexer.Lexer) (*Ast, error) {
	tree := new(Ast)
	var stmt Statement
	var modif Modifier
	inModifier := false
	for _, l := range lex {
		switch l.Type {
		case lexer.StatementDelimiterType:
			if inModifier {
				if stmt == nil {
					return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but modifier finished"))
				}
				if err := stmt.AddModifier(modif); err != nil {
					return nil, err
				}
			}
			if stmt != nil {
				tree.Body = append(tree.Body, stmt)
			}
		case lexer.KeywordType:
			switch l.Literal {
			case "axis":
				stmt = new(Axis)
			default:
				return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported keyword %s", l.Literal))
			}
		case lexer.ModifierDelimiterType:
			if stmt == nil {
				return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but modifier started"))
			}
			if inModifier {
				//TODO: add arguments
				if err := stmt.AddModifier(modif); err != nil {
					return nil, err
				}
			}
			inModifier = true
			modif = nil
		case lexer.IdentifierType:
			if inModifier && modif == nil {
				//TODO: create modifier
			} else {
				return nil, errors.Join(ErrInternal, fmt.Errorf("identifier received but not in modifier or modif is not nil"))
			}
		default:
			err := parseLiteral(l)
			if err != nil {
				return nil, err
			}
		}
	}
	return tree, nil
}

func parseLiteral(lex *lexer.Lexer) error {
	switch lex.Type {
	case lexer.LiteralType:
	case lexer.VariableType:
	case lexer.StringType:
	case lexer.IntType:
	case lexer.FloatType:
	default:
		return errors.Join(ErrUnknownValue, fmt.Errorf("unsupported lex type %s", lex.Type))
	}
	return nil
}
