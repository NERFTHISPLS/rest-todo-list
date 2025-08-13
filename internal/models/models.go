package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
