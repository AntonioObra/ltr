package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
		ltr task new <name>
		ltr task open <id>
		ltr task delete <id>
		ltr task complete <id>`,
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

		taskDefaultContent := []byte("# " + taskName + "\n\n---\n\n- COMPLETED: [FALSE]\n\n---\n")
		taskFile := filepath.Join(ltrDir, "tasks", timestamp, "task.md")

		err = os.WriteFile(
			taskFile,
			taskDefaultContent, 0o644,
		)
		if err != nil {
			return fmt.Errorf("failed to create new task: %w", err)
		}

		if err := openInEditor(taskFile); err != nil {
			return fmt.Errorf("failed to open task: %w", err)
		}

		fmt.Printf("Created new task %s in %s\n", taskName, dirName)
		return nil
	},
}

var (
	listCompleted bool
	listAll       bool
)

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your tasks",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if listCompleted && listAll {
			return fmt.Errorf("cannot use -c and -a together")
		}

		return nil
	},
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

			completed := strings.Contains(string(content), "- COMPLETED: [TRUE]")

			if !listAll {
				if listCompleted && !completed {
					continue
				}

				if !listCompleted && completed {
					continue
				}
			}

			name := strings.TrimPrefix(lines[0], "# ")
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

		t := table.New().
			Headers("ID", "TIMESTAMP", "TASK", "COMPLETED", "CREATED").
			Border(lipgloss.NormalBorder()).
			BorderRow(false)

		for _, task := range tasks {
			row := []string{
				strconv.Itoa(task.ID),
				task.Timestamp,
				task.Name,
				strconv.FormatBool(task.Completed),
				task.CreatedAt.Format("02.01.2006 15:04"),
			}

			if task.Completed {
				t.Row(completedStyle.Render(row[0]),
					completedStyle.Render(row[1]),
					completedStyle.Render(row[2]),
					completedStyle.Render(row[3]),
					completedStyle.Render(row[4]),
				)
			} else {
				t.Row(pendingStyle.Render(row[0]),
					pendingStyle.Render(row[1]),
					pendingStyle.Render(row[2]),
					pendingStyle.Render(row[3]),
					pendingStyle.Render(row[4]),
				)
			}
		}

		fmt.Println(t)
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

var taskOpenCmd = &cobra.Command{
	Use:   "open [id]",
	Short: "Open your task in neovim",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		dirName := filepath.Base(dir)
		taskDir := filepath.Join(dir, ".ltr", "tasks")

		entries, err := os.ReadDir(taskDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if taskID < 0 || taskID >= len(entries) {
			return fmt.Errorf("task ID %d doesn't exist in %s", taskID, dirName)
		}

		entry := entries[taskID]

		if !entry.IsDir() {
			return fmt.Errorf("task ID %d is not a task", taskID)
		}

		taskFile := filepath.Join(taskDir, entry.Name(), "task.md")

		if err := openInEditor(taskFile); err != nil {
			return fmt.Errorf("failed to open task: %w", err)
		}

		return nil
	},
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		dirName := filepath.Base(dir)
		taskDir := filepath.Join(dir, ".ltr", "tasks")

		entries, err := os.ReadDir(taskDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if taskID < 0 || taskID >= len(entries) {
			return fmt.Errorf("task ID %d doesn't exist in %s", taskID, dirName)
		}

		entry := entries[taskID]

		if !entry.IsDir() {
			return fmt.Errorf("task ID %d is not a task", taskID)
		}

		taskFile := filepath.Join(taskDir, entry.Name(), "task.md")

		content, err := os.ReadFile(taskFile)
		if err != nil {
			return fmt.Errorf("failed to read task: %w", err)
		}

		old := "- COMPLETED: [FALSE]"
		new := "- COMPLETED: [TRUE]"

		updatedContent := strings.Replace(string(content), old, new, 1)

		if string(content) == updatedContent {
			return fmt.Errorf("task %d is already completed", taskID)
		}

		if err := os.WriteFile(taskFile, []byte(updatedContent), 0o644); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Printf("Task %d marked as completed!\n", taskID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskNewCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskOpenCmd)
	taskCmd.AddCommand(taskCompleteCmd)

	taskListCmd.Flags().BoolVarP(
		&listCompleted,
		"completed",
		"c",
		false,
		"list completed tasks",
	)

	taskListCmd.Flags().BoolVarP(
		&listAll,
		"all",
		"a",
		false,
		"list all tasks",
	)
}
