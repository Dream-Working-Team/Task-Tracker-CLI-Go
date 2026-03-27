package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"task-cli/internal/auth"
	"task-cli/internal/config"
	"task-cli/internal/storage"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Gestión de autenticación (login, registro, logout)",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Inicia sesión en el sistema",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Usuario: ")
		var username string
		fmt.Scanln(&username)

		fmt.Print("Contraseña: ")
		password := auth.ReadPass()

		dirDatos, _ := config.GetDirectionData()
		routeUsers := filepath.Join(dirDatos, "usuarios.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Login(strings.ToLower(username), password); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("✅ Sesión iniciada como '%s'\n", username)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registra un nuevo usuario",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Nuevo Usuario: ")
		var username string
		fmt.Scanln(&username)

		fmt.Print("Nueva Contraseña: ")
		password := auth.ReadPass()

		dirData, _ := config.GetDirectionData()
		routeUsers := filepath.Join(dirData, "usuarios.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Register(strings.ToLower(username), password); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Registro exitoso. Usa 'task-cli auth login' para acceder.")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Cierra la sesión actual",
	Run: func(cmd *cobra.Command, args []string) {
		auth.CloseSesion()
		fmt.Println("👋 Sesión cerrada.")
	},
}

// La función init() en Go se ejecuta automáticamente al arrancar.
// Aquí conectamos los subcomandos a su padre.
func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd, registerCmd, logoutCmd)
}
