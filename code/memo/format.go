package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: make a real tree printer?
func Format(g lang.Group) string {
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
		extra(&buf, e)
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

func extra(buf *bytes.Buffer, e lang.Expr) {
	switch o := e.(type) {
	case *Scan:
		buf.WriteString(" [")
		for i, c := range o.Cols {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "%d", c)
		}
		buf.WriteByte(']')
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
