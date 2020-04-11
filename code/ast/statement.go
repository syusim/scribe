package ast

import (
	"bytes"
	"fmt"

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

// TODO: we should maybe just use straight up lang.Datums
// everywhere.
type Datum struct {
	d lang.Datum
}

func (d *Datum) Format(buf *bytes.Buffer) {
	switch e := d.d.(type) {
	case lang.DInt:
		fmt.Fprintf(buf, "%d", e)
	case lang.DString:
		fmt.Fprintf(buf, "%q", e)
	case lang.DBool:
		if e {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	}
}

type Row []Datum

func (r Row) Format(buf *bytes.Buffer) {
	buf.WriteByte('[')
	for i, d := range r {
		if i > 0 {
			buf.WriteByte(' ')
		}
		d.Format(buf)
	}
	buf.WriteByte(']')
}

type CreateTable struct {
	Name    string
	Columns []lang.Column
	Data    []Row
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
