package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/justinj/scribe/compiler"
)

func main() {
	build()

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// log.Println("Listening on :3000...")
	// err := http.ListenAndServe(":3000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func build() {
	in := "testbook/"
	out := "build/"

	c := compiler.New(out, in, "code/")

	if err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return c.AddFile(path)
	}); err != nil {
		panic(err)
	}

	if err := c.Build(); err != nil {
		panic(err)
	}
}
