package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func listTasks(args []string, tasks TaskList) error {

	listFlags := flag.NewFlagSet("list", flag.ContinueOnError)

	showAll := listFlags.Bool("all", false, "show completed tasks")

	if err := listFlags.Parse(args); err != nil {
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

func handleAdd(args []string, tasks TaskList) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo add \"task\"")
	}

	tasks = append(tasks, Task{
		ID:   tasks.nextID(),
		Text: args[0],
		Done: false})
	return tasks.Save()
}

func parseID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return -1, fmt.Errorf("invalid task ID: %q", arg)
	}
	return id, nil
}

func handleToggle(args []string, tasks TaskList) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo done [ID]")
	}

	id, err := parseID(args[0])
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

func CommandTask(args []string, tasks TaskList) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: todo task [id]")
	}

	id, err := parseID(args[0])
	if err != nil {
		return err
	}

	index, ok := tasks.indexByID(id)
	if !ok {
		return fmt.Errorf("ID: %v does not exist", id)
	}

	args = args[1:]
	if len(args) != 0 {
		switch args[0] {
		case "show":

			tasks[index].prettyPrint()

		case "edit":
			tasks[index].prettyPrint()

			new := readline("New Task: ")

			tasks[index].Text = new

			tasks.Save()

		case "delete":
			tasks = removeAt(tasks, index)
			tasks.Save()
		}
	} else {
		return fmt.Errorf("Usage: todo task [%v] [show|edit|delete]", id)
	}
	return nil
}

func readline(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)

	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

func removeAt[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}
