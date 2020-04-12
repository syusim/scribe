package ast

import (
	"bytes"
)

type RelExpr interface {
	Format(buf *bytes.Buffer)
}

type TableRef struct {
	Name string
}

func (t *TableRef) Format(buf *bytes.Buffer) {
	buf.WriteString(t.Name)
}

type Join struct {
	Left  RelExpr
	Right RelExpr
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

type Select struct {
	Input     RelExpr
	Predicate Expr
}

func (s *Select) Format(buf *bytes.Buffer) {
	buf.WriteString("(select ")
	s.Input.Format(buf)
	buf.WriteByte(' ')
	s.Predicate.Format(buf)
	buf.WriteByte(')')
}

type Project struct {
	Input RelExpr
	Exprs []Expr
	// TODO: round-trip this field
	Aliases []string
}

func (p *Project) Format(buf *bytes.Buffer) {
	buf.WriteString("(project ")
	p.Input.Format(buf)
	buf.WriteString(" [")
	for i, e := range p.Exprs {
		if i > 0 {
			buf.WriteByte(' ')
		}
		e.Format(buf)
	}
	buf.WriteString("])")
}

type As struct {
	Input    RelExpr
	Name     string
	ColNames []string
}

func (a *As) Format(buf *bytes.Buffer) {
	buf.WriteString("(as ")
	a.Input.Format(buf)
	buf.WriteByte(' ')
	buf.WriteString(a.Name)
	if a.ColNames != nil {
		buf.WriteString(" [")
		for i, n := range a.ColNames {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(n)
		}
		buf.WriteByte(']')
	}
	buf.WriteString(")")
}
