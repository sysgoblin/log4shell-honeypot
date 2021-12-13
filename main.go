package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/sysgoblin/log4shell-honeypot/extractor"
)

func Analyse(text string, request *http.Request) {
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
			log.Printf("Saved payload from %s to file %s\n", request.RemoteAddr, filename)

			// log rull request
			reqDump, err := httputil.DumpRequest(request, true)
			if err != nil {
				log.Printf("Failed to dump request: %v", err)
				continue
			} else {
				log.Printf("Request: %s", reqDump)
			}
		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	// get all headers from the request
	headers := req.Header

	// send all headers to Analyse
	for _, value := range headers {
		for _, v := range value {
			Analyse(v, req)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	address := flag.String("h", ":8080", "HTTP server binding address")
	flag.Parse()

	http.HandleFunc("/", handler)
	http.ListenAndServe(*address, nil)
}
