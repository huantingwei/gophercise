package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	urlshort "./handler"
)

func main() {

	yamlFile := flag.String("yaml", "path.yaml", "path, url pair in yaml")
	jsonFile := flag.String("json", "path.json", "path, url pair in json")
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// read files
	yaml, err := ioutil.ReadFile(*yamlFile)
	check(err)

	json, err := ioutil.ReadFile(*jsonFile)
	check(err)

	// build handlers
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	check(err)
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	check(err)

	fmt.Printf("Starting the server on :%d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
