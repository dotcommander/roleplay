# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2025-05-28

### Added
- **Major User Experience Improvements**
  - `roleplay init` - Interactive setup wizard for first-time configuration
    - Auto-detects local LLM services (Ollama, LM Studio, LocalAI)
    - Guides through provider selection and API configuration
    - Creates example characters automatically
    - Generates proper config.yaml with sensible defaults
  - `roleplay quickstart` - Zero-configuration quick start
    - Automatically detects and uses local LLM services
    - Falls back to environment variables (OPENAI_API_KEY, ANTHROPIC_API_KEY)
    - Creates Rick Sanchez and starts chatting immediately
  - `roleplay config` - Configuration management commands
    - `config list` - Display all active configuration values
    - `config get <key>` - Retrieve specific configuration value
    - `config set <key> <value>` - Update configuration file
    - `config where` - Show configuration file location
    - Masks sensitive values (API keys) when displaying

- **OpenAI-Compatible API Support**
  - Added `base_url` configuration for custom endpoints
  - `--base-url` flag for command-line override
  - Environment variable support: `ROLEPLAY_BASE_URL`, `OPENAI_BASE_URL`, `OLLAMA_HOST`
  - Enables use with Ollama, LM Studio, OpenRouter, and other compatible services
  - Automatic endpoint detection for common local services

- **Model Aliases**
  - Configure friendly names for models in config.yaml
  - Example: `model_aliases: { fast: "gpt-4o-mini", smart: "gpt-4o", local: "llama3" }`
  - Use with `-m fast` instead of full model names

- **Architecture Improvements**
  - TUI componentization - Refactored interactive mode into reusable components
    - Created `internal/tui/components/` with StatusBar, Header, InputArea, MessageList
    - Improved testability with component-level unit tests
    - Reduced cyclomatic complexity significantly
  - Centralized CharacterBot initialization within CharacterManager
    - Eliminated boilerplate provider initialization code
    - Improved consistency across commands

### Changed
- CharacterManager now internally initializes AI providers and UserProfileAgent
- OpenAI provider now supports custom base URLs for compatible services
- Config structure expanded to support new features

### Fixed
- Character creation now properly uses CharacterManager for both bot and persistence

## [0.2.0] - 2025-05-28

### Added
- AI-Powered User Profile Agent - Intelligent system that builds and maintains user profiles
  - Automatic extraction of user information from conversations using LLM analysis
  - Character-specific profiles (each character maintains their own perception of users)
  - Confidence scoring for extracted facts (0.0-1.0)
  - Dynamic profile updates as conversations evolve
  - `profile` command for managing user profiles
    - `show <user-id> <character-id>` - Display a specific user profile
    - `list <user-id>` - List all profiles for a user
    - `delete <user-id> <character-id>` - Delete a user profile
  - Configurable update frequency and analysis depth
  - Privacy-aware design with user control over their data
  - Enriches conversations with learned context about users
- Enhanced `chat` command with session persistence and user profile support
  - Now saves conversation sessions for continuity across chats
  - Loads previous conversation context for more coherent interactions
  - Automatically triggers user profile updates based on configured frequency
  - Includes session ID in output for easy session management
  - Maps character roles correctly for API compatibility
- Scenario Context Cache - New highest-level cache layer for meta-prompts and operational contexts
  - 5-layer cache hierarchy with scenarios at the top (7-day TTL)
  - `scenario` command for managing high-level interaction contexts
    - `create` - Create new scenarios with custom prompts
    - `list` - List all available scenarios
    - `show` - Display scenario details
    - `update` - Update existing scenarios
    - `delete` - Remove scenarios
    - `example` - Show example scenario definitions
  - `--scenario` flag added to `chat` and `interactive` commands
  - Example scenarios included: starship bridge, therapy session, tech support, creative writing
- Command history navigation in interactive mode - use up/down arrows to navigate through previous commands
- `/memories` command to view character's memories about the user (planned)

### Fixed
- Fixed `character show` command to load characters from repository instead of expecting them in memory
- Fixed interactive mode to load all available characters on startup for proper `/list` and `/switch` functionality

### Features
- Initial release of Roleplay character bot system
- OCEAN personality model implementation
- Multi-tier memory system (short, medium, long-term)
- 4-layer caching architecture for 90% cost reduction
- Support for OpenAI and Anthropic providers
- Interactive TUI chat interface
- Character import from markdown files using AI
- Session management and statistics
- Personality evolution with bounded drift
- Emotional state tracking and blending
- Example character files
- Comprehensive documentation

### Features
- `character` command for managing characters
  - `create` - Create character from JSON
  - `list` - List all characters
  - `show` - Show character details
  - `example` - Generate example JSON
- `import` command for AI-powered markdown import
- `chat` command for single message interactions
- `interactive` command for TUI chat interface
- `demo` command for caching demonstration
- `session` command for session management
  - `list` - List all sessions
  - `stats` - Show caching statistics
- `api-test` command for testing API connectivity
- `status` command for configuration status

## [0.1.0] - TBD

Initial public release.