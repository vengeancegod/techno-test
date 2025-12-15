package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "taskmanager",
		Short: "task management cli",
		Long: `A command line interface application for managing your tasks. You can create, list, update, and delete tasks`,
		Version: "1.0.0",
	}
}
