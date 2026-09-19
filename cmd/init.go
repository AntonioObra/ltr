package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init new ltr instance in your project",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		ltrDir := filepath.Join(dir, ".ltr")

		if _, err := os.Stat(ltrDir); err == nil {
			fmt.Printf("ltr already exists in: %s", dirName)
			return nil
		}

		if err := os.MkdirAll(ltrDir, 0o755); err != nil {
			return fmt.Errorf("failed to init new ltr instance: %w", err)
		}

		if err := os.MkdirAll(filepath.Join(ltrDir, "tasks"), 0o755); err != nil {
			return fmt.Errorf("failed to init new ltr instance: %w", err)
		}

		fmt.Printf("Successfully created new ltr instance in %s directory", dirName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
