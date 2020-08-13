package main

import (
	"os"

	"github.com/semka95/gophercises/ex3/cyoa"
)

func main() {
	os.Exit(cyoa.CLI(os.Args[1:]))
}
