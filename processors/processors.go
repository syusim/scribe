package processors

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type lineKind int

const (
	_ lineKind = iota
	bareLine
	snippetRefLine
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
		if strings.HasPrefix(l, "% ") {
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

func unindent(s string) string {
	smallestIndent := len(s)
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		for i := 0; i < len(l); i++ {
			if !unicode.IsSpace(rune(l[i])) {
				if smallestIndent > i {
					smallestIndent = i
					break
				}
			}
		}
		if smallestIndent == 0 {
			return s
		}
	}

	var out bytes.Buffer
	for i, l := range lines {
		if len(l) <= smallestIndent {
			// l is probably just a blank line.
			out.WriteString(l)
		} else {
			out.WriteString(l[smallestIndent:])
		}
		if i < len(lines)-1 {
			out.WriteByte('\n')
		}
	}

	return out.String()
}
