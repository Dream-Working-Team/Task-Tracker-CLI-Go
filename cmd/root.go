package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "Un gestor de tareas multiusuario en la terminal",
	Long:  `Task Tracker CLI Secure. Permite gestionar tareas por usuario de forma aislada y persistente.`,
}

// Execute ejecuta el comando raíz de la aplicación
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
