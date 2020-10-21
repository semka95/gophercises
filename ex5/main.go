package main

import (
	"os"

	"github.com/semka95/gophercises/ex5/sitemap"
)

func main() {
	os.Exit(sitemap.CLI(os.Args[1:]))
}
