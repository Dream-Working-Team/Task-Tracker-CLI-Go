package service

import (
	"fmt"
	"time"

	"task-cli/internal/model"
	"task-cli/internal/storage"
)

// TaskService encapsula la lógica de negocio para gestionar tareas
type TaskService struct {
	// capacity representa la capa de almacenamiento usada por el servicio
	capacity *storage.Storage
}

// NewTaskService crea una instancia del servicio de tareas con su almacenamiento
func NewTaskService(s *storage.Storage) *TaskService {
	return &TaskService{capacity: s}
}

// Add agrega una nueva tarea con ID incremental y la guarda
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

// muteTask busca una tarea por ID, la modifica y guarda los cambios
func (s *TaskService) muteTask(id int, modificar func(*model.Task)) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			// Ejecutamos la modificación específica (estado o descripción)
			modificar(&tasks[i])
			tasks[i].UpdatedAt = time.Now()

			return s.capacity.SaveTask(tasks)
		}
	}
	return fmt.Errorf("tarea %d no encontrada", id)
}

// ChangeStatus actualiza el estado de una tarea por su ID
func (s *TaskService) ChangeStatus(id int, newStatus string) error {
	return s.muteTask(id, func(t *model.Task) {
		t.Status = newStatus
	})
}

// Update cambia la descripción de una tarea por su ID
func (s *TaskService) Update(id int, newDescription string) error {
	return s.muteTask(id, func(t *model.Task) {
		t.Description = newDescription
	})
}

// Delete elimina una tarea por su ID y guarda la lista actualizada
func (s *TaskService) Delete(id int) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {

			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.capacity.SaveTask(tasks)
		}
	}
	return fmt.Errorf("tarea %d no encontrada", id)
}

// List devuelve todas las tareas o solo las que coinciden con un estado
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
