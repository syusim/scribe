package index

// func spool(it *Iterator) []lang.Row {
// 	result := make([]lang.Row, 0)
// 	for r, ok := it.Next(); ok; r, ok = it.Next() {
// 		result = append(result, r)
// 	}
// 	return result
// }

// func TestScan(t *testing.T) {
// 	idx := New([]lang.Row{
// 		{"a", "x"},
// 		{"b", "z"},
// 		{"c", "y"},
// 	},
// 		[]opt.ColOrdinal{1},
// 	)

// 	it := idx.SeekGE(lang.Key{"y"})

// 	res := spool(it)
// 	if !reflect.DeepEqual(res, []lang.Row{{"c", "y"}, {"b", "z"}}) {
// 		t.Fatal("no!")
// 	}
// }
