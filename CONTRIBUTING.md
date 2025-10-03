# Contributing to Anvil

Thank you for your interest in contributing to Anvil! This document provides guidelines and instructions for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- [golangci-lint](https://golangci-lint.run/usage/install/) (recommended for linting)

### Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/osage-io/anvil.git
   cd anvil
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Install Git hooks (recommended):**
   ```bash
   make hooks
   ```
   
   This installs pre-commit and pre-push hooks that will automatically run tests and checks before commits and pushes.

## Git Hooks

We use Git hooks to ensure code quality and prevent broken builds. These hooks are **strongly recommended** for all contributors.

### What the hooks do:

**Pre-commit hook:**
- Checks Go code formatting (`gofmt`)
- Runs `go vet` for static analysis
- Runs `golangci-lint` for comprehensive linting
- Executes fast unit tests (`go test -short`)
- Runtime: ~30-60 seconds

**Pre-push hook:**
- Runs the full test suite with race detection
- Tests cross-compilation for multiple platforms
- Checks test coverage (minimum 70%)
- Runtime: ~2-5 minutes

### Installing hooks:

```bash
# Using make (recommended)
make hooks

# Or manually
./install-git-hooks.sh
```

### Bypassing hooks (emergency only):

If you need to bypass hooks in an emergency (not recommended for regular use):

```bash
# Skip pre-commit checks
git commit --no-verify -m "emergency fix"

# Skip pre-push checks  
git push --no-verify
```

### Uninstalling hooks:

```bash
rm .git/hooks/pre-commit .git/hooks/pre-push
```

### Troubleshooting hooks:

**Hook fails with "golangci-lint not found":**
- Install golangci-lint: https://golangci-lint.run/usage/install/
- Or the hook will skip linting with a warning

**Slow test performance:**
- Pre-commit runs only fast tests (`-short` flag)
- Pre-push runs comprehensive tests (this is intentional)
- Consider running `make test-short` manually during development

**Windows line ending issues:**
- Ensure Git is configured properly: `git config --global core.autocrlf true`
- The hooks should handle cross-platform differences automatically

## Development Workflow

### Making changes:

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Write code and tests:**
   - Follow Go conventions and existing code style
   - Add tests for new functionality
   - Ensure tests pass locally: `make test`

3. **Run quality checks:**
   ```bash
   make lint    # Run linter
   make fmt     # Format code
   make test    # Run all tests
   ```

4. **Commit changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```
   
   The pre-commit hook will automatically run checks.

5. **Push changes:**
   ```bash
   git push origin feature/your-feature-name
   ```
   
   The pre-push hook will run the comprehensive test suite.

### Commit message format:

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]

[optional footer(s)]
```

Examples:
- `feat: add support for Dogecoin derivation`
- `fix: correct BIP32 key derivation for testnet`
- `docs: update installation instructions`
- `test: add benchmarks for crypto operations`

### Testing

**Run all tests:**
```bash
make test
```

**Run fast tests only:**
```bash
make test-short
```

**Run comprehensive pre-push tests:**
```bash
make pre-push
# or directly:
./pre-push-test.sh
```

**Run benchmarks:**
```bash
go test -bench=. -benchmem ./internal/...
```

### Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and reasonably sized
- Write tests for new functionality

### Pull Request Process

1. Ensure all tests pass and hooks are installed
2. Update documentation if needed
3. Create a pull request with a clear description
4. Address any review feedback
5. Ensure CI passes before merging

## Project Structure

```
anvil/
├── cmd/anvil/          # Main application entry point
├── internal/           # Internal packages
│   ├── bitcoin/        # Bitcoin-specific implementations
│   ├── crypto/         # Cryptographic utilities
│   ├── ethereum/       # Ethereum-specific implementations
│   └── tron/          # Tron-specific implementations
├── pkg/               # Public packages (if any)
├── docs/              # Documentation
├── .githooks/         # Git hooks (installed via make hooks)
├── .github/           # GitHub workflows and templates
└── site/              # Website/documentation site
```

## Getting Help

- Check existing [Issues](https://github.com/osage-io/anvil/issues)
- Read the [README.md](README.md) for basic usage
- Look at existing code for patterns and examples
- Ask questions in pull request discussions

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to create a welcoming environment for all contributors.
