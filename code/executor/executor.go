package executor

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/cat"
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
			nil,
			nil,
		)
		return Result{Msg: "ok"}, nil
	default:
		return Result{}, fmt.Errorf("unhandled statement %T", stmt)
	}
}
