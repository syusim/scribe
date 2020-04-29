package fd

import (
	"testing"

	"github.com/justinj/scribe/code/lang"
)

func TestFD(t *testing.T) {
	fd := None.Imply([]FD{{lang.SetFromCols(1), 2}})

	if !fd.Implies(lang.SetFromCols(1), lang.SetFromCols(2)) {
		t.Errorf("1->2 was not true in %s", fd)
	}

	fd2 := fd.Imply([]FD{{lang.SetFromCols(2), 3}})
	if !fd2.Implies(lang.SetFromCols(1), lang.SetFromCols(3)) {
		t.Errorf("1->3 was not true in %s", fd2)
	}
}
