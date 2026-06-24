package main

import (
	"fmt"
	"log"
	"time"
)

type Task struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Status    string
	CreatedAt time.Time
}

func AddTask(name string) {
	task := Task{
		ID:        fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Name:      name,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	result := db.Create(&task)
	if result.Error != nil {
		log.Printf("[QUEUE] Error adding task: %v", result.Error)
		return
	}
	log.Printf("[QUEUE] Task added to DB: %s", task.Name)
}

func ProcessTasks() {
	for {
		var task Task
		// Find first pending task
		result := db.Where("status = ?", "PENDING").First(&task)

		if result.Error == nil {
			log.Printf("[QUEUE] Processing task: %s", task.Name)
			db.Model(&task).Update("Status", "COMPLETED")
		}
		time.Sleep(5 * time.Second)
	}
}

func StartQueueWorker() {
	initDB() // Init DB before queue worker starts
	go ProcessTasks()

	var count int64
	db.Model(&Task{}).Count(&count)
	if count == 0 {
		AddTask("Initial System Sync")
		AddTask("Vectorize Session Data")
	}
}
