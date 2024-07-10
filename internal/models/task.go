package models

import "time"

// Task представляет собой модель задачи
type Task struct {
	ID          string     `json:"id,omitempty"`          // Уникальный идентификатор задачи
	Title       string     `json:"title,omitempty"`       // Заголовок задачи
	Description string     `json:"description,omitempty"` // Описание задачи
	Done        bool       `json:"done,omitempty"`        // Признак завершённости задачи
	CreatedAt   time.Time  `json:"created_at,omitempty"`  // Время создания задачи
	DoneAt      *time.Time `json:"done_at,omitempty"`     // Время завершения задачи (если задача завершена)
	Duration    *float64   `json:"duration,omitempty"`    // Продолжительность выполнения задачи в часах (если указано)
}
