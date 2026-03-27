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
func obtenerServicioTareas() (*service.TaskService, error) {
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

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Gestión de tareas (requiere sesión activa)",
}

var updateCmd = &cobra.Command{
	Use:   "update [id] [nueva_descripcion]",
	Short: "Actualiza la descripción de una tarea existente",
	Args:  cobra.ExactArgs(2), // Validamos que el usuario pase exactamente 2 argumentos
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := obtenerServicioTareas()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Println("❌ Error: El ID debe ser un número entero válido.")
			return
		}

		// args[1] contiene la nueva descripción (Cobra ya maneja las comillas por nosotros)
		if err := svc.Update(id, args[1]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("✅ Tarea %d actualizada exitosamente.\n", id)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Elimina una tarea permanentemente",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := obtenerServicioTareas()
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
		fmt.Printf("✅ Tarea %d eliminada.\n", id)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [descripción]",
	Short: "Agrega una nueva tarea",
	Args:  cobra.MinimumNArgs(1), // Cobra valida automáticamente que haya al menos 1 argumento
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := obtenerServicioTareas()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		// Cobra nos entrega los argumentos limpios. Si hay espacios, el usuario debió usar comillas.
		tarea, err := svc.Add(args[0])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("✅ Tarea agregada (ID: %d)\n", tarea.ID)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista las tareas",
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := obtenerServicioTareas()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		// Obtenemos el flag (bandera) en lugar de argumento posicional
		estado, _ := cmd.Flags().GetString("estado")

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
			fmt.Printf("%d\t[%s]\t%s\n", t.ID, t.Status, t.Description)
		}
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Marca una tarea como completada",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := obtenerServicioTareas()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		id, _ := strconv.Atoi(args[0])
		if err := svc.CambiarEstado(id, model.Complete); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("✅ Tarea %d completada.\n", id)
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	// ¡Agregamos updateCmd y deleteCmd a la lista!
	taskCmd.AddCommand(addCmd, listCmd, completeCmd, updateCmd, deleteCmd)

	listCmd.Flags().StringP("estado", "e", "", "Filtrar por estado (por_hacer, completada)")
}
