package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

func eqConst(s scalar.Group, d lang.Datum) bool {
	if c, ok := s.(*scalar.Constant); ok {
		if c.D.Type() != d.Type() {
			return false
		}
		return lang.Compare(c.D, d) == lang.EQ
	}
	return false
}

func concat(a, b *scalar.Filters) []scalar.Group {
	return append(append(make([]scalar.Group, 0), a.Filters...), b.Filters...)
}

func unfoldFilters(f []scalar.Group) []scalar.Group {
	var newFilters []scalar.Group
	for i, c := range f {
		if a, ok := c.(*scalar.And); ok {
			if newFilters == nil {
				newFilters = make([]scalar.Group, 0)
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
	filters []scalar.Group,
	cols opt.ColSet,
) ([]scalar.Group, []scalar.Group) {
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

	var bound []scalar.Group
	var unbound []scalar.Group
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
	e scalar.Group,
	projs []scalar.Group,
	ids []opt.ColumnID,
) scalar.Group {
	return m.Walk(e, func(in lang.Group) lang.Group {
		if ref, ok := in.(*scalar.ColRef); ok {
			for i, col := range ids {
				if col == ref.Id {
					return projs[i]
				}
			}
		}
		return in
	}).(scalar.Group)
}
