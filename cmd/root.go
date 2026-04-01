package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "A multi-user task manager for the terminal",
	Long:  "Task Tracker CLI Secure. Manage tasks per user with isolated and persistent storage.",
}

func init() {
	rootCmd.SetHelpTemplate(
		`{{with .Long}}{{.}}{{"\n\n"}}{{end}}` +
			`Usage:
  {{.UseLine}}

` +
			`Available Commands:
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}
` +
			`Examples:
  task-cli auth register
  task-cli auth login
  task-cli add "Study Go"
  task-cli list in-progress

` +
			`Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

` +
			`Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
