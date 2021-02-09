package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	workDir string
	port    int = 3000
)

func init() {
	flag.Usage = usage
}

func usage() {
	out := flag.CommandLine.Output()

	fmt.Fprintf(out, "dead-simple application that serves static files from the current directory\n")
	fmt.Fprintf(out, "Usage: serve-go [options] [<work-dir>]\n\n")
	fmt.Fprintf(out, "Options:\n")
	fmt.Fprintf(out, "  -listen: Port to listen to (default %d)\n", port)
	fmt.Fprintf(out, "  <work-dir>: Folder to serve (default to current directory)\n")
}

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
	flag.IntVar(&port, "listen", port, "Port to listen to")
	flag.Parse()

	if flag.NArg() > 0 {
		workDir = flag.Arg(0)
	} else {
		workDir = "."
	}

	fs := SPAFileSystem{http.Dir(workDir)}

	h := http.FileServer(fs)

	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Serving %s on %s\n", workDir, addr)

	return http.ListenAndServe(addr, h)
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
