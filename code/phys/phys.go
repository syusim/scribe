package phys

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

// Physical Props are an object representing a presentational fact about a
// relational expression. They represent anything about a relational expression
// which is not explicitly part of the semantics of the language itself. For
// instance, order is not a part of relational algebra (however, many
// implementations of relational algebra allow you to express a specific
// presentation of the data to you at the end of your computation).
//
// Physical properties have one primary defining characteristic, which is that
// they can be _enforced_. That is, given an expression which does not satisfy
// a given set of physical properties, it is possible to transform its output
// into a form which _does_ satisfy that set of physical properties.
//
// The classical physical property is ordering. Ordering is not a concept in
// relational algebra, but it is useful to exploit it in computations, say, for
// efficiently distinctifying a result set, or to be able to use merge join.
// Ordering is enforceable via a Sort operator.
//
// A secondary characteristic of physical properties is that they form a
// partial order.  If A and B are physical properties, A <= B if any expression
// satisfying B also satisfies A.  In the case of ordering, A <= B if A is a
// prefix of B: a result set (lexicographically) ordered on [x y] is also
// ordered on [x].

// Min is the bottom element of the set of physical properties: it represents
// no demands on its input.
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

// Weaken gives a set of physical properties weaker than p, if possible. If
// that is not possible (meaning p is Min), it returns false.
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
