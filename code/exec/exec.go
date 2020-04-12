package exec

//(imports
import (
	"fmt"

	"github.com/justinj/scribe/code/index"
	"github.com/justinj/scribe/code/lang"
) //)

//(node-interface
type Node interface {
	// Start is called to initialize any state that this node needs to execute.
	Start()

	// Next returns the next row in the Node's result set. If there are no more
	// rows to return, the second return value will be false, otherwise, it will
	// be true.
	Next() (lang.Row, bool)
} //)

//(scan
type scan struct {
	iter *index.Iterator
}

func (s *scan) Start() {}

func (s *scan) Next() (lang.Row, bool) {
	return s.iter.Next()
}

func Scan(iter *index.Iterator) Node {
	return &scan{
		iter: iter,
	}
} //)

//(select1
type select1 struct {
	input Node
	i     int
	d     lang.Datum
}

func (s *select1) Start() {
	s.input.Start()
}

func (s *select1) Next() (lang.Row, bool) {
	next, ok := s.input.Next()
	for ok && lang.Compare(next[s.i], s.d) != lang.EQ {
		next, ok = s.input.Next()
	}
	return next, ok
}

func Select1(in Node, i int, d lang.Datum) Node {
	return &select1{
		input: in,
		i:     i,
		d:     d,
	}
} //)

//(the-rest
type select2 struct {
	input Node

	i int
	j int
}

func (s *select2) Start() {
	s.input.Start()
}

func (s *select2) Next() (lang.Row, bool) {
	next, ok := s.input.Next()
	for ok && next[s.i] != next[s.j] {
		next, ok = s.input.Next()
	}
	return next, ok
}

func Select2(in Node, i, j int) Node {
	return &select2{
		input: in,
		i:     i,
		j:     j,
	}
}

type project struct {
	input Node
	idxs  []int
}

func (p *project) Start() {
	p.input.Start()
}

func (p *project) Next() (lang.Row, bool) {
	next, ok := p.input.Next()
	if !ok {
		return nil, false
	}
	out := make(lang.Row, len(p.idxs))
	for i := range p.idxs {
		out[i] = next[p.idxs[i]]
	}
	return out, true
}

func Project(in Node, idxs []int) Node {
	return &project{
		input: in,
		idxs:  idxs,
	}
}

type cross struct {
	l Node
	r Node

	leftRows []lang.Row
	rightRow lang.Row
	leftIdx  int
	done     bool
}

func (c *cross) Start() {
	c.l.Start()
	c.r.Start()

	// Buffer up all the rows from the left side.
	for next, ok := c.l.Next(); ok; next, ok = c.l.Next() {
		c.leftRows = append(c.leftRows, next)
	}

	// Don't do anything if there are no rows in the left side.
	if len(c.leftRows) == 0 {
		c.done = true
	}

	// And buffer up a single row from the right side.
	var ok bool
	c.rightRow, ok = c.r.Next()
	if !ok {
		c.done = true
	}
}

func (c *cross) Next() (lang.Row, bool) {
	if c.done {
		return nil, false
	}

	// Check if we're done with this rightRow and need a fresh one.
	if c.leftIdx >= len(c.leftRows) {
		c.leftIdx = 0
		var ok bool
		c.rightRow, ok = c.r.Next()
		if !ok {
			c.done = true
			return nil, false
		}
	}

	result := make(lang.Row, 0, len(c.leftRows[c.leftIdx])+len(c.rightRow))
	result = append(result, c.leftRows[c.leftIdx]...)
	result = append(result, c.rightRow...)
	c.leftIdx++

	return result, true
}

func Cross(l, r Node) Node {
	return &cross{
		l: l,
		r: r,
	}
}

// func cols(n Node) []string {
// 	switch e := n.(type) {
// 	case *scan:
// 		// TODO: fix this!
// 		return e.rel.ColNames
// 	case *select1:
// 		return cols(e.input)
// 	case *select2:
// 		return cols(e.input)
// 	case *project:
// 		c := cols(e.input)
// 		out := make([]string, len(e.idxs))
// 		for i := range e.idxs {
// 			out[i] = c[e.idxs[i]]
// 		}
// 		return out
// 	case *cross:
// 		lc := cols(e.l)
// 		rc := cols(e.r)
// 		c := make([]string, len(lc)+len(rc))
// 		copy(c[:len(lc)], lc)
// 		copy(c[len(lc):], rc)
// 		return c
// 	}
// 	panic("unhandled")
// }

// TODO: does this need error handling?
func Spool(n Node) []lang.Row {
	n.Start()
	result := make([]lang.Row, 0)
	for next, ok := n.Next(); ok; next, ok = n.Next() {
		result = append(result, next)
	}
	return result
}

func ChildCount(n Node) int {
	switch n.(type) {
	case *scan:
		return 0
	case *select1, *select2, *project:
		return 1
	case *cross:
		return 2
	default:
		panic(fmt.Sprintf("unhandled node %T", n))
	}
}

func Child(n Node, i int) Node {
	switch e := n.(type) {
	case *select1:
		return e.input
	case *select2:
		return e.input
	case *project:
		return e.input
	case *cross:
		switch i {
		case 0:
			return e.l
		case 1:
			return e.r
		}
	}
	panic("unhandled")
}

// func Explain(n Node) string {
// 	var buf bytes.Buffer
// 	indent := func(depth int) {
// 		for i := 0; i < depth; i++ {
// 			buf.WriteString("  ")
// 		}
// 	}
// 	var p func(n Node, depth int)
// 	p = func(n Node, depth int) {
// 		switch e := n.(type) {
// 		case *scan:
// 			buf.WriteString("scan")
// 		case *select1:
// 			fmt.Fprintf(&buf, "select [%d] = %q", e.i, e.d)
// 		case *select2:
// 			fmt.Fprintf(&buf, "select [%d] = [%d]", e.i, e.j)
// 		case *project:
// 			fmt.Fprintf(&buf, "project %v", e.idxs)
// 		case *cross:
// 			fmt.Fprintf(&buf, "cross")
// 		}

// 		fmt.Fprintf(&buf, " (%s)\n", strings.Join(cols(n), ","))
// 		for i, m := 0, ChildCount(n); i < m; i++ {
// 			indent(depth + 1)
// 			buf.WriteString("-> ")
// 			p(Child(n, i), depth+1)
// 		}
// 	}

// 	p(n, 0)

// 	return buf.String()
// }

// func main() {
// 	r := opt.Relation{
// 		ColNames: []string{"name", "from", "resides"},
// 		Rows: []lang.Row{
// 			{"Jordan", "New York", "New York"},
// 			{"Lauren", "California", "New York"},
// 			{"Justin", "Ontario", "New York"},
// 			{"Devin", "California", "California"},
// 			{"Smudge", "Ontario", "Ontario"},
// 		},
// 	}

// 	c := opt.Relation{
// 		ColNames: []string{"location", "country"},
// 		Rows: []lang.Row{
// 			{"New York", "United States"},
// 			{"California", "United States"},
// 			{"Ontario", "Canada"},
// 		},
// 	}

// 	fmt.Println(
// 		Explain(
// 			Project(
// 				Select2(
// 					Select1(
// 						Cross(
// 							Scan(r),
// 							Scan(c),
// 						),
// 						0,
// 						"Justin",
// 					),
// 					2,
// 					3,
// 				),
// 				[]int{0, 4},
// 			),
// 		),
// 	)
// }

//)
