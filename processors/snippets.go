package processors

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/justinj/scribe/snippets"
)

var snippetRefRegexp = regexp.MustCompile(`^\s*%\s*(\S*)\s*$`)

func MakeSnippetProcessor(
	snippetDir string,
) (func(io.Reader) (io.Reader, error), error) {
	c, err := buildCorpus(snippetDir)
	if err != nil {
		return nil, err
	}

	return func(r io.Reader) (io.Reader, error) {
		var out bytes.Buffer
		rd := bufio.NewReader(r)
		// Probably a way to do this more efficiently, but yolo.
		for line, err := rd.ReadString('\n'); err == nil; line, err = rd.ReadString('\n') {
			matches := snippetRefRegexp.FindStringSubmatch(line)
			if len(matches) == 0 {
				out.WriteString(line)
			} else {
				// We are referencing a snippet!
				referenced := matches[1]
				referencedFile, ok := c.tags[referenced]
				if !ok {
					panic(fmt.Sprintf("bad snippet: %q", referenced))
				}
				file := c.files[referencedFile]
				out.WriteString("```go\n")
				snippets.Extract(file, &out, nil, referenced)
				out.WriteString("```\n")
			}
		}

		return &out, nil
	}, nil
}
