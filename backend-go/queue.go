package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Task struct {
	ID        string
	Name      string
	Status    string
	CreatedAt time.Time
}

type TaskQueue struct {
	mu    sync.Mutex
	Tasks []Task
}

var queue = &TaskQueue{
	Tasks: make([]Task, 0),
}

func AddTask(name string) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	task := Task{
		ID:        fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Name:      name,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}
	queue.Tasks = append(queue.Tasks, task)
	log.Printf("[QUEUE] Task added: %s", task.Name)
}

func ProcessTasks() {
	for {
		queue.mu.Lock()
		if len(queue.Tasks) > 0 {
			for i, t := range queue.Tasks {
				if t.Status == "PENDING" {
					log.Printf("[QUEUE] Processing task: %s", t.Name)
					queue.Tasks[i].Status = "COMPLETED"
					break
				}
			}
		}
		queue.mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func StartQueueWorker() {
	go ProcessTasks()
	AddTask("Initial System Sync")
	AddTask("Vectorize Session Data")
}
