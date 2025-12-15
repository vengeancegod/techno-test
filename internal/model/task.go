package model

import (
	"time"
)

const (
	Open TaskStatus = iota
	Closed
)

type TaskStatus int

type Task struct {
	ID          int
	Title       string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
}

func (s TaskStatus) StringStatus() string {
	switch s {
	case Open:
		return "not done"
	case Closed:
		return "done"
	default:
		return "bug"
	}
}

func ParseTaskStatus(input string) TaskStatus {
	switch input {
	case "0", "not done":
		return Open
	case "1", "done":
		return Closed
	default:
		return 0
	}
}
