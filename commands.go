package main

import (
	"flag"
	"fmt"
	"strconv"
)

func listTasks(args []string) error {

	listFlags := flag.NewFlagSet("list", flag.ContinueOnError)

	showAll := listFlags.Bool("all", false, "show completed tasks")

	if err := listFlags.Parse(args); err != nil {
		return err
	}

	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		task := tasks[i]
		if task.Done && !*showAll {
			continue
		}
		task.prettyPrint()
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
			ID:   tasks.nextID(),
			Text: args[0],
			Done: false})
		return tasks.Save()
	}
}

func parseID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return -1, fmt.Errorf("invalid task ID: %q", arg)
	}
	return id, nil
}

func handleToggle(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo done [ID]")
	}

	id, err := parseID(args[0])
	if err != nil {
		return err
	}

	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	if err := tasks.toggleDone(id); err != nil {
		return err
	}

	for i := range tasks {
		task := tasks[i]
		if task.Done {
			continue
		}
		task.prettyPrint()
	}

	return tasks.Save()
}

func showTask(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo task [id]")
	}

	id, err := parseID(args[0])
	if err != nil {
		return err
	}

	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	index, ok := tasks.indexByID(id)
	if !ok {
		return fmt.Errorf("ID does not exist")
	}

	tasks[index].prettyPrint()
	return nil
}
