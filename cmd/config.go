/**
*
*  /\_/\
* ( o.o )
*  > ^ <
*
* Config cmd
*
* Set configuration values.
*
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `View and modify configuration settings.
	
	Examples:
		ltr config get defaults.priority
		ltr config set defaults.priority high
		ltr config list`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.Get(key)

		if value == nil {
			return fmt.Errorf("configuration key %q not found", key)
		}

		fmt.Printf("%s: %v\n", key, value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		viper.Set(key, value)

		if err := viper.WriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				home, _ := os.UserHomeDir()
				configPath := filepath.Join(home, ".ltr.yaml")

				if err := viper.WriteConfigAs(configPath); err != nil {
					return fmt.Errorf("failed to write config: %w", err)
				}
			} else {
				return fmt.Errorf("failed to write config: %w", err)
			}
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := viper.AllSettings()

		if len(settings) == 0 {
			fmt.Println("No configuration values set")
			return nil
		}

		for key, value := range settings {
			fmt.Printf("%s: %v\n", key, value)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
}
