package model

// Importamos el tiempo para saber la creacion y actualizacion de tareas

import "time"

// Se definen estas palabras como constantes para seguir el principio DRY
// Fomenta la seguridad para la escalabilidad

const (
	ToDo       = "To_do"
	InProgress = "In_progress"
	Complete   = "Done"
)

// Definimos el modelo de tarea como dato estructurado
// Los nombres de los campos deben iniciar en mayusculas para ser exportadoso publicos
// y que el paquete encoding/json pueda acceder a ellos
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"delete_at,omitempty"`
}

// Definimos el modelo de usuario como dato estructurado
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}
