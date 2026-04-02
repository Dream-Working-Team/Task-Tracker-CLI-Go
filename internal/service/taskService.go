package service

import (
	"fmt"
	"strings"
	"task-cli/internal/model"
	"task-cli/internal/storage"
	"time"
)

// TaskService encapsulates business logic for task management.
type TaskService struct {
	// capacity references the storage layer used by the service.
	capacity *storage.Storage
}

// NewTaskService creates a task service instance with its storage dependency.
func NewTaskService(s *storage.Storage) *TaskService {
	return &TaskService{capacity: s}
}

// Add creates a new task with an incremental ID and stores it.
func (s *TaskService) Add(description string) (model.Task, error) {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return model.Task{}, err
	}

	normalizedNew := strings.ToLower(strings.TrimSpace(description))
	for _, t := range tasks {
		if t.DeletedAt != nil {
			continue
		}

		normalizedExisting := strings.ToLower(strings.TrimSpace(t.Description))
		if normalizedExisting == normalizedNew {
			return model.Task{}, fmt.Errorf("a task with this description already exists")
		}
	}

	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	newTask := model.Task{
		ID:          nextID,
		Description: description,
		Status:      model.ToDo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)
	return newTask, s.capacity.SaveTask(tasks)
}

// muteTask finds a task by ID, applies changes, and saves them.
func (s *TaskService) muteTask(id int, modify func(*model.Task) error) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID != id {
			continue
		}
		if tasks[i].DeletedAt != nil {
			return fmt.Errorf("task %d already deleted", id)
		}

		err := modify(&tasks[i])
		if err != nil {
			return err
		}
		tasks[i].UpdatedAt = time.Now()
		return s.capacity.SaveTask(tasks)
	}
	return fmt.Errorf("task %d not found", id)
}

// ChangeStatus updates the status of a task by ID.
func (s *TaskService) ChangeStatus(id int, newStatus string) error {
	return s.muteTask(id, func(t *model.Task) error {
		if t.Status == newStatus {
			return fmt.Errorf("task %d is already in status %q", id, newStatus)
		}
		t.Status = newStatus
		return nil
	})
}

// Update changes a task description by ID.
func (s *TaskService) Update(id int, newDescription string) error {
	return s.muteTask(id, func(t *model.Task) error {
		t.Description = newDescription
		return nil
	})
}

// Delete removes a task by ID and stores the updated list.
func (s *TaskService) Delete(id int) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID != id {
			continue
		}
		if tasks[i].DeletedAt != nil {
			return fmt.Errorf("task %d already deleted", id)
		}

		now := time.Now()
		tasks[i].DeletedAt = &now
		tasks[i].UpdatedAt = now
		return s.capacity.SaveTask(tasks)
	}
	return fmt.Errorf("task %d not found", id)
}

// List returns all tasks or only those that match a status filter.
func (s *TaskService) List(filter string) ([]model.Task, error) {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return nil, err
	}

	filtered := make([]model.Task, 0, len(tasks))
	for _, t := range tasks {
		if t.DeletedAt != nil {
			continue
		}
		if t.Status == filter || filter == "" {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}
