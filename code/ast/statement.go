package ast

import (
	"bytes"

	"github.com/justinj/scribe/code/lang"
)

type Statement interface {
	Format(buf *bytes.Buffer)
}

type RunQuery struct {
	Input RelExpr
}

func (r *RunQuery) Format(buf *bytes.Buffer) {
	buf.WriteString("(run ")
	r.Input.Format(buf)
	buf.WriteByte(')')
}

type CreateTable struct {
	Name    string
	Columns []lang.Column
	Data    []lang.Row
}

func (c *CreateTable) Format(buf *bytes.Buffer) {
	buf.WriteString("(create-table ")
	buf.WriteString(c.Name)
	buf.WriteString(" [")
	for i, col := range c.Columns {
		if i > 0 {
			buf.WriteByte(' ')
		}
		col.Format(buf)
	}
	buf.WriteString("] [")
	for i, row := range c.Data {
		if i > 0 {
			buf.WriteByte(' ')
		}
		row.Format(buf)
	}
	buf.WriteByte(']')
	buf.WriteByte(')')
}
