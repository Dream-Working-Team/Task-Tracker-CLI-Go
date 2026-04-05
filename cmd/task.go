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

// GetServiceTask validates the active session and creates the user's task service
func GetServiceTask() (*service.TaskService, error) {
	userID, err := auth.GetActiveUser()
	if err != nil {
		return nil, err
	}

	dataDir, err := config.GetDataDirectory()
	if err != nil {
		return nil, err
	}

	taskFile := filepath.Join(dataDir, fmt.Sprintf("task_%s.json", userID))
	return service.NewTaskService(&storage.Storage{Route: taskFile}), nil
}

// addCmd creates a new task using the provided description
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

		task, err := svc.Add(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Keep this exact output format to match the expected behavior
		fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	},
}

// updateCmd updates the description of an existing task by ID
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

// deleteCmd removes an existing task permanently by ID
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

// markInProgressCmd changes a task status to in-progress
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

// markDoneCmd changes a task status to done
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

// listCmd shows all tasks or filters them by status
var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "Task list",
	Args:  cobra.MaximumNArgs(1), // Accepts 0 or 1 positional argument
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := GetServiceTask()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		statusFilter := ""
		if len(args) == 1 {
			switch args[0] {
			case "done":
				statusFilter = model.Complete
			case "todo":
				statusFilter = model.ToDo
			case "in-progress":
				statusFilter = model.InProgress
			default:
				fmt.Println("Invalid filter. Use: done, todo, in-progress")
				return
			}
		}

		tasks, err := svc.List(statusFilter)
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

// init registers task commands on the root command
func init() {
	rootCmd.AddCommand(addCmd, updateCmd, deleteCmd, markTodoCmd, markInProgressCmd, markDoneCmd, listCmd)
}
