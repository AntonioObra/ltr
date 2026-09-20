package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

type Snippet struct {
	ID        int
	Timestamp string
	Name      string
	CreatedAt time.Time
}

var snippetCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Manage your snippets in a clean md format",
	Long: `Manage your snippets in a clean md format

	Examples:
		ltr snippet list
		ltr snippet new <name>
		ltr snippet delete <id>`,
}

var snippetNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create new snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		snippetName := args[0]

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

		if err := os.Mkdir(filepath.Join(ltrDir, "snippets", timestamp), 0o755); err != nil {
			return fmt.Errorf("failed to create new snippet: %w", err)
		}

		snippetDefaultContent := []byte("# " + snippetName + "\n\n---\n")

		err = os.WriteFile(
			filepath.Join(ltrDir, "snippets", timestamp, "snippet.md"),
			snippetDefaultContent, 0o644,
		)
		if err != nil {
			return fmt.Errorf("failed to create new snippet: %w", err)
		}

		fmt.Printf("Created new snippet %s in %s\n", snippetName, dirName)
		return nil
	},
}

var snippetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of you  snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		snippetDir := filepath.Join(dir, ".ltr", "snippets")

		entries, err := os.ReadDir(snippetDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var snippets []Snippet

		for id, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			snippetFile := filepath.Join(snippetDir, entry.Name(), "snippet.md")

			content, err := os.ReadFile(snippetFile)
			if err != nil {
				return fmt.Errorf("failed to list snippets: %w", err)
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
				return fmt.Errorf("failed to list snippets: %w", err)
			}

			snippets = append(snippets, Snippet{
				ID:        id,
				Timestamp: entry.Name(),
				Name:      name,
				CreatedAt: createdAt,
			})
		}

		t := table.New().
			Headers("ID", "TIMESTAMP", "SNIPPET", "CREATED").
			Border(lipgloss.NormalBorder()).
			BorderRow(false)

		for _, snippet := range snippets {
			row := []string{
				strconv.Itoa(snippet.ID),
				snippet.Timestamp,
				snippet.Name,
				snippet.CreatedAt.Format("02.01.2006 15:04"),
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

var snippetDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete your snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		snippetId, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		snippetDir := filepath.Join(dir, ".ltr", "snippets")

		entries, err := os.ReadDir(snippetDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if snippetId < 0 || snippetId >= len(entries) {
			return fmt.Errorf("snippet ID %d doesn't exist in %s", snippetId, dirName)
		}

		entry := entries[snippetId]

		if err := os.RemoveAll(filepath.Join(snippetDir, entry.Name())); err != nil {
			return fmt.Errorf("Failed to delete snippet: #%d in %s, %w", snippetId, dirName, err)
		}

		fmt.Printf("Successfully deleted snippet: ID:%d", snippetId)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(snippetCmd)

	snippetCmd.AddCommand(snippetNewCmd)
	snippetCmd.AddCommand(snippetListCmd)
	snippetCmd.AddCommand(snippetDeleteCmd)
}
