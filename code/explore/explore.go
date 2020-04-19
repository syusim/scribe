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

		if sel, ok := next.(*memo.Select); ok {
			// TODO: We have to iterate over each member rather than unwrapping here.
			if scan, ok := sel.Input.Unwrap().(*memo.Scan); ok {
				// Only generate from the canonical one.
				if scan.Index == 0 && scan.Constraint.IsUnconstrained() {
					tab, ok := e.c.TableByName(scan.TableName)
					if !ok {
						panic(fmt.Sprintf("table gone?? %s", scan.TableName))
					}
					for i, n := 0, tab.IndexCount(); i < n; i++ {
						newConstraints, newFilters := constraint.Generate(
							lang.Unwrap(sel.Filter).(*scalar.Filters),
							scan.Cols,
							tab.Index(i),
						)

						// TODO: remove this
						if newConstraints.IsUnconstrained() && i == 0 {
							continue
						}

						// OK this kind of sucks but just do it!!
						newExpr := e.m.Select(
							e.m.Scan(
								scan.TableName,
								scan.Cols,
								i,
								newConstraints,
							),
							e.m.Filters(newFilters),
						).Unwrap()

						// TODO: add to queue!
						r.(*memo.RelGroup).Add(newExpr)
					}
				}
			}
		}
	}
}
