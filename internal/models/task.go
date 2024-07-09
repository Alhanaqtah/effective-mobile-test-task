package models

import "time"

type Task struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Done        bool       `json:"done,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
	Duration    *float64   `json:"duration,omitempty"`
}
