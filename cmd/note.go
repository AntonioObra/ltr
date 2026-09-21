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

type Note struct {
	ID        int
	Timestamp string
	Name      string
	CreatedAt time.Time
}

var noteCmd = &cobra.Command{
	Use:     "note",
	Aliases: []string{"n"},
	Short:   "Manage your notes in a clean md format",
	Long: `Manage your notes in a clean md format

	Examples:
		ltr note list
		ltr note new <name>
		ltr note open <id>
		ltr note delete <id>`,
}

var noteNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create new note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteName := args[0]

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

		if err := os.Mkdir(filepath.Join(ltrDir, "notes", timestamp), 0o755); err != nil {
			return fmt.Errorf("failed to create new note: %w", err)
		}

		noteDefaultContent := []byte("# " + noteName + "\n\nWrite your note here...")
		noteFile := filepath.Join(ltrDir, "notes", timestamp, "note.md")

		err = os.WriteFile(
			noteFile,
			noteDefaultContent, 0o644,
		)
		if err != nil {
			return fmt.Errorf("failed to create new note: %w", err)
		}

		if err := openInEditor(noteFile); err != nil {
			return fmt.Errorf("failed to open note: %w", err)
		}

		fmt.Printf("Created new note %s in %s\n", noteName, dirName)
		return noteListCmd.RunE(taskListCmd, []string{})
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		noteDir := filepath.Join(dir, ".ltr", "notes")

		entries, err := os.ReadDir(noteDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var notes []Note

		for id, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			noteFile := filepath.Join(noteDir, entry.Name(), "note.md")

			content, err := os.ReadFile(noteFile)
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}

			lines := strings.Split(string(content), "\n")

			if len(lines) == 0 {
				continue
			}

			name := strings.TrimPrefix(lines[0], "# ")
			createdAt, err := time.ParseInLocation(
				"20060102150405",
				entry.Name(),
				time.Local,
			)
			if err != nil {
				return fmt.Errorf("failed to list notes: %w", err)
			}

			notes = append(notes, Note{
				ID:        id,
				Timestamp: entry.Name(),
				Name:      name,
				CreatedAt: createdAt,
			})
		}

		t := table.New().
			Headers("ID", "TIMESTAMP", "note", "CREATED").
			Border(lipgloss.NormalBorder()).
			BorderRow(false)

		for _, note := range notes {
			row := []string{
				strconv.Itoa(note.ID),
				note.Timestamp,
				note.Name,
				note.CreatedAt.Format("02.01.2006 15:04"),
			}

			t.Row(normalStyle.Render(row[0]),
				normalStyle.Render(row[1]),
				normalStyle.Render(row[2]),
				normalStyle.Render(row[3]),
			)
		}

		fmt.Println(t)
		return nil
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete your note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteId, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		noteDir := filepath.Join(dir, ".ltr", "notes")

		entries, err := os.ReadDir(noteDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if noteId < 0 || noteId >= len(entries) {
			return fmt.Errorf("note ID %d doesn't exist in %s", noteId, dirName)
		}

		entry := entries[noteId]

		if err := os.RemoveAll(filepath.Join(noteDir, entry.Name())); err != nil {
			return fmt.Errorf("Failed to delete note: #%d in %s, %w", noteId, dirName, err)
		}

		fmt.Printf("Successfully deleted note: ID:%d\n", noteId)
		return noteListCmd.RunE(taskListCmd, []string{})
	},
}

var noteOpenCmd = &cobra.Command{
	Use:   "open [id]",
	Short: "Open your note in neovim",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID: %w", err)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		dirName := filepath.Base(dir)
		noteDir := filepath.Join(dir, ".ltr", "notes")

		entries, err := os.ReadDir(noteDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if noteID < 0 || noteID >= len(entries) {
			return fmt.Errorf("note ID %d doesn't exist in %s", noteID, dirName)
		}

		entry := entries[noteID]

		if !entry.IsDir() {
			return fmt.Errorf("note ID %d is not a note", noteID)
		}

		noteFile := filepath.Join(noteDir, entry.Name(), "note.md")

		if err := openInEditor(noteFile); err != nil {
			return fmt.Errorf("failed to open note: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)

	noteCmd.AddCommand(noteNewCmd)
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	noteCmd.AddCommand(noteOpenCmd)
}
