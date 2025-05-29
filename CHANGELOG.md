# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.8.1] - 2025-05-29

### Added
- **Version Command and Information**
  - New `roleplay version` command shows version details
  - Version flag `--version` or `-V` for quick version check
  - Build information includes Git commit hash and build date
  - Go version and OS/architecture information

### Changed
- **Improved Build System**
  - Makefile now injects version information at build time
  - Automatic version detection from git tags
  - Build date and commit hash embedded in binary

## [0.8.0] - 2025-05-29

### Added
- **Character System Prompt Caching (1,024+ tokens)**
  - Expanded character model with 30+ new fields for richer personas
  - Character prompts now exceed OpenAI's 1,024 token minimum for automatic caching
  - 50% cost reduction on cached character prompts with OpenAI
  - 7-day TTL for character system prompts (configurable)
  - Automatic cache warmup on character creation
  - Cache invalidation support for character updates

- **Quick Character Generation (`quickgen` command)**
  - Generate complete characters from one-line descriptions
  - AI-powered character creation with psychological depth
  - Example: `roleplay character quickgen "A grumpy wizard who loves cats"`
  - Automatic ID generation from character names
  - JSON output option for programmatic use
  - Comprehensive prompt template for consistent generation

### Changed
- **Enhanced Character Model Structure**
  - Added demographics: age, gender, occupation, education, nationality, ethnicity
  - Added physical traits and appearance descriptions
  - Added skills, interests, fears, and goals arrays
  - Added relationships map with detailed connections
  - Added core beliefs, moral code, flaws, and strengths
  - Added catch phrases and dialogue examples
  - Added behavior patterns and emotional triggers
  - Added decision-making style and conflict resolution approach
  - Added worldview, life philosophy, and daily routines
  - Added hobbies, pet peeves, secrets, regrets, and achievements

### Improved
- **Prompt Caching Performance**
  - Character prompts increased from ~500 to 1,940+ tokens
  - 100% prompt cache hit rate after initial request
  - Significant cost savings for high-volume usage
  - Better consistency across conversations

## [0.7.1] - 2025-05-29

### Added
- **Comprehensive Test Suite**
  - Added tests for configuration package
  - Added tests for user profile agent service
  - Added tests for character manager
  - Added tests for all repository implementations
  - Added tests for CLI commands
  - Added tests for TUI model
  - Added integration tests

### Fixed
- Fixed panic in cache cleanup when cleanup interval is zero
- Fixed test compilation errors and import issues
- Updated tests to match current model structures and APIs

### Improved
- Test coverage across all major components
- Error handling validation in tests
- Concurrent access testing for thread safety

## [0.7.0] - 2025-05-29

### Changed
- **Major Command Structure Improvements**
  - Renamed `init` command to `setup` for better clarity
  - Moved `status` command under `config` (now `config status`)
  - Moved `api-test` command under `config` (now `config test`)
  - Moved `import` command under `character` (now `character import`)
  - Hidden advanced commands (`demo`, `scenario`) from main help
  - Disabled shell completion command to reduce clutter

### Added
- **Command Aliases for Common Operations**
  - `i` → `interactive` (start interactive chat)
  - `c` → `chat` (send single message)
  - `ls` → `character list` (list characters)
  
- **Organized Help Display**
  - Custom help template groups commands by category
  - Quick Start section highlights `quickstart` and `setup`
  - Chat Commands section for primary user interactions
  - Character Management section for character operations
  - Configuration section for system management
  - Cleaner, more intuitive command discovery

### Improved
- **User Experience**
  - Commands organized by frequency of use
  - Most common operations are immediately visible
  - Advanced/debug features hidden but still accessible
  - Consistent command naming conventions
  - Better command descriptions

## [0.6.0] - 2025-05-29

### Added
- **Robust JSON Extraction for LLM Outputs**
  - New `internal/utils/json.go` with `ExtractValidJSON` function
  - Handles common LLM issues: prefixes, markdown blocks, truncated JSON
  - Attempts to repair incomplete JSON structures
  - Successfully extracts valid JSON from malformed LLM responses

### Improved
- **Resilient Background User Profile Updates**
  - User profile updates now gracefully handle LLM parsing failures
  - Returns existing profile on failure instead of crashing
  - Better error logging with timestamps and context
  - 30-second timeout for background operations
  - Profile updates get 4000 max tokens (up from 2000)
  
