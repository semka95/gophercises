package urlshort

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
	"net/http"
)

// Link represents path and corresponding link
type Link struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if link, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, link, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

// parseYAML parses YAML and returns array of Links
func parseYAML(yml []byte) ([]Link, error) {
	var link []Link

	err := yaml.Unmarshal(yml, &link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	[
// 		{
//			"path": "path",
//    		"url": "url"
//  	}
//	]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)

	return MapHandler(pathMap, fallback), nil
}

// parseJSON parses JSON and returns array of Links
func parseJSON(jsn []byte) ([]Link, error) {
	var link []Link

	err := json.Unmarshal(jsn, &link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// buildMap converts []Link to map[string]string
func buildMap(parsedYaml []Link) map[string]string {
	pathMap := make(map[string]string, len(parsedYaml))
	for _, v := range parsedYaml {
		pathMap[v.Path] = v.URL
	}

	return pathMap
}

func BoltHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	pathMap := make(map[string]string)

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("links"))

		if err := b.ForEach(func(k, v []byte) error {
			pathMap[string(k)] = string(v)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return MapHandler(pathMap, fallback), nil
}
