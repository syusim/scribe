package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type formatter struct {
}

// TODO: make a real tree printer?
func Format(e Expr) string {
	var buf bytes.Buffer
	depth := 0
	var p func(e Expr)
	p = func(e Expr) {
		for i := 0; i < depth; i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString("-> ")
		spew.Dump(e)
		if r, ok := e.(RelExpr); ok {
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

func extra(buf *bytes.Buffer, e Expr) {
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
	case *Func:
		fmt.Fprintf(buf, " (%s)", o.Op)
	case *Constant:
		fmt.Fprintf(buf, " (%s)", o.D)
	case *ColRef:
		fmt.Fprintf(buf, " (%d)", o.Id)
	}
}
