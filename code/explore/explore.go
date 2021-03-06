package explore

import (
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

		add := func(e lang.Expr) {
			// TODO: I think this is broken right now if we happen to construct this
			// expression via normalization later.
			q = append(q, e)
			r.(*memo.RelGroup).Add(e)
		}

		e.generateIndexScans(next, add)
		e.generateConstrainedIndexScans(next, add)
		e.generateHashJoins(next, add)
	}
}
