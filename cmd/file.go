package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File operations and shortcuts",
	Long: `File operations and shortcuts for common tasks.

Available commands:
  open     - Open file with default editor
  find     - Find files by name or pattern
  grep     - Search for text in files
  backup   - Create backup of file
  diff     - Show differences between files`,
}

var fileOpenCmd = &cobra.Command{
	Use:   "open [file]",
	Short: "Open file with default editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("file path is required")
		}

		filePath := args[0]

		if dryRun {
			color.Yellow("Would open file: %s", filePath)
			return nil
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", filePath)
		}

		// Try to open with default editor
		var cmdExec *exec.Cmd
		switch os := runtime.GOOS; os {
		case "darwin":
			cmdExec = exec.Command("open", filePath)
		case "linux":
			cmdExec = exec.Command("xdg-open", filePath)
		case "windows":
			cmdExec = exec.Command("cmd", "/c", "start", filePath)
		default:
			return fmt.Errorf("unsupported operating system: %s", os)
		}

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		color.Green("Opened file: %s", filePath)
		return nil
	},
}

var fileFindCmd = &cobra.Command{
	Use:   "find [pattern]",
	Short: "Find files by name or pattern",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("search pattern is required")
		}

		pattern := args[0]
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}

		if dryRun {
			color.Yellow("Would search for pattern '%s' in directory '%s'", pattern, dir)
			return nil
		}

		// Use find command
		cmdExec := exec.Command("find", dir, "-name", pattern, "-type", "f")
		output, err := cmdExec.Output()
		if err != nil {
			return fmt.Errorf("failed to find files: %w", err)
		}

		files := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			color.Yellow("No files found matching pattern: %s", pattern)
			return nil
		}

		color.Green("Found %d files:", len(files))
		for _, file := range files {
			if file != "" {
				fmt.Printf("  %s\n", file)
			}
		}

		return nil
	},
}

var fileGrepCmd = &cobra.Command{
	Use:   "grep [pattern] [file]",
	Short: "Search for text in files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("search pattern and file path are required")
		}

		pattern := args[0]
		filePath := args[1]

		if dryRun {
			color.Yellow("Would search for '%s' in file '%s'", pattern, filePath)
			return nil
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", filePath)
		}

		// Use grep command
		cmdExec := exec.Command("grep", "-n", pattern, filePath)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			// grep returns exit code 1 when no matches found
			if strings.Contains(err.Error(), "exit status 1") {
				color.Yellow("No matches found for pattern: %s", pattern)
				return nil
			}
			return fmt.Errorf("failed to search file: %w", err)
		}

		return nil
	},
}

var fileBackupCmd = &cobra.Command{
	Use:   "backup [file]",
	Short: "Create backup of file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("file path is required")
		}

		filePath := args[0]

		if dryRun {
			color.Yellow("Would create backup of file: %s", filePath)
			return nil
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", filePath)
		}

		// Create backup filename
		backupPath := filePath + ".backup"
		if len(args) > 1 {
			backupPath = args[1]
		}

		// Copy file
		cmdExec := exec.Command("cp", filePath, backupPath)
		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}

		color.Green("Created backup: %s", backupPath)
		return nil
	},
}

var fileDiffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "Show differences between files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("two file paths are required")
		}

		file1 := args[0]
		file2 := args[1]

		if dryRun {
			color.Yellow("Would show diff between '%s' and '%s'", file1, file2)
			return nil
		}

		// Check if files exist
		if _, err := os.Stat(file1); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", file1)
		}
		if _, err := os.Stat(file2); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", file2)
		}

		// Use diff command
		cmdExec := exec.Command("diff", file1, file2)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			// diff returns exit code 1 when files are different
			if strings.Contains(err.Error(), "exit status 1") {
				// This is normal for different files
				return nil
			}
			return fmt.Errorf("failed to compare files: %w", err)
		}

		color.Green("Files are identical")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(fileOpenCmd)
	fileCmd.AddCommand(fileFindCmd)
	fileCmd.AddCommand(fileGrepCmd)
	fileCmd.AddCommand(fileBackupCmd)
	fileCmd.AddCommand(fileDiffCmd)
}
