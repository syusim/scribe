package explore

import (
	"fmt"

	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/scalar"
)

func (e *explorer) generateIndexScans(expr lang.Expr, add func(lang.Expr)) {
	if scan, ok := expr.(*memo.Scan); ok {
		// Only generate from the canonical one.
		if scan.Index == 0 && scan.Constraint.IsUnconstrained() {
			tab, ok := e.c.TableByName(scan.TableName)
			if !ok {
				panic(fmt.Sprintf("table gone?? %s", scan.TableName))
			}
			// Skip the first one since we already have it.
			for i, n := 1, tab.IndexCount(); i < n; i++ {
				newExpr := e.m.Scan(
					scan.TableName,
					scan.Cols,
					i,
					// TODO: make this a constant somewhere
					constraint.Constraint{},
				).Unwrap()

				add(newExpr)
			}
		}
	}
}

func (e *explorer) generateConstrainedIndexScans(expr lang.Expr, add func(lang.Expr)) {
	if sel, ok := expr.(*memo.Select); ok {
		for j, m := 0, sel.Input.MemberCount(); j < m; j++ {
			if scan, ok := sel.Input.Member(j).(*memo.Scan); ok {
				// Only generate from the canonical one.
				if scan.Index == 0 && scan.Constraint.IsUnconstrained() {
					tab, ok := e.c.TableByName(scan.TableName)
					if !ok {
						panic(fmt.Sprintf("table gone?? %s", scan.TableName))
					}
					newConstraints, newFilters := constraint.Generate(
						lang.Unwrap(sel.Filter).(*scalar.Filters),
						scan.Cols,
						tab.Index(scan.Index),
					)

					// TODO: remove this
					if newConstraints.IsUnconstrained() && scan.Index == 0 {
						continue
					}

					// OK this kind of sucks but just do it!!
					newExpr := e.m.Select(
						e.m.Scan(
							scan.TableName,
							scan.Cols,
							scan.Index,
							newConstraints,
						),
						e.m.Filters(newFilters),
					).Unwrap()

					add(newExpr)
				}
			}
		}
	}
}
