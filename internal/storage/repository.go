package storage

import "task-cli/internal/model"

type Storage struct {
	Route string
}

func (s *Storage) ReadUsers() ([]model.User, error) {
	return ReadFile[model.User](s.Route)
}

func (s *Storage) SaveUsers(users []model.User) error {
	return SaveFile(s.Route, users)
}

func (s *Storage) ReadTask() ([]model.Task, error) {
	return ReadFile[model.Task](s.Route)
}

func (s *Storage) SaveTask(task []model.Task) error {
	return SaveFile(s.Route, task)
}
