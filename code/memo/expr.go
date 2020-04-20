package memo

import (
	"github.com/justinj/scribe/code/lang"
)

type relExpr interface {
	lang.Expr
}

// TODO: this is actually a group, but not sure how to name them appropriately.
// Also not sure how to structure this right so that other packages and inspect it sanely.
// TODO/idea: what if every expr had its children directly, but also was a thing in a map to a metadata thing?
type RelGroup struct {
	Es []relExpr

	Props Props
}

// this breaks EVERYTHING!!! but i guess it's fine
func (r *RelGroup) SetBest(e lang.Expr) {
	r.Es[0] = e
}

func (r *RelGroup) Unwrap() relExpr {
	return r.Es[0]
}

func (g *RelGroup) Add(r relExpr) {
	g.Es = append(g.Es, r)
}

func (r *RelGroup) MemberCount() int {
	return len(r.Es)
}

func (r *RelGroup) Member(i int) lang.Expr {
	return r.Es[i]
}
