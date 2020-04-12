package index

import (
	"reflect"
	"testing"

	"github.com/justinj/scribe/code/opt"
)

func spool(it *Iterator) []opt.Row {
	result := make([]opt.Row, 0)
	for r, ok := it.Next(); ok; r, ok = it.Next() {
		result = append(result, r)
	}
	return result
}

func TestScan(t *testing.T) {
	idx := New([]opt.Row{
		{"a", "x"},
		{"b", "z"},
		{"c", "y"},
	},
		[]opt.ColOrdinal{1},
	)

	it := idx.SeekGE(opt.Key{"y"})

	res := spool(it)
	if !reflect.DeepEqual(res, []opt.Row{{"c", "y"}, {"b", "z"}}) {
		t.Fatal("no!")
	}
}
