package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/sysgoblin/log4shell-honeypot/extractor"
)

func Analyse(text string, remote string) {
	log.Printf("Testing: %s\n", text)

	pattern := regexp.MustCompile(`\${jndi:(.*)}`)
	finder := extractor.NewFinder(pattern)

	injections := finder.FindInjections(text)
	for _, url := range injections {
		log.Printf("Fetching payload for: jndi:%s", url.String())

		files, err := extractor.FetchFromLdap(url)
		if err != nil {
			log.Printf("Failed to fetch class from %s", url)
			continue
		}

		for _, filename := range files {
			log.Printf("Saved payload from %s to file %s\n", remote, filename)
		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	// get the ua/ref strings from the request
	useragent := req.Header.Get("User-Agent")
	referer := req.Header.Get("Referer")
	url := req.URL.String()

	// log details of the request
	log.Printf("Request from %s: %s, %s, %s", req.RemoteAddr, url, useragent, referer)

	// send ua to Analyse
	Analyse(useragent, req.RemoteAddr)
	// send referer if it exists
	if referer != "" {
		Analyse(referer, req.RemoteAddr)
	}

	fmt.Fprintf(w, "thanks lol\n")
}

func main() {
	address := flag.String("h", ":8080", "HTTP server binding address")
	flag.Parse()

	f, err := os.OpenFile("http.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/", handler)
	http.ListenAndServe(*address, nil)
}
