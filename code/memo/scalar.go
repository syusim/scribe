package memo

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

type ScalarExpr interface {
	Expr

	Type() lang.Type
}

type ColRef struct {
	Id opt.ColumnID

	// We need to store this in here because the alternative is
	// passing in a context when we ask for the type and I find that
	// distasteful.
	Typ lang.Type
}

func (c *ColRef) Type() lang.Type {
	return c.Typ
}

// TODO: remove this wrapper
type Constant struct {
	D lang.Datum
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

// TODO: Should these each be their own ops (probably)?
type Func struct {
	Op   lang.Func
	Args []ScalarExpr
}

func (f *Func) Type() lang.Type {
	// TODO: make this a method on Datum
	switch f.Op {
	case lang.Eq:
		return lang.Bool
	default:
		panic(fmt.Sprintf("unhandled: %v", f.Op))
	}
}
