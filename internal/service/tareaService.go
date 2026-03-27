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

func NuevoTareaService(s *storage.Storage) *TaskService {
	return &TaskService{capacity: s}
}

func (s *TaskService) Agregar(descripcion string) (model.Task, error) {
	tareas, err := s.capacity.ReadTask()
	if err != nil {
		return model.Task{}, err
	}

	id := 1
	if len(tareas) > 0 {
		id = tareas[len(tareas)-1].ID + 1
	}

	nueva := model.Task{
		ID: id, Description: descripcion, Status: model.ToDo,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tareas = append(tareas, nueva)
	return nueva, s.capacity.SaveTask(tareas)
}

func (s *TaskService) CambiarEstado(id int, estado string) error {
	tareas, err := s.capacity.ReadTask()
	if err != nil {
		return err
	}

	for i := range tareas {
		if tareas[i].ID == id {
			tareas[i].Status = estado
			tareas[i].UpdatedAt = time.Now()
			return s.capacity.SaveTask(tareas)
		}
	}
	return fmt.Errorf("tarea %d no encontrada", id)
}

func (s *TaskService) Listar(filtro string) ([]model.Task, error) {
	tareas, err := s.capacity.ReadTask()
	if err != nil {
		return nil, err
	}

	if filtro == "" {
		return tareas, nil
	}

	var filtradas []model.Task
	for _, t := range tareas {
		if t.Status == filtro {
			filtradas = append(filtradas, t)
		}
	}
	return filtradas, nil
}
