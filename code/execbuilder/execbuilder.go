package execbuilder

import (
	"fmt"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/memo"
)

type builder struct {
	cat *cat.Catalog
}

func New(cat *cat.Catalog) *builder {
	return &builder{
		cat: cat,
	}
}

func (b *builder) Build(e memo.RelExpr) (exec.Node, error) {
	switch o := e.E.(type) {
	case *memo.Scan:
		tab, ok := b.cat.TableByName(o.TableName)
		if !ok {
			// This should have been verified already.
			panic(fmt.Sprintf("table %q not found", o.TableName))
		}
		if tab.IndexCount() == 0 {
			// TODO: this should be ensured to be impossible.
			panic("no indexes buddy!")
		}
		// TODO: pass in which one to use
		idx := tab.Index(0)
		iter := idx.Scan()

		return exec.Scan(iter), nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e.E))
	}
}
