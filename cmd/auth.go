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

// authCmd groups the system authentication subcommands
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management (login, register, logout)",
}

// loginCmd asks for credentials and starts a session when they are valid
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

		dataDir, _ := config.GetDataDirectory()
		routeUsers := filepath.Join(dataDir, "users.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Login(strings.ToLower(username), password); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Session started as '%s'\n", username)
	},
}

// registerCmd asks for input data and registers a new user
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

		dataDir, _ := config.GetDataDirectory()
		routeUsers := filepath.Join(dataDir, "users.json")

		svc := auth.NewAuthService(&storage.Storage{Route: routeUsers})
		if err := svc.Register(strings.ToLower(username), password); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Registration successful. Use 'task-cli auth login' to log in.")
	},
}

// logoutCmd closes the current active user session
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close the current session",
	Run: func(cmd *cobra.Command, args []string) {
		auth.CloseSession()
		fmt.Println("Session closed.")
	},
}

// init registers the auth command and its subcommands
func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd, registerCmd, logoutCmd)
}
