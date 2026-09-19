package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		ltrDir := filepath.Join(dir, ".ltr")

		if _, err := os.Stat(ltrDir); err != nil {
			fmt.Printf("ltr doesnt exists in: %s", dirName)
			return nil
		}

		t := time.Now()
		timestamp := t.Format("20060102150405")

		if err := os.Mkdir(filepath.Join(ltrDir, "tasks", timestamp), 0o755); err != nil {
			return fmt.Errorf("failed to create new task: %w", err)
		}

		taskDefaultContent := []byte("# " + taskName + "\n\n ---\n\n - COMPLETED: [FALSE] \n\n---\n\n")

		err = os.WriteFile(
			filepath.Join(ltrDir, "tasks", timestamp, "task.md"),
			taskDefaultContent, 0o644,
		)
		if err != nil {
			return fmt.Errorf("failed to create new task: %w", err)
		}

		fmt.Printf("Created new task %s in %s\n", taskName, dirName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskNewCmd)
}
