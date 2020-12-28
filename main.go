package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
)

var site = flag.String("site", "", "Website for sitemap creation")

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNs   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}
type URL struct {
	XMLName xml.Name `xml:"url"`
	LOC     string   `xml:"loc"`
}

const (
	XMLPrefix = `<?xml version="1.0" encoding="UTF-8"?>`
	XMLNs     = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

func main() {
	flag.Parse()
	if *site == "" {
		log.Fatal("site cannot be null")
	}

	links, err := BuildSiteMap(*site, 0)
	if err != nil {
		log.Fatal(err)
	}

	var lu []URL
	for _, li := range links {
		item := URL{LOC: li}
		lu = append(lu, item)
	}

	xmlStr, err := xml.MarshalIndent(URLSet{XMLNs: XMLNs, URLs: lu}, " ", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(XMLPrefix + string(xmlStr))
}
