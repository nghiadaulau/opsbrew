package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"opsbrew/internal/config"
	"opsbrew/internal/kubernetes"
)

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes operations and shortcuts",
	Long: `Kubernetes operations and shortcuts for common kubectl workflows.

Available commands:
  kctx     - Switch kubectl context with fuzzy finder
  kns      - Switch kubectl namespace with fuzzy finder
  klogs    - Get pod logs with fuzzy finder
  kpods    - List pods with fuzzy finder
  ksvc     - List services
  kingress - List ingress resources
  kexec    - Execute command in pod with fuzzy finder`,
}

var kctxCmd = &cobra.Command{
	Use:   "kctx [context]",
	Short: "Switch kubectl context with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var targetContext string

		if len(args) > 0 {
			targetContext = args[0]
			// Check if it's an alias
			if alias, exists := cfg.Kubernetes.ContextAliases[targetContext]; exists {
				targetContext = alias
			}
		} else {
			// Use fuzzy finder to select context
			contexts, err := kubernetes.GetContexts()
			if err != nil {
				return fmt.Errorf("failed to get contexts: %w", err)
			}

			selected, err := kubernetes.SelectContext(contexts)
			if err != nil {
				return fmt.Errorf("failed to select context: %w", err)
			}
			targetContext = selected
		}

		if dryRun {
			color.Yellow("Would run: kubectl config use-context %s", targetContext)
			return nil
		}

		// Switch context
		cmdExec := exec.Command("kubectl", "config", "use-context", targetContext)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to switch context: %w", err)
		}

		color.Green("Switched to context: %s", targetContext)
		return nil
	},
}

var knsCmd = &cobra.Command{
	Use:   "kns [namespace]",
	Short: "Switch kubectl namespace with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetRepoConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var targetNamespace string

		if len(args) > 0 {
			targetNamespace = args[0]
			// Check if it's an alias
			if alias, exists := cfg.Kubernetes.NamespaceAliases[targetNamespace]; exists {
				targetNamespace = alias
			}
		} else {
			// Use fuzzy finder to select namespace
			namespaces, err := kubernetes.GetNamespaces()
			if err != nil {
				return fmt.Errorf("failed to get namespaces: %w", err)
			}

			selected, err := kubernetes.SelectNamespace(namespaces)
			if err != nil {
				return fmt.Errorf("failed to select namespace: %w", err)
			}
			targetNamespace = selected
		}

		if dryRun {
			color.Yellow("Would run: kubectl config set-context --current --namespace=%s", targetNamespace)
			return nil
		}

		// Switch namespace
		cmdExec := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+targetNamespace)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to switch namespace: %w", err)
		}

		color.Green("Switched to namespace: %s", targetNamespace)
		return nil
	},
}

var klogsCmd = &cobra.Command{
	Use:   "klogs [pod]",
	Short: "Get pod logs with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetPod string

		if len(args) > 0 {
			targetPod = args[0]
		} else {
			// Use fuzzy finder to select pod
			pods, err := kubernetes.GetPods()
			if err != nil {
				return fmt.Errorf("failed to get pods: %w", err)
			}

			selected, err := kubernetes.SelectPod(pods)
			if err != nil {
				return fmt.Errorf("failed to select pod: %w", err)
			}
			targetPod = selected
		}

		// Get additional flags
		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetInt("tail")

		if dryRun {
			cmdStr := fmt.Sprintf("kubectl logs %s", targetPod)
			if follow {
				cmdStr += " -f"
			}
			if tail > 0 {
				cmdStr += fmt.Sprintf(" --tail=%d", tail)
			}
			color.Yellow("Would run: %s", cmdStr)
			return nil
		}

		// Build kubectl logs command
		kubectlArgs := []string{"logs", targetPod}
		if follow {
			kubectlArgs = append(kubectlArgs, "-f")
		}
		if tail > 0 {
			kubectlArgs = append(kubectlArgs, fmt.Sprintf("--tail=%d", tail))
		}

		cmdExec := exec.Command("kubectl", kubectlArgs...)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to get logs: %w", err)
		}

		return nil
	},
}

var kpodsCmd = &cobra.Command{
	Use:   "kpods",
	Short: "List pods with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		pods, err := kubernetes.GetPods()
		if err != nil {
			return fmt.Errorf("failed to get pods: %w", err)
		}

		kubernetes.DisplayPods(pods)
		return nil
	},
}

var ksvcCmd = &cobra.Command{
	Use:   "ksvc",
	Short: "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			color.Yellow("Would run: kubectl get services")
			return nil
		}

		cmdExec := exec.Command("kubectl", "get", "services")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to get services: %w", err)
		}

		return nil
	},
}

var kingressCmd = &cobra.Command{
	Use:   "kingress",
	Short: "List ingress resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			color.Yellow("Would run: kubectl get ingress")
			return nil
		}

		cmdExec := exec.Command("kubectl", "get", "ingress")
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to get ingress: %w", err)
		}

		return nil
	},
}

var kexecCmd = &cobra.Command{
	Use:   "kexec [pod] [command]",
	Short: "Execute command in pod with fuzzy finder",
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetPod string
		var command string

		if len(args) > 0 {
			targetPod = args[0]
		} else {
			// Use fuzzy finder to select pod
			pods, err := kubernetes.GetPods()
			if err != nil {
				return fmt.Errorf("failed to get pods: %w", err)
			}

			selected, err := kubernetes.SelectPod(pods)
			if err != nil {
				return fmt.Errorf("failed to select pod: %w", err)
			}
			targetPod = selected
		}

		if len(args) > 1 {
			command = args[1]
		} else {
			command = "/bin/bash"
		}

		if dryRun {
			color.Yellow("Would run: kubectl exec -it %s -- %s", targetPod, command)
			return nil
		}

		// Execute command in pod
		kubectlArgs := []string{"exec", "-it", targetPod, "--"}
		kubectlArgs = append(kubectlArgs, strings.Split(command, " ")...)

		cmdExec := exec.Command("kubectl", kubectlArgs...)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(k8sCmd)
	k8sCmd.AddCommand(kctxCmd)
	k8sCmd.AddCommand(knsCmd)
	k8sCmd.AddCommand(klogsCmd)
	k8sCmd.AddCommand(kpodsCmd)
	k8sCmd.AddCommand(ksvcCmd)
	k8sCmd.AddCommand(kingressCmd)
	k8sCmd.AddCommand(kexecCmd)

	// Add flags for klogs
	klogsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	klogsCmd.Flags().IntP("tail", "t", 0, "Number of lines to show from the end of the logs")
}
