package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/numtide/serve-go/spa"
)

var (
	workDir   string
	port      int    = 3000
	oEmbedUrl string = os.Getenv("SERVEGO_OEMBED_URL")
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
	fmt.Fprintf(out, "  -oembed-url: Sets the oEmbed Link header if set (env: $SERVEGO_OEMBED_URL) (default %s)\n", oEmbedUrl)
	fmt.Fprintf(out, "  <work-dir>: Folder to serve (default to current directory)\n")
}

func run() error {
	flag.IntVar(&port, "listen", port, "Port to listen to")
	flag.StringVar(&oEmbedUrl, "oembed-url", oEmbedUrl, "Sets the oEmbed Link header if set")
	flag.Parse()

	if flag.NArg() > 0 {
		workDir = flag.Arg(0)
	} else {
		workDir = "."
	}

	fs := http.Dir(workDir)

	h := spa.FileServer(fs)

	if oEmbedUrl != "" {
		var err error
		h, err = spa.NewOembedMiddleware(h, oEmbedUrl)
		if err != nil {
			panic(err)
		}
	}

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
