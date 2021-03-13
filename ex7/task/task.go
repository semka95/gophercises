package task

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// CLI runs the task command line app and returns its exit status.
func CLI() int {
	var app appEnv
	err := app.fromArgs()
	if err != nil {
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

// Command represents subcommand of task command
type Command int

const (
	_ Command = iota
	// Add task
	Add
	// Do task
	Do
	// List task
	List
)

// appEnv represents parsed command line arguments
type appEnv struct {
	command Command
	task    string
	doNums  []int
	store   Storage
}

// fromArgs parses command line arguments into appEnv struct
func (app *appEnv) fromArgs() error {
	var rootCmd = &cobra.Command{
		Use:   "task",
		Short: "task is a CLI for managing your TODOs.",
	}

	var cmdAdd = &cobra.Command{
		Use:   "add [task to add]",
		Short: "Add a new task to your TODO list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a task")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			app.command = Add
			app.task = strings.Join(args, " ")
		},
	}

	var cmdDo = &cobra.Command{
		Use:   "do [number of task to complete]",
		Short: "Mark a task on your TODO list as complete",
		Args: func(cmd *cobra.Command, args []string) error {
			var ids []int
			for _, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					fmt.Println("Failed to parse the argument:", arg)
					continue
				}
				ids = append(ids, id)
			}
			if len(ids) == 0 {
				return fmt.Errorf("nothing to delete")
			}

			app.doNums = ids

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			app.command = Do
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list [no options!]",
		Short: "List all of your incomplete tasks",
		Run: func(cmd *cobra.Command, args []string) {
			app.command = List
		},
	}

	rootCmd.AddCommand(cmdAdd)
	rootCmd.AddCommand(cmdDo)
	rootCmd.AddCommand(cmdList)
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func (app *appEnv) run() error {
	storage, err := NewBoltStore("tasks.db", "tasks")
	if err != nil {
		return err
	}
	app.store = storage

	switch app.command {
	case Add:
		err = app.Add()
	case List:
		err = app.List()
	case Do:
		err = app.Do()
	default:
		return nil
	}

	return err
}

// Task represents task
type Task struct {
	ID    int
	Value string
}

// Add adds new task to Storage
func (app appEnv) Add() error {
	_, err := app.store.Store(app.task)
	if err != nil {
		return err
	}
	fmt.Printf("Added \"%v\" to your task list.\n", app.task)

	return nil
}

// List displays all tasks
func (app appEnv) List() error {
	tasks, err := app.store.GetAll()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("You have no tasks to complete! Why not take a vacation? 🏖")
		return nil
	}

	fmt.Println("You have the following tasks:")
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task.Value)
	}

	return nil
}

// Do completes task by given number
func (app appEnv) Do() error {
	tasks, err := app.store.GetAll()
	if err != nil {
		return err
	}

	for _, id := range app.doNums {
		if id <= 0 || id > len(tasks) {
			fmt.Println("invalid task number: ", id)
			continue
		}
		task := tasks[id-1]
		err := app.store.Delete(task.ID)
		if err != nil {
			fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
			continue
		}
		fmt.Printf("Marked \"%d\" as completed.\n", id)
	}

	return nil
}
