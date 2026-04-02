package storage

// model contains the domain entities used by the repository.
import "task-cli/internal/model"

// Storage represents a file-based persistence repository.
type Storage struct {
	// Route is the file path where data is stored.
	Route string
}

// ReadUsers retrieves the user list from the configured file.
func (s *Storage) ReadUsers() ([]model.User, error) {
	return ReadFile[model.User](s.Route)
}

// SaveUsers stores the user list in the configured file.
func (s *Storage) SaveUsers(users []model.User) error {
	return SaveFile(s.Route, users)
}

// ReadTask retrieves the task list from the configured file.
func (s *Storage) ReadTask() ([]model.Task, error) {
	return ReadFile[model.Task](s.Route)
}

// SaveTask stores the task list in the configured file.
func (s *Storage) SaveTask(task []model.Task) error {
	return SaveFile(s.Route, task)
}
