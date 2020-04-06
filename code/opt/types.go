package opt

import "bytes"

//(relational-types
type Row []string
type Relation struct {
	ColNames []string
	Rows     []Row
} //)
//(key-type
type Key Row //)

//(col-ordinal-type
type ColOrdinal int //)

//(relation.string
func (t Relation) String() string {
	widest := make([]int, len(t.ColNames))

	for i, n := range t.ColNames {
		if widest[i] < len(n) {
			widest[i] = len(n)
		}
	}

	for i := range t.Rows {
		for j := range t.Rows[i] {
			l := len(t.Rows[i][j])
			if widest[j] < l {
				widest[j] = l
			}
		}
	}

	var buf bytes.Buffer
	for i, n := range t.ColNames {
		if i > 0 {
			buf.WriteString(" | ")
		}
		for k := 0; k < (widest[i]-len(n))/2; k++ {
			buf.WriteByte(' ')
		}
		buf.WriteString(n)
		for k := 0; k < (widest[i]-len(n))/2; k++ {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte('\n')
	for i := range widest {
		if i > 0 {
			buf.WriteString("-+-")
		}
		for j := 0; j < widest[i]; j++ {
			buf.WriteByte('-')
		}
	}
	buf.WriteByte('\n')
	for i := range t.Rows {
		for j := range t.Rows[i] {
			d := t.Rows[i][j]
			if j > 0 {
				buf.WriteString(" | ")
			}
			buf.WriteString(d)
			for k := 0; k < widest[j]-len(d); k++ {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('\n')
	}

	return buf.String()
} //)
