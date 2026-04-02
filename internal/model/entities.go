package model

// Import time to track creation and update timestamps for tasks.

import "time"

// Define statuses as constants to follow DRY principles
// and keep status values consistent across the project.

const (
	ToDo       = "To_do"
	InProgress = "In_progress"
	Complete   = "Done"
)

// Task defines the structured model for task data.
// Field names must start with uppercase letters to be exported
// so encoding/json can read and write them.
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"delete_at,omitempty"`
}

// User defines the structured model for user data.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}
