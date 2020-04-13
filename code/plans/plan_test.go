package plan_test

import (
	"fmt"
	"testing"

	"github.com/justinj/bitwise/datadriven"
	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/builder"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
)

func TestNorm(t *testing.T) {
	datadriven.Walk(t, "testdata/norm", func(t *testing.T, path string) {
		catalog := cat.New()
		datadriven.RunTest(t, path, func(td *datadriven.TestData) string {
			switch td.Cmd {
			case "ddl":
				stmt, err := ast.ParseStatement(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				switch c := stmt.(type) {
				case *ast.CreateTable:
					catalog.AddTable(
						c.Name,
						c.Columns,
						c.Data,
						// Just have one empty index.
						[][]opt.ColOrdinal{{}},
					)
					return "ok\n"
				default:
					panic("unhandled")
				}
			case "plan-scalar":
				expr, err := ast.ParseExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				b := builder.New(catalog)
				e, err := b.BuildScalar(expr, nil)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				return memo.Format(e)
			case "plan":
				expr, err := ast.ParseRelExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				b := builder.New(catalog)
				e, _, err := b.Build(expr)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				return memo.Format(e)
			default:
				panic("unhandled")
			}
		})
	})
}
