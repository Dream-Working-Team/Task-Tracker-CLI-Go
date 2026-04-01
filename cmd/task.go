package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"task-cli/internal/auth"
	"task-cli/internal/config"
	"task-cli/internal/model"
	"task-cli/internal/service"
	"task-cli/internal/storage"

	"github.com/spf13/cobra"
)

// GetServiceTask valida la sesión activa y crea el servicio de tareas del usuario
func GetServiceTask() (*service.TaskService, error) {
	userID, err := auth.GetActiveUser()
	if err != nil {
		return nil, err
	}

	dirDatos, err := config.GetDirectionData()
	if err != nil {
		return nil, err
	}

	archivo := filepath.Join(dirDatos, fmt.Sprintf("task_%s.json", userID))
	return service.NewTaskService(&storage.Storage{Route: archivo}), nil
}

// addCmd crea una nueva tarea con la descripción indicada
var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		tarea, err := svc.Add(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Salida exacta según requerimiento
		fmt.Printf("Task added successfully (ID: %d)\n", tarea.ID)
	},
}

// updateCmd actualiza la descripción de una tarea existente por ID
var updateCmd = &cobra.Command{
	Use:   "update [id] [new_description]",
	Short: "Update the description of an existing task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("Error: The ID must be a valid integer.")
			return
		}

		if err := svc.Update(id, args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d updated successfully.\n", id)
	},
}

// deleteCmd elimina una tarea existente de forma permanente por ID
var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Remove a task permanently",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("Error: The ID must be a valid integer.")
			return
		}

		if err := svc.Delete(id); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d deleted successfully.\n", id)
	},
}

var markTodoCmd = &cobra.Command{
	Use:   "mark-todo [id]",
	Short: "Mark a task as to-do",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("Error: The ID must be a valid integer.")
			return
		}

		if err := svc.ChangeStatus(id, model.ToDo); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Task %d marked as to-do.\n", id)

	},
}

// markInProgressCmd cambia el estado de una tarea a en progreso
var markInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress [id]",
	Short: "Mark a task as in-progress",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("Error: The ID must be a valid integer.")
			return
		}

		if err := svc.ChangeStatus(id, model.InProgress); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d marked as in-progress.\n", id)
	},
}

// markDoneCmd cambia el estado de una tarea a completada
var markDoneCmd = &cobra.Command{
	Use:   "mark-done [id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("Error: The ID must be a valid integer.")
			return
		}

		if err := svc.ChangeStatus(id, model.Complete); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d marked as done.\n", id)
	},
}

// listCmd muestra todas las tareas o las filtra por estado
var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "Task list",
	Args:  cobra.MaximumNArgs(1), // Acepta 0 o 1 argumento posicional
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		estado := ""
		if len(args) == 1 {
			switch args[0] {
			case "done":
				estado = model.Complete
			case "todo":
				estado = model.ToDo
			case "in-progress":
				estado = model.InProgress
			default:
				fmt.Println("Invalid filter. Use: done, todo, in-progress")
				return
			}
		}

		tasks, err := svc.List(estado)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("There are no tasks.")
			return
		}

		fmt.Printf("%-4s %-15s %s\n", "ID", "Status", "Description")
		for _, t := range tasks {
			fmt.Printf("%-4d %-15s %s\n", t.ID, "["+t.Status+"]", t.Description)
		}
	},
}

// init registra los comandos de tareas en el comando raíz.
func init() {
	rootCmd.AddCommand(addCmd, updateCmd, deleteCmd, markTodoCmd, markInProgressCmd, markDoneCmd, listCmd)
}
