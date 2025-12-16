package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/numtide/serve-go/spa"
)

var (
	workDir               string
	port                  int    = 3000
	oEmbedUrl             string = os.Getenv("SERVEGO_OEMBED_URL")
	hstsSeconds           uint64 = 0
	hstsIncludeSubDomains bool   = false
	hstsPreload           bool   = false
)

func init() {
	flag.Usage = usage
}

func usage() {
	out := flag.CommandLine.Output()

	_, _ = fmt.Fprintf(out, "dead-simple application that serves static files from the current directory\n")
	_, _ = fmt.Fprintf(out, "Usage: serve-go [options] [<work-dir>]\n\n")
	_, _ = fmt.Fprintf(out, "Options:\n")
	_, _ = fmt.Fprintf(out, "  -listen: Port to listen to (default %d)\n", port)
	_, _ = fmt.Fprintf(out, "  -oembed-url: Sets the oEmbed Link header if set (env: $SERVEGO_OEMBED_URL) (default %s)\n", oEmbedUrl)
	_, _ = fmt.Fprintf(out, "  -hstsSeconds: Sets the HTTP Strict Transport Security max-age header if larger than 0 (env: $SERVEGO_HSTS_SECONDS) (default %d)\n", hstsSeconds)
	_, _ = fmt.Fprintf(out, "  -hstsIncludeSubDomains: Sets the HTTP Strict Transport Security includeSubDomains header if set (env: $SERVEGO_HSTS_INCLUDE_SUBDOMAINS) (default %v)\n", hstsIncludeSubDomains)
	_, _ = fmt.Fprintf(out, "  -hstsPreload: Sets the HTTP Strict Transport Security preload header if set (env: $SERVEGO_HSTS_PRELOAD) (default %v)\n", hstsPreload)
	_, _ = fmt.Fprintf(out, "  <work-dir>: Folder to serve (default to current directory)\n")
}

func run() error {
	flag.IntVar(&port, "listen", port, "Port to listen to")
	flag.StringVar(&oEmbedUrl, "oembed-url", oEmbedUrl, "Sets the oEmbed Link header if set")
	flag.Uint64Var(&hstsSeconds, "hstsSeconds", hstsSeconds, "Sets the HTTP Strict Transport Security max-age if set")
	flag.Bool("hstsIncludeSubDomains", hstsIncludeSubDomains, "Sets the HTTP Strict Transport Security includeSubDomains if set")
	flag.Bool("hstsPreload", hstsPreload, "Sets the HTTP Strict Transport Security includeSubDomains if set")
	flag.Parse()

	if flag.NArg() > 0 {
		workDir = flag.Arg(0)
	} else {
		workDir = "."
	}

	val, ok := os.LookupEnv("SERVEGO_HSTS_SECONDS")
	if ok {
		var err error
		hstsSeconds, err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	_, ok = os.LookupEnv("SERVEGO_HSTS_INCLUDE_SUBDOMAINS")
	if ok {
		hstsIncludeSubDomains = true
	}
	_, ok = os.LookupEnv("SERVEGO_HSTS_PRELOAD")
	if ok {
		hstsPreload = true
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

	if hstsSeconds > 0 {
		var err error
		h, err = spa.NewHSTSMiddleware(h, hstsSeconds, hstsIncludeSubDomains, hstsPreload)
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
