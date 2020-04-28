package exec

//(imports
import (
	"bytes"
	"sort"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/index"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
) //)

type ScalarExpr scalar.Group

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
	// TODO: need this to be a disjunction.
	constraint constraint.Constraint
	ordering   []opt.ColOrdinal
}

func (s *scan) Start() {}

func (s *scan) Next() (lang.Row, bool) {
	next, ok := s.iter.Next()
	if ok && s.constraint.End != nil {
		cmp := opt.KeyCompare(next, s.constraint.End, s.ordering)
		if cmp == lang.GT || cmp == lang.EQ && !s.constraint.InclusiveEnd {
			return nil, false
		}
	}
	return next, ok
}

func Scan(idx *cat.Index, constraint constraint.Constraint) Node {
	var iter *index.Iterator
	if constraint.Start != nil {
		if constraint.InclusiveStart {
			iter = idx.ScanGE(constraint.Start)
		} else {
			iter = idx.ScanGT(constraint.Start)
		}
	} else {
		iter = idx.Scan()
	}
	return &scan{
		iter:       iter,
		constraint: constraint,
		ordering:   idx.Ordering,
	}
} //)

//(select1
type select1 struct {
	input Node
	p     ScalarExpr
}

func (s *select1) Start() {
	s.input.Start()
}

func (s *select1) Next() (lang.Row, bool) {
	var evaled lang.Datum = lang.DBool(false)
	var next lang.Row
	for evaled != lang.DBool(true) {
		var ok bool
		next, ok = s.input.Next()
		if !ok {
			return nil, false
		}
		var err error
		evaled, err = scalar.Eval(s.p, next)
		if err != nil {
			// TODO: fixme
			panic(err)
		}
	}
	return next, true
}

func Select(in Node, pred ScalarExpr) Node {
	return &select1{
		input: in,
		p:     pred,
	}
} //)

//(project
type project struct {
	input Node
	exprs []ScalarExpr
}

func (p *project) Start() {
	p.input.Start()
}

func (p *project) Next() (lang.Row, bool) {
	next, ok := p.input.Next()
	if !ok {
		return nil, false
	}
	row := make(lang.Row, len(p.exprs))
	for i := range p.exprs {
		evaled, err := scalar.Eval(p.exprs[i], next)
		if err != nil {
			// TODO: fixme
			panic("no good chief")
		}
		row[i] = evaled
	}
	return row, true
}

func Project(in Node, exprs []ScalarExpr) Node {
	return &project{
		input: in,
		exprs: exprs,
	}
} //)

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

	result := appendRows(c.leftRows[c.leftIdx], c.rightRow)
	c.leftIdx++

	return result, true
}

func Cross(l, r Node) Node {
	return &cross{
		l: l,
		r: r,
	}
}

func appendRows(l, r lang.Row) lang.Row {
	result := make(lang.Row, 0, len(l)+len(r))
	result = append(result, l...)
	result = append(result, r...)
	return result
}

type merge struct {
	l         Node
	r         Node
	leftIdxs  []opt.ColOrdinal
	rightIdxs []opt.ColOrdinal

	lr    lang.Row
	rr    lang.Row
	queue []lang.Row
}

func Merge(l, r Node, leftIdxs, rightIdxs []opt.ColOrdinal) Node {
	return &merge{
		l:         l,
		r:         r,
		leftIdxs:  leftIdxs,
		rightIdxs: rightIdxs,
	}
}

func (m *merge) Start() {
	m.l.Start()
	m.r.Start()
}

func (m *merge) left() lang.Row {
	if m.lr == nil {
		m.lr, _ = m.l.Next()
	}
	return m.lr
}

func (m *merge) popLeft() lang.Row {
	r := m.left()
	m.lr = nil
	return r
}

func (m *merge) right() lang.Row {
	if m.rr == nil {
		m.rr, _ = m.r.Next()
	}
	return m.rr
}

func (m *merge) popRight() lang.Row {
	r := m.right()
	m.rr = nil
	return r
}

