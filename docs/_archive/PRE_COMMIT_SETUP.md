# Pre-commit Hooks Setup Guide

Pre-commit hooks automatically run code quality checks before each commit, preventing bad code from being committed and ensuring all contributions meet project standards.

## Table of Contents

- [Quick Start](#quick-start)
- [What Hooks Do](#what-hooks-do)
- [Installation](#installation)
- [Usage](#usage)
- [Troubleshooting](#troubleshooting)
- [Configuration](#configuration)

## Quick Start

```bash
# One-time setup
pre-commit install

# That's it! Hooks now run automatically on git commit
git commit -m "feat: add awesome feature"
```

## What Hooks Do

Pre-commit hooks automatically run in this order before each commit:

1. **trailing-whitespace** - Removes trailing spaces from files
2. **end-of-file-fixer** - Ensures all files end with newline
3. **check-yaml** - Validates YAML syntax (if applicable)
4. **check-added-large-files** - Prevents accidentally committing large files
5. **check-merge-conflict** - Detects merge conflict markers
6. **detect-private-key** - Prevents committing secrets and credentials
7. **mixed-line-ending** - Enforces consistent line endings (LF)
8. **gofumpt** - Formats Go code (stricter than gofmt)
9. **goimports** - Organizes imports with local package prefix
10. **go-mod-tidy** - Tidies go.mod and go.sum files
11. **golangci-lint** - Full Go linting suite

If any hook **fails**, the commit is **aborted**. You must fix the issues and commit again.

## Installation

### Prerequisites

Ensure you have the required tools installed:

```bash
# Go 1.19+ (for gofumpt, goimports, golangci-lint)
go version

# pre-commit framework
pip install pre-commit
# or
brew install pre-commit  # macOS
```

### Setup

Navigate to the project root and install hooks:

```bash
cd /path/to/s9s
pre-commit install
```

Verify installation:

```bash
# Check that hook is installed
cat .git/hooks/pre-commit

# You should see a reference to .pre-commit-config.yaml
```

### One-time Setup for Existing Repos

If you already have the repository cloned:

```bash
# Install hooks
pre-commit install

# Run hooks on all files to catch any existing issues
pre-commit run --all-files

# Fix any issues found
make fmt
make lint

# Then commit
git add .
git commit -m "fix: address pre-commit issues"
```

## Usage

### Automatic Hooks (Default)

Hooks run automatically before each commit:

```bash
git commit -m "feat: add new feature"
# Hooks run automatically
# If they pass: commit succeeds
# If they fail: commit is aborted, fix issues and try again
```

### Manual Hook Execution

Sometimes you want to run hooks manually without committing:

```bash
# Run all hooks on changed files
pre-commit run

# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run gofumpt --all-files

# Run multiple specific hooks
pre-commit run gofumpt goimports golangci-lint --all-files
```

### Skipping Hooks (Emergency Only!)

In exceptional circumstances, you can skip hooks:

```bash
# Skip hooks for this commit
git commit --no-verify -m "emergency: skip hooks"

# BUT you must fix issues immediately after!
make fmt
make lint
git add .
git commit -m "fix: address pre-commit issues"
```

## Recommended Workflow

Follow this workflow to work smoothly with pre-commit hooks:

```bash
# 1. Make changes to code
# Edit files...

# 2. Test changes
make test

# 3. Fix formatting automatically
make fmt

# 4. Check for linting issues
make lint
# Fix any remaining issues manually

# 5. Stage changes
git add .

# 6. Commit (hooks run automatically)
git commit -m "feat: add awesome feature"
# If hooks fail, fix issues and commit again

# 7. If hooks pass, you're done!
# If they fail, go back to step 3 or 4
```

## Troubleshooting

### "pre-commit: command not found"

You haven't installed pre-commit framework:

```bash
pip install pre-commit
# or
brew install pre-commit  # macOS
```

### "Hook failed but I need to commit now"

Use `--no-verify` (only in emergencies!):

```bash
git commit --no-verify -m "emergency: skip hooks"
# But immediately fix issues:
make fmt && make lint
```

### "Hooks modified my files"

Some hooks modify files (gofumpt, goimports, go-mod-tidy). This is expected:

```bash
# Review changes
git diff

# Stage modified files
git add .

# Try commit again
git commit -m "feat: add awesome feature"
```

### "golangci-lint hook takes too long"

First commit will be slower as tools are installed. Subsequent commits are faster.

If golangci-lint consistently times out:

```bash
# Run it locally to debug
golangci-lint run

# Check what linters are slow
golangci-lint run -v

# See if there are specific slow files
# You may need to refactor or disable certain linters temporarily
```

### "Hook failed but I didn't change that code"

Pre-commit hooks check changed files. If a file you modified affects other files:

```bash
# Run hooks on all files to see full scope
pre-commit run --all-files

# Fix all issues
make fmt && make lint

# Commit
git add . && git commit -m "fix: address pre-commit issues"
```

### "Different results locally vs. CI"

Ensure your tools are up to date:

```bash
# Update pre-commit framework
pip install --upgrade pre-commit

# Update Go tools
go install github.com/mvdan/gofumpt@latest
go install golang.org/x/tools/cmd/goimports@latest
golangci-lint --version  # Should be v1.55.2 or later

# Update hooks configuration
pre-commit autoupdate

# Run hooks again
pre-commit run --all-files
```

## Configuration

### Hook Configuration

Hooks are configured in `.pre-commit-config.yaml`:

```yaml
# Specific Go formatting tool
- repo: https://github.com/dnephin/pre-commit-golang
  hooks:
    - id: go-fmt
      name: gofumpt
      entry: gofumpt
      args: [-w]  # Write in-place

# golangci-lint with timeout
- repo: https://github.com/golangci/golangci-lint
  rev: v1.55.2
  hooks:
    - id: golangci-lint
      args: [--new-from-rev=HEAD~1, --timeout=5m]
```

### Customizing Hooks

To disable a specific hook temporarily:

```yaml
# In .pre-commit-config.yaml
- repo: https://github.com/golangci/golangci-lint
  hooks:
    - id: golangci-lint
      stages: [manual]  # Run only with pre-commit run golangci-lint
```

Then run manually:

```bash
pre-commit run golangci-lint --all-files
```

## Best Practices

1. **Run hooks locally before pushing**
   ```bash
   pre-commit run --all-files
   ```

2. **Don't ignore hook failures**
   - They indicate real issues
   - Fix them properly, don't just skip hooks

3. **Keep hooks updated**
   ```bash
   pre-commit autoupdate
   ```

4. **Use hooks during development**
   - They prevent rework when you push to CI
   - Faster feedback loop

5. **Involve team in hook changes**
   - If you want to change hooks, discuss first
   - Document why in `.pre-commit-config.yaml`

## Integration with IDE

### VSCode

Install extensions:
- [Pre-commit](https://marketplace.visualstudio.com/items?itemName=pustelto.pre-commit-vscode) - Pre-commit framework support

Then configure in `.vscode/settings.json`:

```json
{
  "python.linting.enabled": false,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.fixAll": true
    }
  }
}
```

### GoLand / IntelliJ IDEA

Enable Git Hooks integration:
1. File → Settings → Tools → Git → Enable git hooks
2. Or use: Settings → Languages & Frameworks → Go → On Save (Ctrl+S)

## Related Documentation

- [CONTRIBUTING.md](../CONTRIBUTING.md#-linting-and-code-quality) - Contributing guidelines
- [docs/LINTING.md](./LINTING.md) - Linting standards and best practices
- [.pre-commit-config.yaml](../.pre-commit-config.yaml) - Hook configuration

## Questions?

- Check [LINTING.md](./LINTING.md) for linting questions
- See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines
- Open an issue on GitHub for help
