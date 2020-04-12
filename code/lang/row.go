package lang

import "bytes"

type Row []Datum
type Key Row

func (r Row) Format(buf *bytes.Buffer) {
	buf.WriteByte('[')
	for i, d := range r {
		if i > 0 {
			buf.WriteByte(' ')
		}
		d.Format(buf)
	}
	buf.WriteByte(']')
}
