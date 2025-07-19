# 🎭 Roleplay - AI Character Bot with Psychological Modeling

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](https://github.com/dotcommander/roleplay/actions)

**Roleplay** is a sophisticated character bot system that creates psychologically-realistic AI characters with personality evolution, emotional intelligence, and multi-layered memory systems. 

🚀 **90% Cost Reduction** through advanced 4-layer prompt caching  
🧠 **OCEAN Personality Model** with dynamic trait evolution  
💭 **Emotional Intelligence** with real-time state tracking  
🔄 **Universal Provider Support** - works with any OpenAI-compatible API

## ✨ Key Features

### 🎭 Character System
- **OCEAN Personality Model** - Scientific psychological framework with 5 core traits
- **Emotional Intelligence** - Real-time emotional state tracking and blending
- **Personality Evolution** - Characters learn and adapt with bounded drift protection
- **Memory System** - Three-tier memory (short/medium/long-term) with emotional weighting

### 💬 Chat Experience  
- **Interactive TUI** - Beautiful terminal interface with real-time metrics
- **AI Character Import** - Convert any markdown file into a character using AI
- **Session Management** - Persistent conversations with performance tracking
- **User Profiling** - Characters automatically learn about users over time

### ⚡ Performance & Cost Optimization
- **4-Layer Caching** - Advanced prompt caching achieving 90% cost reduction
- **Adaptive TTL** - Dynamic cache duration based on conversation patterns  
- **Rate Limiting** - Smart request throttling to maximize cache effectiveness
- **Token Optimization** - Automatic prompt structuring for optimal caching

### 🔄 Universal Provider Support
- **OpenAI** (GPT-4o, GPT-4o-mini, o1-mini)
- **Anthropic** (Claude 3.5 Sonnet, Haiku)  
- **Google Gemini** (1.5 Flash, Pro)
- **Local Models** (Ollama, LM Studio)
- **Cloud Services** (Groq, OpenRouter)
- **Custom Endpoints** (Any OpenAI-compatible API)

## 🚀 Quick Start

### Prerequisites
- **Go 1.23+** (for building from source)
- **API Key** for your chosen provider, or local LLM service like Ollama

### Installation

#### Option 1: Install from Release (Recommended)
```bash
# Download latest release for your platform
curl -L https://github.com/dotcommander/roleplay/releases/latest/download/roleplay-$(uname -s)-$(uname -m).tar.gz | tar xz
chmod +x roleplay
sudo mv roleplay /usr/local/bin/
```

#### Option 2: Build from Source
```bash
git clone https://github.com/dotcommander/roleplay.git
cd roleplay
go install  # Installs to $GOPATH/bin/roleplay
```

### Getting Started

#### 1. Initial Setup
```bash
# Interactive setup wizard - guides you through configuration
roleplay setup

# Or use the quick demo with built-in Rick Sanchez
roleplay demo
```

#### 2. Start Chatting
```bash
# Interactive TUI mode (recommended)
roleplay interactive

# Single message chat
roleplay chat "Hello, how are you?"

# Chat with specific character
roleplay chat "What's your latest invention?" --character rick-c137
```

## ⚙️ Configuration

### Quick Setup
```bash
# Guided setup wizard (recommended for first-time users)
roleplay setup

# Check current configuration
roleplay status

# Test your API connection
roleplay api-test
```

### Configuration Hierarchy
Settings are resolved in this order (highest to lowest priority):

1. **Command Flags** - `--api-key`, `--model`, `--provider`
2. **Environment Variables** - `ROLEPLAY_API_KEY`, `OPENAI_API_KEY`, etc.
3. **Config File** - `~/.config/roleplay/config.yaml`
4. **Defaults** - Sensible fallbacks

### Manual Configuration
**Config Location:** `~/.config/roleplay/config.yaml`

```yaml
# Basic setup
provider: openai
api_key: sk-your-key-here
model: gpt-4o-mini

# Advanced caching settings
cache:
  default_ttl: 10m
  adaptive_ttl: true
  max_entries: 10000

# Memory system
memory:
  short_term_window: 20
  medium_term_duration: 24h
  consolidation_rate: 0.1
```

## 📖 Usage Guide

### Character Management

```bash
# List all available characters
roleplay character list

# Quick-generate a character from description
roleplay character quickgen "A wise old wizard with a mysterious past"

# Import from any markdown file using AI
roleplay import ~/Documents/character-bio.md

# Import from Characters system format
roleplay import ~/path/to/character.json --source characters

# Create from structured JSON
roleplay character create wizard.json

# View character details
roleplay character show gandalf-123

# Generate example template
roleplay character example > my-character.json
```

### Chat Modes

```bash
# 🎭 Interactive TUI (recommended) - full-featured chat interface
roleplay interactive

# 💬 Quick chat - single message and response
roleplay chat "What's your greatest fear?"

# 🚀 Demo mode - see caching performance in action
roleplay demo

# 🎲 Scenario mode - chat within a specific scenario
roleplay chat "Ready for the mission?" --scenario starship-bridge
```

### Performance & Analytics

```bash
# View cache performance metrics
roleplay session stats

# List conversation history
roleplay session list

# Monitor user profiles (if enabled)
roleplay profile show alice
```

## 🎭 Creating Characters

### Built-in Characters
- **Rick Sanchez** - Genius scientist from Rick & Morty (demo character)

### Ready-to-Use Characters
Check `examples/characters/` for pre-made characters:
- **Sophia the Philosopher** - Thoughtful guide through life's questions
- **Captain Rex Thunderbolt** - Bold adventurer and sky pirate  
- **Dr. Luna Quantum** - Meticulous quantum physicist

### AI-Powered Import & Cross-System Compatibility
Transform any text into a character or import from other systems:

```bash
# Import from markdown, text, or any document
roleplay import ~/Documents/character-bio.md
roleplay import ~/Downloads/novel-character.txt

# Import from Characters system (cross-system compatibility)
roleplay import ~/path/to/characters-format.json --source characters

# AI automatically extracts:
# ✅ Name and personality traits (OCEAN model)
# ✅ Speech patterns and mannerisms  
# ✅ Background story and motivations
# ✅ Emotional baseline and quirks
```

### 🌉 Cross-System Compatibility

Roleplay features a **universal character bridge** that enables seamless character transfer between different AI systems:

```bash
# Import characters from the Characters differential system
roleplay import ~/characters-export.json --source characters

# Auto-detection of format (works with most JSON character files)
roleplay import ~/any-character.json --source auto

# Verbose import to see conversion details
roleplay import ~/character.json --source characters --verbose
```

**Bridge Features:**
- **Universal Format**: Common character representation across systems
- **Intelligent Conversion**: Automatically maps attributes, traits, and personality
- **OCEAN Model Mapping**: Converts any personality system to scientific OCEAN traits
- **Preservation**: Maintains original data while adapting to Roleplay's format
- **Warning System**: Alerts about potential data loss during conversion

### Quick Character Generation
```bash
# Generate from a simple description
roleplay character quickgen "A stoic samurai warrior from feudal Japan"
roleplay character quickgen "A cheerful robot chef who loves cooking"
```

### Manual Character Creation
Generate a template and customize:

```bash
# Create template
roleplay character example > my-character.json

# Edit with your favorite editor
vim my-character.json

# Import the character
roleplay character create my-character.json
```

### Character JSON Structure
```json
{
  "name": "Character Name",
  "backstory": "Detailed background story...",
  "personality": {
    "openness": 0.8,        // Creativity, curiosity (0-1)
    "conscientiousness": 0.6, // Organization, discipline (0-1)  
    "extraversion": 0.7,     // Social energy, assertiveness (0-1)
    "agreeableness": 0.8,    // Cooperation, trust (0-1)
    "neuroticism": 0.3       // Emotional stability (0-1)
  },
  "speech_style": "How they communicate...",
  "quirks": ["Specific mannerisms", "Unique habits"],
  "skills": ["Area of expertise", "Special abilities"],
  "goals": ["Primary motivation", "Life ambitions"]
}
```

## 🔧 Provider Setup

### Supported Providers

| Provider | Models | Notes |
|----------|--------|-------|
| **OpenAI** | GPT-4o, GPT-4o-mini, o1-mini | Official API with prompt caching |
| **Anthropic** | Claude 3.5 Sonnet, Haiku | Via OpenAI-compatible proxy |
| **Google Gemini** | 1.5 Flash, 1.5 Pro | Via OpenAI-compatible proxy |
| **Ollama** | Llama, Mistral, etc. | Local models (no API key) |
| **Groq** | Llama3, Mixtral | Fast inference cloud |
| **OpenRouter** | 100+ models | Unified access to many providers |

### Quick Provider Setup

#### OpenAI
```bash
export OPENAI_API_KEY=sk-your-key-here
roleplay chat "Hello!" --provider openai --model gpt-4o-mini
```

#### Anthropic (Claude)
```bash
export ANTHROPIC_API_KEY=sk-ant-your-key
roleplay chat "Hello!" --provider anthropic --model claude-3-5-sonnet-20241022
```

#### Ollama (Local)
```bash
# Start Ollama with your model
ollama serve
ollama pull llama3

# No API key needed
roleplay chat "Hello!" --provider ollama --model llama3
```

#### Gemini
```bash
export GEMINI_API_KEY=your-key
roleplay chat "Hello!" \
  --provider openai \
  --base-url https://generativelanguage.googleapis.com/v1beta/openai/ \
  --model models/gemini-1.5-flash-latest
```

### Environment Variables

```bash
# Universal settings
export ROLEPLAY_API_KEY=your-key      # Works with any provider
export ROLEPLAY_PROVIDER=openai       # Default provider
export ROLEPLAY_MODEL=gpt-4o-mini     # Default model

# Provider-specific (auto-detected)
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
export GEMINI_API_KEY=...
export GROQ_API_KEY=gsk-...

# Local services
export OLLAMA_HOST=http://localhost:11434
```

## 🏗️ Architecture Overview

### 🌉 Universal Character Bridge
Roleplay's bridge system enables cross-platform character compatibility:

**Bridge Architecture:**
- **Universal Format**: JSON-based character representation with OCEAN personality model
- **Converter Registry**: Pluggable converters for different AI systems (Characters, CharacterAI, etc.)
- **Intelligent Mapping**: Automatic trait analysis and personality conversion
- **Data Preservation**: Original format metadata maintained for round-trip compatibility

**Workflow Example:**
```bash
# 1. Characters system exports character
characters export abc123 --format roleplay --output-dir ./exports

# 2. Roleplay imports the character  
roleplay import ./exports/vampire_companion_roleplay.json --source characters

# 3. Character is now available in Roleplay with full personality mapping
roleplay chat "Tell me about your immortal experiences" --character vampire-companion-123
```

### 4-Layer Prompt Caching System
Roleplay achieves 90% cost reduction through strategic prompt caching:

| Layer | Content | TTL | Cache Efficiency |
|-------|---------|-----|------------------|
| **System** | Global instructions, safety guidelines | 24h+ | 95%+ hit rate |
| **Character** | Core personality, backstory, traits | 6-12h | 90%+ hit rate |
| **User Context** | User-specific memories, relationships | 1-3h | 70%+ hit rate |
| **Conversation** | Recent chat history | 5-15m | 40%+ hit rate |

### Core Components

#### 🧠 Character System  
- **OCEAN Model**: Scientific personality framework with 5 traits
- **Emotional States**: 6-dimensional emotion tracking (joy, anger, fear, etc.)
- **Memory System**: 3-tier memory with automatic consolidation
- **Evolution**: Bounded personality drift based on interactions

#### ⚡ Performance Engine
- **Dual Caching**: Prompt cache + response cache for maximum efficiency
- **Adaptive TTL**: Dynamic cache duration based on activity patterns
- **Rate Limiting**: Smart throttling to maximize cache effectiveness  
- **Background Workers**: Automatic memory consolidation and cache cleanup

#### 🔌 Provider Abstraction
- **Universal Interface**: Single OpenAI-compatible API for all providers
- **Factory Pattern**: Centralized provider initialization
- **Smart Routing**: Automatic model selection and endpoint configuration
- **Failure Handling**: Graceful degradation and error recovery

### Performance Metrics
- **90% cost reduction** through intelligent caching
- **80% latency reduction** for cached responses  
- **Thread-safe** operations with comprehensive mutex usage
- **Background processing** maintains UI responsiveness

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