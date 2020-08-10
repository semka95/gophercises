package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/semka95/gophercises/ex2/urlshort"
	bolt "go.etcd.io/bbolt"
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

	// Create BoltDB database
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())
	defer db.Close()

	// Fill in data
	err = fillBoltDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// Listen for interrupt signal to close and cleanup database
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		dbPath := db.Path()
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(dbPath); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

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
	// Build the BoltHandler using the JSONHandler as the
	// fallback
	boltHandler, err := urlshort.BoltHandler(db, jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", boltHandler))
}

func fillBoltDB(db *bolt.DB) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("links"))
		if err != nil {
			return err
		}

		if err := b.Put([]byte("/bolt"), []byte("https://pkg.go.dev/go.etcd.io/bbolt")); err != nil {
			return err
		}
		if err := b.Put([]byte("/yandex"), []byte("https://yandex.ru")); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
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
