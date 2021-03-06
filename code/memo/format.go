package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: make a real tree printer?
func Format(m *Memo, g lang.Group) string {
	var buf bytes.Buffer
	depth := 0
	var p func(g lang.Group)
	p = func(g lang.Group) {
		for i := 0; i < depth; i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString("-> ")
		e := lang.Unwrap(g)
		buf.WriteString(reflect.TypeOf(e).Elem().Name())
		m.extra(&buf, e)
		buf.WriteByte('\n')
		depth++
		for i, n := 0, e.ChildCount(); i < n; i++ {
			p(e.Child(i))
		}
		depth--
	}

	p(g)

	return buf.String()
}

// TODO: make a real tree printer?
func (m *Memo) Format(g lang.Group) string {
	queue := []lang.Group{g}

	ids := make(map[lang.Group]int)
	nextId := 1
	getId := func(g lang.Group) int {
		if id, ok := ids[g]; ok {
			return id
		}
		ids[g] = nextId
		nextId++
		return ids[g]
	}
	seen := make(map[lang.Group]struct{})
	enqueue := func(g lang.Group) {
		if _, ok := seen[g]; ok {
			return
		}
		seen[g] = struct{}{}
		queue = append(queue, g)
	}
	deque := func() (lang.Group, bool) {
		if len(queue) == 0 {
			return nil, false
		}
		var next lang.Group
		next, queue = queue[0], queue[1:]
		return next, true
	}

	var buf bytes.Buffer

	for next, ok := deque(); ok; next, ok = deque() {
		fmt.Fprintf(&buf, "G%d\n", getId(next))
		for i, n := 0, next.MemberCount(); i < n; i++ {
			expr := next.Member(i)
			fmt.Fprintf(&buf, "  - %s", reflect.TypeOf(expr).Elem().Name())

			for j, m := 0, expr.ChildCount(); j < m; j++ {
				c := expr.Child(j)
				fmt.Fprintf(&buf, " G%d", getId(c))
				enqueue(c)
			}
			m.extra(&buf, expr)
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

func (m *Memo) extra(buf *bytes.Buffer, e lang.Expr) {
	switch o := e.(type) {
	case *Scan:
		buf.WriteString(" [")
		for i, c := range o.Cols {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "%d", c)
		}
		buf.WriteString("] ")

		tab, ok := m.catalog.TableByName(o.TableName)
		if !ok {
			panic("no table")
		}
		idx := tab.Index(o.Index)

		fmt.Fprintf(buf, "@%s", idx.Name)
		if o.Constraint.Start != nil || o.Constraint.End != nil {
			buf.WriteByte(' ')
			o.Constraint.Format(buf)
		}
	case *HashJoin:
		buf.WriteString(" [")
		for i := range o.LeftCols {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "%d=%d", o.LeftCols[i], o.RightCols[i])
		}
		buf.WriteString("]")
	case *Project:
		buf.WriteString(" [")
		for i, c := range o.ColIDs {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "%d", c)
		}
		buf.WriteString("] ")
		buf.WriteString(o.PassthroughCols.String())
	case *Sort:
		if len(o.Ordering) > 0 {
			buf.WriteString(" (ordering: [")
			for i, c := range o.Ordering {
				if i > 0 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(buf, "%d", c)
			}
			buf.WriteString("]) ")
		}
	case *Root:
		if len(o.Ordering) > 0 {
			buf.WriteString(" (required ordering: [")
			for i, c := range o.Ordering {
				if i > 0 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(buf, "%d", c)
			}
			buf.WriteString("]) ")
		}
	case *scalar.Func:
		fmt.Fprintf(buf, " (%s)", o.Op)
	case *scalar.Constant:
		fmt.Fprintf(buf, " (%s)", o.D)
	case *scalar.ColRef:
		fmt.Fprintf(buf, " (%d)", o.Id)
	}
}
