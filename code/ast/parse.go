package ast

import (
	"fmt"
	"strconv"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/sexp"
)

func ParseExpr(s string) (Expr, error) {
	t, err := sexp.Parse(s)
	if err != nil {
		return nil, err
	}
	return parseExpr(t)
}

func parseExpr(s sexp.Sexp) (Expr, error) {
	switch e := s.(type) {
	case sexp.Atom:
		switch e {
		case "true":
			return lang.DBool(true), nil
		case "false":
			return lang.DBool(false), nil
		}

		if '0' <= e[0] && e[0] <= '9' {
			n, err := strconv.Atoi(string(e))
			if err != nil {
				return nil, err
			}
			return lang.DInt(n), nil
		}
		return ColumnReference(e), nil
	case sexp.String:
		return lang.DString(e), nil
	case sexp.List:
		if len(e) == 0 {
			return nil, fmt.Errorf("empty list is not an expression")
		}
		head, ok := e[0].(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("head must be function name")
		}
		op, err := lang.GetFunc(string(head))
		if err != nil {
			return nil, err
		}
		args := make([]Expr, len(e)-1)
		for i, arg := range e[1:] {
			parsed, err := parseExpr(arg)
			if err != nil {
				return nil, err
			}
			args[i] = parsed
		}
		return ScalarFunc{op, args}, nil
	}
	panic(fmt.Sprintf("unexpected type: %T", s))
}

func parseExprList(s sexp.Sexp) ([]Expr, error) {
	l, ok := s.(sexp.List)
	if !ok {
		return nil, fmt.Errorf("expected expression list")
	}
	result := make([]Expr, len(l))
	for i := range l {
		next, err := parseExpr(l[i])
		if err != nil {
			return nil, err
		}
		result[i] = next
	}
	return result, nil
}

func ParseQuery(s string) (Node, error) {
	t, err := sexp.Parse(s)
	if err != nil {
		return nil, err
	}
	return parseQuery(t)
}

func parseQuery(s sexp.Sexp) (Node, error) {
	switch e := s.(type) {
	case sexp.Atom:
		// An atom in a query position is a table reference.
		return &TableRef{string(e)}, nil
	case sexp.String:
		return nil, fmt.Errorf("expected relational expression, found string (%q)", e)
	case sexp.List:
		if len(e) == 0 {
			return nil, fmt.Errorf("empty list is not a relational expression")
		}
		head, ok := e[0].(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("head of operator must be atom")
		}
		switch head {
		case "join":
			if len(e) != 4 {
				return nil, fmt.Errorf("join takes 3 arguments")
			}
			left, err := parseQuery(e[1])
			if err != nil {
				return nil, err
			}
			right, err := parseQuery(e[2])
			if err != nil {
				return nil, err
			}
			on, err := parseExpr(e[3])
			if err != nil {
				return nil, err
			}
			return &Join{left, right, on}, nil
		default:
			return nil, fmt.Errorf("unrecognized relational operator %s", e[0])
		}
	default:
		panic(fmt.Sprintf("unhandled type %T", s))
	}
}
