package snippets

import (
	"bytes"
	"strings"
	"unicode"
)

type extractionType int

const (
	_ extractionType = iota
	textExtraction
	straddlerExtraction
)

type subextraction struct {
	k extractionType
	// text
	t string

	// straddler
	l string
	e string
	r string
}

// TODO: make this not alloc?
func trimRight(s string) string {
	lines := strings.Split(s, "\n")
	var out bytes.Buffer

	// Go until we find a non-blank line.
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		out.WriteString(lines[start])
		out.WriteByte('\n')
		start++
	}

	if start >= len(lines) {
		return ""
	}

	out.WriteString(lines[start])
	level := indentationLen(lines[start])
	for i := start + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" || indentationLen(lines[i]) != level {
			break
		}
		out.WriteByte('\n')
		out.WriteString(lines[i])
	}
	return out.String()
}

func trimLeft(s string) string {
	// TODO: make this not suck
	lines := strings.Split(strings.TrimRight(s, " "), "\n")
	revd := make([]string, len(lines))
	for i := range lines {
		revd[len(lines)-1-i] = lines[i]
	}
	right := trimRight(strings.Join(revd, "\n"))
	revd = strings.Split(right, "\n")
	lines = make([]string, len(revd))
	for i := range revd {
		lines[len(lines)-1-i] = revd[i]
	}

	return strings.Join(lines, "\n")
}

func indentationLen(s string) int {
	for i := 0; i < len(s); i++ {
		if !unicode.IsSpace(rune(s[i])) {
			return i
		}
	}
	return len(s)
}

func appendExtractions(a, b subextraction) subextraction {
	switch a.k {
	case textExtraction:
		switch b.k {
		case textExtraction:
			return subextraction{
				k: textExtraction,
				t: a.t + b.t,
			}
		case straddlerExtraction:
			return subextraction{
				k: straddlerExtraction,
				l: a.t + b.l,
				e: b.e,
				r: b.r,
			}
		}
	case straddlerExtraction:
		switch b.k {
		case textExtraction:
			return subextraction{
				k: straddlerExtraction,
				l: a.l,
				e: a.e,
				r: a.r + b.t,
			}
		}
	}
	panic("invalid")
}

type Extraction struct {
	Pre      string
	Contents string
	Post     string
}

func ExtractCtx(b Block, f FlagSet, tag string) Extraction {
	e := extractCtx(b, f, tag)
	if e.k != straddlerExtraction {
		panic("I think this means we didn't find it")
	}
	return Extraction{
		Pre:      trimLeft(e.l),
		Contents: e.e,
		Post:     trimRight(e.r) + "\n",
	}
}

func extractCtx(b Block, f FlagSet, tag string) subextraction {
	switch e := b.(type) {
	case *literal:
		return subextraction{
			k: textExtraction,
			t: e.contents,
		}
	case *pair:
		left := extractCtx(e.l, f, tag)
		right := extractCtx(e.r, f, tag)
		return appendExtractions(left, right)
	case *fence:
		if e.tag == tag {
			// We found it! Render according to the flagSet.
			var buf bytes.Buffer
			e.contents.Render(&buf, f)
			return subextraction{
				k: straddlerExtraction,
				e: buf.String(),
			}
		} else if _, ok := f[e.tag]; ok {
			// This is not the block we were looking for. Keep digging.
			return extractCtx(e.contents, f, tag)
		} else {
			return subextraction{k: textExtraction}
		}
	}
	panic("no good bub")
}
