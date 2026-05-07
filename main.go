package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TaskList []Task

const path string = ".todo"

func (t TaskList) Save() error {
	data, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		return fmt.Errorf("Saving Error: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write: %w", err)
	}
	return nil
}

func (tasks TaskList) nextId() int {
	maxId := 0

	for _, task := range tasks {
		if task.ID > maxId {
			maxId = task.ID
		}
	}
	return maxId + 1
}

func (tasks TaskList) indexById(id int) (int, bool) {
	for i, task := range tasks {
		if task.ID == id {
			return i, true
		}
	}
	return -1, false
}

func (tasks TaskList) toggleDone(id int) error {
	i, ok := tasks.indexById(id)
	if !ok {
		return fmt.Errorf("task not found by ID")
	}

	tasks[i].Done = true
	return nil
}

func LoadTasks() (TaskList, error) {
	tasks := make(TaskList, 0)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return tasks, fmt.Errorf("unable to read .todo: %w", err)
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return TaskList{}, fmt.Errorf("Unable to Parse .todo: %w", err)
	}
	return tasks, nil
}

func listTasks() error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		fmt.Println(tasks[i])
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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: todo [add|list|done]")
		return
	}

	command := args[0]

	switch command {
	case "add":
		if err := handleAdd(args[1:]); err != nil {
			fmt.Println(err)
		}

	case "list":
		if err := listTasks(); err != nil {
			fmt.Println(err)
		}

	case "done":
		if err := handleToggle(args[1:]); err != nil {
			fmt.Println(err)
		}
	}
}
