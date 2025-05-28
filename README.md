# Roleplay - Advanced AI Character Bot with Psychological Modeling

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

A sophisticated character bot system that implements psychologically-realistic AI characters with personality evolution, emotional states, and multi-layered memory systems. Features advanced prompt caching strategies that achieve 90% cost reduction in LLM API usage.

## ✨ Features

- 🎭 **Interactive TUI Chat**: Beautiful terminal interface with real-time chat, personality display, and performance metrics
- 🧠 **OCEAN Personality Model**: Characters with dynamic personality traits (Openness, Conscientiousness, Extraversion, Agreeableness, Neuroticism)
- 💭 **Emotional Intelligence**: Real-time emotional state tracking and blending
- 🗂️ **Multi-Tier Memory System**: Short-term, medium-term, and long-term memory with emotional weighting
- 🌱 **Personality Evolution**: Characters learn and adapt based on interactions with bounded drift
- ⚡ **4-Layer Caching Architecture**: Sophisticated caching system for optimal performance (90% cost reduction)
- 🔄 **Universal OpenAI-Compatible Support**: Works with any OpenAI-compatible API (OpenAI, Anthropic, Ollama, Groq, etc.)
- 📊 **Adaptive TTL**: Dynamic cache duration based on conversation patterns
- 📥 **Character Import**: Import characters from unstructured markdown files using AI

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- API key for your chosen provider (or local LLM service like Ollama)

### Installation

#### Option 1: Install from source

```bash
# Clone the repository
git clone https://github.com/dotcommander/roleplay.git
cd roleplay

# Install globally
go install

# Or build locally
go build -o roleplay
```

#### Option 2: Install from release

```bash
# Download the latest release for your platform
curl -L https://github.com/dotcommander/roleplay/releases/latest/download/roleplay-$(uname -s)-$(uname -m).tar.gz | tar xz
chmod +x roleplay
sudo mv roleplay /usr/local/bin/
```

### First Run

```bash
# Run the setup wizard (recommended)
roleplay init

# Or manually set your API key
export OPENAI_API_KEY="your-api-key"  # For OpenAI
export ANTHROPIC_API_KEY="your-key"   # For Anthropic
export OLLAMA_HOST="http://localhost:11434"  # For Ollama

# Quick start with built-in Rick Sanchez character
roleplay demo

# Or start interactive chat with any character
roleplay interactive --character rick-c137 --user your-name
```

## 📖 Usage

### Character Management

```bash
# List all characters
roleplay character list

# Create a character from JSON
roleplay character create character.json

# Import character from markdown (AI-powered)
roleplay import ~/Documents/my-character.md

# Show character details
roleplay character show character-id

# Generate example character JSON
roleplay character example > my-character.json
```

### Chat Commands

```bash
# Interactive mode (recommended) - Beautiful TUI
roleplay interactive --character rick-c137 --user your-name

# Single message chat
roleplay chat "Hello!" --character rick-c137 --user your-name

# Demo mode - Shows caching performance
roleplay demo
```

### Session Management

```bash
# List all sessions
roleplay session list

# Show session statistics (cache performance)
roleplay session stats
```

## 🎭 Example Characters

The system includes Rick Sanchez as a built-in demo character. You can import many more characters from markdown files or create your own!

### Example Characters Available

Check the `examples/characters/` directory for ready-to-use character files:
- **Sophia the Philosopher** - Thoughtful thinker who guides through questions
- **Captain Rex Thunderbolt** - Bold adventurer and sky pirate
- **Dr. Luna Quantum** - Meticulous quantum physicist

### Importing Characters from Markdown

You can import characters from unstructured markdown files using AI:

```bash
# Import a character from any markdown file
roleplay import ~/Documents/character-description.md

# The AI will analyze the file and extract:
# - Character name and personality
# - OCEAN personality traits
# - Speech patterns and quirks
# - Background story
```

### Creating Your Own Character

Create a JSON file with this structure:

```json
{
  "name": "Example Character",
  "backstory": "Character's background story...",
  "personality": {
    "openness": 0.8,
    "conscientiousness": 0.6,
    "extraversion": 0.7,
    "agreeableness": 0.8,
    "neuroticism": 0.3
  },
  "speech_style": "How the character speaks...",
  "quirks": ["quirk1", "quirk2"],
  "current_mood": {
    "joy": 0.7,
    "surprise": 0.3,
    "anger": 0.1,
    "fear": 0.2,
    "sadness": 0.1,
    "disgust": 0.1
  }
}
```

## ⚙️ Configuration

### Setup Wizard (Recommended)

The easiest way to configure roleplay is using the interactive setup wizard:

```bash
roleplay init
```

This will guide you through:
- Choosing your LLM provider (OpenAI, Anthropic, Ollama, etc.)
- Configuring API endpoints and keys
- Setting default models
- Creating example characters

### Manual Configuration

Create `~/.config/roleplay/config.yaml`:

