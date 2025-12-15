package cli

import (
	"context"
	"fmt"
	"strconv"
	"techno/internal/model"
	"techno/internal/service"

	"github.com/spf13/cobra"
)

type TaskCommands struct {
	taskService service.TaskService
}

func NewTaskCommands(taskService service.TaskService) *TaskCommands {
	return &TaskCommands{
		taskService: taskService,
	}
}

func (tc *TaskCommands) RegisterCommands(rootCmd *cobra.Command) {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
		Long:  "Create, read, update, delete tasks",
	}

	taskCmd.AddCommand(tc.createCmd())
	taskCmd.AddCommand(tc.listCmd())
	taskCmd.AddCommand(tc.getCmd())
	taskCmd.AddCommand(tc.updateCmd())
	taskCmd.AddCommand(tc.deleteCmd())

	rootCmd.AddCommand(taskCmd)
}

func (tc *TaskCommands) createCmd() *cobra.Command {
	var title, description string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new task",
		Long:    "Create a new task with title and optional description",
		Example: `  taskmanager task create -t "Buy groceries" -d "Milk, bread, eggs" taskmanager task create --title "Meeting" --description "Team sync at 3pm"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			task := &model.Task{
				Title:       title,
				Description: description,
			}

			if err := tc.taskService.CreateTask(context.Background(), task); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("Task created successfull\n")
			fmt.Printf("ID: %d\n", task.ID)
			fmt.Printf("Title: %s\n", task.Title)
			fmt.Printf("Status: %s\n", task.Status.StringStatus())
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Task title (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Task description")
	cmd.MarkFlagRequired("title")

	return cmd
}

func (tc *TaskCommands) listCmd() *cobra.Command {
	var statusStr string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all tasks",
		Long:    "List all tasks or filter by status",
		Example: `  taskmanager task list taskmanager task list -s pending taskmanager task list --status completed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var tasks []*model.Task
			var err error

			if statusStr != "" {
				taskStatus := model.ParseTaskStatus(statusStr)
				tasks, err = tc.taskService.GetByStatus(context.Background(), taskStatus)
				if err != nil {
					return fmt.Errorf("failed to get tasks by status: %w", err)
				}
			} else {
				tasks, err = tc.taskService.GetAll(context.Background())
				if err != nil {
					return fmt.Errorf("failed to get all tasks: %w", err)
				}
			}
			if len(tasks) == 0 {
				fmt.Println("No tasks found")
				return nil
			}

			fmt.Printf("\n%-5s %-40s %-30s %-15s %-20s\n", "ID", "Title", "Description", "Status", "Created At")
			for _, task := range tasks {
				title := task.Title
				if len(title) > 40 {
					title = title[:37] + "..."
				}
				fmt.Printf("%-5d %-40s %-15s %-20s\n", task.ID, title, task.Status.StringStatus(), task.CreatedAt.Format("2006-01-02 15:04"))
			}
			fmt.Printf("\nTotal: %d task(s)\n\n", len(tasks))
			return nil
		},
	}

	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "Filter by status (done/not_done)")
	return cmd
}

func (tc *TaskCommands) getCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get [id]",
		Short:   "Get task by ID",
		Long:    "Display detailed information about a specific task",
		Example: `  taskmanager task get 1 taskmanager task get 42`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %w", err)
			}

			task, err := tc.taskService.GetByID(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}
			fmt.Printf("ID:          %d\n", task.ID)
			fmt.Printf("Title:       %s\n", task.Title)
			fmt.Printf("Description: %s\n", task.Description)
			fmt.Printf("Status:      %s\n", task.Status.StringStatus())
			fmt.Printf("Created At:  %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
			return nil
		},
	}
}

func (tc *TaskCommands) updateCmd() *cobra.Command {
	var title, description, statusStr string

	cmd := &cobra.Command{
		Use:     "update [id]",
		Short:   "Update a task",
		Long:    "Update task title, description, or status",
		Example: `  taskmanager task update 1 -t "New title" taskmanager task update 1 -s completed taskmanager task update 1 -t "New title" -d "New description" -s in_progress`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %w", err)
			}

			existingTask, err := tc.taskService.GetByID(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			if title != "" {
				existingTask.Title = title
			}
			if description != "" {
				existingTask.Description = description
			}
			if statusStr != "" {
				existingTask.Status = model.ParseTaskStatus(statusStr)
			}

			if err := tc.taskService.UpdateTask(context.Background(), existingTask); err != nil {
				return fmt.Errorf("failed to update task: %w", err)
			}

			fmt.Printf("Task %d updated successfull\n", id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "New task title")
	cmd.Flags().StringVarP(&description, "description", "d", "", "New task description")
	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "New task status (done/not_done)")

	return cmd
}

func (tc *TaskCommands) deleteCmd() *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:     "delete [id]",
		Short:   "Delete a task",
		Long:    "Delete a task by its ID",
		Example: `  taskmanager task delete 1 taskmanager task delete 1 -y`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %w", err)
			}

			if !confirm {
				fmt.Printf("Are you sure you want to delete task %d? [y/N]: ", id)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			if err := tc.taskService.DeleteTask(context.Background(), id); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			fmt.Printf("Task %d deleted successfull\n", id)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&confirm, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
