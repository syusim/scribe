package index

import (
	"sort"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

type order []opt.ColOrdinal

type T struct {
	data    []lang.Row
	orderBy order
}

//(index.iterator-def
type Iterator struct {
	index *T
	pos   int
} //)

type cmpResult int

const (
	lt cmpResult = -1
	eq           = 0
	gt           = 1
)

// TODO: this needs to use lang.Compare
func compare(a, b lang.Row, orderBy order) cmpResult {
	for _, idx := range orderBy {
		cmp := lang.Compare(a[idx], b[idx])
		if cmp == lang.LT {
			return lt
		} else if cmp == lang.GT {
			return gt
		}
	}
	return eq
}

func compareKey(a lang.Row, key lang.Key, orderBy order) cmpResult {
	for i, idx := range orderBy {
		cmp := lang.Compare(a[idx], key[i])
		if cmp == lang.LT {
			return lt
		} else if cmp == lang.GT {
			return gt
		}
	}
	return eq
}

//(index.header
func New(data []lang.Row, order []opt.ColOrdinal) *T { //)
	//(index.make-a-copy
	d := make([]lang.Row, len(data))
	copy(d, data)
	//)

	//(index.sort-it
	sort.Slice(d, func(i, j int) bool {
		return compare(d[i], d[j], order) == lt
	}) //)

	//(index.closer
	return &T{
		data:    d,
		orderBy: order,
	}
} //)

//(index.iter
func (idx *T) Iter() *Iterator {
	return &Iterator{
		index: idx,
		pos:   0,
	}
} //)

func (idx *T) SeekGE(key lang.Key) *Iterator {
	start := sort.Search(len(idx.data), func(i int) bool {
		return compareKey(idx.data[i], key, idx.orderBy) != lt
	})

	return &Iterator{
		index: idx,
		pos:   start,
	}
}

func (idx *T) SeekGT(key lang.Key) *Iterator {
	start := sort.Search(len(idx.data), func(i int) bool {
		return compareKey(idx.data[i], key, idx.orderBy) == gt
	})

	return &Iterator{
		index: idx,
		pos:   start,
	}
}

//(index.it.next
func (it *Iterator) Next() (lang.Row, bool) {
	if it.pos >= len(it.index.data) {
		return nil, false
	}
	it.pos++

	return it.index.data[it.pos-1], true
} //)

//(index.it.prev
func (it *Iterator) Prev() (lang.Row, bool) {
	if it.pos <= 1 {
		return nil, false
	}
	it.pos--

	return it.index.data[it.pos], true
} //)