```yaml
# Provider profile name (used for config resolution)
provider: openai  # or anthropic, ollama, gemini, etc.

# API Configuration
api_key: your-api-key-here
base_url: https://api.openai.com/v1  # Optional, for custom endpoints
model: gpt-4o-mini

# Caching Configuration
cache:
  max_entries: 10000
  cleanup_interval: 5m
  default_ttl: 10m
  adaptive_ttl: true

# Memory System
memory:
  short_term_window: 20
  medium_term_duration: 24h
  consolidation_rate: 0.1

# Personality Evolution
personality:
  evolution_enabled: true
  max_drift_rate: 0.02
  stability_threshold: 10
```

### Supported Providers

Roleplay uses a unified OpenAI-compatible provider that works with:

- **OpenAI** - Official OpenAI API
- **Anthropic** - Claude models via OpenAI-compatible endpoint
- **Google Gemini** - Via OpenAI-compatible proxy
- **Ollama** - Local models (no API key required)
- **LM Studio** - Local models (no API key required)
- **Groq** - Fast inference cloud service
- **OpenRouter** - Access multiple providers
- **Any OpenAI-compatible API** - Custom endpoints

### Environment Variables

```bash
# General configuration
export ROLEPLAY_PROVIDER=openai
export ROLEPLAY_API_KEY=your-api-key
export ROLEPLAY_BASE_URL=https://api.custom.com/v1
export ROLEPLAY_MODEL=gpt-4o-mini

# Provider-specific API keys (auto-detected)
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
export GEMINI_API_KEY=...
export GROQ_API_KEY=gsk-...

# Local services
export OLLAMA_HOST=http://localhost:11434
export ROLEPLAY_CACHE_DEFAULT_TTL=10m
export ROLEPLAY_CACHE_ADAPTIVE_TTL=true
```

## 🏗️ Architecture

### 4-Layer Caching System

1. **Admin/System Layer** - Global system prompts (24h+ TTL)
2. **Character Personality Layer** - Core traits and backstory (6-12h TTL)
3. **User Memory Layer** - User-specific relationships (1-3h TTL)
4. **Current Chat History** - Recent conversation (5-15m TTL)

### Key Components

- **Character System**: OCEAN personality model with emotional states
- **Memory System**: Three-tier memory with emotional weighting
- **Cache System**: Dual caching (prompt + response) with adaptive TTL
- **Provider Factory**: Centralized AI provider initialization and management
- **Provider Abstraction**: Supports multiple AI providers (OpenAI, Anthropic)

## 📊 Performance

- **90% cost reduction** through intelligent caching
- **Adaptive TTL** extends cache duration for active conversations
- **Background workers** for cache cleanup and memory consolidation
- **Thread-safe** operations throughout

## 🗺️ Roadmap

### Near-term (Q1 2025)
- [ ] **Streaming Support**: Real-time streaming of AI responses for more fluid conversations
- [ ] **Voice Input**: Add voice-to-text capabilities for natural conversation
- [ ] **Text-to-Speech**: AI responses read aloud with character-appropriate voices
- [ ] **Web UI**: Browser-based interface alongside the terminal UI
- [ ] **Character Marketplace**: Share and download community-created characters

### Mid-term (Q2-Q3 2025)
- [ ] **Multi-modal Characters**: Support for image generation and visual responses
- [ ] **Group Conversations**: Multiple characters interacting in the same chat
- [ ] **Character Relationships**: Characters remember relationships with each other
- [ ] **Mobile Support**: iOS and Android apps with full feature parity
- [ ] **Plugin System**: Extensible architecture for custom behaviors

### Long-term (Q4 2025+)
- [ ] **Advanced Emotional Modeling**: More nuanced emotional states and reactions
- [ ] **Character Learning**: Long-term personality evolution across users
- [ ] **Multi-language Support**: Characters speaking in different languages
- [ ] **Game Integration**: SDK for integrating characters into games
- [ ] **Enterprise Features**: Team management, analytics, and compliance tools

### Community Requested
- [ ] Discord bot integration
- [ ] Twitch chat integration
- [ ] Custom memory strategies
- [ ] Character mood visualization
- [ ] Export conversations to various formats

Want to contribute to these features? Check out our [Contributing Guidelines](CONTRIBUTING.md) or start a discussion in our [GitHub Discussions](https://github.com/dotcommander/roleplay/discussions).

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the beautiful TUI
- Uses [Cobra](https://github.com/spf13/cobra) for CLI management
- Inspired by advances in conversational AI and personality modeling

## 📞 Support

- 🐛 Issues: [GitHub Issues](https://github.com/dotcommander/roleplay/issues)
- 💡 Discussions: [GitHub Discussions](https://github.com/dotcommander/roleplay/discussions)
- 📚 Wiki: [GitHub Wiki](https://github.com/dotcommander/roleplay/wiki)

---

Made with ❤️ by the Roleplay team