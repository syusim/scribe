package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/justinj/scribe/processors"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	build()
}

func build() {
	in := "testbook/"
	out := "build/"

	os.RemoveAll(out)

	err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		res, err := processors.Build(f, "code/")
		if err != nil {
			return err
		}

		postPath := path[len(in):]
		newPath := filepath.Join(out, postPath) + ".html"

		dir, _ := filepath.Split(newPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		out, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer out.Close()
		io.Copy(out, res)

		return nil
	})

	if err != nil {
		panic(err)
	}
}
