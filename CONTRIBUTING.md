# 🤝 Contributing to Roleplay

Thank you for considering contributing to Roleplay! Your contributions help make this AI character bot system better for everyone. Whether you're fixing bugs, adding features, improving documentation, or helping with testing, every contribution is valuable.

## 📋 Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Code Style Guidelines](#code-style-guidelines)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project is governed by our Code of Conduct. By participating, you agree to uphold this code. Please report unacceptable behavior to the project maintainers.

## 🚀 Ways to Contribute

### 🐛 Bug Reports
Found a bug? Help us fix it! Before reporting:

1. **Search existing issues** to avoid duplicates
2. **Use the latest version** to ensure the bug still exists
3. **Provide detailed information**:
   - Clear, descriptive title
   - Exact steps to reproduce
   - Expected vs. actual behavior
   - Configuration details (`roleplay status`)
   - Environment info (OS, Go version)
   - Relevant logs or screenshots

**Example Bug Report:**
```
Title: Character personality evolution not working with Ollama provider

Steps to reproduce:
1. Configure roleplay with Ollama provider
2. Chat with character for 20+ messages
3. Check personality traits with `roleplay character show character-id`

Expected: Personality traits should show slight changes
Actual: Personality traits remain exactly the same

Environment:
- OS: macOS 14.0
- Go: 1.23.1
- Provider: ollama (llama3)
- Config: ~/.config/roleplay/config.yaml (attached)
```

### ✨ Feature Requests
Have an idea to improve Roleplay? We'd love to hear it!

- **Start with a GitHub Discussion** for large features
- **Use GitHub Issues** for smaller enhancements
- **Explain the problem** your feature would solve
- **Describe your proposed solution** in detail
- **Consider alternative approaches**

### 📖 Documentation
Help make Roleplay easier to use:

- Fix typos or unclear explanations
- Add examples and use cases
- Improve code comments
- Write tutorials or guides
- Translate documentation

### 🧪 Testing
Help ensure Roleplay works reliably:

- Test new features and report issues
- Add unit tests for uncovered code
- Create integration tests
- Test on different platforms and providers

## 🛠️ Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.23+ | Main programming language |
| **golangci-lint** | Latest | Code linting and formatting |
| **API Key** | - | Testing with AI providers (OpenAI, Anthropic, or Ollama for local) |

### Quick Start

```bash
# 1. Fork and clone
git clone https://github.com/your-username/roleplay.git
cd roleplay

# 2. Add upstream remote
git remote add upstream https://github.com/dotcommander/roleplay.git

# 3. Install dependencies
go mod download

# 4. Verify setup
make test        # Run all tests
make lint        # Check code style
make build       # Build binary

# 5. Test your changes
./roleplay demo  # Quick functionality test
```

### Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Make your changes**
   - Follow the [code style guidelines](#code-style-guidelines)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test thoroughly**
   ```bash
   # Run full test suite
   make test

   # Test with different providers
   export OPENAI_API_KEY=sk-...
   ./roleplay api-test --provider openai

   # Manual testing
   ./roleplay interactive
   ```

4. **Commit with clear messages**
   ```bash
   git commit -m "feat: add streaming response support

   - Implement streaming for OpenAI provider
   - Add progress indicators to TUI
   - Update documentation with streaming examples
   
   Fixes #123"
   ```

## 📏 Code Style Guidelines

### Go Best Practices
```go
// ✅ Good: Clear function name and documentation
// ProcessCharacterRequest handles conversation requests for a specific character
func (cb *CharacterBot) ProcessCharacterRequest(ctx context.Context, req *ConversationRequest) (*AIResponse, error) {
    if req.CharacterID == "" {
        return nil, fmt.Errorf("character ID is required")
    }
    // ... implementation
}

// ❌ Bad: Unclear name and no documentation
func (cb *CharacterBot) process(r interface{}) interface{} {
    // ... implementation
}
```

### Code Organization
- **Package structure**: Follow the established `internal/` structure
- **Interface design**: Define interfaces in packages that use them
- **Error handling**: Use clear, actionable error messages
- **Concurrency**: Use mutexes appropriately for thread safety

### Documentation Standards
```go
// ✅ Good: Complete documentation
// CacheEntry represents a cached prompt entry with metadata for TTL management.
// It tracks access patterns and hit counts for adaptive caching strategies.
type CacheEntry struct {
    Breakpoints []CacheBreakpoint // Cache layers with individual TTLs
    CreatedAt   time.Time         // Entry creation timestamp
    LastAccess  time.Time         // Most recent access for TTL calculation
    HitCount    int               // Number of cache hits for popularity tracking
}

// ❌ Bad: Missing or unclear documentation
type CacheEntry struct {
    Breakpoints []CacheBreakpoint
    CreatedAt   time.Time
    LastAccess  time.Time
    HitCount    int
}
```

### Testing Guidelines

```bash
# Run comprehensive tests
make test                    # All tests with coverage
make test-integration        # Integration tests only
make test-cache             # Cache system tests
go test -race ./...         # Race condition detection

# Test specific functionality
go test ./internal/cache -v
go test ./internal/providers -run TestOpenAI
```

### Commit Message Format
Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat:` New features
- `fix:` Bug fixes  
- `docs:` Documentation changes
- `style:` Code style changes
- `refactor:` Code refactoring
- `test:` Test additions/changes
- `chore:` Maintenance tasks

**Examples:**
```bash
feat(cache): implement adaptive TTL for character prompts

- Add complexity-based TTL calculation
- Extend cache duration for active conversations  
- Add background cleanup worker

Closes #145

fix(tui): resolve character list scrolling issue

The character list wasn't properly updating scroll position
when new characters were added during active sessions.

Fixes #123
```

## 🏗️ Project Architecture

Understanding the codebase structure helps you contribute effectively:

```
roleplay/
├── cmd/                     # CLI command implementations
│   ├── root.go             # Main command setup and config
│   ├── chat.go             # Chat command logic
│   ├── interactive.go      # TUI implementation
│   └── character.go        # Character management
├── internal/               # Internal packages (not exported)
│   ├── cache/              # 4-layer caching system
│   ├── config/             # Configuration structures
│   ├── factory/            # Provider factory pattern
│   ├── models/             # Domain models (Character, Memory, etc.)
│   ├── providers/          # AI provider implementations
│   ├── services/           # Core business logic (CharacterBot)
│   ├── repository/         # Data persistence layer
│   └── tui/                # Terminal UI components
├── prompts/                # LLM prompt templates
├── examples/               # Example characters and configs
└── docs/                   # Technical documentation
```

## 🔧 Common Contribution Patterns

### Adding a New AI Provider

1. **Implement the interface** in `internal/providers/`:
   ```go
   type NewProvider struct {
       apiKey string
       baseURL string
   }

   func (p *NewProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
       // Your implementation here
   }
   ```

2. **Add to factory** in `internal/factory/provider.go`:
   ```go
   case "newprovider":
       return providers.NewProvider(apiKey, baseURL), nil
   ```

3. **Update documentation** and add tests

### Enhancing the Cache System

The cache system is in `internal/cache/`. Key areas:

- **`cache.go`**: Core cache logic and TTL management
- **`response_cache.go`**: Response-level caching
- **`types.go`**: Cache layer definitions

### Adding New Commands

1. Create command file in `cmd/`
2. Add to root command in `cmd/root.go`
3. Follow existing patterns for flags and validation

## 📤 Pull Request Process

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] All tests pass (`make test`)
- [ ] Code is properly documented
- [ ] Manual testing completed
- [ ] CHANGELOG.md updated (if applicable)

### PR Checklist

- [ ] **Clear title** following conventional commit format
- [ ] **Detailed description** explaining the change
- [ ] **Link to related issues** (Fixes #123)
- [ ] **Screenshots/demos** for UI changes
- [ ] **Breaking changes** clearly documented

### Review Process

1. **Automated checks** must pass (tests, linting)
2. **Code review** by maintainers
3. **Manual testing** for significant changes
4. **Documentation review** for user-facing changes

## 🎯 Good First Issues

Look for issues labeled:
- `good-first-issue` - Perfect for newcomers
- `help-wanted` - We'd love community help
- `documentation` - Improve docs and examples

## 💬 Getting Help

- **GitHub Discussions** - Questions and design discussions
- **GitHub Issues** - Bug reports and feature requests  
- **Code Comments** - Inline documentation and examples

## 🎉 Recognition

Contributors are recognized in:
- README.md contributors section
- GitHub contributors graph
- Release notes for significant contributions

Thank you for making Roleplay better! 🚀