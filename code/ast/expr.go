package ast

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

type Expr interface {
	Node
}

func ExprStr(e Expr) string {
	var buf bytes.Buffer
	e.Format(&buf)
	return buf.String()
}

type ColumnReference string

func (c ColumnReference) Format(buf *bytes.Buffer) {
	buf.WriteString(string(c))
}

type ScalarFunc struct {
	Op   lang.Func
	Args []Expr
}

func (f ScalarFunc) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "(%s", f.Op.String())
	for _, arg := range f.Args {
		buf.WriteByte(' ')
		arg.Format(buf)
	}
	buf.WriteByte(')')
}
