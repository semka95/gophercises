package main

import (
	"github.com/semka95/gophercises/ex1/quiz"
	"os"
)

func main() {
	os.Exit(quiz.CLI(os.Args[1:]))
}
