package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/sysgoblin/log4shell-honeypot/extractor"
	"github.com/sysgoblin/log4shell-honeypot/responses"
)

var (
	service        *string
	serviceOptions = []string{"elastic", "apache"}
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
		} else {
			for _, filename := range files {
				log.Printf("Saved payload from %s to file %s\n", request.RemoteAddr, filename)
			}
		}

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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

	// if request is get method
	if req.Method == "GET" {
		// check if get contains query parameters
		if len(req.URL.Query()) > 0 {
			// unescape query and send all parameters to Analyse
			for _, value := range req.URL.Query() {
				for _, v := range value {
					ve, _ := url.QueryUnescape(v)
					Analyse(ve, req)
				}
			}
		}
		// generate response
		switch *service {
		case "elastic":
			responses.CreateElasticResponse(w)
		case "apache":
			responses.CreateApacheResponse(w)
		}
	} else if req.Method == "POST" {
		switch *service {
		case "apache":
			b, _ := ioutil.ReadAll(req.Body)
			// remove url encoding
			body := string(b)
			bodyEscaped, _ := url.QueryUnescape(body)
			// add body back to request as we've consumed it
			req.Body = ioutil.NopCloser(strings.NewReader(bodyEscaped))
			Analyse(bodyEscaped, req)
		default:
			log.Printf("Service %s does not support POST method", *service)
		}
	}
}

func main() {
	address := flag.String("h", ":8080", "HTTP server binding address")
	service = flag.String("s", "apache", "Service to emulate response headers for")
	flag.Parse()

	// check if service is in serviceOptions
	if !contains(serviceOptions, *service) {
		log.Fatalf("Service %s is not supported", *service)
	}

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
