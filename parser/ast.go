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
	ErrUnknownType     = errors.New("unknown type")
	ErrInternal        = errors.New("internal error")
	ErrInvalidModifier = errors.New("invalid modifier")
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
					return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but Modifier finished"))
				}
				if err := stmt.AddModifier(modif); err != nil {
					return nil, err
				}
			}
			if stmt != nil {
				tree.Body = append(tree.Body, stmt)
			}
		case lexer.KeywordType:
			//TODO: create statement
		case lexer.ModifierDelimiterType:
			if stmt == nil {
				return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but Modifier started"))
			}
			if inModifier {
				if !modif.ValidStatement(stmt) {
					return nil, errors.Join(ErrInvalidModifier, fmt.Errorf("cannot apply modifier %s to statement %s", modif, stmt))
				}
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
			}
		default:
			return nil, errors.Join(ErrUnknownType, fmt.Errorf("unsupported lex type %s", l.Type))
		}
	}
	return tree, nil
}
