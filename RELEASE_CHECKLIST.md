# Release Checklist

This checklist helps ensure a smooth release process for Roleplay.

## Pre-Release

- [ ] Update module path in `go.mod` to your GitHub username:
  ```
  module github.com/YOUR-USERNAME/roleplay
  ```
- [ ] Run the import update script:
  ```bash
  ./scripts/update-imports.sh
  ```
- [ ] Update version in CHANGELOG.md
- [ ] Run all tests: `make test`
- [ ] Run linter: `make lint`
- [ ] Build for all platforms: `make build-all`
- [ ] Test the binary locally
- [ ] Update README.md with any new features
- [ ] Review and update documentation

## GitHub Setup

1. Create a new repository on GitHub:
   - Name: `roleplay`
   - Description: "Advanced AI Character Bot with Psychological Modeling"
   - Public repository
   - Don't initialize with README (we have one)

2. Push the code:
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin git@github.com:YOUR-USERNAME/roleplay.git
   git push -u origin main
   ```

3. Configure repository settings:
   - Add topics: `go`, `ai`, `chatbot`, `cli`, `llm`, `openai`, `anthropic`
   - Set up branch protection for `main`
   - Enable GitHub Actions

## Creating a Release

1. Update version tag:
   ```bash
   git tag -a v0.1.0 -m "Initial release"
   git push origin v0.1.0
   ```

2. GitHub Actions will automatically:
   - Run tests
   - Build binaries for all platforms
   - Create a GitHub release with artifacts

3. After release is created:
   - [ ] Edit release notes if needed
   - [ ] Announce on social media
   - [ ] Update any documentation sites

## Post-Release

- [ ] Monitor issues for bug reports
- [ ] Create milestone for next version
- [ ] Update development branch

## Release Notes Template

```markdown
## What's New

- 🎭 Interactive TUI chat interface
- 🧠 OCEAN personality model with dynamic evolution
- ⚡ 4-layer caching for 90% cost reduction
- 📥 AI-powered character import from markdown
- 🔄 Support for OpenAI and Anthropic

## Installation

See the [README](README.md) for detailed installation instructions.

## Acknowledgments

Thanks to all contributors and testers!
```