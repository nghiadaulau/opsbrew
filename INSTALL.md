# Installation Guide

## Quick Install

### Using Go Install (Recommended)

```bash
# Install latest version
go install github.com/nghiadaulau/opsbrew@latest

# Add Go bin to PATH (if not already added)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Or for zsh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Verify installation
opsbrew --version
```

### Using Pre-built Binaries

#### Linux (x86_64)
```bash
# Download latest release
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Linux_x86_64.tar.gz

# Extract and install
tar -xzf opsbrew_Linux_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Verify installation
opsbrew --version
```

#### macOS (x86_64)
```bash
# Download latest release
curl -L -o opsbrew_Darwin_x86_64.tar.gz https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Darwin_x86_64.tar.gz

# Extract and install
tar -xzf opsbrew_Darwin_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Verify installation
opsbrew --version
```

#### macOS (Apple Silicon)
```bash
# Download latest release
curl -L -o opsbrew_Darwin_arm64.tar.gz https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Darwin_arm64.tar.gz

# Extract and install
tar -xzf opsbrew_Darwin_arm64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Verify installation
opsbrew --version
```

#### Windows
```bash
# Download latest release
curl -L -o opsbrew_Windows_x86_64.zip https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Windows_x86_64.zip

# Extract and add to PATH
# Or use PowerShell:
Expand-Archive -Path opsbrew_Windows_x86_64.zip -DestinationPath C:\opsbrew
# Add C:\opsbrew to your PATH environment variable
```

### From Source

```bash
# Clone repository
git clone https://github.com/nghiadaulau/opsbrew.git
cd opsbrew

# Build from source
go build -o opsbrew .

# Install to system
sudo mv opsbrew /usr/local/bin/

# Verify installation
opsbrew --version
```

## Upgrade

### Using Go Install (Recommended)

```bash
# Upgrade to latest version
go install github.com/nghiadaulau/opsbrew@latest

# Verify new version
opsbrew --version
```

### Using Pre-built Binaries

#### Linux/macOS
```bash
# Download latest release
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Linux_x86_64.tar.gz
# or for macOS: opsbrew_Darwin_x86_64.tar.gz

# Stop any running opsbrew processes
pkill opsbrew

# Backup current version (optional)
sudo cp /usr/local/bin/opsbrew /usr/local/bin/opsbrew.backup

# Extract and replace
tar -xzf opsbrew_Linux_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Verify new version
opsbrew --version
```

#### Windows
```bash
# Download latest release
curl -L -o opsbrew_Windows_x86_64.zip https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Windows_x86_64.zip

# Stop any running opsbrew processes
taskkill /F /IM opsbrew.exe

# Extract and replace
Expand-Archive -Path opsbrew_Windows_x86_64.zip -DestinationPath C:\opsbrew -Force
```

### From Source

```bash
# Pull latest changes
cd opsbrew
git pull origin main

# Rebuild
go build -o opsbrew .

# Replace existing installation
sudo mv opsbrew /usr/local/bin/

# Verify new version
opsbrew --version
```

## Uninstall

### Remove Binary
```bash
# Remove the binary
sudo rm /usr/local/bin/opsbrew

# Verify removal
which opsbrew
```

### Remove Configuration (Optional)
```bash
# Remove global config
rm ~/.opsbrew.yaml

# Remove per-repo configs (if any)
find . -name "opsbrew.yaml" -delete
```

## Shell Completions

### Bash
```bash
# Generate bash completion
opsbrew completion bash > ~/.local/share/bash-completion/completions/opsbrew

# Or add to ~/.bashrc
echo 'source <(opsbrew completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh
```bash
# Generate zsh completion
opsbrew completion zsh > ~/.zsh/completions/_opsbrew

# Or add to ~/.zshrc
echo 'source <(opsbrew completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish
```bash
# Generate fish completion
opsbrew completion fish > ~/.config/fish/completions/opsbrew.fish
```

### PowerShell
```powershell
# Generate PowerShell completion
opsbrew completion powershell | Out-String | Invoke-Expression

# Add to profile for persistence
opsbrew completion powershell > $PROFILE
```

## Configuration

### Global Configuration
```bash
# Create global config
cat > ~/.opsbrew.yaml << 'EOF'
ui:
  confirm: false
  verbose: false

brew:
  recipes: {}

git:
  default_branch: main
  sync_strategy: rebase

kubernetes:
  default_namespace: default
  context_timeout: 30s
EOF
```

### Per-Repository Configuration
```bash
# Create repo-specific config
cat > opsbrew.yaml << 'EOF'
ui:
  confirm: true
  verbose: true

brew:
  recipes:
    deploy:
      description: "Deploy to production"
      commands:
        - "kubectl apply -f k8s/"
        - "kubectl rollout status deployment/app"

git:
  default_branch: develop
  sync_strategy: merge
EOF
```

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Fix permissions
sudo chmod +x /usr/local/bin/opsbrew
```

#### Command Not Found
```bash
# Check if binary exists
ls -la /usr/local/bin/opsbrew

# Check PATH
echo $PATH

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin
```

#### Version Issues
```bash
# Check current version
opsbrew --version

# Check available versions
git tag -l | grep "^v" | sort -V

# Install specific version
go install github.com/nghiadaulau/opsbrew@v1.0.0
```

#### Build Issues
```bash
# Clean and rebuild
go clean
go mod download
go build -o opsbrew .
```

### Verification

After installation, verify everything works:

```bash
# Check version
opsbrew --version

# Check help
opsbrew --help

# Test git commands
opsbrew git status

# Test k8s commands
opsbrew k8s kctx

# Test file commands
opsbrew file backup README.md

# Test brew commands
opsbrew brew list
```

## Package Managers (Future)

### Homebrew (macOS/Linux)
```bash
# When Homebrew tap is available
brew tap nghiadaulau/opsbrew
brew install opsbrew
```

### Chocolatey (Windows)
```powershell
# When Chocolatey package is available
choco install opsbrew
```

### Snap (Linux)
```bash
# When Snap package is available
sudo snap install opsbrew
```

## Development Installation

For developers who want to work on opsbrew:

```bash
# Clone repository
git clone https://github.com/nghiadaulau/opsbrew.git
cd opsbrew

# Install dependencies
go mod download

# Build
go build -o opsbrew .

# Run tests
go test ./...

# Install development version
go install .
```

## Security

### Verify Checksums
```bash
# Download checksums
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/checksums.txt

# Verify downloaded binary
sha256sum -c checksums.txt
```

### GPG Verification (Future)
When GPG signatures are available:
```bash
# Import GPG key
gpg --keyserver keyserver.ubuntu.com --recv-keys <KEY_ID>

# Verify signature
gpg --verify opsbrew_Linux_x86_64.tar.gz.sig
```
