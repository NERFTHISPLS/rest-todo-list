package models

import "time"

const DefaultTaskStatus = "new"

// Task представляет задачу в системе
// swagger:model Task
type Task struct {
	// ID задачи (только в ответе)
	// example: 1
	ID int `json:"id"`

	// Заголовок задачи
	// required: true
	// example: Купить молоко
	Title string `json:"title"`

	// Описание задачи
	// required: false
	// example: Взять 2 литра и хлеб
	Description string `json:"description"`

	// Статус задачи
	// required: true
	// enum: new,in_progress,done
	// example: new
	Status string `json:"status"`

	// Дата создания (только в ответе)
	// example: 2025-08-13T14:52:00Z
	CreatedAt time.Time `json:"created_at"`

	// Дата последнего обновления (только в ответе)
	// example: 2025-08-13T15:12:00Z
	UpdatedAt time.Time `json:"updated_at"`
}