- **Character List Display**
  - Shows full character backstories without truncation
  - Beautiful visual formatting with emojis and separators
  - Displays speech style and quirks for each character
  - Proper text wrapping for readability
  - Total character count at the bottom

### Fixed
- **User Profile Agent**
  - Fixed JSON parsing errors from truncated LLM responses
  - Background profile updates no longer affect main chat flow
  - Improved error messages distinguish between parsing and save failures

## [0.5.0] - 2025-05-29

### Added
- **Enhanced Configuration Management**
  - New `roleplay config list` shows configuration values with their sources (config file, env var, flag, or default)
  - New `roleplay config get <key>` retrieves specific configuration values
  - New `roleplay config set <key> <value>` updates config.yaml directly from CLI
  - New `roleplay config where` shows the configuration file location
  - Configuration sources are now transparent and debuggable

### Improved
- **Configuration User Experience**
  - `roleplay init` now shows where config will be saved and provides a summary after completion
  - Actionable error messages for missing API keys with 4 clear fix options
  - Updated documentation emphasizes `~/.config/roleplay/config.yaml` as the primary config location
  - Clear configuration precedence: flags > env vars > config file > defaults
  - Better guidance in init wizard pointing to config management commands

- **Documentation**
  - README now includes "Configuration Made Simple" section
  - Clear explanation of configuration precedence
  - Emphasized the "golden path" of using `roleplay init` and config file

## [0.4.2] - 2025-05-29

### Fixed
- **Configuration Generation**
  - Fixed `roleplay init` to correctly set `provider: openai` for Gemini and Anthropic
  - Fixed Gemini default model to include required `models/` prefix
  - Fixed Gemini base URL to point to correct OpenAI-compatible endpoint
  - Ensured all OpenAI-compatible providers use `provider: openai` in config

### Improved
- **Configuration Clarity**
  - Simplified provider configuration for OpenAI-compatible endpoints
  - Made it clear that Gemini, Anthropic, and others use `provider: openai`
  - Better alignment between init wizard choices and actual configuration

## [0.4.1] - 2025-05-29

### Fixed
- **Gemini API Support**
  - Fixed URL construction bug that automatically appended `/v1` to all base URLs
  - Provider now respects the exact base URL provided by users
  - Gemini's `/v1beta/openai/` endpoint now works correctly
  
### Added
- **HTTP Debug Logging**
  - Added `DEBUG_HTTP=true` environment variable for troubleshooting API requests
  - Shows exact URLs, HTTP methods, and response status codes
  - Helpful for debugging OpenAI-compatible endpoint issues

### Improved
- **Documentation**
  - Added Gemini configuration example to CLAUDE.md
  - Documented correct endpoint URLs for various providers
  - Added troubleshooting guidance for OpenAI-compatible endpoints

## [0.4.0] - 2025-05-28

### Changed
- **BREAKING: Unified OpenAI-Compatible Provider**
  - Replaced multiple provider implementations with a single OpenAI-compatible provider
  - All LLM interactions now use the `sashabaranov/go-openai` SDK
  - The `provider` setting is now a profile name for configuration resolution
  - Removed provider-specific code in favor of universal implementation

### Added
- **Expanded Provider Support**
  - Support for any OpenAI-compatible API endpoint
  - Pre-configured profiles for Anthropic, Google Gemini, Groq, and more
  - Better support for local models (Ollama, LM Studio) without API keys
  - Custom endpoint configuration via `base_url`

### Improved
- **Configuration Resolution**
  - Enhanced priority system: Flags > Config File > Environment Variables
  - Automatic API key detection from provider-specific environment variables
  - Smarter base URL resolution (including OLLAMA_HOST auto-detection)
  - Model defaults based on provider profile

- **Setup Wizard Enhancement**
  - Added more provider presets (8 total)
  - Auto-detection of running local services
  - Better guidance for OpenAI-compatible services

### Removed
- `SupportsBreakpoints()` and `MaxBreakpoints()` methods from provider interface
- Provider-specific implementations (AnthropicProvider)
- Complex provider type switching logic

### Technical
- Simplified provider interface to just `SendRequest()`, `SendStreamRequest()`, and `Name()`
- Consistent "openai_compatible" provider name for all endpoints
- Cleaner factory pattern implementation
- Updated all tests for the unified provider approach

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

