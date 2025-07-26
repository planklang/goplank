package parser

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"strconv"
)

var (
	ErrUnknownValue      = errors.New("unknown value")
	ErrInternal          = errors.New("internal error")
	ErrInvalidModifier   = errors.New("invalid modifier")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInvalidLiteral    = errors.New("invalid literal")
	ErrMissingLiteral    = errors.New("missing literal")
	ErrUnexpectedToken   = errors.New("unexpected token")
	ErrDelimiterExcepted = errors.Join(ErrUnexpectedToken, errors.New("delimiter excepted"))
)

func Parse(lex *lexer.TokenList) (*Ast, error) {
	// top-level = [ figure, [{ figure-delimiter, [figure] }] ];

	tree := new(Ast)
	tree.Type = AstTypeDefault

	for lex.Next() {
		fig, err := parseFigure(lex)
		if err != nil {
			return nil, err
		}
		tree.Body = append(tree.Body, fig)

		if !lex.Next() {
			return tree, nil
		}
		if lex.Current().Type != lexer.FigureDelimiterType {
			return nil, errors.Join(ErrDelimiterExcepted, fmt.Errorf("expected figure delimiter, not %s", lex.Current()))
		}
	}

	return tree, nil
}

func parseFigure(lex *lexer.TokenList) (*Figure, error) {
	// figure = [ statement, [{ statement-delimiter, [ statement ] }] ];

	fig := new(Figure)

	for lex.Next() {
		stmt, err := parseStatement(lex)
		if err != nil {
			return fig, err
		}
		fig.Stmts = append(fig.Stmts, stmt)

		if !lex.Next() {
			return fig, nil
		}
		if lex.Current().Type != lexer.StatementDelimiterType {
			return nil, errors.Join(ErrDelimiterExcepted, fmt.Errorf("expected statement delimiter, not %s", lex.Current()))
		}
	}

	return fig, nil
}

func parseStatement(lex *lexer.TokenList) (*Statement, error) {
	// statement = keyword, [ arguments ], [{ property-delimiter, property }];

	if lex.Current().Type != lexer.KeywordType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected keyword, not %s", lex.Current()))
	}

	stmt := new(Statement)
	stmt.Keyword = lex.Current().Literal

	if !lex.Next() {
		return stmt, nil
	}

	if lex.Current().Type != lexer.ModifierDelimiterType {
		args, err := parseArgument(lex)
		if err != nil {
			return nil, err
		}
		stmt.Arguments = args
	}

	for lex.Current().Type == lexer.ModifierDelimiterType {
		if !lex.Next() {
			return nil, errors.Join(lexer.ErrInvalidExpression, fmt.Errorf("expected modifier definition after modifier delimiter"))
		}

		mod, err := parseProperty(lex)
		if err != nil {
			return nil, err
		}
		stmt.Modifiers = append(stmt.Modifiers, mod)
	}

	return stmt, nil
}

func parseProperty(lex *lexer.TokenList) (*Modifier, error) {
	// property = ? identifier ?, [ arguments ]

	if lex.Current().Type != lexer.IdentifierType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected modifier name, not %v", lex.Current()))
	}

	mod := new(Modifier)
	mod.Name = lex.Current().Literal

	if !lex.Next() {
		return mod, nil
	}

	args, err := parseArgument(lex)
	if err != nil {
		return nil, err
	}
	mod.Arguments = args
	return mod, nil
}

func parseArgument(lex *lexer.TokenList) (*types.Tuple, error) {
	tuple := new(types.Tuple)

	for lex.Current().Type != lexer.StatementDelimiterType &&
		lex.Current().Type != lexer.ModifierDelimiterType &&
		lex.Current().Type != lexer.FigureDelimiterType { // do not call [TokenList.Next] because argument does not require anything
		val, err := parseWeakDelimiters(lex)
		if err != nil {
			return nil, err
		}
		tuple.AddValues(val)

		if !lex.Next() { // call [TokenList.Next] here because parseWeakDelimiters never skips the last one
			return tuple, nil
		}
	}

	// handle optional parenthesis for argument
	if len(tuple.GetValues()) == 1 {
		if t, ok := tuple.GetValues()[0].(*types.Tuple); ok {
			return t, nil
		}
	}

	return tuple, nil
}

func parseWeakDelimiters(lex *lexer.TokenList) (types.Value, error) {
	if lex.Current().Type != lexer.WeakDelimiterType {
		return parseLiteral(lex.Current())
	}
	fn := func(c types.ValueContainer, end string) error {
		for lex.Next() && lex.Current().Type != lexer.WeakDelimiterType && lex.Current().Literal != end {
			if lex.Current().Type == lexer.ModifierDelimiterType ||
				lex.Current().Type == lexer.FigureDelimiterType ||
				lex.Current().Type == lexer.StatementDelimiterType {
				return errors.Join(ErrMissingLiteral, fmt.Errorf("unfinished container %v", c))
			}
			val, err := parseWeakDelimiters(lex)
			if err != nil {
				return err
			}
			if !c.CanContain(val) {
				return errors.Join(ErrInvalidLiteral, fmt.Errorf("container cannot contain %v", val))
			}
			c.AddValues(val)
		}
		return nil
	}
	switch lex.Current().Literal {
	case "(":
		tuple := new(types.Tuple)
		return tuple, fn(tuple, ")") // valid because tuple is a pointer
	case "[":
		list := new(types.List)
		return list, fn(list, "]") // valid because list is a pointer
	case "]", ")":
		return nil, errors.Join(ErrInvalidLiteral, fmt.Errorf("cannot close a container with %s", lex.Current().Literal))
	}
	return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported weak delimiters %s", lex.Current().Type))
}

func parseLiteral(lex *lexer.Lexer) (types.Value, error) {
	switch lex.Type {
	case lexer.IdentifierType:
		return types.NewDefaultLiteral(lex.Literal), nil
	case lexer.VariableType:
		//TODO: handle
	case lexer.StringType:
		return types.String(lex.Literal), nil
	case lexer.IntType:
		i, err := strconv.ParseInt(lex.Literal, 10, 64)
		if err != nil {
			return nil, err
		}
		return types.Int(i), nil
	case lexer.FloatType:
		f, err := strconv.ParseFloat(lex.Literal, 64)
		if err != nil {
			return nil, err
		}
		return types.Float(f), nil
	}
	return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported literal lex types %s", lex.Type))
}
