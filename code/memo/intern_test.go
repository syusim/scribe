package memo

import (
	"testing"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

func TestHash(t *testing.T) {
	tcs := []struct {
		val interface{}
		key string
	}{
		{
			Scan{
				TableName: "abc",
				Cols:      []opt.ColumnID{1, 2, 3},
			},
			"Scan/abc/1,2,3",
		}, {
			scalar.ColRef{
				Id:  1,
				Typ: lang.Int,
			},
			"ColRef/1/int",
		}, {
			scalar.Func{
				Op:   lang.Plus,
				Args: []scalar.Expr{},
			},
			"Func/+/",
		}, {
			lang.DInt(3),
			"3",
		}, {
			lang.DString("foo"),
			"\"foo\"",
		}, {
			lang.DBool(true),
			"true",
		},
	}
	for _, tc := range tcs {
		k := hash(tc.val)
		if k != tc.key {
			t.Errorf("expected %q, got %q", tc.key, k)
		}
	}
}

func TestIntern(t *testing.T) {
	m := New()
	j1 := m.internScan(Scan{
		TableName: "abc",
		Cols:      []opt.ColumnID{1, 2, 3},
	})

	j2 := m.internScan(Scan{
		TableName: "abc",
		Cols:      []opt.ColumnID{1, 2, 3},
	})

	if j1 != j2 {
		t.Fatalf("j1 (%p) != j2 (%p)", j1, j2)
	}
}
