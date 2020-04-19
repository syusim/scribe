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

func (k Key) Format(buf *bytes.Buffer) {
	for i, d := range k {
		if i > 0 {
			buf.WriteByte('/')
		}
		d.Format(buf)
	}
}
