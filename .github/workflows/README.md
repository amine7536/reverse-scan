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

Handles all builds and releases using GoReleaser.

**Features:**

- Runs tests before release
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
  - FreeBSD (amd64)
- Generates SHA256 checksums
- Creates GitHub releases with release notes
- Builds and publishes multi-platform Docker images (amd64, arm64) to Docker Hub
- Uses snapshot mode for non-tag builds (no push)

**Triggers:**

- Push to main, master, or develop branches (snapshot build, no push)
- Push tags matching `v*` (full release with GitHub release and Docker push)
- Pull requests (snapshot build, no push)
- Manual workflow dispatch

**Usage:**

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

**Container Registry:**

- Docker Hub: `amine7536/reverse-scan`
- Requires `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets for pushing images

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
