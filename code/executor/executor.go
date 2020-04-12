package executor

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/builder"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/execbuilder"
	"github.com/justinj/scribe/code/opt"
)

type executor struct {
	catalog *cat.Catalog
}

func New() *executor {
	return &executor{
		catalog: cat.New(),
	}
}

type Result struct {
	Msg string
}

func (e *executor) Run(cmd string) (Result, error) {
	stmt, err := ast.ParseStatement(cmd)
	if err != nil {
		return Result{}, err
	}

	switch c := stmt.(type) {
	case *ast.CreateTable:
		e.catalog.AddTable(
			c.Name,
			c.Columns,
			c.Data,
			// Just have one empty index.
			[][]opt.ColOrdinal{{}},
		)
		return Result{Msg: "ok"}, nil
	case *ast.RunQuery:
		b := builder.New(e.catalog)
		rel, _, err := b.Build(c.Input)
		if err != nil {
			return Result{}, err
		}

		eb := execbuilder.New(e.catalog)
		plan, _, err := eb.Build(rel)
		if err != nil {
			return Result{}, err
		}

		rows := exec.Spool(plan)

		var buf bytes.Buffer
		for i, row := range rows {
			buf.WriteByte('[')
			for j, d := range row {
				if j > 0 {
					buf.WriteByte(' ')
				}
				d.Format(&buf)
			}
			buf.WriteByte(']')
			if i != len(rows)-1 {
				buf.WriteByte('\n')
			}
		}

		return Result{Msg: buf.String()}, nil
	default:
		return Result{}, fmt.Errorf("unhandled statement %T", stmt)
	}
}
