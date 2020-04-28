package scalar

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

type Group interface {
	lang.Expr
	lang.Group

	Type() lang.Type
}

// type Expr interface {
// 	// Scalar Exprs are _both_. They're groups of just themselves.
// 	Group
// }

type ColRef struct {
	Id lang.ColumnID

	// We need to store this in here because the alternative is
	// passing in a context when we ask for the type and I find that
	// distasteful.
	Typ lang.Type
}

func (c *ColRef) ChildCount() int {
	return 0
}

func (c *ColRef) Child(i int) lang.Group {
	panic("no children")
}

func (c *ColRef) Type() lang.Type {
	return c.Typ
}

func (c *ColRef) MemberCount() int {
	return 1
}

func (c *ColRef) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type ExecColRef struct {
	Idx int
	Typ lang.Type
}

func (c *ExecColRef) ChildCount() int {
	return 0
}

func (c *ExecColRef) Child(i int) lang.Group {
	panic("no children")
}

func (c *ExecColRef) Type() lang.Type {
	return c.Typ
}

func (c *ExecColRef) MemberCount() int {
	return 1
}

func (c *ExecColRef) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

// TODO: remove this wrapper?
type Constant struct {
	D lang.Datum
}

func (c *Constant) ChildCount() int {
	return 0
}

func (c *Constant) Child(i int) lang.Group {
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

func (c *Constant) MemberCount() int {
	return 1
}

func (c *Constant) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type Plus struct {
	Left  Group
	Right Group
}

func (e *Plus) ChildCount() int {
	return 2
}

func (e *Plus) Child(i int) lang.Group {
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

func (c *Plus) MemberCount() int {
	return 1
}

func (c *Plus) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type Times struct {
	Left  Group
	Right Group
}

func (e *Times) ChildCount() int {
	return 2
}

func (e *Times) Child(i int) lang.Group {
	switch i {
	case 0:
		return e.Left
	case 1:
		return e.Right
	default:
		panic("out of bounds")
	}
}

func (e *Times) Type() lang.Type {
	return lang.Int
}

func (c *Times) MemberCount() int {
	return 1
}

func (c *Times) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type And struct {
	Left  Group
	Right Group
}

func (e *And) ChildCount() int {
	return 2
}

func (e *And) Child(i int) lang.Group {
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

func (c *And) MemberCount() int {
	return 1
}

func (c *And) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type Eq struct {
	Left  Group
	Right Group
}

func (e *Eq) ChildCount() int {
	return 2
}

func (e *Eq) Child(i int) lang.Group {
	switch i {
	case 0:
		return e.Left
	case 1:
		return e.Right
	default:
		panic("out of bounds")
	}
}

func (e *Eq) Type() lang.Type {
	return lang.Bool
}

func (c *Eq) MemberCount() int {
	return 1
}

func (c *Eq) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

type Filters struct {
	Filters []Group
}

func (e *Filters) ChildCount() int {
	return len(e.Filters)
}

func (e *Filters) Child(i int) lang.Group {
	return e.Filters[i]
}

func (e *Filters) Type() lang.Type {
	return lang.Bool
}

func (c *Filters) MemberCount() int {
	return 1
}

func (c *Filters) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}

// TODO: Should these each be their own ops (probably)?
type Func struct {
	Op   lang.Func
	Args []Group
}

func (f *Func) ChildCount() int {
	return len(f.Args)
}

func (f *Func) Child(i int) lang.Group {
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

func (c *Func) MemberCount() int {
	return 1
}

func (c *Func) Member(i int) lang.Expr {
	if i != 0 {
		panic("out of bounds")
	}
	return c
}
