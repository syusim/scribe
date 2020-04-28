package optimize

import (
	"math"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/phys"
)

// TODO: this needs a better name, it's like,
// the pair of (Group, required phys props) to be used for
// hashing.
type pair struct {
	g lang.Group
	p *phys.Props
}

type optResult struct {
	expr lang.Expr
	idx  int
	cost Cost
}

type optimizer struct {
	catalog *cat.Catalog
	// TODO: kill exprCosts? subsumed by optimized
	exprCosts map[lang.Expr]Cost
	optimized map[pair]optResult
	props     map[string]*phys.Props
	m         *memo.Memo
}

// special case phys.Min so we can just pass that around.
func (o *optimizer) internPhys(props phys.Props) *phys.Props {
	h := phys.Hash(props)
	if p, ok := o.props[h]; ok {
		return p
	}
	x := &props
	o.props[h] = x
	return o.props[h]
}

func Optimize(g lang.Group, catalog *cat.Catalog, m *memo.Memo) lang.Group {
	o := &optimizer{
		exprCosts: make(map[lang.Expr]Cost),
		optimized: make(map[pair]optResult),
		props:     make(map[string]*phys.Props),
		catalog:   catalog,
		m:         m,
	}

	o.OptimizeGroup(g, o.internPhys(phys.Min))

	return o.twiddle(g, o.internPhys(phys.Min))
}

type Cost float64

func (o *optimizer) OptimizeGroup(g lang.Group, reqd *phys.Props) {
	p := pair{g, reqd}
	if _, ok := o.optimized[p]; ok {
		return
	}

	bestCost := Cost(math.Inf(+1))
	bestIdx := -1
	var bestExpr lang.Expr
	for i, n := 0, g.MemberCount(); i < n; i++ {
		expr := g.Member(i)
		// TODO: I smell a way this could be cleaner.
		next, ok := reqd, true
		for ok {
			if o.CanProvide(expr, next) {
				o.OptimizeExpr(expr, next)

				curExpr := expr

				if reqd != next {
					// TODO we actually need to compute the "difference" operator between the two
					// sets of physical props.
					// TODO: tidy this up a little, extract it. make it nice, y'know?
					// So the wonkiness here I don't understand yet is: the Sort operator
					// technically belongs in ths same group as its input, as it has the
					// same logical properties. That's no good, however, since it means
					// it refers to itself and contains a cycle, and invaldiates a thing
					// I thought was true, which was that every relgroup shows up in a
					// final tree at most once. I think the answer is that enforcers must live
					// outside the memo, or in their own group, and then when we twiddle
					// something that is a parent of an enforcer we have to redirect its
					// pointer to the enforced version. I'm not sure how to reconcile
					// this with the definitional fact of "everything with the same
					// logical props lives in the same group."
					curExpr = &memo.Sort{
						Input:    g.(*memo.RelGroup),
						Ordering: reqd.Ordering,
					}
				}

				cost := o.ComputeCost(curExpr, next)

				if cost < bestCost {
					bestCost = cost
					bestIdx = i
					bestExpr = curExpr
				}
			}
			pro, ok := next.Weaken()
			if !ok {
				break
			}
			next = o.internPhys(pro)

			// TODO: I'm not sure this is correct in general.
			o.OptimizeGroup(g, next)
		}
	}
	if bestIdx == -1 {
		panic("didn't find a valid expression")
	}

	o.optimized[p] = optResult{
		expr: bestExpr,
		idx:  bestIdx,
	}
}

func (o *optimizer) twiddle(g lang.Group, reqd *phys.Props) lang.Group {
	res, ok := o.optimized[pair{g, reqd}]
	if !ok {
		panic("group was not optimized")
	}

	e := res.expr
	newChildren := make([]lang.Group, e.ChildCount())
	// Twiddle each subexpression.
	for i, n := 0, e.ChildCount(); i < n; i++ {
		reqdChildProps := o.ReqdPhys(e, reqd, i)
		newChildren[i] = o.twiddle(e.Child(i), reqdChildProps)
	}

	return o.m.Render(e, newChildren)
}

func (o *optimizer) OptimizeExpr(e lang.Expr, props *phys.Props) {
	// Ensure each child is optimized.
	for i, n := 0, e.ChildCount(); i < n; i++ {
		reqdChildProps := o.ReqdPhys(e, props, i)
		o.OptimizeGroup(e.Child(i), reqdChildProps)
	}
}

func (o *optimizer) ComputeCost(e lang.Expr, reqd *phys.Props) Cost {
	// We must obey Bellman!
	var cost Cost
	for i, n := 0, e.ChildCount(); i < n; i++ {
		reqdChild := o.ReqdPhys(e, reqd, i)
		childResult, ok := o.optimized[pair{e.Child(i), reqdChild}]
		if !ok {
			panic("child was not optimized")
		}
		cost += childResult.cost
	}

	if _, ok := e.(*memo.Sort); ok {
		cost += 100
	}
	if _, ok := e.(*memo.Join); ok {
		cost += 50
	}
	if _, ok := e.(*memo.Select); ok {
		cost += 10
	}
	cost += 1

	o.exprCosts[e] = cost
	return o.exprCosts[e]
}
