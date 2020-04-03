package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/justinj/scribe/compiler"
	"github.com/radovskyb/watcher"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var watch = flag.Bool("watch", false, "whether to watch for changes")

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

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func build() {
	in := "testbook/"
	out := "build/"

	c := compiler.New(out, in, "code/")

	pipe := c.Run()

	if err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		pipe <- compiler.FileAdded{Path: path}

		return nil
	}); err != nil {
		panic(err)
	}

	w := watcher.New()
	if err := w.AddRecursive(in); err != nil {
		panic(err)
	}

	if *watch {
		go func() {
			for {
				select {
				case event := <-w.Event:
					if event.IsDir() {
						continue
					}
					wd, err := os.Getwd()
					if err != nil {
						panic(err)
					}
					path, err := filepath.Rel(wd, event.Path)
					if err != nil {
						panic(err)
					}

					pipe <- compiler.FileAdded{Path: path}
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}()

		go func() {
			// Start the watching process - it'll check for changes every 100ms.
			if err := w.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()
	} else {
		pipe <- compiler.Stop
	}
}
