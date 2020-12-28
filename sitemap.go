package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/belkale/link"
)

type siteMap struct {
	site     *url.URL
	maxDepth int
	visited  map[string]bool
}

func fetchLinks(webpage *url.URL) ([]string, error) {
	log.Printf("Crawling %s", webpage)
	resp, err := http.Get(webpage.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	links, err := link.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, li := range links {
		result = append(result, li.HREF)
	}
	return result, nil
}

func getSameHostLinks(hostURL *url.URL, urls []string) []string {
	var shost []string

	set := make(map[string]bool)
	hostStr := hostURL.String()
	hostStrTrimmed := strings.TrimRight(hostStr, "/")
	for _, r := range urls {
		r = strings.TrimRight(r, "/")
		if strings.HasPrefix(r, "/") {
			set[hostStrTrimmed+r] = true
		} else if strings.HasPrefix(r, hostStr) {
			set[r] = true
		}
	}

	for k, _ := range set {
		shost = append(shost, k)
	}
	return shost
}

func BuildSiteMap(site string, maxDepth int) ([]string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	sm := siteMap{site: u, maxDepth: maxDepth, visited: visited}
	links := sm.crawl([]string{u.String()}, 0)

	var siteURLs []string
	for st, _ := range sm.visited {
		siteURLs = append(siteURLs, st)
	}
	siteURLs = append(siteURLs, links...)
	return siteURLs, nil
}

func (s siteMap) crawl(webpages []string, depth int) []string {
	if depth > s.maxDepth {
		return nil
	}

	var nextLinks []string
	for _, w := range webpages {
		if s.visited[w] {
			continue
		}

		s.visited[w] = true
		wurl, err := url.Parse(w)
		if err != nil {
			log.Print(err)
			continue
		}

		urls, err := fetchLinks(wurl)
		if err != nil {
			log.Print(err)
			continue
		}

		surls := getSameHostLinks(s.site, urls)
		nextLinks = append(nextLinks, surls...)
	}

	if depth == s.maxDepth {
		return nextLinks
	}

	return s.crawl(nextLinks, depth+1)
}
