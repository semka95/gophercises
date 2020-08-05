package quiz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	limit     int64
	csv       string
	csvFile   *os.File
	questions [][]string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("quiz", flag.ContinueOnError)
	fl.StringVar(
		&app.csv, "csv", defaultProblems, "a csv file in the format of \"question,answer\"",
	)
	fl.Int64Var(
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
	score := 0

	timer := time.NewTimer(time.Second * time.Duration(app.limit))
	defer timer.Stop()

	done := make(chan struct{}, 1)
	errCh := make(chan error)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	problemNum, err := app.readQuestions()
	if err != nil {
		return err
	}

	go func() {
		for i, question := range app.questions {
			select {
			case <-done:
				return
			default:
			}

			fmt.Printf("Problem #%v: %v = ", i, question[0])

			var answer string
			n, err := fmt.Scan(&answer)
			if err != nil {
				errCh <- err
				return
			}
			if n != 1 {
				errCh <- fmt.Errorf("wrong number of arguments: %v instead of 1", n)
				return
			}
			if answer == question[1] {
				score++
			}
		}
		done <- struct{}{}
	}()

	for {
		select {
		case <-quit:
			fmt.Printf("\nYou scored %v out of %v.\n", score, problemNum)
			return nil
		case err := <-errCh:
			return err
		case <-done:
			fmt.Printf("\nYou scored %v out of %v.\n", score, problemNum)
			return nil
		case <-timer.C:
			done <- struct{}{}
			fmt.Printf("\nYou scored %v out of %v.\n", score, problemNum)
			return nil
		}
	}
}

func (app *appEnv) readQuestions() (int, error) {
	r := csv.NewReader(app.csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return 0, err
	}

	app.questions = records
	problemNum := len(app.questions)
	return problemNum, nil
}
