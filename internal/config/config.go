package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the opsbrew configuration structure
type Config struct {
	Git struct {
		DefaultBranch string            `yaml:"default_branch"`
		Aliases       map[string]string `yaml:"aliases"`
		AutoFetch     bool              `yaml:"auto_fetch"`
	} `yaml:"git"`

	Kubernetes struct {
		DefaultContext  string            `yaml:"default_context"`
		DefaultNamespace string            `yaml:"default_namespace"`
		ContextAliases  map[string]string `yaml:"context_aliases"`
		NamespaceAliases map[string]string `yaml:"namespace_aliases"`
	} `yaml:"kubernetes"`

	Brew struct {
		Recipes map[string]Recipe `yaml:"recipes"`
	} `yaml:"brew"`

	Templates struct {
		Path string `yaml:"path"`
	} `yaml:"templates"`

	UI struct {
		Colors    bool `yaml:"colors"`
		Verbose   bool `yaml:"verbose"`
		Confirm   bool `yaml:"confirm"`
		DryRun    bool `yaml:"dry_run"`
	} `yaml:"ui"`
}

// Recipe represents a saved command recipe
type Recipe struct {
	Description string   `yaml:"description"`
	Commands    []string `yaml:"commands"`
	Tags        []string `yaml:"tags"`
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	var cfg Config

	// Read config from viper
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(cfg *Config) error {
	// Marshal config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Get config file path
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(home, ".opsbrew.yaml")
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig() error {
	cfg := &Config{}

	// Set default Git configuration
	cfg.Git.DefaultBranch = "main"
	cfg.Git.AutoFetch = true
	cfg.Git.Aliases = map[string]string{
		"st":   "status",
		"co":   "checkout",
		"br":   "branch",
		"cm":   "commit",
		"pl":   "pull",
		"ps":   "push",
		"lg":   "log --oneline --graph",
		"sync": "pull --rebase",
	}

	// Set default Kubernetes configuration
	cfg.Kubernetes.DefaultContext = ""
	cfg.Kubernetes.DefaultNamespace = "default"
	cfg.Kubernetes.ContextAliases = map[string]string{
		"prod": "production-cluster",
		"dev":  "development-cluster",
		"stg":  "staging-cluster",
	}
	cfg.Kubernetes.NamespaceAliases = map[string]string{
		"app": "application",
		"db":  "database",
		"mon": "monitoring",
	}

	// Set default Brew configuration
	cfg.Brew.Recipes = map[string]Recipe{
		"daily-sync": {
			Description: "Daily development workflow",
			Commands: []string{
				"git fetch --all",
				"git pull origin main",
				"git checkout -b feature/$(date +%Y%m%d)",
			},
			Tags: []string{"daily", "git"},
		},
		"deploy-check": {
			Description: "Pre-deployment checks",
			Commands: []string{
				"kubectl get pods",
				"kubectl get services",
				"kubectl get ingress",
			},
			Tags: []string{"deploy", "k8s"},
		},
	}

	// Set default Templates configuration
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cfg.Templates.Path = filepath.Join(home, ".opsbrew", "templates")

	// Set default UI configuration
	cfg.UI.Colors = true
	cfg.UI.Verbose = false
	cfg.UI.Confirm = false
	cfg.UI.DryRun = false

	return SaveConfig(cfg)
}

// GetRepoConfig loads repository-specific configuration
func GetRepoConfig() (*Config, error) {
	// Check for .opsbrew.yaml in current directory
	if _, err := os.Stat(".opsbrew.yaml"); err == nil {
		viper.SetConfigFile(".opsbrew.yaml")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read repo config: %w", err)
		}
		return LoadConfig()
	}

	// Fall back to global config
	return LoadConfig()
}
