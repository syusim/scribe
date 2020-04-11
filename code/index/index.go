package index

import (
	"sort"

	"github.com/justinj/scribe/code/opt"
)

type order []opt.ColOrdinal

type T struct {
	data    []opt.Row
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

func compare(a, b opt.Row, orderBy order) cmpResult {
	for _, idx := range orderBy {
		if a[idx] < b[idx] {
			return lt
		} else if a[idx] > b[idx] {
			return gt
		}
	}
	return eq
}

func compareKey(a opt.Row, key opt.Key, orderBy order) cmpResult {
	for i, idx := range orderBy {
		if a[idx] < key[i] {
			return lt
		} else if a[idx] > key[i] {
			return gt
		}
	}
	return eq
}

//(index.header
func New(data []opt.Row, order []opt.ColOrdinal) *T { //)
	//(index.make-a-copy
	d := make([]opt.Row, len(data))
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

//(index.seekge
func (idx *T) SeekGE(key opt.Key) *Iterator {
	//[index.seekge-slow
	//start := 0
	//for start < len(idx.data) && compareKey(idx.data[start], key, idx.orderBy) == lt {
	//	start++
	//}
	//]
	//(index.seekge-binsearch
	start := sort.Search(len(idx.data), func(i int) bool {
		return compareKey(idx.data[i], key, idx.orderBy) != lt
	}) //)

	return &Iterator{
		index: idx,
		pos:   start,
	}
} //)

//(index.it.next
func (it *Iterator) Next() (opt.Row, bool) {
	if it.pos >= len(it.index.data) {
		return nil, false
	}
	it.pos++

	return it.index.data[it.pos-1], true
} //)

//(index.it.prev
func (it *Iterator) Prev() (opt.Row, bool) {
	if it.pos <= 1 {
		return nil, false
	}
	it.pos--

	return it.index.data[it.pos], true
} //)
