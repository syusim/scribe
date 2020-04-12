package builder

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/cat"
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

func (b *builder) Build(e ast.RelExpr) (memo.RelExpr, error) {
	switch a := e.(type) {
	case *ast.TableRef:
		// TODO: look it up in the catalog.
		return memo.Wrap(&memo.Scan{
			TableName: a.Name,
		}), nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
