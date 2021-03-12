package task

import (
	"fmt"
	"strconv"
)

// Task represents data needed for task CLI
type Task struct {
	args  string
	store Storage
}

// Add adds new task to Storage
func (t Task) Add() error {
	err := t.store.Store(t.args)
	if err != nil {
		return err
	}
	fmt.Printf("Added \"%v\" to your task list.\n", t.args)

	return nil
}

// List displays all tasks
func (t Task) List() error {
	tasks, err := t.store.GetAll()
	if err != nil {
		return err
	}
	fmt.Printf("You have the following tasks:\n%v", tasks)

	return nil
}

// Do completes task by given number
func (t Task) Do() error {
	id, err := strconv.Atoi(t.args)
	if err != nil {
		return err
	}

	deleted, err := t.store.Delete(id)
	if err != nil {
		return err
	}

	fmt.Printf("You have completed the \"%v\" task.\n", deleted)
	return nil
}
