package kubernetes

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

// Context represents a kubectl context
type Context struct {
	Name    string
	Current bool
}

// Namespace represents a kubernetes namespace
type Namespace struct {
	Name    string
	Current bool
	Status  string
}

// Pod represents a kubernetes pod
type Pod struct {
	Name      string
	Ready     string
	Status    string
	Restarts  string
	Age       string
	Namespace string
}

// GetContexts returns all available kubectl contexts
func GetContexts() ([]Context, error) {
	output, err := exec.Command("kubectl", "config", "get-contexts", "--no-headers", "-o", "name").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}

	currentOutput, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}
	currentContext := strings.TrimSpace(string(currentOutput))

	var contexts []Context
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		contexts = append(contexts, Context{
			Name:    strings.TrimSpace(line),
			Current: strings.TrimSpace(line) == currentContext,
		})
	}

	return contexts, nil
}

// SelectContext uses fuzzy finder to select a context
func SelectContext(contexts []Context) (string, error) {
	idx, err := fuzzyfinder.Find(
		contexts,
		func(i int) string {
			ctx := contexts[i]
			if ctx.Current {
				return fmt.Sprintf("  * %s", ctx.Name)
			}
			return fmt.Sprintf("    %s", ctx.Name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			ctx := contexts[i]
			return fmt.Sprintf("Context: %s\nCurrent: %t", ctx.Name, ctx.Current)
		}),
	)
	if err != nil {
		return "", err
	}

	return contexts[idx].Name, nil
}

// GetNamespaces returns all available namespaces
func GetNamespaces() ([]Namespace, error) {
	output, err := exec.Command("kubectl", "get", "namespaces", "--no-headers", "-o", "custom-columns=NAME:.metadata.name,STATUS:.status.phase").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	currentOutput, err := exec.Command("kubectl", "config", "view", "--minify", "-o", "jsonpath={..namespace}").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current namespace: %w", err)
	}
	currentNamespace := strings.TrimSpace(string(currentOutput))
	if currentNamespace == "" {
		currentNamespace = "default"
	}

	var namespaces []Namespace
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			namespaces = append(namespaces, Namespace{
				Name:    parts[0],
				Status:  parts[1],
				Current: parts[0] == currentNamespace,
			})
		}
	}

	return namespaces, nil
}

// SelectNamespace uses fuzzy finder to select a namespace
func SelectNamespace(namespaces []Namespace) (string, error) {
	idx, err := fuzzyfinder.Find(
		namespaces,
		func(i int) string {
			ns := namespaces[i]
			if ns.Current {
				return fmt.Sprintf("  * %s (%s)", ns.Name, ns.Status)
			}
			return fmt.Sprintf("    %s (%s)", ns.Name, ns.Status)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			ns := namespaces[i]
			return fmt.Sprintf("Namespace: %s\nStatus: %s\nCurrent: %t", ns.Name, ns.Status, ns.Current)
		}),
	)
	if err != nil {
		return "", err
	}

	return namespaces[idx].Name, nil
}

// GetPods returns all pods in the current namespace
func GetPods() ([]Pod, error) {
	output, err := exec.Command("kubectl", "get", "pods", "--no-headers", "-o", "custom-columns=NAME:.metadata.name,READY:.status.containerStatuses[*].ready,STATUS:.status.phase,RESTARTS:.status.containerStatuses[*].restartCount,AGE:.metadata.creationTimestamp").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var pods []Pod
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 5 {
			pods = append(pods, Pod{
				Name:     parts[0],
				Ready:    parts[1],
				Status:   parts[2],
				Restarts: parts[3],
				Age:      parts[4],
			})
		}
	}

	return pods, nil
}

// SelectPod uses fuzzy finder to select a pod
func SelectPod(pods []Pod) (string, error) {
	idx, err := fuzzyfinder.Find(
		pods,
		func(i int) string {
			pod := pods[i]
			return fmt.Sprintf("%s (%s) - %s", pod.Name, pod.Status, pod.Ready)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			pod := pods[i]
			return fmt.Sprintf("Pod: %s\nStatus: %s\nReady: %s\nRestarts: %s\nAge: %s", 
				pod.Name, pod.Status, pod.Ready, pod.Restarts, pod.Age)
		}),
	)
	if err != nil {
		return "", err
	}

	return pods[idx].Name, nil
}

// DisplayPods displays pods with formatting
func DisplayPods(pods []Pod) {
	fmt.Println("=== Pods ===")
	for _, pod := range pods {
		statusColor := getStatusColor(pod.Status)
		statusColor.Printf("  %s (%s) - %s\n", pod.Name, pod.Status, pod.Ready)
	}
}

// getStatusColor returns the appropriate color for pod status
func getStatusColor(status string) *color.Color {
	switch strings.ToLower(status) {
	case "running":
		return color.New(color.FgGreen)
	case "pending":
		return color.New(color.FgYellow)
	case "failed", "error":
		return color.New(color.FgRed)
	case "succeeded":
		return color.New(color.FgBlue)
	default:
		return color.New(color.FgWhite)
	}
}
