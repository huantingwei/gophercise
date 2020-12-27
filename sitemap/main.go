package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type urlset struct {
	Base string `xml:"xmlns,attr"`
	Urls []loc  `xml:"url"`
}

type loc struct {
	Value string `xml:"loc"`
}

func main() {

	urlFlag := flag.String("url", "https://www.calhoun.io/", "url of the target website (only https, ends with /)")
	maxDepthFlag := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	outputFlag := flag.String("out", "output.xml", "the name of the output xml file")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepthFlag)

	var smap urlset
	smap.Base = *urlFlag
	for _, p := range pages {
		smap.Urls = append(smap.Urls, loc{Value: p})
	}

	output, err := os.Create(*outputFlag)
	check(err, "Fail to create file")

	w := bufio.NewWriter(output)

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err := enc.Encode(smap); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}

func hrefs(body io.Reader, base string) []string {
	links, _ := ParseLink(body)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, base):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func get(urlStr string) []string {

	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}

	base := baseURL.String()
	return hrefs(resp.Body, base)
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	// use map as a set
	// use struct{} as value
	// because an empty struct does not allocate any memory
	seen := make(map[string]empty)

	// current queue
	var q map[string]empty
	// new queue
	nq := map[string]empty{
		urlStr: empty{},
	}
	// bfs
	for i := 0; i <= maxDepth; i++ {
		// assign nq to be q
		q, nq = nq, make(map[string]empty)

		for url, _ := range q {
			if _, ok := seen[url]; ok {
				// already seen
				continue
			}
			// if not already seen

			// mark as seen
			seen[url] = empty{}

			// (push) in new links to nq
			for _, link := range get(url) {
				nq[link] = empty{}
			}
		}
	}

	// return seen links

	// preallocate
	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}
