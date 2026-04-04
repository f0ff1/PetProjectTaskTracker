package model

import "time"

type Task struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	UserTaskID  int        `json:"user_task_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Tags        []string   `json:"tags"`
}
