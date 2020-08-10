package main

import (
	"os"

	"github.com/semka95/gophercises/ex2/urlshort"
)

func main() {
	os.Exit(urlshort.CLI(os.Args[1:]))
}
