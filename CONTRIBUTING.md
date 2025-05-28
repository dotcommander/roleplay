# Contributing to Roleplay

First off, thank you for considering contributing to Roleplay! It's people like you that make Roleplay such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to support@example.com.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include your configuration and environment details**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain which behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Development Process

### Prerequisites

- Go 1.23 or higher
- golangci-lint (for linting)
- An OpenAI or Anthropic API key for testing

### Setting Up Your Development Environment

```bash
# Clone your fork
git clone https://github.com/your-username/roleplay.git
cd roleplay

# Add upstream remote
git remote add upstream https://github.com/original-owner/roleplay.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run

# Build the project
go build -o roleplay
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Write idiomatic Go code
- Add comments for exported functions and types
- Keep functions focused and small
- Write unit tests for new functionality

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/cache

# Run tests with race detection
go test -race ./...
```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add character import functionality

- Implement AI-powered markdown parser
- Add import command to CLI
- Support for various character formats
- Add comprehensive error handling

Fixes #123
```

### Project Structure

```
roleplay/
├── cmd/                    # CLI commands
├── internal/              # Internal packages
│   ├── cache/            # Caching system
│   ├── config/           # Configuration
│   ├── importer/         # Character importer
│   ├── manager/          # Character manager
│   ├── models/           # Data models
│   ├── providers/        # AI providers
│   ├── repository/       # Data persistence
│   ├── services/         # Core services
│   └── utils/            # Utilities
├── prompts/              # LLM prompt templates
├── examples/             # Example files
└── scripts/              # Build and utility scripts
```

### Adding a New AI Provider

1. Create a new file in `internal/providers/`
2. Implement the `AIProvider` interface
3. Add provider initialization in `cmd/root.go`
4. Update documentation

Example:
```go
type MyProvider struct {
    apiKey string
}

func NewMyProvider(apiKey string) *MyProvider {
    return &MyProvider{apiKey: apiKey}
}

func (p *MyProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
    // Implementation
}

func (p *MyProvider) SupportsBreakpoints() bool {
    return false
}

func (p *MyProvider) MaxBreakpoints() int {
    return 0
}

func (p *MyProvider) Name() string {
    return "myprovider"
}
```

### Documentation

- Update the README.md if you change functionality
- Add godoc comments to all exported types and functions
- Include examples in documentation where appropriate
- Update CHANGELOG.md for notable changes

## Release Process

We use GitHub Actions for automated releases. To create a new release:

1. Update version in appropriate files
2. Update CHANGELOG.md
3. Create a git tag: `git tag -a v1.2.3 -m "Release version 1.2.3"`
4. Push the tag: `git push origin v1.2.3`
5. GitHub Actions will automatically build and create the release

## Questions?

Feel free to open an issue with your question or reach out on our Discord server.

## Recognition

Contributors will be recognized in our README.md file. Thank you for your contributions!