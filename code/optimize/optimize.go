package optimize

import (
	"fmt"
	"math"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
)

type optimizer struct {
	exprCosts map[lang.Expr]Cost
	optimized map[lang.Group]struct{}
}

func Optimize(g lang.Group) {
	o := &optimizer{
		exprCosts: make(map[lang.Expr]Cost),
		optimized: make(map[lang.Group]struct{}),
	}

	o.OptimizeGroup(g)
}

type Cost float64

func (o *optimizer) OptimizeGroup(g lang.Group) {
	if _, ok := o.optimized[g]; ok {
		return
	}
	o.optimized[g] = struct{}{}

	bestCost := Cost(math.Inf(+1))
	bestIdx := -1
	for i, n := 0, g.MemberCount(); i < n; i++ {
		o.OptimizeExpr(g.Member(i))

		cost := o.ComputeCost(g.Member(i))
		if cost < bestCost {
			bestCost = cost
			bestIdx = i
		}
	}
	if bestIdx == -1 {
		panic("didn't find a valid expression")
	}

	// Now ratchet the group.
	if rel, ok := g.(*memo.RelGroup); ok {
		rel.Ratchet(bestIdx)
	} else if g.MemberCount() > 1 {
		panic(fmt.Sprintf("don't know how to ratchet a %T", g))
	}
}

func (o *optimizer) OptimizeExpr(e lang.Expr) {
	// Ensure each child is optimized.
	for i, n := 0, e.ChildCount(); i < n; i++ {
		o.OptimizeGroup(e.Child(i))
	}
}

func (o *optimizer) ComputeCost(e lang.Expr) Cost {
	// We must obey Bellman!
	var cost Cost
	// Invariant: the children here have been optimized, so the Unwrap expression
	// is the best one.
	for i, n := 0, e.ChildCount(); i < n; i++ {
		cost += o.exprCosts[lang.Unwrap(e.Child(i))]
	}

	cost += 1

	o.exprCosts[e] = cost
	return o.exprCosts[e]
}
