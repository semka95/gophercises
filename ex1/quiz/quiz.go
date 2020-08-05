package quiz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func CLI(args []string) int {
	var app appEnv

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

const (
	defaultTimeLimit = 30
	defaultProblems  = "./problems.csv"
)

type appEnv struct {
	limit   int
	csv     string
	csvFile *os.File
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("quiz", flag.ContinueOnError)
	fl.StringVar(
		&app.csv, "csv", defaultProblems, "a csv file in the format of \"question,answer\"",
	)
	fl.IntVar(
		&app.limit, "limit", defaultTimeLimit, "the time limit for the quiz in seconds",
	)

	if err := fl.Parse(args); err != nil {
		return err
	}

	file, err := os.Open(app.csv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "got bad output type: %v\n", app.csv)
		fl.Usage()
		return flag.ErrHelp
	}
	app.csvFile = file

	return nil
}

func (app *appEnv) run() error {
	// read csv and ask questions
	r := csv.NewReader(app.csvFile)
	problemNum := 0
	score := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		problemNum++
		fmt.Printf("Problem #%v: %v = ", problemNum, record[0])

		var answer string
		n, err := fmt.Scan(&answer)
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("wrong number of arguments: %v instead of 1", n)
		}
		if answer == record[1] {
			score++
		}
	}

	//output final result
	fmt.Printf("You scored %v out of %v.\n", score, problemNum)
	return nil
}
