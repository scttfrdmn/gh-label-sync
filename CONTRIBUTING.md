# Contributing to gh-label-sync

Thank you for your interest in contributing! This document provides guidelines for contributing to gh-label-sync.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/gh-label-sync.git
   cd gh-label-sync
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Install the extension locally** for testing:
   ```bash
   gh extension install .
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   ```
   Or for bug fixes:
   ```bash
   git checkout -b fix/bug-description
   ```

2. **Make your changes** following our coding standards

3. **Test your changes**:
   ```bash
   # Build the extension
   go build

   # Run tests
   go test ./...

   # Test the extension manually
   gh label-sync --help
   ```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add CSV export support`
- `fix: handle empty label descriptions`
- `docs: update installation instructions`
- `test: add tests for diff algorithm`
- `refactor: simplify parser logic`
- `chore: update dependencies`

### Pull Requests

1. **Push your branch** to your fork:
   ```bash
   git push origin feat/your-feature-name
   ```

2. **Open a pull request** on GitHub with:
   - Clear title following conventional commits format
   - Description of what changed and why
   - Reference to related issues (if any)
   - Screenshots/examples for UI changes

3. **Respond to feedback** from maintainers

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Run `go vet` to catch common issues
- Keep functions focused and small
- Add comments for exported functions

### Testing

- Add tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Documentation

- Update README.md if adding features
- Add inline comments for complex logic
- Update command help text if adding flags
- Include examples for new functionality

## Project Structure

```
gh-label-sync/
├── main.go              # Entry point
├── cmd/                 # Command implementations
│   ├── sync.go
│   ├── export.go
│   └── clone.go
├── pkg/
│   ├── api/            # GitHub API wrapper
│   ├── parser/         # YAML/JSON/CSV parsing
│   ├── diff/           # Label diff algorithm
│   └── format/         # Output formatting
└── .github/
    └── workflows/      # CI/CD workflows
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/parser
```

### Manual Testing

```bash
# Build and install locally
go build && gh extension install .

# Test sync command
gh label-sync sync --file examples/labels.yml --dry-run

# Test export
gh label-sync export

# Test with different formats
gh label-sync sync --file examples/labels.json
gh label-sync sync --file examples/labels.csv
```

## Release Process

Maintainers handle releases:

1. Update version in code
2. Update CHANGELOG (if we add one)
3. Create and push tag: `git tag v1.0.0 && git push --tags`
4. GitHub Actions builds and creates release automatically

## Getting Help

- **Questions?** Open a [Discussion](https://github.com/scttfrdmn/gh-label-sync/discussions)
- **Bug?** Open an [Issue](https://github.com/scttfrdmn/gh-label-sync/issues/new?template=bug_report.md)
- **Feature idea?** Open an [Issue](https://github.com/scttfrdmn/gh-label-sync/issues/new?template=feature_request.md)

## Code of Conduct

Be respectful and constructive. We're all here to make better tools together.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
