# gh-label-sync

A GitHub CLI extension for bulk label management and synchronization from YAML/JSON files.

## Why This Extension?

Setting up labels in new repositories or syncing label changes across multiple repos is tedious with the GitHub web UI or requires writing custom scripts. This extension provides declarative label management with support for multiple file formats.

## Installation

```bash
gh extension install scttfrdmn/gh-label-sync
```

## Commands

### Sync Labels from File

```bash
gh label-sync --file .github/labels.yml
gh label-sync --file labels.json --repo owner/repo
gh label-sync --file labels.yml --dry-run
```

**Flags:**
- `--file` / `-f` (required): Path to label definition file (YAML or JSON)
- `--dry-run`: Show what would change without applying
- `--repo` / `-R`: Target repository (default: current repo)
- `--force`: Update existing labels even if they differ (default: skip)
- `--delete-unmanaged`: Remove labels not in file (dangerous, default: false)

### Export Labels

```bash
gh label-sync export > .github/labels.yml
gh label-sync export --format json > labels.json
gh label-sync export --repo source/repo | gh label-sync sync --file -
```

**Flags:**
- `--format`: Output format (`yaml` [default] or `json`)
- `--repo` / `-R`: Source repository

### Clone Labels Between Repos

```bash
gh label-sync clone source/repo --repo target/repo
gh label-sync clone source/repo --repo target/repo --force
```

Quick way to copy all labels from one repository to another.

## File Formats

### YAML Format (Recommended)

```yaml
# .github/labels.yml
labels:
  - name: "bug"
    color: "d73a4a"
    description: "Something isn't working"

  - name: "enhancement"
    color: "a2eeef"
    description: "New feature or request"

  - name: "type: feature"
    color: "1d76db"
    description: "New feature implementation"

  - name: "priority: critical"
    color: "b60205"
    description: "Needs immediate attention"
```

### JSON Format

```json
{
  "labels": [
    {
      "name": "bug",
      "color": "d73a4a",
      "description": "Something isn't working"
    },
    {
      "name": "enhancement",
      "color": "a2eeef",
      "description": "New feature or request"
    }
  ]
}
```

### CSV Format

```csv
name,color,description
bug,d73a4a,Something isn't working
enhancement,a2eeef,New feature or request
type: feature,1d76db,New feature implementation
```

**Field Requirements:**
- `name` (required): Label name
- `color` (required): 6-character hex color (with or without `#`)
- `description` (optional): Label description

## Behavior

### Default Sync Behavior

1. **Create missing labels**: Labels in file but not in repo → create
2. **Skip differing labels**: Labels exist but differ → skip (unless `--force`)
3. **Keep unmanaged labels**: Labels in repo but not in file → keep (unless `--delete-unmanaged`)

### Example Output

```bash
$ gh label-sync sync --file .github/labels.yml

Analyzing labels...
  ✓ bug - matches
  ✓ enhancement - matches
  + type: feature - will create
  + type: bug - will create
  ~ priority: high - exists but differs (color: ff9800 → fb8c00)
  ⚠ help wanted - exists but not in file

Summary:
  2 labels match
  2 labels to create
  1 label to update (use --force to apply)
  1 unmanaged label (use --delete-unmanaged to remove)

? Apply changes? (Y/n)
```

## Use Cases

### New Repository Setup

```bash
# Export labels from a template repository
gh label-sync export --repo myorg/template > .github/labels.yml

# Apply to new repository
cd new-project
gh label-sync sync --file .github/labels.yml
```

### Organization-Wide Label Standards

```bash
# Define labels once
cat > org-labels.yml <<EOF
labels:
  - name: "priority: critical"
    color: "b60205"
    description: "Needs immediate attention"
  - name: "priority: high"
    color: "d93f0b"
    description: "High priority"
  # ... more labels
EOF

# Apply to all repos
for repo in $(gh repo list myorg --json nameWithOwner -q '.[].nameWithOwner'); do
  gh label-sync sync --file org-labels.yml --repo $repo --force
done
```

### Copy Labels Between Repos

```bash
# Quick clone
gh label-sync clone source/repo --repo target/repo

# Or using export/sync
gh label-sync export --repo source/repo | \
  gh label-sync sync --file - --repo target/repo
```

## Development

### Prerequisites

- Go 1.21 or later
- GitHub CLI (`gh`) installed

### Local Installation

```bash
git clone https://github.com/scttfrdmn/gh-label-sync.git
cd gh-label-sync
go build
gh extension install .
```

### Testing

```bash
# Test the extension locally
gh label-sync export

# Run Go tests
go test ./...
```

### Building

```bash
# Build for current platform
go build -o gh-label-sync

# Cross-compile for release (done automatically by GitHub Actions)
GOOS=linux GOARCH=amd64 go build -o gh-label-sync-linux-amd64
```

## Architecture

```
gh-label-sync/
├── main.go              # Entry point and CLI setup
├── cmd/                 # Command implementations
│   ├── sync.go
│   ├── export.go
│   └── clone.go
├── pkg/
│   ├── api/            # GitHub API client wrapper
│   ├── parser/         # YAML/JSON/CSV parsing
│   ├── diff/           # Label diff algorithm
│   └── format/         # Output formatting
└── .github/
    └── workflows/
        └── release.yml  # Automated cross-platform builds
```

## Release Process

1. Update version in code
2. Commit changes: `git commit -am "Release v1.0.0"`
3. Create tag: `git tag v1.0.0`
4. Push tag: `git push --tags`
5. GitHub Actions automatically builds cross-platform binaries and creates release

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Related

- [gh-milestone](https://github.com/scttfrdmn/gh-milestone) - Companion extension for milestone management
- [GitHub CLI](https://cli.github.com/)
- [go-gh](https://github.com/cli/go-gh) - Go library for building GitHub CLI extensions

## Background

This extension was created because GitHub CLI maintainers [decided not to include](https://github.com/cli/cli/issues/9180) file-based label import in the core CLI, preferring users create extensions to avoid "format debates" between JSON, YAML, CSV, etc.

By being an extension, we can support **all** formats! 🎉

The design is based on comprehensive [real-world experience](../research/cargoship-setup-feedback.md) setting up GitHub projects.
