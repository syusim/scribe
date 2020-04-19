package explore

import (
	"fmt"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
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

		if s, ok := next.(*memo.Scan); ok {
			// Only generate from the canonical one.
			if s.Index == 0 {
				tab, ok := e.c.TableByName(s.TableName)
				if !ok {
					panic(fmt.Sprintf("table gone?? %s", s.TableName))
				}
				for i, n := 1, tab.IndexCount(); i < n; i++ {
					// OK this kind of sucks but just do it!!
					newExpr := e.m.Scan(
						s.TableName,
						s.Cols,
						i,
					).Unwrap()
					// TODO: add to queue?

					r.(*memo.RelGroup).Add(newExpr)
				}
			}
		}
	}
}
