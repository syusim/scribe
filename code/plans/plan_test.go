package plan_test

import (
	"fmt"
	"testing"

	"github.com/justinj/bitwise/datadriven"
	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/builder"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/explore"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/optimize"
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
					if err := catalog.AddTable(c); err != nil {
						return fmt.Sprintf("error: %s\n", err)
					}
					return "ok\n"
				default:
					panic("unhandled")
				}
			case "plan-scalar":
				expr, err := ast.ParseExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				m := memo.New(catalog)
				b := builder.New(catalog, m)
				e, err := b.BuildScalar(expr, nil)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				return memo.Format(m, e)
			case "plan":
				expr, err := ast.ParseRelExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				m := memo.New(catalog)
				b := builder.New(catalog, m)
				e, _, err := b.Build(expr)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				return memo.Format(m, e)
			case "plan-memo":
				expr, err := ast.ParseRelExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				m := memo.New(catalog)
				b := builder.New(catalog, m)
				e, _, err := b.Build(expr)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}

				explore.Explore(m, catalog, e)

				return m.Format(e)
			case "plan-full":
				expr, err := ast.ParseRelExpr(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				m := memo.New(catalog)
				b := builder.New(catalog, m)
				e, _, err := b.Build(expr)

				explore.Explore(m, catalog, e)
				g := optimize.Optimize(e, catalog, m)

				if err != nil {
					return fmt.Sprintf("error: %s\n", err)
				}
				return memo.Format(m, g)
			default:
				panic("unhandled")
			}
		})
	})
}
