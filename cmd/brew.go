package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/nghiadaulau/opsbrew/internal/config"
)

var brewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Manage and run command recipes/macros",
	Long: `Brew allows you to save and run command recipes for daily workflows.

Available commands:
  save     - Save a new recipe
  list     - List all saved recipes
  run      - Run a saved recipe
  delete   - Delete a saved recipe
  edit     - Edit a saved recipe`,
}

var brewSaveCmd = &cobra.Command{
	Use:   "save [name]",
	Short: "Save a new recipe",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("recipe name is required")
		}

		name := args[0]
		description, _ := cmd.Flags().GetString("description")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// Get commands from user
		fmt.Printf("Enter commands for recipe '%s' (one per line, empty line to finish):\n", name)
		var commands []string
		for {
			fmt.Print("> ")
			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				color.Red("Error reading input: %v", err)
				return err
			}
			if input == "" {
				break
			}
			commands = append(commands, input)
		}

		if len(commands) == 0 {
			return fmt.Errorf("no commands provided")
		}

		// Load current config
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Add recipe
		cfg.Brew.Recipes[name] = config.Recipe{
			Description: description,
			Commands:    commands,
			Tags:        tags,
		}

		// Save config
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save recipe: %w", err)
		}

		color.Green("Recipe '%s' saved successfully", name)
		return nil
	},
}

var brewListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved recipes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(cfg.Brew.Recipes) == 0 {
			color.Yellow("No recipes found")
			return nil
		}

		fmt.Println("=== Saved Recipes ===")
		for name, recipe := range cfg.Brew.Recipes {
			color.Cyan("  %s", name)
			if recipe.Description != "" {
				fmt.Printf("    Description: %s\n", recipe.Description)
			}
			fmt.Printf("    Commands: %d\n", len(recipe.Commands))
			if len(recipe.Tags) > 0 {
				fmt.Printf("    Tags: %s\n", strings.Join(recipe.Tags, ", "))
			}
			fmt.Println()
		}

		return nil
	},
}

var brewRunCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run a saved recipe",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("recipe name is required")
		}

		name := args[0]
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		recipe, exists := cfg.Brew.Recipes[name]
		if !exists {
			return fmt.Errorf("recipe '%s' not found", name)
		}

		if dryRun {
			color.Yellow("Would run recipe '%s':", name)
			for i, command := range recipe.Commands {
				color.Yellow("  %d. %s", i+1, command)
			}
			return nil
		}

		// Check if we need confirmation
		if !confirm && !cfg.UI.Confirm {
			fmt.Printf("Run recipe '%s'? (y/N): ", name)
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				color.Red("Error reading input: %v", err)
				return err
			}
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				color.Yellow("Operation cancelled")
				return nil
			}
		}

		color.Green("Running recipe: %s", name)
		if recipe.Description != "" {
			fmt.Printf("Description: %s\n", recipe.Description)
		}
		fmt.Println()

		// Execute commands
		for i, command := range recipe.Commands {
			color.Cyan("Executing command %d/%d: %s", i+1, len(recipe.Commands), command)

			// Split command into parts
			parts := strings.Fields(command)
			if len(parts) == 0 {
				continue
			}

			cmdExec := exec.Command(parts[0], parts[1:]...)
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Stdin = os.Stdin

			if err := cmdExec.Run(); err != nil {
				color.Red("Command failed: %s", command)
				return fmt.Errorf("recipe execution failed: %w", err)
			}

			fmt.Println()
		}

		color.Green("Recipe '%s' completed successfully", name)
		return nil
	},
}

var brewDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a saved recipe",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("recipe name is required")
		}

		name := args[0]
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if _, exists := cfg.Brew.Recipes[name]; !exists {
			return fmt.Errorf("recipe '%s' not found", name)
		}

		if dryRun {
			color.Yellow("Would delete recipe: %s", name)
			return nil
		}

		// Check if we need confirmation
		if !confirm && !cfg.UI.Confirm {
			fmt.Printf("Delete recipe '%s'? (y/N): ", name)
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				color.Red("Error reading input: %v", err)
				return err
			}
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				color.Yellow("Operation cancelled")
				return nil
			}
		}

		delete(cfg.Brew.Recipes, name)

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to delete recipe: %w", err)
		}

		color.Green("Recipe '%s' deleted successfully", name)
		return nil
	},
}

var brewEditCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a saved recipe",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("recipe name is required")
		}

		name := args[0]
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		recipe, exists := cfg.Brew.Recipes[name]
		if !exists {
			return fmt.Errorf("recipe '%s' not found", name)
		}

		// Show current recipe
		fmt.Printf("Current recipe '%s':\n", name)
		fmt.Printf("Description: %s\n", recipe.Description)
		fmt.Printf("Tags: %s\n", strings.Join(recipe.Tags, ", "))
		fmt.Println("Commands:")
		for i, command := range recipe.Commands {
			fmt.Printf("  %d. %s\n", i+1, command)
		}
		fmt.Println()

		// Get new description
		fmt.Print("New description (press Enter to keep current): ")
		var newDescription string
		if _, err := fmt.Scanln(&newDescription); err != nil {
			color.Red("Error reading input: %v", err)
			return err
		}
		if newDescription != "" {
			recipe.Description = newDescription
		}

		// Get new tags
		fmt.Print("New tags (comma-separated, press Enter to keep current): ")
		var newTags string
		if _, err := fmt.Scanln(&newTags); err != nil {
			color.Red("Error reading input: %v", err)
			return err
		}
		if newTags != "" {
			recipe.Tags = strings.Split(newTags, ",")
			for i, tag := range recipe.Tags {
				recipe.Tags[i] = strings.TrimSpace(tag)
			}
		}

		// Get new commands
		fmt.Println("Enter new commands (one per line, empty line to finish):")
		var newCommands []string
		for {
			fmt.Print("> ")
			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				color.Red("Error reading input: %v", err)
				return err
			}
			if input == "" {
				break
			}
			newCommands = append(newCommands, input)
		}

		if len(newCommands) > 0 {
			recipe.Commands = newCommands
		}

		// Save updated recipe
		cfg.Brew.Recipes[name] = recipe

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save recipe: %w", err)
		}

		color.Green("Recipe '%s' updated successfully", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.AddCommand(brewSaveCmd)
	brewCmd.AddCommand(brewListCmd)
	brewCmd.AddCommand(brewRunCmd)
	brewCmd.AddCommand(brewDeleteCmd)
	brewCmd.AddCommand(brewEditCmd)

	// Add flags for brew save
	brewSaveCmd.Flags().StringP("description", "d", "", "Recipe description")
	brewSaveCmd.Flags().StringSliceP("tags", "t", []string{}, "Recipe tags")
}
