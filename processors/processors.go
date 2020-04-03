package processors

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type lineKind int

const (
	_ lineKind = iota
	bareLine
	snippetRefLine
	metaLine
)

type line struct {
	kind lineKind
	lex  string
}

func tokenize(in io.Reader) []line {
	result := make([]line, 0)
	r := bufio.NewReader(in)

	var buf bytes.Buffer

	flush := func() {
		if buf.Len() > 0 {
			result = append(result, line{
				kind: bareLine,
				lex:  buf.String(),
			})
			buf.Reset()
		}
	}

	for l, err := r.ReadString('\n'); err == nil; l, err = r.ReadString('\n') {
		if l == "---\n" {
			flush()
			var meta bytes.Buffer
			for l, err := r.ReadString('\n'); err == nil; l, err = r.ReadString('\n') {
				if l == "---\n" {
					break
				}
				meta.WriteString(l)
			}
			result = append(result, line{
				kind: metaLine,
				lex:  meta.String(),
			})
		} else if strings.HasPrefix(l, "% ") {
			flush()
			result = append(result, line{
				kind: snippetRefLine,
				lex:  strings.Trim(l[2:], " \n"),
			})
		} else {
			buf.WriteString(l)
		}
	}

	flush()
	return result
}
