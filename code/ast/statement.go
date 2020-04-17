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

type IndexDef struct {
	Name string
	Cols []string
}

func (d *IndexDef) Format(buf *bytes.Buffer) {
	buf.WriteByte('[')
	buf.WriteString(d.Name)
	buf.WriteString(" [")
	for i, c := range d.Cols {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(c)
	}
	buf.WriteString("]]")
}

type CreateTable struct {
	Name    string
	Columns []lang.Column
	Data    []lang.Row
	Indexes []IndexDef
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
	if len(c.Indexes) > 0 {
		buf.WriteString(" [")
		for i, idx := range c.Indexes {
			if i > 0 {
				buf.WriteByte(' ')
			}
			idx.Format(buf)
		}
		buf.WriteByte(']')
	}
	buf.WriteByte(')')
}
