package snippets

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type FlagSet map[string]struct{}

type Block interface {
	Render(*bytes.Buffer, FlagSet)
}

type pair struct {
	l Block
	r Block
}

func (b *pair) Render(buf *bytes.Buffer, f FlagSet) {
	b.l.Render(buf, f)
	b.r.Render(buf, f)
}

type literal struct {
	contents string
}

func (b *literal) Render(buf *bytes.Buffer, _ FlagSet) {
	buf.WriteString(b.contents)
}

type fence struct {
	tag      string
	contents Block
}

func (b *fence) Render(buf *bytes.Buffer, f FlagSet) {
	if _, ok := f[b.tag]; ok {
		b.contents.Render(buf, f)
	}
}

type tokStream struct {
	buffered token
	r        *bufio.Reader
}

type kind int

const (
	noKind kind = iota
	lineKind
	startBlockKind
	endBlockKind
)

type token struct {
	k   kind
	lex string
}

var openRegexp = regexp.MustCompile(`^\s*//\((\w+)\s*$`)
var closeRegexp = regexp.MustCompile(`^\s*//\)\s*$`)

func (t *tokStream) populate() bool {
	if t.buffered.k == noKind {
		line, err := t.r.ReadString('\n')
		if err != nil {
			return false
		}

		matches := openRegexp.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.buffered = token{
				k:   startBlockKind,
				lex: matches[1],
			}
			return true
		}

		if closeRegexp.MatchString(line) {
			t.buffered = token{
				k:   endBlockKind,
				lex: "",
			}
			return true
		}

		t.buffered = token{
			k:   lineKind,
			lex: line,
		}
		return true
	}

	return true
}

func (t *tokStream) Next() (tok token, ok bool) {
	if !t.populate() {
		return token{}, false
	}
	n := t.buffered
	t.buffered = token{}
	return n, true
}

func (t *tokStream) Peek() (tok token, ok bool) {
	if !t.populate() {
		return token{}, false
	}
	return t.buffered, true
}

func build(in *tokStream) (Block, error) {
	tok, ok := in.Next()
	if !ok {
		return &literal{""}, io.EOF
	}
	switch tok.k {
	case lineKind:
		return &literal{tok.lex}, nil
	case startBlockKind:
		var result Block
		result = &literal{""}

		r, _ := in.Peek()
		if r.k == endBlockKind {
			_, _ = in.Next()
			return result, nil
		}

		for n, err := build(in); err == nil; n, err = build(in) {
			result = &pair{result, n}
			r, _ := in.Peek()
			if r.k == endBlockKind {
				// Skip over it.
				_, _ = in.Next()
				return &fence{
					contents: result,
					tag:      tok.lex,
				}, nil
			}
		}
	case endBlockKind:
		panic("hit end block")
	}

	panic(fmt.Sprintf("no good chief: %v", tok.k))
}

func New(in io.Reader) (Block, error) {
	r := &tokStream{r: bufio.NewReader(in)}
	var result Block
	result = &literal{""}

	for n, err := build(r); err == nil; n, err = build(r) {
		result = &pair{result, n}
	}
	return result, nil
}
