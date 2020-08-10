package urlshort

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	yamlLinks = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	jsonLinks = `
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
)

type appEnv struct {
	yamlPath  string
	jsonPath  string
	yamlLinks []byte
	jsonLinks []byte
}

func CLI(args []string) int {
	app := appEnv{
		yamlLinks: []byte(yamlLinks),
		jsonLinks: []byte(jsonLinks),
	}

	err := app.fromArgs(args)
	if err != nil {
		return 2
	}

	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}

	return 0
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("urlshort", flag.ContinueOnError)
	fl.StringVar(&app.yamlPath, "yaml", "", "a yaml file in the format of \"- path: path  url: url\"")
	fl.StringVar(&app.jsonPath, "json", "", "a json file in the format of \"[{\"path\": \"path\", \"link\": \"link\"}]\"")

	if err := fl.Parse(args); err != nil {
		return err
	}

	// loading data from file if it was specified
	if app.yamlPath != "" {
		res, err := readFile(app.yamlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "got bad output type: %v\n", app.yamlPath)
			fl.Usage()
			return flag.ErrHelp
		}
		app.yamlLinks = res
	}

	if app.jsonPath != "" {
		res, err := readFile(app.jsonPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "got bad output type: %v\n", app.jsonPath)
			fl.Usage()
			return flag.ErrHelp
		}
		app.jsonLinks = res
	}

	return nil
}

func readFile(fileName string) ([]byte, error) {
	result, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (app *appEnv) run() error {
	mux := defaultMux()

	// Create BoltDB database
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Properly close database
	defer func() {
		dbPath := db.Path()
		if err := db.Close(); err != nil {
			log.Println(err)
		}

		if err := os.Remove(dbPath); err != nil {
			log.Println(err)
		}
	}()

	// Fill in data
	err = fillBoltDB(db)
	if err != nil {
		return err
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := YAMLHandler(app.yamlLinks, mapHandler)
	if err != nil {
		return err
	}
	// Build the JSONHandler using the YAMLHandler as the
	// fallback
	jsonHandler, err := JSONHandler(app.jsonLinks, yamlHandler)
	if err != nil {
		return err
	}
	// Build the BoltHandler using the JSONHandler as the
	// fallback
	boltHandler, err := BoltHandler(db, jsonHandler)
	if err != nil {
		return err
	}

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: boltHandler,
	}

	// Launch server
	go func() {
		log.Printf("Starting the server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	// Listen for interrupt signal to close and cleanup database, close http server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Program interrupted")

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

// fillBoltDB insert some data in BoltDB database
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
