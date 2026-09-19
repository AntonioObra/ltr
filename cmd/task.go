package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage your tasks in a clean md format",
	Long: `Manage your tasks in a clean md format

	Examples:
		ltr task list
		ltr task new <name>`,
}

var taskNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskName := args[0]

		fmt.Printf("Created new task %s in \n", taskName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskNewCmd)
}
