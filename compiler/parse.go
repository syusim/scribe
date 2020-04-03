package compiler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/justinj/scribe/snippets"
)

type Section interface {
}

type Prose struct {
	Contents string
}

type Code struct {
	Pre  string
	Code string
	Post string
}

type DocMeta struct {
	Title    string
	OrderBy  string
	Contents []Section
}

func ParseDocument(snippetDir, path string) (DocMeta, error) {
	// TODO: cache the corpus on the Compiler.
	c, err := buildCorpus(snippetDir)
	if err != nil {
		return DocMeta{}, err
	}

	in, err := os.Open(path)
	if err != nil {
		return DocMeta{}, err
	}

	var meta DocMeta

	// TODO: this should be global to all documents?
	seenFlags := make(snippets.FlagSet)
	for _, tok := range tokenize(in) {
		switch tok.kind {
		case metaLine:
			json.Unmarshal([]byte(tok.lex), &meta)
		case bareLine:
			meta.Contents = append(meta.Contents, Prose{tok.lex})
		case snippetRefLine:
			seenFlags[tok.lex] = struct{}{}
			pre, mid, post := c.getSnip(seenFlags, tok.lex)
			meta.Contents = append(meta.Contents, Code{pre, mid, post})
		}
	}
	return meta, nil
}

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
