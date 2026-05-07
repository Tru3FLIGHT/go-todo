package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID   int
	Text string
	Done bool
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
			return tasks, nil
		}
		return tasks, fmt.Errorf("unable to read .todo: %w", err)
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func main() {
	tasks, err := LoadTasks()
	if err != nil {
		fmt.Println(err)
	}
	if len(tasks) == 0 {
		fmt.Println("There is no .todo")
	}
	for i := 0; i < len(tasks); i++ {
		fmt.Println(tasks[i])
	}
}
