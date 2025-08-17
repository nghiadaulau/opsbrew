package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"opsbrew/internal/config"
)

var (
	cfgFile string
	verbose bool
	dryRun  bool
	confirm bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opsbrew",
	Short: "A CLI tool to simplify and shorten repetitive DevOps terminal commands",
	Long: `opsbrew is a powerful CLI tool designed to streamline DevOps workflows.

It provides shortcuts for common Git and kubectl operations, with features like:
- Fuzzy finder for branches, contexts, and namespaces
- Command recipes and macros
- Safe defaults with dry-run and confirmation modes
- Project template initialization

Examples:
  opsbrew git status
  opsbrew git sync
  opsbrew kctx
  opsbrew kns
  opsbrew klogs
  opsbrew init go-service
  opsbrew brew save my-workflow`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opsbrew.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would be done without executing")
	rootCmd.PersistentFlags().BoolVar(&confirm, "confirm", false, "skip confirmation prompts")

	// Local flags
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".opsbrew" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".opsbrew")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			color.Green("Using config file: %s", viper.ConfigFileUsed())
		}
	} else {
		// Create default config if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := config.CreateDefaultConfig(); err != nil {
				color.Red("Error creating default config: %v", err)
			} else {
				color.Green("Created default config file: %s", viper.ConfigFileUsed())
			}
		}
	}
}