func (m *merge) Next() (lang.Row, bool) {
	for len(m.queue) == 0 {
		if m.left() == nil || m.right() == nil {
			return nil, false
		}

		cmp := opt.RowCompare2(m.left(), m.right(), m.leftIdxs, m.rightIdxs)
		switch cmp {
		case lang.LT:
			m.lr = nil
		case lang.GT:
			m.rr = nil
		case lang.EQ:
			leftRows := []lang.Row{m.popLeft()}
			for m.left() != nil && opt.RowCompare(leftRows[0], m.left(), m.leftIdxs) == lang.EQ {
				leftRows = append(leftRows, m.popLeft())
			}

			rightRows := []lang.Row{m.popRight()}
			for m.right() != nil && opt.RowCompare(rightRows[0], m.right(), m.rightIdxs) == lang.EQ {
				rightRows = append(rightRows, m.popRight())
			}

			for _, l := range leftRows {
				for _, r := range rightRows {
					m.queue = append(m.queue, appendRows(l, r))
				}
			}
		}
	}

	var next lang.Row
	next, m.queue = m.queue[0], m.queue[1:]
	return next, true
}

func hashRow(r lang.Row, key []int) string {
	var buf bytes.Buffer
	for i, idx := range key {
		if i > 0 {
			buf.WriteByte('/')
		}
		r[idx].Format(&buf)
	}
	return buf.String()
}

type hash struct {
	l Node
	r Node

	leftIdxs  []int
	rightIdxs []int
	queue     []lang.Row

	table map[string][]lang.Row
}

func Hash(l, r Node, leftIdxs, rightIdxs []int) Node {
	return &hash{
		l:         l,
		r:         r,
		leftIdxs:  leftIdxs,
		rightIdxs: rightIdxs,
	}
}

func (h *hash) Start() {
	h.l.Start()
	h.r.Start()

	h.table = make(map[string][]lang.Row)

	for next, ok := h.l.Next(); ok; next, ok = h.l.Next() {
		key := hashRow(next, h.leftIdxs)
		h.table[key] = append(h.table[key], next)
	}
}

func (h *hash) Next() (lang.Row, bool) {
	for len(h.queue) == 0 {
		next, ok := h.r.Next()
		if !ok {
			return nil, false
		}
		key := hashRow(next, h.rightIdxs)
		xs := h.table[key]
		for _, x := range xs {
			h.queue = append(h.queue, appendRows(x, next))
		}
	}

	var next lang.Row
	next, h.queue = h.queue[0], h.queue[1:]
	return next, true
}

//(sort
// TODO: naming this sort conflicts with the sortRows package. think of a better name (maybe suffix Node to all names)
type sortRows struct {
	input    Node
	ordering []opt.ColOrdinal

	rows []lang.Row
	idx  int
}

func (s *sortRows) Start() {
	s.rows = Spool(s.input)
	sort.Slice(s.rows, func(i, j int) bool {
		return opt.RowCompare(s.rows[i], s.rows[j], s.ordering) == lang.LT
	})
}

func (s *sortRows) Next() (lang.Row, bool) {
	if s.idx >= len(s.rows) {
		return nil, false
	}
	s.idx++
	return s.rows[s.idx-1], true
}

func Sort(in Node, ordering []opt.ColOrdinal) Node {
	return &sortRows{
		input:    in,
		ordering: ordering,
	}
} //)

type constant struct {
	rows []lang.Row
	idx  int
}

func (s *constant) Start() {}

func (s *constant) Next() (lang.Row, bool) {
	if s.idx >= len(s.rows) {
		return nil, false
	}

	s.idx += 1
	return s.rows[s.idx-1], true
}

func Constant(rows []lang.Row) Node {
	return &constant{
		rows: rows,
	}
}

// TODO: does this need error handling?
func Spool(n Node) []lang.Row {
	n.Start()
	result := make([]lang.Row, 0)
	for next, ok := n.Next(); ok; next, ok = n.Next() {
		result = append(result, next)
	}
	return result
}
