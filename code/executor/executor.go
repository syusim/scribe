package executor

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/builder"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/execbuilder"
	"github.com/justinj/scribe/code/explore"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/optimize"
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
		err := e.catalog.AddTable(c)
		if err != nil {
			return Result{}, fmt.Errorf("error: %s", err)
		}
		return Result{Msg: "ok"}, nil
	case *ast.RunQuery:
		mem := memo.New(e.catalog)

		b := builder.New(e.catalog, mem)
		rel, scope, err := b.Build(c.Input)
		if err != nil {
			return Result{}, err
		}

		explore.Explore(mem, e.catalog, rel)
		optimize.Optimize(rel, e.catalog, mem)

		// The relational representation of the plan doesn't have a notion of the
		// ordering of columns, however, that information is encoded in the order
		// of the columns stored in the final outScope.

		eb := execbuilder.New(e.catalog, mem)
		plan, err := eb.Build(rel, scope.OutCols())
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
