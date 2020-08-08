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
	limit    int64
	csv      string
	csvFile  *os.File
	problems []problem
}

type problem struct {
	question string
	answer   string
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

	errCh := make(chan error)
	answerCh := make(chan string)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	problemNum, err := app.parseProblems()
	if err != nil {
		return err
	}

	timer := time.NewTimer(time.Second * time.Duration(app.limit))
	defer timer.Stop()

	for i, problem := range app.problems {
		fmt.Printf("Problem #%v: %v = ", i+1, problem.question)

		//TODO: use single goroutine and in/out channels instead creating goroutine for every problem
		go getAnswer(answerCh, errCh)

		select {
		case answer := <-answerCh:
			if answer == problem.answer {
				score++
			}
		case <-quit:
			fmt.Printf("\nProgram interrupted. You scored %v out of %v.\n", score, problemNum)
			fmt.Print("Press enter to quit...")
			<-errCh
			return nil
		case err := <-errCh:
			return err
		case <-timer.C:
			fmt.Printf("\nTime ran out (%vs). You scored %v out of %v.\n", app.limit, score, problemNum)
			fmt.Print("Press enter to quit...")
			<-errCh
			return nil
		}
	}

	fmt.Printf("You scored %v out of %v.\n", score, problemNum)
	fmt.Print("Press enter to quit...")
	fmt.Scanf("\n")

	return nil
}

func getAnswer(answerCh chan string, errCh chan error) {
	var answer string
	n, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		errCh <- err
		return
	}
	if n != 1 {
		errCh <- fmt.Errorf("wrong number of arguments: %v instead of 1", n)
		return
	}
	answerCh <- answer
}

func (app *appEnv) parseProblems() (int, error) {
	r := csv.NewReader(app.csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return 0, err
	}

	problemNum := len(records)
	app.problems = make([]problem, problemNum)

	for i, record := range records {
		app.problems[i] = problem{question: record[0], answer: record[1]}
	}

	return problemNum, nil
}
