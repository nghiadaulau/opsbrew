# Release Guide

This project uses [GoReleaser](https://goreleaser.com/) for automated releases.

## How to Create a Release

### 1. Create and Push a Tag

```bash
# Create a new tag
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### 2. Automated Release Process

When you push a tag starting with `v*` (e.g., `v1.0.0`, `v1.2.3`), the GitHub Actions workflow will automatically:

1. **Build binaries** for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64) 
   - Windows (amd64)

2. **Create archives**:
   - `.tar.gz` for Linux/macOS
   - `.zip` for Windows

3. **Generate checksums** for verification

4. **Create a GitHub Release** with:
   - Release notes from git commits
   - All binary artifacts
   - Checksums file

### 3. Release Artifacts

Each release will include:

```
opsbrew_Linux_x86_64.tar.gz
opsbrew_Linux_arm64.tar.gz
opsbrew_Darwin_x86_64.tar.gz
opsbrew_Darwin_arm64.tar.gz
opsbrew_Windows_x86_64.zip
checksums.txt
```

### 4. Installation from Release

Users can download and install from releases:

```bash
# Download for Linux
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Linux_x86_64.tar.gz
tar -xzf opsbrew_Linux_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/

# Download for macOS
wget https://github.com/nghiadaulau/opsbrew/releases/latest/download/opsbrew_Darwin_x86_64.tar.gz
tar -xzf opsbrew_Darwin_x86_64.tar.gz
sudo mv opsbrew /usr/local/bin/
```

## Configuration

The release process is configured in `.goreleaser.yml`:

- **Builds**: Multi-platform binary builds
- **Archives**: Tar.gz and zip formats
- **Checksums**: SHA256 checksums for verification
- **Changelog**: Auto-generated from git commits
- **Release**: GitHub release creation

## Future Enhancements

The following features are commented out in `.goreleaser.yml` and can be enabled when needed:

### Homebrew Tap
```yaml
brews:
  - name: opsbrew
    homepage: "https://github.com/nghiadaulau/opsbrew"
    repository:
      owner: nghiadaulau
      name: homebrew-tap
```

### Package Managers (DEB/RPM)
```yaml
nfpms:
  - maintainer: "your-email@example.com"
    formats:
      - deb
      - rpm
```

### Docker Images
```yaml
dockers:
  - image_templates:
      - "ghcr.io/{{ .Env.GITHUB_REPOSITORY }}:{{ .Version }}"
```

## Manual Release (if needed)

If you need to create a release manually:

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Test the release process (dry run)
goreleaser release --snapshot --clean --skip-publish

# Create a real release
goreleaser release --clean
```

## Version Management

- Use semantic versioning: `v1.0.0`, `v1.1.0`, `v2.0.0`
- Version information is embedded in binaries via ldflags
- Changelog is auto-generated from conventional commits
