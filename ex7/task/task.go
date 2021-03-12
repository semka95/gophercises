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
	args    string
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
			app.args = strings.Join(args, " ")
		},
	}

	var cmdDo = &cobra.Command{
		Use:   "do [number of task to complete]",
		Short: "Mark a task on your TODO list as complete",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a number of task")
			}

			if _, err := strconv.Atoi(args[0]); err != nil {
				return err
			}
			app.args = args[0]
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
			app.args = ""
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

	task := Task{
		app.args,
		storage,
	}

	switch app.command {
	case Add:
		err = task.Add()
	case List:
		err = task.List()
	case Do:
		err = task.Do()
	default:
		return nil
	}

	return err
}
