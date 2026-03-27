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

// Helper: Verifica la sesión e instancia el servicio dinámicamente
func GetServiceTask() (*service.TaskService, error) {
	usuarioID, err := auth.GetActiveUser()
	if err != nil {
		return nil, err
	}

	dirDatos, err := config.GetDirectionData()
	if err != nil {
		return nil, err
	}

	archivo := filepath.Join(dirDatos, fmt.Sprintf("tareas_%s.json", usuarioID))
	return service.NewTaskService(&storage.Storage{Route: archivo}), nil
}

// ---------------------------------------------------------
// COMANDOS PLANOS (Raíz del CLI)
// ---------------------------------------------------------

var addCmd = &cobra.Command{
	Use:   "add [descripción]",
	Short: "Agrega una nueva tarea",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		tarea, err := svc.Add(args[0])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		// Salida exacta según requerimiento
		fmt.Printf("Task added successfully (ID: %d)\n", tarea.ID)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [id] [nueva_descripcion]",
	Short: "Actualiza la descripción de una tarea existente",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("❌ Error: El ID debe ser un número entero válido.")
			return
		}

		if err := svc.Update(id, args[1]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d updated successfully.\n", id)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Elimina una tarea permanentemente",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("❌ Error: El ID debe ser un número entero válido.")
			return
		}

		if err := svc.Delete(id); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d deleted successfully.\n", id)
	},
}

var markInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress [id]",
	Short: "Marca una tarea como en progreso",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("❌ Error: El ID debe ser un número entero válido.")
			return
		}

		if err := svc.ChangeStatus(id, model.InProgress); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d marked as in-progress.\n", id)
	},
}

var markDoneCmd = &cobra.Command{
	Use:   "mark-done [id]",
	Short: "Marca una tarea como completada",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("❌ Error: El ID debe ser un número entero válido.")
			return
		}

		if err := svc.ChangeStatus(id, model.Complete); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("Task %d marked as done.\n", id)
	},
}

var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "Lista las tareas",
	Args:  cobra.MaximumNArgs(1), // Acepta 0 o 1 argumento posicional
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		estado := ""
		// Evaluamos si el usuario pasó un filtro (ej. list done)
		if len(args) == 1 {
			switch args[0] {
			case "done":
				estado = model.Complete
			case "todo":
				estado = model.ToDo
			case "in-progress":
				estado = model.InProgress
			default:
				fmt.Println("❌ Filtro inválido. Usa: done, todo, in-progress")
				return
			}
		}

		tareas, err := svc.List(estado)
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		if len(tareas) == 0 {
			fmt.Println("📂 No hay tareas.")
			return
		}

		fmt.Println("ID\tEstado\t\tDescripción")
		for _, t := range tareas {
			// Usando t.Status y t.Description según la convención de tu código
			fmt.Printf("%d\t[%s]\t%s\n", t.ID, t.Status, t.Description)
		}
	},
}

func init() {
	// APLANAMIENTO ARQUITECTÓNICO:
	// Ya no usamos 'taskCmd'. Conectamos todos los comandos directamente al comando raíz ('rootCmd').
	rootCmd.AddCommand(addCmd, updateCmd, deleteCmd, markInProgressCmd, markDoneCmd, listCmd)
}
