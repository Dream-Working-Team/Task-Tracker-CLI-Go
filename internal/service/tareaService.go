package service

import (
	"fmt"
	"time"

	"task-cli/internal/model"
	"task-cli/internal/storage"
)

type TaskService struct {
	capacity *storage.Storage
}

func NuewTaskService(s *storage.Storage) *TaskService {
	return &TaskService{capacity: s}
}

func (s *TaskService) Add(descripcion string) (model.Task, error) {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return model.Task{}, err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	newTask := model.Task{
		ID: id, Description: descripcion, Status: model.ToDo,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tasks = append(tasks, newTask)
	return newTask, s.capacity.SaveTask(tasks)
}

func (s *TaskService) ChangeStatus(id int, estado string) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = estado
			tasks[i].UpdatedAt = time.Now()
			return s.capacity.SaveTask(tasks)
		}
	}
	return fmt.Errorf("tarea %d no encontrada", id)
}

func (s *TaskService) List(filter string) ([]model.Task, error) {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return nil, err
	}

	if filter == "" {
		return tasks, nil
	}

	var filtered []model.Task
	for _, t := range tasks {
		if t.Status == filter {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}
