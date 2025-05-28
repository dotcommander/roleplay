# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a sophisticated Go-based character bot architecture that implements psychologically-realistic AI characters with personality evolution, emotional states, and multi-layered memory systems. The codebase demonstrates advanced caching strategies to achieve 90% cost reduction in LLM API usage.

## Architecture

### Core Components

1. **Character System**
   - OCEAN personality model (Openness, Conscientiousness, Extraversion, Agreeableness, Neuroticism)
   - Emotional states with dynamic blending
   - Three-tier memory system (short-term, medium-term, long-term)
   - Personality evolution with bounded drift

2. **4-Layer Prompt Caching Architecture**
   Our sophisticated caching system implements 4 strategic layers for maximum token savings:
   
   **Layer 1: Admin/System Layer** - Global system prompts and admin instructions (longest TTL)
   **Layer 2: Character Personality Layer** - Core character traits, backstory, personality (long TTL)
   **Layer 3: User Memory Layer** - User-specific memories, relationships, context (medium TTL)
   **Layer 4: Current Chat History** - Recent conversation context (short TTL/no cache)

3. **Dual Caching System**
   - **Response Cache**: Stores complete API responses to avoid duplicate requests
   - **Prompt Cache**: Layers prompts with strategic breakpoints for provider caching
   - Automatic cache hit detection and metrics tracking
   - Adaptive TTL based on conversation activity and character complexity

4. **Provider Abstraction**
   - Interface-based design supporting multiple AI providers
   - Anthropic implementation with prompt caching (4 breakpoints)
   - OpenAI implementation with response caching and parameter optimization
   - Smart routing based on features, cost, or latency

5. **Performance Optimizations**
   - Adaptive TTL: 50% extension for active conversations, 20% for complex characters
   - Background workers for cache cleanup and memory consolidation
   - Thread-safe operations with proper mutex usage
   - Token tracking and optimization
   - Response deduplication for identical requests

## Development Commands

```bash
# Build the application
go build -o roleplay

# Run commands directly
go run main.go character example
go run main.go character create thorin.json
go run main.go chat "Hello!" --character warrior-123 --user user-789

# Install globally
go install

# Format code
go fmt ./...

# Download dependencies
go mod download
go mod tidy
```

## Key Design Patterns

- **Clean Architecture**: Separation between domain models, business logic, and external providers
- **Dependency Injection**: Providers registered at runtime
- **Interface-First Design**: All major components defined as interfaces
- **Concurrent Design**: Thread-safe operations throughout

## Important Implementation Details

### 4-Layer Cache Implementation
The caching system uses strategic breakpoints aligned with our 4-layer architecture:

**Layer 1: Admin/System Layer**
- Global system instructions and safety guidelines
- Administrative prompts and framework instructions
- Longest TTL (24+ hours) - rarely changes

**Layer 2: Character Personality Layer** 
- Character backstory, personality traits (OCEAN model)
- Core behavioral patterns and speech style
- Character-specific quirks and mannerisms
- Long TTL (6-12 hours) - stable character traits

**Layer 3: User Memory Layer**
- User-specific relationship dynamics
- Conversation history and shared memories
- User preferences and interaction patterns
- Medium TTL (1-3 hours) - evolves with relationship

**Layer 4: Current Chat History**
- Recent conversation turns and immediate context
- Current emotional state and active topics
- Short TTL (5-15 minutes) or no caching for real-time responses

### Memory Consolidation
- Automatic consolidation when short-term memory exceeds 10 entries
- Emotional weighting preserves important memories
- Background process runs every 5 minutes

### Personality Evolution
- Bounded drift prevents radical personality changes
- Learning rate of 0.1 for gradual adaptation
- Trait changes capped at ±0.2 from baseline

## Project Structure

The codebase follows clean Go CLI architecture with global configuration:

```
roleplay/
├── main.go                 # Entry point (<20 lines)
├── cmd/                    # Command definitions
│   ├── root.go            # Root command + shared config
│   ├── chat.go            # Chat command handler
│   ├── character.go       # Character management commands
│   ├── demo.go            # Caching demonstration
│   ├── interactive.go     # TUI chat interface
│   ├── session.go         # Session management
│   ├── status.go          # Configuration status
│   └── apitest.go         # API connectivity testing
├── internal/              # Private packages
│   ├── cache/             # Dual caching system (prompt + response)
│   ├── config/            # Configuration structures
│   ├── models/            # Domain models (Character, Memory, etc.)
│   ├── providers/         # AI provider implementations
│   ├── services/          # Core bot service and business logic
│   ├── repository/        # Character and session persistence
│   ├── manager/           # High-level character management
│   └── utils/             # Shared utilities (text wrapping, etc.)
├── prompts/               # LLM prompt templates (externalized)
├── migrate-config.sh      # Configuration migration script
├── chat-with-rick.sh      # Quick Rick Sanchez demo script
└── go.mod

### Global Configuration
- Config directory: `~/.config/roleplay/`
- Character storage: `~/.config/roleplay/characters/`
- Session storage: `~/.config/roleplay/sessions/`
- Cache storage: `~/.config/roleplay/cache/`
- Global binary: `~/go/bin/roleplay` (symlinked)
```

## Command Structure

```bash
roleplay
├── character              # Character management
│   ├── create            # Create from JSON file  
│   ├── list              # List all available characters
│   ├── show              # Display character details
│   └── example           # Generate example JSON
├── session                # Session management
│   ├── list              # List sessions for character(s)
│   └── stats             # Show caching performance metrics
├── interactive            # TUI chat interface (auto-creates Rick)
├── chat                   # Single message chat
├── demo                   # Caching demonstration (uses Rick by default)
├── api-test               # Test API connectivity
└── status                 # Show current configuration
```

## Cache Performance Features

### Demo Mode
- `roleplay demo` - Interactive demonstration of 4-layer caching
- Shows cache hits/misses in real-time with visual feedback
- Demonstrates token savings and cost reduction
- Uses Rick Sanchez character for engaging demo experience

### Session Persistence
- All conversations saved with cache metrics
- `roleplay session stats` shows aggregate caching performance
- Tracks hit rates, tokens saved, and cost savings across sessions
- Session data persists between application runs

### Cache Metrics Tracking
- Real-time cache hit/miss tracking
- Token usage optimization
- Cost savings calculations
- Performance latency measurements

## Usage Example

```go
// Initialize bot
config := Config{
    MaxShortTermMemory: 10,
    MaxMediumTermMemory: 50,
    MaxLongTermMemory: 200,
    CacheTTL: 5 * time.Minute,
}
bot := NewCharacterBot(config)

// Register providers
bot.RegisterProvider("anthropic", NewAnthropicProvider(apiKey))

// Create character
character := Character{
    ID: "warrior-maiden",
    Name: "Lyra",
    Personality: Personality{
        Openness: 0.7,
        Conscientiousness: 0.8,
        Extraversion: 0.6,
        Agreeableness: 0.5,
        Neuroticism: 0.3,
    },
    // ... other fields
}
bot.CreateCharacter(character)

// Process conversation
request := ConversationRequest{
    CharacterID: "warrior-maiden",
    UserID: "user123",
    Message: "Tell me about your adventures",
}
response, err := bot.ProcessRequest(ctx, request)
```

## Prompt Caching Strategy

Our goal is to implement prompt-caching in 4 layers:
- Admin layer
- System character prompt layer
- User memory layer
- Current chat history layer