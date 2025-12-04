# GitHub Actions Workflows

This directory contains GitHub Actions workflows for CI/CD automation.

## Workflows

### CI Workflow (`ci.yml`)
Runs on every push and pull request to main branches.

**Features:**
- Tests on multiple platforms (Linux, macOS, Windows)
- Tests with Go 1.24 and 1.25
- Runs tests with race detection
- Generates code coverage reports
- Uploads coverage to Codecov
- Builds binaries for multiple platforms
- Runs golangci-lint for code quality

**Triggers:**
- Push to main, master, or develop branches
- Pull requests to main, master, or develop branches

### Release Workflow (`release.yml`)
Automatically creates releases when version tags are pushed.

**Features:**
- Runs tests before release
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
  - FreeBSD (amd64)
- Generates SHA256 checksums
- Creates GitHub releases with release notes

**Triggers:**
- Push tags matching `v*` (e.g., v1.0.0)

**Usage:**
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Docker Workflow (`docker.yml`)
Builds and publishes Docker images.

**Features:**
- Builds multi-platform Docker images (amd64, arm64)
- Pushes to GitHub Container Registry (ghcr.io)
- Optionally pushes to Docker Hub on release tags
- Uses layer caching for faster builds
- Tags images based on branch/tag names

**Triggers:**
- Push to main, master, or develop branches
- Push tags matching `v*`
- Pull requests (build only, no push)

**Container Registries:**
- GitHub Container Registry: `ghcr.io/amine7536/reverse-scan`
- Docker Hub (on releases): Requires `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets

## Configuration Files

### Dependabot (`dependabot.yml`)
Automatically creates PRs to update:
- GitHub Actions versions
- Go module dependencies
- Docker base images

Updates are checked weekly.

### golangci-lint (`.golangci.yml`)
Configures code linting with multiple linters for code quality, including:
- errcheck, govet, staticcheck
- gofmt, goimports
- gocritic, revive
- And more

## Secrets Required

For Docker Hub publishing (optional):
- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token

GitHub token (`GITHUB_TOKEN`) is automatically provided by GitHub Actions.
