// https://github.com/Adikso/minecraft-log4j-honeypot/tree/master/extractor
package extractor

import (
	"log"
	"net/url"
	"regexp"
)

type Finder struct {
	RegexPattern *regexp.Regexp
}

func NewFinder(pattern *regexp.Regexp) *Finder {
	return &Finder{RegexPattern: pattern}
}

func (f *Finder) FindInjections(text string) []*url.URL {
	var urls []*url.URL

	res := f.RegexPattern.FindAllStringSubmatch(text, -1)
	for i := range res {
		address, err := url.Parse(res[i][1])
		if err != nil {
			continue
		}

		urls = append(urls, address)
		log.Printf("Found injection: %s", address)
	}

	return urls
}
