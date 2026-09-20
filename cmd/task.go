package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type Task struct {
	ID        int
	Timestamp string
	Name      string
	Completed bool
	CreatedAt time.Time
}

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

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of you  tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		taskDir := filepath.Join(dir, ".ltr", "tasks")

		entries, err := os.ReadDir(taskDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var tasks []Task

		for id, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			taskFile := filepath.Join(taskDir, entry.Name(), "task.md")

			content, err := os.ReadFile(taskFile)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			lines := strings.Split(string(content), "\n")

			if len(lines) == 0 {
				continue
			}

			name := strings.TrimPrefix(lines[0], "# ")
			completed := strings.Contains(string(content), "- COMPLETED: [TRUE]")
			createdAt, err := time.ParseInLocation(
				"20060102150405",
				entry.Name(),
				time.Local,
			)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			tasks = append(tasks, Task{
				ID:        id,
				Timestamp: entry.Name(),
				Name:      name,
				Completed: completed,
				CreatedAt: createdAt,
			})
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIMESTAMP\tTASK\tCOMPLETED\tCREATED")

		for _, task := range tasks {
			fmt.Fprintf(
				w,
				"%d\t%s\t%s\t%t\t%s\n",
				task.ID,
				task.Timestamp,
				task.Name,
				task.Completed,
				task.CreatedAt.Format("02.01.2006 15:04"),
			)
		}

		w.Flush()

		return nil
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete your task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		taskDir := filepath.Join(dir, ".ltr", "tasks")

		entries, err := os.ReadDir(taskDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if taskId < 0 || taskId >= len(entries) {
			return fmt.Errorf("task ID %d doesn't exist in %s", taskId, dirName)
		}

		entry := entries[taskId]

		if err := os.RemoveAll(filepath.Join(taskDir, entry.Name())); err != nil {
			return fmt.Errorf("Failed to delete task: #%d in %s, %w", taskId, dirName, err)
		}

		fmt.Printf("Successfully deleted task: ID:%d", taskId)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskNewCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskDeleteCmd)
}
