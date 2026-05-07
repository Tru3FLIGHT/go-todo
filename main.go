package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

const path string = ".todo"

func SaveTasks(t []Task) error {
	data, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		return fmt.Errorf("Saving Error: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write: %w", err)
	}
	fmt.Println("Write success")
	return nil
}

func LoadTasks() ([]Task, error) {
	tasks := []Task{}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, fmt.Errorf("there is no .todo")
		}
		return tasks, fmt.Errorf("unable to read .todo: %w", err)
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func testWrite() error {
	return SaveTasks([]Task{{1, "task 1", false}, {2, "task 2", true}})
}

func listTasks() error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := 0; i < len(tasks); i++ {
		fmt.Println(tasks[i])
	}
	return nil
}

func updateId(tasks []Task) []Task {
	for i, task := range tasks {
		task.Id = i + 1
	}
	return tasks
}

func handleAdd(args []string) error {
	if tasks, err := LoadTasks(); err != nil {
		return err
	} else {
		tasks = append(tasks, Task{len(tasks) + 1, args[0], false})
		return SaveTasks(tasks)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: todo [add|list|id]")
		return
	}

	command := args[0]

	switch command {
	case "add":
		if len(args) == 1 {
			fmt.Println("todo add [task] [flags]")
			return
		}
		handleAdd(args[1:])
	case "list":
		if err := listTasks(); err != nil {
			fmt.Println(err)
		}
	case "id":
		if tasks, err := LoadTasks(); err != nil {
			fmt.Println(err)
		} else {
			tasks = updateId(tasks)
			if err = SaveTasks(tasks); err != nil {
				fmt.Println(err)
			}
		}
	}
}
