package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Snippet struct {
	ID        int
	Timestamp string
	Name      string
	Language  string
	CreatedAt time.Time
}

var snippetCmd = &cobra.Command{
	Use:     "snippet",
	Aliases: []string{"s"},
	Short:   "Manage your snippets in a clean md format",
	Long: `Manage your snippets in a clean md format

	Examples:
		ltr snippet list
		ltr snippet new <name>
		ltr snippet open <id>
		ltr snippet delete <id>`,
}

var snippetLanguage string

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

		language := "NONE"

		if snippetLanguage != "" {
			language = strings.ToUpper(snippetLanguage)
		}

		snippetDefaultContent := []byte("# " + snippetName + "\n\n---\n\n- LANGUAGE: [" + language + "]\n\n---\n\n```" + language + "\nfunc main() {}\n```\n")
		snippetFile := filepath.Join(ltrDir, "snippets", timestamp, "snippet.md")

		err = os.WriteFile(
			snippetFile,
			snippetDefaultContent, 0o644,
		)
		if err != nil {
			return fmt.Errorf("failed to create new snippet: %w", err)
		}

		if err := openInEditor(snippetFile); err != nil {
			return fmt.Errorf("failed to open snippet: %w", err)
		}

		fmt.Printf("Created new snippet %s in %s\n", snippetName, dirName)
		return snippetListCmd.RunE(taskListCmd, []string{})
	},
}

var snippetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your snippets",
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

		if len(entries) == 0 {
			fmt.Println(redStyle.Render("No snippets found..."))
			return nil
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

			var language string

			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(line, "- LANGUAGE: [") {
					language = strings.TrimPrefix(line, "- LANGUAGE: [")
					language = strings.TrimSuffix(language, "]")
				}
			}

			snippets = append(snippets, Snippet{
				ID:        id,
				Timestamp: entry.Name(),
				Name:      name,
				Language:  language,
				CreatedAt: createdAt,
			})
		}

		slices.Reverse(snippets)

		t := table.New().
			Headers("ID", "TIMESTAMP", "SNIPPET", "LANGUAGE", "CREATED").
			Border(lipgloss.NormalBorder()).
			BorderRow(false)

		for index, snippet := range snippets {
			row := []string{
				strconv.Itoa(snippet.ID),
				snippet.Timestamp,
				snippet.Name,
				snippet.Language,
				snippet.CreatedAt.Format("02.01.2006 15:04"),
			}

			if index%2 == 0 {
				t.Row(blueStyle.Render(row[0]),
					blueStyle.Render(row[1]),
					blueStyle.Render(row[2]),
					blueStyle.Render(row[3]),
					blueStyle.Render(row[4]),
				)
			} else {
				t.Row(redStyle.Render(row[0]),
					redStyle.Render(row[1]),
					redStyle.Render(row[2]),
					redStyle.Render(row[3]),
					redStyle.Render(row[4]),
				)
			}

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

		fmt.Printf("Successfully deleted snippet: ID:%d\n", snippetId)
		return snippetListCmd.RunE(taskListCmd, []string{})
	},
}

var snippetOpenCmd = &cobra.Command{
	Use:   "open [id]",
	Short: "Open your snippet in neovim",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		snippetID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid snippet ID: %w", err)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		dirName := filepath.Base(dir)
		snippetDir := filepath.Join(dir, ".ltr", "snippets")

		entries, err := os.ReadDir(snippetDir)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		if snippetID < 0 || snippetID >= len(entries) {
			return fmt.Errorf("snippet ID %d doesn't exist in %s", snippetID, dirName)
		}

		entry := entries[snippetID]

		if !entry.IsDir() {
			return fmt.Errorf("snippet ID %d is not a snippet", snippetID)
		}

		snippetFile := filepath.Join(snippetDir, entry.Name(), "snippet.md")

		if err := openInEditor(snippetFile); err != nil {
			return fmt.Errorf("failed to open snippet: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(snippetCmd)

	snippetCmd.AddCommand(snippetNewCmd)
	snippetCmd.AddCommand(snippetListCmd)
	snippetCmd.AddCommand(snippetDeleteCmd)
	snippetCmd.AddCommand(snippetOpenCmd)

	snippetNewCmd.Flags().StringVarP(
		&snippetLanguage,
		"lang",
		"l",
		"",
		"snippet language",
	)
}
