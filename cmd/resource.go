package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

type Resource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var resourceCmd = &cobra.Command{
	Use:     "resource",
	Aliases: []string{"r"},
	Short:   "Manage your resources in a clean json format",
	Long: `Manage your resources in a clean json format

	Examples:
		ltr resource list
		ltr resource new <name> <url>
		ltr resource delete <id>`,
}

var resourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		resourcesFile := filepath.Join(dir, ".ltr", "resources.json")

		content, err := os.ReadFile(resourcesFile)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var resources []Resource

		if err := json.Unmarshal(content, &resources); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		t := table.New().
			Headers("ID", "NAME", "URL").
			Border(lipgloss.NormalBorder()).
			BorderRow(false)

		for index, resource := range resources {
			row := []string{
				strconv.Itoa(index),
				resource.Name,
				resource.URL,
			}

			if index%2 == 0 {
				t.Row(completedStyle.Render(row[0]),
					completedStyle.Render(row[1]),
					blueStyle.Render(row[2]),
				)
			} else {
				t.Row(pendingStyle.Render(row[0]),
					pendingStyle.Render(row[1]),
					pendingStyle.Render(row[2]),
				)
			}
		}

		fmt.Println(t)
		return nil
	},
}

var resourceNewCmd = &cobra.Command{
	Use:   "new [name] [url]",
	Short: "Add new resource",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceName := args[0]
		resourceURL := args[1]

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		resourcesFile := filepath.Join(dir, ".ltr", "resources.json")

		content, err := os.ReadFile(resourcesFile)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var resources []Resource

		if err := json.Unmarshal(content, &resources); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		resources = append(resources, Resource{
			Name: resourceName,
			URL:  resourceURL,
		})

		data, err := json.MarshalIndent(resources, "", "	")
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		err = os.WriteFile(resourcesFile, data, 0o644)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		fmt.Print("Successfully added new resource!\n")
		return resourceListCmd.RunE(resourceListCmd, []string{})
	},
}

var resourceDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete resource",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirName := filepath.Base(dir)
		resourcesFile := filepath.Join(dir, ".ltr", "resources.json")

		content, err := os.ReadFile(resourcesFile)
		if err != nil {
			return fmt.Errorf("ltr doesnt exists in: %s, %w", dirName, err)
		}

		var resources []Resource

		if err := json.Unmarshal(content, &resources); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		if resourceID < 0 || resourceID >= len(resources) {
			return fmt.Errorf("resource ID %d doesn't exist", resourceID)
		}

		resources = append(resources[:resourceID], resources[resourceID+1:]...)

		data, err := json.MarshalIndent(resources, "", "	")
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		err = os.WriteFile(resourcesFile, data, 0o644)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		fmt.Printf("Successfully deleted resource: ID:%d\n", resourceID)
		return resourceListCmd.RunE(resourceListCmd, []string{})
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	resourceCmd.AddCommand(resourceListCmd)
	resourceCmd.AddCommand(resourceNewCmd)
	resourceCmd.AddCommand(resourceDeleteCmd)
}
