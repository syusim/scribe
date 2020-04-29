package phys

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

var Min = &Props{}

type Props struct {
	Ordering lang.Ordering
}

func Hash(p Props) string {
	var buf bytes.Buffer
	for i := range p.Ordering {
		if i > 0 {
			buf.WriteString(",")
		}
		fmt.Fprintf(&buf, "%d", p.Ordering[i])
	}
	return buf.String()
}

// Weaken climbs down the partial order. Perhaps there's a
// smarter way we could do it than hopping immediately to the
// bottom once we have more execution things.
func (p Props) Weaken() (*Props, bool) {
	if len(p.Ordering) == 0 {
		return Min, false
	}

	return Min, true

	// When we support segmented sort, rather than returning min,
	// do this:
	// return Props{
	// 	Ordering: p.Ordering[:len(p.Ordering)-1],
	// }, true
}
