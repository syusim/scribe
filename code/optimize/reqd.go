package optimize

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/phys"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: re-evaluate if these should be public. i feel like they could just be
// their own package even, or live in memo?
func (o *optimizer) CanProvide(e lang.Expr, p *phys.Props) bool {
	switch e := e.(type) {
	case *memo.Select:
		// Can provide it by requiring it of the input.
		return true
	case *memo.Project:
		// Can provide it by requiring it of the input.
		return true
	case *memo.HashJoin:
		// TODO: we can provide the right's ordering!
		// Can't provide anything!
		return len(p.Ordering) == 0
	case *memo.Join:
		// Can't provide anything!
		return len(p.Ordering) == 0
	case *memo.Scan:
		// It can if the required ordering is a prefix of the indexed columns.
		// TODO: make this neater
		tab, ok := o.catalog.TableByName(e.TableName)
		if !ok {
			panic("no good chief")
		}
		idx := tab.Index(e.Index)
		for i, col := range p.Ordering {
			// TODO: use a better algorithm here
			ord := -1
			for j := range e.Cols {
				if e.Cols[j] == col {
					ord = j
					break
				}
			}
			if i >= len(idx.Ordering) || idx.Ordering[i] != lang.ColOrdinal(ord) {
				return false
			}
		}
		return true
	default:
		// skeptical here...
		return true
	}
}

func (o *optimizer) ReqdPhys(e lang.Expr, p *phys.Props, i int) *phys.Props {
	switch e := e.(type) {
	case *memo.Root:
		// We expect the requried props here to be nothing.
		inputProps := o.internPhys(phys.Props{
			Ordering: e.Ordering,
		})

		return inputProps
	case *memo.Join:
		return o.internPhys(phys.Min)
	case *memo.HashJoin:
		return o.internPhys(phys.Min)
	case *memo.Sort:
		return o.internPhys(phys.Min)
	case *memo.Select:
		if i == 0 {
			return p
		} else {
			return o.internPhys(phys.Min)
		}
	case *memo.Project:
		if i == 0 {
			return p
		} else {
			return o.internPhys(phys.Min)
		}
	case scalar.Group:
		return o.internPhys(phys.Min)
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
