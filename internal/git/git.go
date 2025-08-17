package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

// FileStatus represents the status of a git file
type FileStatus struct {
	Path   string
	Status string
	Type   string
}

// GitStatus represents the overall git status
type GitStatus struct {
	Modified   []FileStatus
	Staged     []FileStatus
	Untracked  []FileStatus
	Deleted    []FileStatus
	Renamed    []FileStatus
	Conflicted []FileStatus
}

// Branch represents a git branch
type Branch struct {
	Name   string
	Current bool
	Remote bool
}

// ParseStatus parses git status output
func ParseStatus(output string) *GitStatus {
	status := &GitStatus{}
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse porcelain format: XY PATH
		if len(line) < 3 {
			continue
		}

		xy := line[:2]
		path := line[3:]

		fileStatus := FileStatus{
			Path:   path,
			Status: xy,
		}

		switch {
		case strings.HasPrefix(xy, "M"):
			if xy[1] == 'M' {
				status.Modified = append(status.Modified, fileStatus)
			} else {
				status.Staged = append(status.Staged, fileStatus)
			}
		case strings.HasPrefix(xy, "A"):
			status.Staged = append(status.Staged, fileStatus)
		case strings.HasPrefix(xy, "D"):
			status.Deleted = append(status.Deleted, fileStatus)
		case strings.HasPrefix(xy, "R"):
			status.Renamed = append(status.Renamed, fileStatus)
		case strings.HasPrefix(xy, "??"):
			status.Untracked = append(status.Untracked, fileStatus)
		case strings.HasPrefix(xy, "UU"), strings.HasPrefix(xy, "AA"), strings.HasPrefix(xy, "DD"):
			status.Conflicted = append(status.Conflicted, fileStatus)
		}
	}

	return status
}

// DisplayStatus displays git status with colors
func DisplayStatus(status *GitStatus, useColors bool) {
	if useColors {
		color.Green("=== Git Status ===")
	} else {
		fmt.Println("=== Git Status ===")
	}

	// Show current branch
	branch, err := getCurrentBranch()
	if err == nil {
		if useColors {
			color.Cyan("On branch: %s", branch)
		} else {
			fmt.Printf("On branch: %s\n", branch)
		}
	}

	fmt.Println()

	// Display staged changes
	if len(status.Staged) > 0 {
		if useColors {
			color.Green("Changes to be committed:")
		} else {
			fmt.Println("Changes to be committed:")
		}
		for _, file := range status.Staged {
			if useColors {
				color.Green("  %s", file.Path)
			} else {
				fmt.Printf("  %s\n", file.Path)
			}
		}
		fmt.Println()
	}

	// Display modified files
	if len(status.Modified) > 0 {
		if useColors {
			color.Yellow("Changes not staged for commit:")
		} else {
			fmt.Println("Changes not staged for commit:")
		}
		for _, file := range status.Modified {
			if useColors {
				color.Yellow("  %s", file.Path)
			} else {
				fmt.Printf("  %s\n", file.Path)
			}
		}
		fmt.Println()
	}

	// Display untracked files
	if len(status.Untracked) > 0 {
		if useColors {
			color.Red("Untracked files:")
		} else {
			fmt.Println("Untracked files:")
		}
		for _, file := range status.Untracked {
			if useColors {
				color.Red("  %s", file.Path)
			} else {
				fmt.Printf("  %s\n", file.Path)
			}
		}
		fmt.Println()
	}

	// Display conflicted files
	if len(status.Conflicted) > 0 {
		if useColors {
			color.Red("Unmerged paths:")
		} else {
			fmt.Println("Unmerged paths:")
		}
		for _, file := range status.Conflicted {
			if useColors {
				color.Red("  %s", file.Path)
			} else {
				fmt.Printf("  %s\n", file.Path)
			}
		}
		fmt.Println()
	}

	// Summary
	totalChanges := len(status.Staged) + len(status.Modified) + len(status.Untracked) + len(status.Deleted) + len(status.Renamed) + len(status.Conflicted)
	if totalChanges == 0 {
		if useColors {
			color.Green("Working tree clean")
		} else {
			fmt.Println("Working tree clean")
		}
	}
}

// GetBranches returns all available branches
func GetBranches() ([]Branch, error) {
	// Get local branches
	localOutput, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get local branches: %w", err)
	}

	// Get current branch
	currentOutput, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	currentBranch := strings.TrimSpace(string(currentOutput))

	// Get remote branches
	remoteOutput, err := exec.Command("git", "branch", "-r", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	var branches []Branch

	// Add local branches
	localBranches := strings.Split(strings.TrimSpace(string(localOutput)), "\n")
	for _, branch := range localBranches {
		if branch == "" {
			continue
		}
		branches = append(branches, Branch{
			Name:    branch,
			Current: branch == currentBranch,
			Remote:  false,
		})
	}

	// Add remote branches
	remoteBranches := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")
	for _, branch := range remoteBranches {
		if branch == "" {
			continue
		}
		// Skip HEAD reference
		if strings.Contains(branch, "HEAD") {
			continue
		}
		branches = append(branches, Branch{
			Name:    branch,
			Current: false,
			Remote:  true,
		})
	}

	return branches, nil
}

// SelectBranch uses fuzzy finder to select a branch
func SelectBranch(branches []Branch) (string, error) {
	idx, err := fuzzyfinder.Find(
		branches,
		func(i int) string {
			branch := branches[i]
			if branch.Current {
				return fmt.Sprintf("  * %s", branch.Name)
			}
			if branch.Remote {
				return fmt.Sprintf("    %s (remote)", branch.Name)
			}
			return fmt.Sprintf("    %s", branch.Name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			branch := branches[i]
			return fmt.Sprintf("Branch: %s\nType: %s", branch.Name, branchType(branch))
		}),
	)
	if err != nil {
		return "", err
	}

	return branches[idx].Name, nil
}

// DisplayBranches displays branches with formatting
func DisplayBranches(branches []Branch) {
	fmt.Println("=== Branches ===")
	for _, branch := range branches {
		if branch.Current {
			color.Cyan("  * %s", branch.Name)
		} else if branch.Remote {
			fmt.Printf("    %s (remote)\n", branch.Name)
		} else {
			fmt.Printf("    %s\n", branch.Name)
		}
	}
}

// getCurrentBranch returns the current branch name
func getCurrentBranch() (string, error) {
	output, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// branchType returns a human-readable branch type
func branchType(branch Branch) string {
	if branch.Current {
		return "Current"
	}
	if branch.Remote {
		return "Remote"
	}
	return "Local"
}
