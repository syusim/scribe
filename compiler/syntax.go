package compiler

import (
	"bytes"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type highlighter struct {
	formatter *html.Formatter
	lexer     chroma.Lexer
	style     *chroma.Style
	cachedCss string
}

func newHighlighter() *highlighter {
	lexer := lexers.Get("go")
	style := styles.Get("friendly")
	formatter := html.New(html.WithClasses(true))

	return &highlighter{
		lexer:     lexer,
		style:     style,
		formatter: formatter,
	}
}

func (h *highlighter) highlight(out *bytes.Buffer, s string) error {
	it, err := h.lexer.Tokenise(nil, s)
	if err != nil {
		return err
	}

	h.formatter.Format(out, h.style, it)
	return nil
}

func (h *highlighter) css() string {
	if h.cachedCss == "" {
		var buf bytes.Buffer
		if err := h.formatter.WriteCSS(&buf, h.style); err != nil {
			return ""
		}
		h.cachedCss = buf.String()
	}
	return h.cachedCss
}
