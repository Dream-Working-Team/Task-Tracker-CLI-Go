package storage

// model contiene las entidades del dominio utilizadas por el repositorio
import "task-cli/internal/model"

// Storage representa el repositorio de persistencia basado en archivo
type Storage struct {
	// Route es la ruta del archivo donde se almacenan los datos
	Route string
}

// ReadUsers obtiene la lista de usuarios desde el archivo configurado
func (s *Storage) ReadUsers() ([]model.User, error) {
	return ReadFile[model.User](s.Route)
}

// SaveUsers guarda la lista de usuarios en el archivo configurado
func (s *Storage) SaveUsers(users []model.User) error {
	return SaveFile(s.Route, users)
}

// ReadTask obtiene la lista de tareas desde el archivo configurado
func (s *Storage) ReadTask() ([]model.Task, error) {
	return ReadFile[model.Task](s.Route)
}

// SaveTask guarda la lista de tareas en el archivo configurado
func (s *Storage) SaveTask(task []model.Task) error {
	return SaveFile(s.Route, task)
}
