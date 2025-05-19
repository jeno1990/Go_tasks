package main

import "time"

type Task struct {
	ID          int
	Name        string
	Description string
	Status      bool
	CreatedAt   time.Time
}
