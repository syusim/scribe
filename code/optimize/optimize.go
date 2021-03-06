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
	catalog   *cat.Catalog
	optimized map[pair]optResult
	props     map[string]*phys.Props
	m         *memo.Memo
}

// special case phys.Min so we can just pass that around.
func (o *optimizer) internPhys(props phys.Props) *phys.Props {
	if len(props.Ordering) == 0 {
		return phys.Min
	}
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
		optimized: make(map[pair]optResult),
		props:     make(map[string]*phys.Props),
		catalog:   catalog,
		m:         m,
	}

	o.OptimizeGroup(g, phys.Min)

	return o.twiddle(g, phys.Min)
}

type Cost float64

// Enforce computes the enforcer operator which brings an expression satisfying
// physical props from to physical props to.
func enforce(from, to *phys.Props, e *memo.RelGroup) lang.Expr {
	return &memo.Sort{
		Input:    e,
		Ordering: to.Ordering,
	}
}

// OptimizeGroup optimizes group g relative to physical props reqd.
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

				withProps := expr
				if reqd != next {
					withProps = enforce(next, reqd, g.(*memo.RelGroup))
				}

				cost := o.ComputeCost(withProps, next)

				if cost < bestCost {
					bestCost = cost
					bestIdx = i
					bestExpr = withProps
				}
			}
			weaker, ok := next.Weaken()
			if !ok {
				break
			}
			next = weaker

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

	return cost
}
