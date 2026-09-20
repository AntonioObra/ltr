/**
*
*  /\_/\
* ( o.o )
*  > ^ <
*
* Root cmd
*
* Root command for the Last Trace CLI application.
*
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var (
	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3"))

	normalStyle = lipgloss.NewStyle()
)

var rootCmd = &cobra.Command{
	Use:   "ltr",
	Short: "A CLI version of Last Trace",
	Long:  "CLI version of Last Trace",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// * Execute adds all child commands to the root command and sets flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// * Called before any command runs
	cobra.OnInitialize(initConfig)

	// * Define persistent flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.ltr.yaml)")

	// * Set up environment variable prefix
	// * LTR_DEFAULTS_PRIORITY will map to defaults.priority
	viper.SetEnvPrefix("LTR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// * Look for config file in multiple locations
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ltr")
	}

	// * Read env variables
	viper.AutomaticEnv()

	// * Set defaults
	// viper.SetDefault("test", "idegasbre")

	// * Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		}
	}
}
