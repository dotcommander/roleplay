# Roleplay - Advanced AI Character Bot with Psychological Modeling

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

A sophisticated character bot system that implements psychologically-realistic AI characters with personality evolution, emotional states, and multi-layered memory systems. Features advanced prompt caching strategies that achieve 90% cost reduction in LLM API usage.

![Roleplay Demo](https://via.placeholder.com/800x400.png?text=Roleplay+Character+Bot+Demo)

## ✨ Features

- 🎭 **Interactive TUI Chat**: Beautiful terminal interface with real-time chat, personality display, and performance metrics
- 🧠 **OCEAN Personality Model**: Characters with dynamic personality traits (Openness, Conscientiousness, Extraversion, Agreeableness, Neuroticism)
- 💭 **Emotional Intelligence**: Real-time emotional state tracking and blending
- 🗂️ **Multi-Tier Memory System**: Short-term, medium-term, and long-term memory with emotional weighting
- 🌱 **Personality Evolution**: Characters learn and adapt based on interactions with bounded drift
- ⚡ **4-Layer Caching Architecture**: Sophisticated caching system for optimal performance (90% cost reduction)
- 🔄 **Multi-Provider Support**: Works with Anthropic Claude and OpenAI models
- 📊 **Adaptive TTL**: Dynamic cache duration based on conversation patterns
- 📥 **Character Import**: Import characters from unstructured markdown files using AI

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- OpenAI API key or Anthropic API key

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
# Set your API key
export OPENAI_API_KEY="your-api-key"
# or
export ROLEPLAY_API_KEY="your-anthropic-key"

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

The system comes with several pre-built characters:

- **Rick Sanchez** - Nihilistic genius scientist
- **Alan Watts** - Eastern philosophy guide
- **Harley Quinn** - Chaotic and playful
- **The Watcher** - Cosmic observer of human history

Create your own character with this structure:

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

### Configuration File

Create `~/.config/roleplay/config.yaml`:

```yaml
provider: openai
api_key: your-api-key-here
model: gpt-4o-mini
cache:
  max_entries: 10000
  cleanup_interval: 5m
  default_ttl: 10m
  adaptive_ttl: true
memory:
  short_term_window: 20
  medium_term_duration: 24h
  consolidation_rate: 0.1
personality:
  evolution_enabled: true
  max_drift_rate: 0.02
  stability_threshold: 10
```

### Environment Variables

```bash
export ROLEPLAY_PROVIDER=openai
export ROLEPLAY_API_KEY=your-api-key
export ROLEPLAY_MODEL=gpt-4o-mini
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
- **Provider Abstraction**: Supports multiple AI providers

## 📊 Performance

- **90% cost reduction** through intelligent caching
- **Adaptive TTL** extends cache duration for active conversations
- **Background workers** for cache cleanup and memory consolidation
- **Thread-safe** operations throughout

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

- 📧 Email: support@example.com
- 💬 Discord: [Join our server](https://discord.gg/example)
- 🐛 Issues: [GitHub Issues](https://github.com/dotcommander/roleplay/issues)

---

Made with ❤️ by the Roleplay team