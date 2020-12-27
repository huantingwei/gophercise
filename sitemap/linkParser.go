package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseLink(r io.Reader) (links []Link, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node) []Link
	f = func(n *html.Node) []Link {
		var links []Link
		var l Link
		if n.Type == html.ElementNode && n.Data == "a" {

			valid := false
			for _, a := range n.Attr {
				if a.Key == "href" {
					// check if it's a tag (ex: "#something")
					if strings.HasPrefix(a.Val, "#") {
						l.Href = strings.TrimSpace(a.Val)
						valid = true
						break
					}
				}
			}
			c := n.FirstChild
			if c.Type == html.TextNode {
				l.Text = strings.TrimSpace(c.Data)
			}

			if valid {
				links = append(links, l)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			links = append(links, f(c)...)
		}
		return links
	}
	links = f(doc)
	return
}
