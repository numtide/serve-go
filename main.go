package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type SPAFileSystem struct {
	http.FileSystem
}

func (fs SPAFileSystem) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)
	if err != nil && os.IsNotExist(err) {
		// Fallback to index.html if the file doesn't exist
		return fs.FileSystem.Open("index.html")
	}
	return file, err
}

func run() error {
	var (
		workDir string
		err     error
		port    int = 3000
	)

	flag.IntVar(&port, "listen", port, "Port to listen to")
	flag.Parse()

	if flag.NArg() > 0 {
		workDir = flag.Arg(0)
	} else {
		workDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fs := SPAFileSystem{http.Dir(workDir)}

	h := http.FileServer(fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), h)
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
