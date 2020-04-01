package pipeline

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// the pipeline package provides a handful of interfaces and primitives for
// defining the process by which to build the finished product.

type Node interface {
	Start()
	Next() (row []string, ok bool)
}

type dirReader struct {
	dir    string
	fnames []string
}

func NewDirReader(dir string) Node {
	return &dirReader{dir: dir}
}

func (d *dirReader) Start() {
	filepath.Walk(d.dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			d.fnames = append(d.fnames, path)
		}
		return nil
	})
}

func (d *dirReader) Next() ([]string, bool) {
	if len(d.fnames) == 0 {
		return nil, false
	}
	var fname string
	fname, d.fnames = d.fnames[0], d.fnames[1:]

	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(fmt.Sprintf("error reading file: %q", fname))
	}

	return []string{fname, string(dat)}, true
}

func Spool(n Node) string {
	var buf bytes.Buffer
	n.Start()
	r, ok := n.Next()
	for ok {
		for i := range r {
			if i > 0 {
				buf.WriteByte('\t')
			}
			fmt.Fprintf(&buf, "%q", r[i])
		}
		buf.WriteByte('\n')
		r, ok = n.Next()
	}
	return buf.String()
}
