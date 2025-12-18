# Release Process

This project uses [goreleaser](https://goreleaser.com/) and [svu](https://github.com/caarlos0/svu) for automated releases.

## Overview

- **goreleaser**: Handles building binaries for multiple platforms, creating archives, checksums, Docker images, and GitHub releases
- **svu**: Semantic version utility that determines the next version based on git tags and conventional commits

## Release Workflow

### Option 1: Manual Tag Creation

1. Create and push a tag manually:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. The release workflow will automatically trigger and:
   - Build binaries for Linux, macOS, Windows, and FreeBSD (amd64 and arm64)
   - Create archives and checksums
   - Build multi-arch Docker images
   - Create a GitHub release with auto-generated release notes

### Option 2: Automated Versioning with svu

1. Go to Actions → Tag Version workflow
2. Click "Run workflow"
3. Select version type:
   - **auto**: Determines version automatically based on conventional commits
   - **major**: Bumps major version (e.g., 1.0.0 → 2.0.0)
   - **minor**: Bumps minor version (e.g., 1.0.0 → 1.1.0)
   - **patch**: Bumps patch version (e.g., 1.0.0 → 1.0.1)

4. The workflow will:
   - Calculate the next version using svu
   - Create and push the tag
   - Automatically trigger the release workflow

## Conventional Commits

When using `auto` versioning, svu determines the version bump based on commit messages:

- `feat:` or `feat(scope):` → Minor version bump
- `fix:` or `fix(scope):` → Patch version bump
- `BREAKING CHANGE:` or `!` → Major version bump

Examples:
```
feat: add support for IPv6 scanning
fix: resolve timeout issue in port scanner
feat!: redesign configuration file format (breaking change)
```

## Release Artifacts

Each release includes:

- **Binaries**: For multiple platforms and architectures
- **Archives**: `.tar.gz` (Unix) and `.zip` (Windows)
- **Checksums**: SHA256 checksums in `checksums.txt`
- **Docker Images**: Multi-arch images pushed to Docker Hub `amine7536/reverse-scan`
  - `amine7536/reverse-scan:latest`
  - `amine7536/reverse-scan:v1.0.0`

## Local Testing

Test the goreleaser configuration locally:

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Check configuration
goreleaser check

# Build snapshot (without publishing)
goreleaser release --snapshot --clean --skip=publish

# Test Docker build
goreleaser release --snapshot --clean --skip=publish,sign
```

## Version Information

The version is injected at build time using ldflags. When building manually:

```bash
go build -ldflags="-X main.Version=v1.0.0"
```

When using goreleaser, the version is automatically injected from the git tag.

## Docker Images

Docker images are built for both amd64 and arm64 architectures and published to Docker Hub:

```bash
# Pull latest
docker pull amine7536/reverse-scan:latest

# Pull specific version
docker pull amine7536/reverse-scan:v1.0.0
```

## Changelog

Changelogs are automatically generated from commit messages and organized by type:
- New Features
- Bug Fixes
- Enhancements
- Documentation
- Other Changes

Commits starting with `docs:`, `test:`, `chore:`, or merge commits are excluded from the changelog.
