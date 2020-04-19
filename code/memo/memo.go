package memo

import "github.com/justinj/scribe/code/scalar"

type Memo struct {
	hashes      map[string]interface{}
	scalarProps map[scalar.Group]ScalarProps
}

func New() *Memo {
	return &Memo{
		hashes:      make(map[string]interface{}),
		scalarProps: make(map[scalar.Group]ScalarProps),
	}
}
