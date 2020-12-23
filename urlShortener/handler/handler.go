package handler

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// http.HandlerFunc = a function
// http.Handler has a method ServeHTTP(ResponseWriter, *Request)
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	// parse yaml file
	type pathUrl struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}

	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert to map[string]string
	pathUrlMap := map[string]string{}
	for _, v := range pathUrls {
		pathUrlMap[v.Path] = v.URL
	}

	return MapHandler(pathUrlMap, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse json file
	type pathUrl struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	}
	var pathUrls []pathUrl
	err := json.Unmarshal(jsonData, &pathUrls)
	if err != nil {
		return nil, err
	}
	// convert to map[string]string
	pathUrlMap := map[string]string{}
	for _, v := range pathUrls {
		pathUrlMap[v.Path] = v.URL
	}
	return MapHandler(pathUrlMap, fallback), nil
}
