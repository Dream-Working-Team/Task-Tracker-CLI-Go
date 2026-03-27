package service

import (
	// fmt se utiliza para formatear mensajes de error
	"fmt"
	// time se usa para registrar fechas de creación y actualización
	"time"

	// model contiene las entidades del dominio, como Task y sus estados
	"task-cli/internal/model"
	// storage provee el acceso a la persistencia de tareas
	"task-cli/internal/storage"
)

// TaskService encapsula la lógica de negocio para gestionar tareas
type TaskService struct {
	// capacity representa la capa de almacenamiento usada por el servicio
	capacity *storage.Storage
}

// NuewTaskService crea el servicio de tareas con su dependencia de almacenamiento
func NewTaskService(s *storage.Storage) *TaskService {
	return &TaskService{capacity: s}
}

// Add crea una nueva tarea, asigna un ID incremental y la persiste
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

func (s *TaskService) ChangeStatus(id int, newStatus string) error {
	return s.muteTask(id, func(t *model.Task) {
		t.Status = newStatus
	})
}

// Actualizar modifica solo la descripción. (NUEVA)
func (s *TaskService) Update(id int, newDescription string) error {
	return s.muteTask(id, func(t *model.Task) {
		t.Description = newDescription
	})
}

// List devuelve todas las tareas o solo las que coinciden con un estado filtro
func (s *TaskService) Delete(id int) error {
	tasks, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			// Truco idiomático: unimos todo lo que está antes del índice 'i'
			// con todo lo que está después del índice 'i'.
			tasks = append(tasks[:i], tasks[i+1:]...)
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
