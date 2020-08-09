package main

import (
	"flag"
	"fmt"
	"github.com/semka95/gophercises/ex2/urlshort"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	yamlFile string
	jsonFile string
)

func init() {
	flag.StringVar(&yamlFile, "yaml", "", "a yaml file in the format of \"- path: path  url: url\"")
	flag.StringVar(&jsonFile, "json", "", "a json file in the format of \"[{\"path\": \"path\", \"link\": \"link\"}]\"")
}

func main() {
	flag.Parse()
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	jsonURL := `
[
  {
    "path": "/telegram",
    "url": "https://telegram.org"
  },
  {
    "path": "/matrix",
    "url": "https://matrix.org"
  }
]
`
	// loading data from file if it was specified
	if yamlFile != "" {
		res, err := readFile(yamlFile)
		if err != nil {
			log.Fatal(err)
		}
		yaml = res
	}

	if jsonFile != "" {
		res, err := readFile(jsonFile)
		if err != nil {
			log.Fatal(err)
		}
		jsonURL = res
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	// Build the JSONHandler using the YAMLHandler as the
	// fallback
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonURL), yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", jsonHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFile(fileName string) (string, error) {
	result, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
