package main

import (
	"encoding/json"
	"os"
)

type LogEntry struct {
	Time       string   `json:"time"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Success    bool     `json:"success"`
	Error      string   `json:"error,omitempty"`
	DurationMS int64    `json:"duration_ms"`
}

func startLogger(path string) (chan<- LogEntry, <-chan error) {
	logs := make(chan LogEntry, 1)
	done := make(chan error, 1)

	go func() {
		entry := <-logs
		data, err := json.MarshalIndent(entry, "", "   ")
		if err == nil {
			err = os.WriteFile(path, data, 0644)
		}

		done <- err
	}()

	return logs, done
}
