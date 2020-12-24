package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	htm, err := ioutil.ReadFile("ex5.html")
	check(err)

	r := strings.NewReader(string(htm))
	links, err := ParseLink(r)
	check(err)

	err = WriteFile("output", links)
	check(err)
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
					if string(a.Val[0]) != "#" {
						l.Href = a.Val
						valid = true
						break
					}
				}
			}
			c := n.FirstChild
			if c.Type == html.TextNode {
				l.Text = c.Data
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

func WriteFile(fileName string, links []Link) (err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, l := range links {
		f.WriteString(l.Text + "\n")
		f.WriteString(l.Href + "\n")
	}
	return nil
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
