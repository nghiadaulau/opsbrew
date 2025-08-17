package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"opsbrew/internal/config"
	"opsbrew/internal/git"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git operations and shortcuts",
	Long: `Git operations and shortcuts for common workflows.

Available commands:
  status    - Show git status with enhanced formatting
  sync      - Pull with rebase (git pull --rebase)
  checkout  - Checkout branch with fuzzy finder
  branch    - List branches with fuzzy finder
  fetch     - Fetch all remotes
  pull      - Pull from current branch
  push      - Push to current branch`,
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status with enhanced formatting",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if dryRun {
			color.Yellow("Would run: git status")
			return nil
		}

		// Run git status
		output, err := exec.Command("git", "status", "--porcelain").Output()
		if err != nil {
			return fmt.Errorf("failed to get git status: %w", err)
		}

		// Parse and display status
		status := git.ParseStatus(string(output))
		git.DisplayStatus(status, cfg.UI.Colors)

		return nil
	},
}

var gitSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pull with rebase (git pull --rebase)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if dryRun {
			color.Yellow("Would run: git pull --rebase")
			return nil
		}

		// Check if we need confirmation
		if !confirm && !cfg.UI.Confirm {
			fmt.Print("Pull with rebase? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				color.Yellow("Operation cancelled")
				return nil
			}
		}

		// Get current branch
		branchOutput, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(branchOutput))

		color.Green("Syncing branch: %s", currentBranch)

		// Run git pull --rebase
		cmdExec := exec.Command("git", "pull", "--rebase")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to sync: %w", err)
		}

		color.Green("Sync completed successfully")
		return nil
	},
}

var gitCheckoutCmd = &cobra.Command{
	Use:   "checkout [branch]",
	Short: "Checkout branch with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var targetBranch string

		if len(args) > 0 {
			targetBranch = args[0]
		} else {
			// Use fuzzy finder to select branch
			branches, err := git.GetBranches()
			if err != nil {
				return fmt.Errorf("failed to get branches: %w", err)
			}

			selected, err := git.SelectBranch(branches)
			if err != nil {
				return fmt.Errorf("failed to select branch: %w", err)
			}
			targetBranch = selected
		}

		if dryRun {
			color.Yellow("Would run: git checkout %s", targetBranch)
			return nil
		}

		// Check if branch exists locally
		_, err = exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+targetBranch).Output()
		if err != nil {
			// Branch doesn't exist locally, try to checkout from remote
			color.Yellow("Branch %s not found locally, checking out from remote...", targetBranch)
			cmdExec := exec.Command("git", "checkout", "-b", targetBranch, "origin/"+targetBranch)
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			if err := cmdExec.Run(); err != nil {
				return fmt.Errorf("failed to checkout branch %s: %w", targetBranch, err)
			}
		} else {
			// Branch exists locally
			cmdExec := exec.Command("git", "checkout", targetBranch)
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			if err := cmdExec.Run(); err != nil {
				return fmt.Errorf("failed to checkout branch %s: %w", targetBranch, err)
			}
		}

		color.Green("Switched to branch: %s", targetBranch)
		return nil
	},
}

var gitBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List branches with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		branches, err := git.GetBranches()
		if err != nil {
			return fmt.Errorf("failed to get branches: %w", err)
		}

		git.DisplayBranches(branches)
		return nil
	},
}

var gitFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all remotes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			color.Yellow("Would run: git fetch --all")
			return nil
		}

		color.Green("Fetching all remotes...")
		cmdExec := exec.Command("git", "fetch", "--all")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to fetch: %w", err)
		}

		color.Green("Fetch completed successfully")
		return nil
	},
}

var gitPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull from current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			color.Yellow("Would run: git pull")
			return nil
		}

		color.Green("Pulling from current branch...")
		cmdExec := exec.Command("git", "pull")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}

		color.Green("Pull completed successfully")
		return nil
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push to current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			color.Yellow("Would run: git push")
			return nil
		}

		color.Green("Pushing to current branch...")
		cmdExec := exec.Command("git", "push")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}

		color.Green("Push completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitSyncCmd)
	gitCmd.AddCommand(gitCheckoutCmd)
	gitCmd.AddCommand(gitBranchCmd)
	gitCmd.AddCommand(gitFetchCmd)
	gitCmd.AddCommand(gitPullCmd)
	gitCmd.AddCommand(gitPushCmd)
}
