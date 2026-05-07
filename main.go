package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func (t TaskList) nextId() int {
	maxId := 0

	for _, task := range t {
		if task.ID > maxId {
			maxId = task.ID
		}
	}
	return maxId + 1
}

func (t TaskList) indexById(id int) (int, bool) {
	for i, task := range t {
		if task.ID == id {
			return i, true
		}
	}
	return -1, false
}

func (t TaskList) toggleDone(id int) error {
	i, ok := t.indexById(id)
	if !ok {
		return fmt.Errorf("task not found by ID")
	}

	t[i].Done = true
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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: todo [add|list|done|task]")
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
