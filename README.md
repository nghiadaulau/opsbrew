# opsbrew

A powerful CLI tool designed to simplify and shorten repetitive DevOps terminal commands. opsbrew provides shortcuts for common Git and kubectl operations, with features like fuzzy finders, command recipes, and project templates.

**Source**: [github.com/nghiadaulau/opsbrew](https://github.com/nghiadaulau/opsbrew)

## Features

- **Git Operations**: Enhanced Git commands with fuzzy finder for branches
- **Kubernetes Management**: kubectl shortcuts with context/namespace switching, HPA management, and scaling
- **File Operations**: Common file operations like backup, diff, find, and grep
- **Command Recipes**: Save and run command macros for daily workflows
- **Project Templates**: Bootstrap common project structures
- **Safe Defaults**: Built-in `--dry-run` and `--confirm` flags
- **Configuration**: YAML-based configuration (global + per-repo)
- **Shell Completions**: Full shell completion support

## Installation

### Quick Install (Recommended)

```bash
# Install latest version
go install github.com/nghiadaulau/opsbrew@latest

# Verify installation
opsbrew --version
```

### Pre-built Binaries

```bash
# Linux/macOS
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Linux_x86_64.tar.gz
tar -xzf opsbrew_Linux_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Windows
curl -L -o opsbrew_Windows_x86_64.zip https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Windows_x86_64.zip
# Extract and add to PATH
```

### From Source

```bash
git clone https://github.com/nghiadaulau/opsbrew.git
cd opsbrew
go build -o opsbrew .
sudo mv opsbrew /usr/local/bin/
```

📖 **Detailed Installation Guide**: See [INSTALL.md](INSTALL.md) for complete instructions including upgrade, uninstall, and troubleshooting.

## Releases

This project uses [GoReleaser](https://goreleaser.com/) for automated releases. When you push a tag starting with `v*`, it automatically:

- Builds binaries for Linux, macOS, and Windows
- Creates GitHub releases with all artifacts
- Generates checksums for verification

### Creating a Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

See [RELEASE.md](RELEASE.md) for detailed release instructions.

## Documentation

This project uses [mdBook](https://rust-lang.github.io/mdBook/) to generate beautiful documentation sites.

- **Live Documentation**: [https://nghiadaulau.github.io/opsbrew](https://nghiadaulau.github.io/opsbrew)
- **Documentation Setup**: See [DOCUMENTATION.md](DOCUMENTATION.md) for details

The documentation site includes:
- Complete command reference
- Installation guides
- Release information
- Interactive search
- Dark/light theme toggle

### Using Go Install

```bash
go install github.com/nghiadaulau/opsbrew@latest
```

## Quick Start

1. **Initialize configuration**:
   ```bash
   opsbrew --help
   # This will create ~/.opsbrew.yaml with default settings
   ```

2. **Git operations**:
   ```bash
   opsbrew git status          # Enhanced git status
   opsbrew git sync           # git pull --rebase
   opsbrew git checkout       # Fuzzy finder for branches
   ```

3. **Kubernetes operations**:
   ```bash
   opsbrew k8s kctx           # Switch kubectl context
   opsbrew k8s kns            # Switch namespace
   opsbrew k8s klogs          # Get pod logs with fuzzy finder
   opsbrew k8s khpa set-min my-hpa 2  # Set HPA min replicas
   opsbrew k8s kscale deployment my-app 5  # Scale deployment
   ```

4. **Command recipes**:
   ```bash
   opsbrew brew save daily-sync  # Save a new recipe
   opsbrew brew run daily-sync   # Execute saved recipe
   opsbrew brew list            # List all recipes
   ```

5. **File operations**:
   ```bash
   opsbrew file backup config.yaml  # Create backup
   opsbrew file find "*.yaml"       # Find files by pattern
   opsbrew file diff file1.yaml file2.yaml  # Compare files
   ```

6. **Project templates**:
   ```bash
   opsbrew init github-actions my-app    # Create GitHub Actions workflow
   opsbrew init k8s-deployment my-app    # Create K8s Deployment
   opsbrew init k8s-service my-app       # Create K8s Service
   ```

## Configuration

opsbrew uses YAML configuration files. The global config is located at `~/.opsbrew.yaml`, and you can have per-repository configs in `.opsbrew.yaml`.

### Example Configuration

```yaml
# Git configuration
git:
  default_branch: "main"
  auto_fetch: true
  aliases:
    st: "status"
    co: "checkout"
    sync: "pull --rebase"

# Kubernetes configuration
kubernetes:
  default_context: "production"
  context_aliases:
    prod: "production-cluster"
    dev: "development-cluster"

# Command recipes
brew:
  recipes:
    daily-sync:
      description: "Daily development workflow"
      commands:
        - "git fetch --all"
        - "git pull origin main"
        - "git checkout -b feature/$(date +%Y%m%d)"
      tags: ["daily", "git"]

# UI settings
ui:
  colors: true
  verbose: false
  confirm: false
  dry_run: false
```

## Commands

### Git Commands

- `opsbrew git status` - Enhanced git status with colors
- `opsbrew git sync` - Pull with rebase
- `opsbrew git checkout [branch]` - Checkout branch with fuzzy finder
- `opsbrew git branch` - List branches with fuzzy finder
- `opsbrew git fetch` - Fetch all remotes
- `opsbrew git pull` - Pull from current branch
- `opsbrew git push` - Push to current branch

### Kubernetes Commands

- `opsbrew k8s kctx [context]` - Switch kubectl context with fuzzy finder
- `opsbrew k8s kns [namespace]` - Switch namespace with fuzzy finder
- `opsbrew k8s klogs [pod]` - Get pod logs with fuzzy finder
- `opsbrew k8s kpods` - List pods with fuzzy finder
- `opsbrew k8s ksvc` - List services
- `opsbrew k8s kingress` - List ingress resources
- `opsbrew k8s kexec [pod] [command]` - Execute command in pod
- `opsbrew k8s khpa list` - List all HPAs
- `opsbrew k8s khpa get [name]` - Get HPA details
- `opsbrew k8s khpa set-min [name] [value]` - Set minimum replicas
- `opsbrew k8s khpa set-max [name] [value]` - Set maximum replicas
- `opsbrew k8s khpa set-target [name] [value]` - Set target CPU percentage
- `opsbrew k8s kscale [type] [name] [replicas]` - Scale deployment/replicaset/statefulset

### File Commands

- `opsbrew file open [file]` - Open file with default editor
- `opsbrew file find [pattern] [dir]` - Find files by name or pattern
- `opsbrew file grep [pattern] [file]` - Search for text in files
- `opsbrew file backup [file] [backup-path]` - Create backup of file
- `opsbrew file diff [file1] [file2]` - Show differences between files

### Brew Commands (Command Recipes)

- `opsbrew brew save [name]` - Save a new recipe
- `opsbrew brew list` - List all saved recipes
- `opsbrew brew run [name]` - Execute a saved recipe
- `opsbrew brew delete [name]` - Delete a recipe
- `opsbrew brew edit [name]` - Edit a recipe

### Init Commands (Project Templates)

- `opsbrew init github-actions [name]` - Create GitHub Actions workflow
- `opsbrew init k8s-deployment [name]` - Create Kubernetes Deployment manifest
- `opsbrew init k8s-service [name]` - Create Kubernetes Service manifest
- `opsbrew init k8s-pod [name]` - Create Kubernetes Pod manifest
- `opsbrew init k8s-configmap [name]` - Create Kubernetes ConfigMap manifest
- `opsbrew init dockerfile [name]` - Create multi-stage Dockerfile
- `opsbrew init list` - List available templates

### Global Flags

- `--config` - Specify config file path
- `--verbose, -v` - Enable verbose output
- `--dry-run` - Show what would be done without executing
- `--confirm` - Skip confirmation prompts

## Shell Completions

Generate shell completions:

```bash
# Bash
opsbrew completion bash > ~/.bash_completion.d/opsbrew

# Zsh
opsbrew completion zsh > "${fpath[1]}/_opsbrew"

# Fish
opsbrew completion fish > ~/.config/fish/completions/opsbrew.fish

# PowerShell
opsbrew completion powershell > opsbrew.ps1
```

## Examples

### Daily Development Workflow

```bash
# Start the day
opsbrew brew run daily-sync

# Check git status
opsbrew git status

# Switch to a feature branch
opsbrew git checkout

# Switch to development context
opsbrew k8s kctx dev

# Check pod logs
opsbrew k8s klogs

# Manage HPA settings
opsbrew k8s khpa set-min my-app 2
opsbrew k8s khpa set-max my-app 10

# Scale deployment
opsbrew k8s kscale deployment my-app 5

# File operations
opsbrew file backup config.yaml
opsbrew file find "*.yaml"
```

### Deployment Workflow

```bash
# Pre-deployment checks
opsbrew brew run deploy-check

# Switch to production context
opsbrew k8s kctx prod

# Check application status
opsbrew k8s kpods
opsbrew k8s ksvc
```

### Project Setup

```bash
# Create GitHub Actions workflow
opsbrew init github-actions my-api

# Create Kubernetes manifests
opsbrew init k8s-deployment my-api
opsbrew init k8s-service my-api
opsbrew init k8s-configmap my-api

# Create Dockerfile
opsbrew init dockerfile my-api
```

## Development

### Prerequisites

- Go 1.24 or later
- Git
- kubectl (for Kubernetes commands)

### Building from Source

```bash
git clone https://github.com/nghiadaulau/opsbrew.git
cd opsbrew
go mod download
go build -o opsbrew .
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [ktr0731/go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) - Fuzzy finder
- [fatih/color](https://github.com/fatih/color) - Color output
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
