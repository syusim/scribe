package scalar

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

type Expr interface {
	lang.Expr

	Type() lang.Type
}

type ColRef struct {
	Id opt.ColumnID

	// We need to store this in here because the alternative is
	// passing in a context when we ask for the type and I find that
	// distasteful.
	Typ lang.Type
}

func (c *ColRef) ChildCount() int {
	return 0
}

func (c *ColRef) Child(i int) lang.Expr {
	panic("no children")
}

func (c *ColRef) Type() lang.Type {
	return c.Typ
}

type ExecColRef struct {
	Idx int
	Typ lang.Type
}

func (c *ExecColRef) ChildCount() int {
	return 0
}

func (c *ExecColRef) Child(i int) lang.Expr {
	panic("no children")
}

func (c *ExecColRef) Type() lang.Type {
	return c.Typ
}

// TODO: remove this wrapper?
type Constant struct {
	D lang.Datum
}

func (c *Constant) ChildCount() int {
	return 0
}

func (c *Constant) Child(i int) lang.Expr {
	panic("no children")
}

func (c *Constant) Type() lang.Type {
	// TODO: make this a method on Datum
	switch c.D.(type) {
	case lang.DInt:
		return lang.Int
	case lang.DString:
		return lang.String
	case lang.DBool:
		return lang.Bool
	default:
		panic(fmt.Sprintf("unhandled: %T", c.D))
	}
}

type Plus struct {
	Left  Expr
	Right Expr
}

func (e *Plus) ChildCount() int {
	return 2
}

func (e *Plus) Child(i int) lang.Expr {
	switch i {
	case 0:
		return e.Left
	case 1:
		return e.Right
	default:
		panic("out of bounds")
	}
}

func (e *Plus) Type() lang.Type {
	return lang.Int
}

type And struct {
	Left  Expr
	Right Expr
}

func (e *And) ChildCount() int {
	return 2
}

func (e *And) Child(i int) lang.Expr {
	switch i {
	case 0:
		return e.Left
	case 1:
		return e.Right
	default:
		panic("out of bounds")
	}
}

func (e *And) Type() lang.Type {
	return lang.Bool
}

type Filters struct {
	Filters []Expr
}

func (e *Filters) ChildCount() int {
	return len(e.Filters)
}

func (e *Filters) Child(i int) lang.Expr {
	return e.Filters[i]
}

func (e *Filters) Type() lang.Type {
	return lang.Bool
}

// TODO: Should these each be their own ops (probably)?
type Func struct {
	Op   lang.Func
	Args []Expr
}

func (f *Func) ChildCount() int {
	return len(f.Args)
}

func (f *Func) Child(i int) lang.Expr {
	return f.Args[i]
}

func (f *Func) Type() lang.Type {
	// TODO: make this a method on Datum
	switch f.Op {
	case lang.Eq:
		return lang.Bool
	case lang.Plus, lang.Minus, lang.Times:
		return lang.Int
	default:
		panic(fmt.Sprintf("unhandled: %v", f.Op))
	}
}
