package explore

import (
	"fmt"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/scalar"
)

type explorer struct {
	m *memo.Memo
	c *cat.Catalog

	explored map[lang.Group]struct{}
}

func Explore(m *memo.Memo, c *cat.Catalog, r *memo.RelGroup) {
	e := &explorer{
		m:        m,
		c:        c,
		explored: make(map[lang.Group]struct{}),
	}

	e.ExploreGroup(r)
}

func (e *explorer) ExploreGroup(r lang.Group) {
	if _, ok := e.explored[r]; ok {
		return
	}
	e.explored[r] = struct{}{}

	q := []lang.Expr{}
	for i, n := 0, r.MemberCount(); i < n; i++ {
		q = append(q, r.Member(i))
	}
	for len(q) > 0 {
		var next lang.Expr
		next, q = q[0], q[1:]

		for i, n := 0, next.ChildCount(); i < n; i++ {
			e.ExploreGroup(next.Child(i))
		}

		add := func(e lang.Expr) {
			// TODO: I think this is broken right now if we happen to construct this
			// expression via normalization later.
			q = append(q, e)
			r.(*memo.RelGroup).Add(e)
		}

		if scan, ok := next.(*memo.Scan); ok {
			// Only generate from the canonical one.
			if scan.Index == 0 && scan.Constraint.IsUnconstrained() {
				tab, ok := e.c.TableByName(scan.TableName)
				if !ok {
					panic(fmt.Sprintf("table gone?? %s", scan.TableName))
				}
				// Skip the first one since we already have it.
				for i, n := 1, tab.IndexCount(); i < n; i++ {
					newExpr := e.m.Scan(
						scan.TableName,
						scan.Cols,
						i,
						// TODO: make this a constant somewhere
						constraint.Constraint{},
					).Unwrap()

					add(newExpr)
				}
			}
		}

		if sel, ok := next.(*memo.Select); ok {
			for j, m := 0, sel.Input.MemberCount(); j < m; j++ {
				if scan, ok := sel.Input.Member(j).(*memo.Scan); ok {
					// Only generate from the canonical one.
					if scan.Index == 0 && scan.Constraint.IsUnconstrained() {
						tab, ok := e.c.TableByName(scan.TableName)
						if !ok {
							panic(fmt.Sprintf("table gone?? %s", scan.TableName))
						}
						newConstraints, newFilters := constraint.Generate(
							lang.Unwrap(sel.Filter).(*scalar.Filters),
							scan.Cols,
							tab.Index(scan.Index),
						)

						// TODO: remove this
						if newConstraints.IsUnconstrained() && scan.Index == 0 {
							continue
						}

						// OK this kind of sucks but just do it!!
						newExpr := e.m.Select(
							e.m.Scan(
								scan.TableName,
								scan.Cols,
								scan.Index,
								newConstraints,
							),
							e.m.Filters(newFilters),
						).Unwrap()

						add(newExpr)
					}
				}
			}
		}
	}
}
