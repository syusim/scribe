package memo

type Expr interface {
	ChildCount() int
	Child(i int) Expr
}

// TODO: this is actually a group, but not sure how to name them appropriately.
// Also not sure how to structure this right so that other packages and inspect it sanely.
type RelExpr struct {
	E relExpr
}

// TODO: these types seem extremely wonky.
// should Child/ChildCount just be methods on
// RelExpr which then defer to a big switch for
// the op? regardless it seems they shouldn't
// BOTH be like this.
func (r RelExpr) ChildCount() int {
	return r.E.ChildCount()
}

func (r RelExpr) Child(i int) Expr {
	return r.E.Child(i)
}

type relExpr interface {
	Expr
}
