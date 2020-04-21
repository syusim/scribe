package memo

import (
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/scalar"
)

type Memo struct {
	hashes      map[string]interface{}
	scalarProps map[scalar.Group]ScalarProps
	catalog     *cat.Catalog
}

func New(c *cat.Catalog) *Memo {
	return &Memo{
		hashes:      make(map[string]interface{}),
		scalarProps: make(map[scalar.Group]ScalarProps),
		catalog:     c,
	}
}
