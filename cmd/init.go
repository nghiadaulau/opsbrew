package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nghiadaulau/opsbrew/internal/config"
	"github.com/nghiadaulau/opsbrew/internal/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [template] [project-name]",
	Short: "Initialize a new project from template",
	Long: `Initialize a new project from available templates.

Available templates:
  github-actions - GitHub Actions workflow template
  k8s-deployment - Kubernetes Deployment manifest
  k8s-service    - Kubernetes Service manifest
  k8s-pod        - Kubernetes Pod manifest
  k8s-configmap  - Kubernetes ConfigMap manifest
  dockerfile     - Multi-stage Dockerfile template`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("template name is required")
		}

		templateName := args[0]
		projectName := ""
		if len(args) > 1 {
			projectName = args[1]
		}

		// Get additional flags
		outputDir, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")

		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if dryRun {
			color.Yellow("Would initialize template: %s", templateName)
			if projectName != "" {
				color.Yellow("Project name: %s", projectName)
			}
			if outputDir != "" {
				color.Yellow("Output directory: %s", outputDir)
			}
			return nil
		}

		// Initialize template
		if err := templates.InitializeTemplate(templateName, projectName, outputDir, force, cfg); err != nil {
			return fmt.Errorf("failed to initialize template: %w", err)
		}

		color.Green("Project initialized successfully!")
		return nil
	},
}

var initListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates := templates.GetAvailableTemplates()

		fmt.Println("=== Available Templates ===")
		for _, template := range templates {
			color.Cyan("  %s", template.Name)
			fmt.Printf("    Description: %s\n", template.Description)
			fmt.Printf("    Files: %d\n", len(template.Files))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initListCmd)

	// Add flags for init
	initCmd.Flags().StringP("output", "o", "", "Output directory (default: current directory)")
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite existing files")
}
