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
	ErrUnknownType = errors.New("unknown type")
)

type Ast struct {
	Type AstType
	Body []statement
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
	var stmt statement
	for _, l := range lex {
		switch l.Type {
		case lexer.StatementDelimiterType:
			if stmt != nil {
				tree.Body = append(tree.Body, stmt)
			}
		default:
			return nil, errors.Join(ErrUnknownType, fmt.Errorf("unsupported lex type %s", l.Type))
		}
	}
	return tree, nil
}
