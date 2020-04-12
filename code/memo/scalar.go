package memo

import "github.com/justinj/scribe/code/lang"

type ScalarExpr interface {
	Expr
}

type Constant struct {
	D lang.Datum
}

// TODO: Should these each be their own ops (probably)?
type Func struct {
	Op   lang.Func
	Args []ScalarExpr
}
