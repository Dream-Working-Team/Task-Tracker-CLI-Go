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

// authCmd agrupa los subcomandos de autenticación del sistema
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Gestión de autenticación (login, registro, logout)",
}

// loginCmd solicita credenciales e inicia sesión si son válidas
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the system",
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := auth.GetActiveUser(); err == nil {
			fmt.Println("Error: You are already logged in. Use 'task-cli auth logout' first.")
			return
		}

		fmt.Print("Username: ")
		var username string
		fmt.Scanln(&username)

		fmt.Print("Password: ")
		password := auth.ReadPass()

		dirDatos, _ := config.GetDirectionData()
		routeUsers := filepath.Join(dirDatos, "users.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Login(strings.ToLower(username), password); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Session started as '%s'\n", username)
	},
}

// registerCmd pide datos y registra un nuevo usuario en el sistema
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := auth.GetActiveUser(); err == nil {
			fmt.Println("Error: Log out of your current session with 'task-cli auth logout' to register a new user.")
			return
		}

		fmt.Print("New Username: ")
		var username string
		fmt.Scanln(&username)

		fmt.Print("New Password: ")
		password := auth.ReadPass()

		dirData, _ := config.GetDirectionData()
		routeUsers := filepath.Join(dirData, "users.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Register(strings.ToLower(username), password); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Registration successful. Use 'task-cli auth login' to log in.")
	},
}

// logoutCmd cierra la sesión actual del usuario activo
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close the current session",
	Run: func(cmd *cobra.Command, args []string) {
		auth.CloseSesion()
		fmt.Println("Session closed.")
	},
}

// init registra el comando auth y sus subcomandos
func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd, registerCmd, logoutCmd)
}
