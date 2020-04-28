package explore

import (
	"fmt"

	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
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

func (e *explorer) generateHashJoins(expr lang.Expr, add func(lang.Expr)) {
	if j, ok := expr.(*memo.Join); ok {
		extraConds := make([]scalar.Group, 0)
		leftCols := make([]opt.ColumnID, 0)
		rightCols := make([]opt.ColumnID, 0)
		for _, f := range lang.Unwrap(j.On).(*scalar.Filters).Filters {
			added := false
			if eq, ok := lang.Unwrap(f).(*scalar.Eq); ok {
				if l, ok := lang.Unwrap(eq.Left).(*scalar.ColRef); ok {
					if r, ok := lang.Unwrap(eq.Right).(*scalar.ColRef); ok {
						// TODO: check that they straddle the two?
						leftCols = append(leftCols, l.Id)
						rightCols = append(rightCols, r.Id)
						added = true
					}
				}
			}
			if !added {
				extraConds = append(extraConds, f)
			}
		}
		// TODO: add the commuted version too for the ordering
		if len(leftCols) > 0 {
			add(e.m.Select(
				e.m.HashJoin(
					j.Left,
					j.Right,
					leftCols,
					rightCols,
				),
				e.m.Filters(extraConds),
			).Unwrap())
		}
	}
}
