package parser

import (
	"encoding/json"
	"errors"
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
	ErrUnexpectedToken = errors.New("unexpected token")
)

type Ast struct {
	Type AstType
	Body []*Figure
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
