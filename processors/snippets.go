package processors

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/justinj/scribe/snippets"
)

type corpus struct {
	files map[string]snippets.Block
	tags  map[string]string
}

func buildCorpus(dir string) (*corpus, error) {
	c := &corpus{
		files: make(map[string]snippets.Block),
		tags:  make(map[string]string),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		b, err := snippets.New(f)
		if err != nil {
			return err
		}
		tags, err := snippets.Tags(b)
		if err != nil {
			return err
		}

		c.files[path] = b
		for tag, _ := range tags {
			if _, ok := c.tags[tag]; ok {
				return fmt.Errorf("duplicate tag: %q", tag)
			}
			c.tags[tag] = path
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

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
				out.WriteString("```\n")
				snippets.Extract(file, &out, nil, referenced)
				out.WriteString("```\n")
			}
		}

		return &out, nil
	}, nil
}
