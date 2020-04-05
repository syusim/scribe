package main

import (
	"sort"

	"github.com/justinj/scribe/code/opt"
)

type order []opt.ColOrdinal

type index struct {
	data    []opt.Row
	orderBy order
}

//(index.iterator-def
type iterator struct {
	index *index
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

func compareKey(a, key opt.Row, orderBy order) cmpResult {
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
func New(data opt.Relation, order []opt.ColOrdinal) *index { //)
	//(index.make-a-copy
	d := make([]opt.Row, len(data.Rows))
	copy(d, data.Rows)
	//)

	//(index.sort-it
	sort.Slice(d, func(i, j int) bool {
		return compare(d[i], d[j], order) == lt
	}) //)

	//(index.closer
	return &index{
		data:    d,
		orderBy: order,
	}
} //)

//(index.seekge
func (idx *index) SeekGE(key opt.Row) *iterator {
	//[index.seekge-slow
	//start := 0
	//for compareKey(idx.data[start], key, idx.orderBy) == lt {
	//	start++
	//}
	//]
	//(index.seekge-binsearch
	start := sort.Search(len(idx.data), func(i int) bool {
		return compareKey(idx.data[i], key, idx.orderBy) != lt
	})

	//)

	return &iterator{
		index: idx,
		pos:   start,
	}
} //)

//(index.it.next
func (it *iterator) Next() (opt.Row, bool) {
	s := it.index.data
	if it.pos >= len(s) {
		return nil, false
	}
	it.pos++

	return it.index.data[it.pos-1], true
} //)
