package memo

import (
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

func eqConst(s scalar.Expr, d lang.Datum) bool {
	if c, ok := s.(*scalar.Constant); ok {
		if c.D.Type() != d.Type() {
			return false
		}
		return lang.Compare(c.D, d) == lang.EQ
	}
	return false
}

func concat(a, b *scalar.Filters) []scalar.Expr {
	return append(append(make([]scalar.Expr, 0), a.Filters...), b.Filters...)
}

func unfoldFilters(f []scalar.Expr) []scalar.Expr {
	var newFilters []scalar.Expr
	for i, c := range f {
		if a, ok := c.(*scalar.And); ok {
			if newFilters == nil {
				newFilters = make([]scalar.Expr, 0)
				newFilters = append(newFilters, f[:i]...)
			}
			// TODO: this is sort of inefficient, since we're relying on the rule to
			// re-trigger, possibly many times. better would be to walk the And tree
			// and extract all the conjuncts.
			newFilters = append(newFilters, a.Left, a.Right)
		} else if newFilters != nil {
			newFilters = append(newFilters, c)
		}
	}

	return newFilters
}

func extractBoundUnbound(
	m *Memo,
	filters []scalar.Expr,
	cols opt.ColSet,
) ([]scalar.Expr, []scalar.Expr) {
	canPush := false
	for _, f := range filters {
		freeVars := m.GetScalarProps(f).FreeVars
		if freeVars.SubsetOf(cols) {
			canPush = true
			break
		}
	}

	if !canPush {
		return nil, filters
	}

	var bound []scalar.Expr
	var unbound []scalar.Expr
	for _, f := range filters {
		freeVars := m.GetScalarProps(f).FreeVars
		if freeVars.SubsetOf(cols) {
			bound = append(bound, f)
		} else {
			unbound = append(unbound, f)
		}
	}
	return bound, unbound
}

func inlineIn(
	m *Memo,
	e scalar.Expr,
	projs []scalar.Expr,
	ids []opt.ColumnID,
) scalar.Expr {
	return exec.ScalarExpr(m.Walk(e, func(in lang.Expr) lang.Expr {
		if ref, ok := in.(*scalar.ColRef); ok {
			for i, col := range ids {
				if col == ref.Id {
					return projs[i]
				}
			}
		}
		return in
	}).(scalar.Expr))
}