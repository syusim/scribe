package ast

import (
	"bytes"
)

type Node interface {
	Format(buf *bytes.Buffer)
}

type TableRef struct {
	Name string
}

func (t *TableRef) Format(buf *bytes.Buffer) {
	buf.WriteString(t.Name)
}

type Join struct {
	Left  Node
	Right Node
	On    Expr
}

func (j *Join) Format(buf *bytes.Buffer) {
	buf.WriteString("(join ")
	j.Left.Format(buf)
	buf.WriteByte(' ')
	j.Right.Format(buf)
	buf.WriteByte(' ')
	j.On.Format(buf)
	buf.WriteString(")")
}
