package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

type formatter struct {
}

// TODO: make a real tree printer?
func Format(e lang.Expr) string {
	var buf bytes.Buffer
	depth := 0
	var p func(e lang.Expr)
	p = func(e lang.Expr) {
		for i := 0; i < depth; i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString("-> ")
		if r, ok := e.(*RelExpr); ok {
			e = r.E
		}
		buf.WriteString(reflect.TypeOf(e).Elem().Name())
		extra(&buf, e)
		buf.WriteByte('\n')
		depth++
		for i, n := 0, e.ChildCount(); i < n; i++ {
			p(e.Child(i))
		}
		depth--
	}

	p(e)

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
	case *scalar.Func:
		fmt.Fprintf(buf, " (%s)", o.Op)
	case *scalar.Constant:
		fmt.Fprintf(buf, " (%s)", o.D)
	case *scalar.ColRef:
		fmt.Fprintf(buf, " (%d)", o.Id)
	}
}
