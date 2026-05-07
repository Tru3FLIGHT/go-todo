package main

import (
	"fmt"
	"strconv"
)

func listTasks() error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		task := tasks[i]
		var comp string
		if task.Done {
			comp = "[X]"
		} else {
			comp = "[ ]"
		}

		fmt.Println(comp, task.Text)
	}
	return nil
}

func handleAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo add \"task\"")
	}

	if tasks, err := LoadTasks(); err != nil {
		return err
	} else {
		tasks = append(tasks, Task{
			ID:   tasks.nextId(),
			Text: args[0],
			Done: false})
		return tasks.Save()
	}
}

func handleToggle(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo done [ID]")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %q", args[0])
	}

	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	if err := tasks.toggleDone(id); err != nil {
		return err
	}

	return tasks.Save()
}
