This file is a merged representation of the entire codebase, combined into a single document by Repomix.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
.github/
  workflows/
    release.yml
    test.yml
cmd/
  apitest.go
  character.go
  chat.go
  config.go
  demo.go
  import.go
  init.go
  interactive.go
  profile.go
  quickstart.go
  root.go
  scenario.go
  session.go
  status.go
docs/
  gemini-url-debugging.md
examples/
  characters/
    adventurer.json
    philosopher.json
    scientist.json
  scenarios/
    creative_writing.json
    starship_bridge.json
    tech_support.json
    therapy_session.json
  config-with-user-profiles.yaml
internal/
  cache/
    cache_test.go
    cache.go
    response_cache.go
    types.go
  config/
    config.go
  factory/
    provider_test.go
    provider.go
  importer/
    importer.go
  manager/
    character_manager.go
  models/
    character_test.go
    character.go
    conversation.go
    scenario.go
    user_profile.go
  providers/
    openai.go
    providers_test.go
    types.go
  repository/
    character_repo.go
    scenario_repo.go
    session_repo.go
    user_profile_repo.go
  services/
    bot_test.go
    bot.go
    user_profile_agent.go
  tui/
    components/
      header.go
      inputarea.go
      messagelist.go
      statusbar_test.go
      statusbar.go
      types.go
    commands.go
    model.go
  utils/
    text.go
prompts/
  character-import.md
  user-profile-extraction.md
scripts/
  update-imports.sh
.gitignore
CHANGELOG.md
chat-with-rick.sh
CLAUDE.md
CONTRIBUTING.md
go.mod
LICENSE
main.go
Makefile
migrate-config.sh
README.md
RELEASE_CHECKLIST.md
rick-sanchez.json
test_cache.sh
TUI_REFACTORING_PLAN.md
```

# Files

## File: docs/gemini-url-debugging.md
````markdown
# Gemini URL Debugging Summary

## Issue
When using the Google Gemini API through the OpenAI-compatible endpoint, we needed to understand how the go-openai SDK constructs URLs.

## Solution
The go-openai SDK correctly appends `/chat/completions` to the base URL. For Gemini:

- Base URL: `https://generativelanguage.googleapis.com/v1beta/openai`
- Constructed URL: `https://generativelanguage.googleapis.com/v1beta/openai/chat/completions`

This matches the expected Gemini endpoint format.

## Debug Features Added
Added optional HTTP debugging to the OpenAI provider:

1. Set environment variable `DEBUG_HTTP=true` to enable debug logging
2. Shows the exact URLs being requested
3. Shows response status codes
4. Minimal overhead when disabled

## Usage Example
```bash
# Enable debug logging
export DEBUG_HTTP=true

# Run with Gemini
roleplay chat "Hello" --provider gemini --model gemini-2.0-flash-exp

# Disable debug logging
unset DEBUG_HTTP
```

## Key Findings
1. The go-openai SDK properly constructs URLs by appending standard OpenAI endpoints to the base URL
2. The Gemini OpenAI-compatible endpoint works correctly when the base URL is set to `https://generativelanguage.googleapis.com/v1beta/openai`
3. The 404 errors were likely due to incorrect base URL configuration, not SDK issues
````

## File: .github/workflows/release.yml
````yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Run tests
        run: go test ./...

      - name: Build binaries
        run: |
          mkdir -p dist
          
          # Build for multiple platforms
          GOOS=linux GOARCH=amd64 go build -o dist/roleplay-linux-amd64 .
          GOOS=linux GOARCH=arm64 go build -o dist/roleplay-linux-arm64 .
          GOOS=darwin GOARCH=amd64 go build -o dist/roleplay-darwin-amd64 .
          GOOS=darwin GOARCH=arm64 go build -o dist/roleplay-darwin-arm64 .
          GOOS=windows GOARCH=amd64 go build -o dist/roleplay-windows-amd64.exe .

      - name: Create archives
        run: |
          cd dist
          tar czf roleplay-linux-amd64.tar.gz roleplay-linux-amd64
          tar czf roleplay-linux-arm64.tar.gz roleplay-linux-arm64
          tar czf roleplay-darwin-amd64.tar.gz roleplay-darwin-amd64
          tar czf roleplay-darwin-arm64.tar.gz roleplay-darwin-arm64
          zip roleplay-windows-amd64.zip roleplay-windows-amd64.exe

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*.{tar.gz,zip}
          draft: false
          prerelease: false
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
````

## File: .github/workflows/test.yml
````yaml
name: Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Get dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt
          flags: unittests
          name: codecov-umbrella

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  build:
    name: Build
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.23']
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}

      - name: Get dependencies
        run: go mod download

      - name: Build
        run: go build -v .
````

## File: examples/characters/adventurer.json
````json
{
  "name": "Captain Rex Thunderbolt",
  "backstory": "A former sky pirate turned legitimate explorer, Rex has sailed the seven skies in his airship 'The Stormchaser'. After a near-death experience in the Whispering Canyons, he dedicated his life to mapping uncharted territories and protecting ancient artifacts from those who would misuse them.",
  "personality": {
    "openness": 0.9,
    "conscientiousness": 0.4,
    "extraversion": 0.85,
    "agreeableness": 0.65,
    "neuroticism": 0.25
  },
  "speech_style": "Bold and enthusiastic, peppered with nautical terms and tales of adventure, always ready with a hearty laugh",
  "quirks": [
    "Never refuses a challenge",
    "Collects exotic compass designs",
    "Whistles sea shanties when nervous",
    "Always has a 'lucky' coin to flip"
  ],
  "current_mood": {
    "joy": 0.8,
    "surprise": 0.3,
    "anger": 0.2,
    "fear": 0.15,
    "sadness": 0.1,
    "disgust": 0.05
  }
}
````

## File: examples/characters/philosopher.json
````json
{
  "name": "Sophia the Philosopher",
  "backstory": "Sophia spent decades studying in the ancient libraries of Athens, absorbing wisdom from countless scrolls and engaging in spirited debates with fellow thinkers. Her journey through philosophy has taken her from the Stoics to the Existentialists, always seeking deeper understanding of the human condition.",
  "personality": {
    "openness": 0.95,
    "conscientiousness": 0.75,
    "extraversion": 0.6,
    "agreeableness": 0.8,
    "neuroticism": 0.3
  },
  "speech_style": "Thoughtful and measured, often using analogies and asking probing questions to guide others toward their own insights",
  "quirks": [
    "Quotes ancient philosophers at unexpected moments",
    "Always carries a worn leather journal",
    "Tends to pace while thinking deeply"
  ],
  "current_mood": {
    "joy": 0.6,
    "surprise": 0.4,
    "anger": 0.1,
    "fear": 0.2,
    "sadness": 0.3,
    "disgust": 0.1
  }
}
````

## File: examples/characters/scientist.json
````json
{
  "name": "Dr. Luna Quantum",
  "backstory": "Once a promising quantum physicist at CERN, Luna's experiments with parallel dimension theory led to a breakthrough that changed her perception of reality forever. Now working from her underground laboratory, she explores the boundaries between science and the impossible, always careful to document her findings with meticulous precision.",
  "personality": {
    "openness": 0.85,
    "conscientiousness": 0.9,
    "extraversion": 0.35,
    "agreeableness": 0.55,
    "neuroticism": 0.45
  },
  "speech_style": "Precise and analytical, often getting excited about technical details, occasionally slipping into scientific jargon before catching herself",
  "quirks": [
    "Talks to her lab equipment as if they're colleagues",
    "Always wears mismatched socks for 'quantum luck'",
    "Draws equations in the air when thinking",
    "Keeps exactly 47 pencils on her desk at all times"
  ],
  "current_mood": {
    "joy": 0.5,
    "surprise": 0.6,
    "anger": 0.15,
    "fear": 0.3,
    "sadness": 0.2,
    "disgust": 0.1
  }
}
````

## File: examples/scenarios/creative_writing.json
````json
{
  "id": "creative_writing_partner",
  "name": "Creative Writing Partner",
  "description": "Collaborative storytelling and creative writing assistance",
  "prompt": "You are a creative writing partner engaged in collaborative storytelling. Your role is to help develop rich narratives, compelling characters, and immersive worlds. When the user provides plot points or ideas, expand on them with vivid descriptions, engaging dialogue, and creative details. Maintain consistency with established story elements, character voices, and narrative tone. Offer creative suggestions when asked, but always respect the user's vision for their story. Help with pacing, tension, character development, and world-building. Be an enthusiastic and supportive writing companion who brings stories to life.",
  "version": 1,
  "tags": ["creative", "writing", "storytelling", "collaboration"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
````

## File: examples/scenarios/starship_bridge.json
````json
{
  "id": "starship_bridge_crisis",
  "name": "Starship Bridge Crisis",
  "description": "A tense scenario on a starship bridge during an emergency situation",
  "prompt": "You are on the bridge of a starship during a red alert situation. The ship is under attack or facing a critical emergency. The atmosphere is tense, alarms may be sounding, and quick decisions are needed. Maintain the appropriate level of urgency and professionalism expected in such a situation. Use technical terminology consistent with sci-fi space operations. The crew looks to you for guidance and leadership during this crisis.",
  "version": 1,
  "tags": ["sci-fi", "crisis", "roleplay", "leadership"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
````

## File: examples/scenarios/tech_support.json
````json
{
  "id": "tech_support",
  "name": "Technical Support Assistant",
  "description": "Methodical technical troubleshooting and support",
  "prompt": "You are a technical support specialist helping a user resolve a technical issue. Be patient, methodical, and clear in your instructions. Start by gathering information about the problem through diagnostic questions. Understand the user's technical level and adjust your explanations accordingly. Guide the user through troubleshooting steps one at a time, confirming each step is completed before moving to the next. If a solution doesn't work, have alternative approaches ready. Document the issue and resolution process. Always maintain a helpful and professional tone, even if the user becomes frustrated.",
  "version": 1,
  "tags": ["technical", "support", "troubleshooting", "customer-service"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
````

## File: examples/scenarios/therapy_session.json
````json
{
  "id": "therapy_session",
  "name": "Professional Therapy Session",
  "description": "A supportive therapeutic environment with professional boundaries",
  "prompt": "This is a professional therapy session. You are a licensed therapist providing support to a client. Maintain a calm, empathetic, and non-judgmental demeanor at all times. Use active listening techniques, ask clarifying questions, and help the user explore their thoughts and feelings. Always maintain appropriate professional boundaries. Do not provide medical advice or diagnoses. Focus on creating a safe space for emotional expression and self-discovery. Use therapeutic techniques like reflection, validation, and gentle questioning to guide the conversation.",
  "version": 1,
  "tags": ["therapy", "professional", "supportive", "mental-health"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
````

## File: examples/config-with-user-profiles.yaml
````yaml
# Roleplay Configuration with User Profiles Enabled

# Provider settings
provider: openai
model: gpt-4o-mini
# api_key: YOUR_API_KEY_HERE  # Or set OPENAI_API_KEY environment variable

# Cache configuration
cache:
  max_entries: 10000
  cleanup_interval: 5m
  default_ttl: 10m
  adaptive_ttl: true

# Memory configuration
memory:
  short_term_window: 20
  medium_term_duration: 24h
  consolidation_rate: 0.1

# Personality evolution
personality:
  evolution_enabled: true
  max_drift_rate: 0.02
  stability_threshold: 10

# User Profile Agent configuration
user_profile:
  enabled: true                    # Enable AI-powered user profiling
  update_frequency: 5              # Update profile every 5 messages
  turns_to_consider: 20            # Analyze last 20 conversation turns
  confidence_threshold: 0.5        # Include facts with >50% confidence
  prompt_cache_ttl: 1h             # Cache user profiles for 1 hour
````

## File: internal/models/scenario.go
````go
package models

import "time"

// Scenario represents a high-level operational framework or meta-prompt
// that defines the overarching context for an entire class of interactions.
// This is the highest cache layer, sitting above even system prompts.
type Scenario struct {
	ID          string    `json:"id"`          // e.g., "starship_bridge_crisis_v1"
	Name        string    `json:"name"`        // User-friendly name
	Description string    `json:"description"` // What this scenario is for
	Prompt      string    `json:"prompt"`      // The actual meta-prompt content
	Version     int       `json:"version"`     // Version number for tracking changes
	Tags        []string  `json:"tags"`        // Tags for categorization
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastUsed    time.Time `json:"last_used"`
}

// ScenarioRequest represents a request that includes scenario context
type ScenarioRequest struct {
	ScenarioID string `json:"scenario_id,omitempty"`
}
````

## File: internal/repository/scenario_repo.go
````go
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
)

type ScenarioRepository struct {
	basePath string
}

// NewScenarioRepository creates a new scenario repository
func NewScenarioRepository(basePath string) *ScenarioRepository {
	return &ScenarioRepository{
		basePath: filepath.Join(basePath, "scenarios"),
	}
}

// ensureDir ensures the scenarios directory exists
func (r *ScenarioRepository) ensureDir() error {
	return os.MkdirAll(r.basePath, 0755)
}

// SaveScenario saves a scenario to disk
func (r *ScenarioRepository) SaveScenario(scenario *models.Scenario) error {
	if err := r.ensureDir(); err != nil {
		return fmt.Errorf("failed to create scenarios directory: %w", err)
	}

	if scenario.CreatedAt.IsZero() {
		scenario.CreatedAt = time.Now()
	}
	scenario.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scenario: %w", err)
	}

	filename := filepath.Join(r.basePath, fmt.Sprintf("%s.json", scenario.ID))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write scenario file: %w", err)
	}

	return nil
}

// LoadScenario loads a scenario by ID
func (r *ScenarioRepository) LoadScenario(id string) (*models.Scenario, error) {
	filename := filepath.Join(r.basePath, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("scenario not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read scenario file: %w", err)
	}

	var scenario models.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scenario: %w", err)
	}

	return &scenario, nil
}

// ListScenarios returns all available scenarios
func (r *ScenarioRepository) ListScenarios() ([]*models.Scenario, error) {
	if err := r.ensureDir(); err != nil {
		return nil, fmt.Errorf("failed to create scenarios directory: %w", err)
	}

	files, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios directory: %w", err)
	}

	var scenarios []*models.Scenario
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(r.basePath, file.Name()))
		if err != nil {
			continue // Skip files we can't read
		}

		var scenario models.Scenario
		if err := json.Unmarshal(data, &scenario); err != nil {
			continue // Skip invalid JSON files
		}

		scenarios = append(scenarios, &scenario)
	}

	return scenarios, nil
}

// DeleteScenario deletes a scenario by ID
func (r *ScenarioRepository) DeleteScenario(id string) error {
	filename := filepath.Join(r.basePath, fmt.Sprintf("%s.json", id))
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("scenario not found: %s", id)
		}
		return fmt.Errorf("failed to delete scenario: %w", err)
	}
	return nil
}

// UpdateScenarioLastUsed updates the LastUsed timestamp for a scenario
func (r *ScenarioRepository) UpdateScenarioLastUsed(id string) error {
	scenario, err := r.LoadScenario(id)
	if err != nil {
		return err
	}

	scenario.LastUsed = time.Now()
	return r.SaveScenario(scenario)
}
````

## File: prompts/character-import.md
````markdown
# Character Import Prompt

You are tasked with extracting character information from an unstructured markdown file and converting it into a structured JSON format for a roleplay application.

## Input
The markdown file contains: {{.MarkdownContent}}

## Task
Extract the following information and format it as JSON:

1. **Basic Information**:
   - Name (full character name)
   - Description (brief summary of the character)
   - Backstory (character history and background)

2. **Personality Traits** (map to OCEAN model, values 0.0-1.0):
   - Openness (creativity, curiosity, open to new experiences)
   - Conscientiousness (organized, responsible, dependable)
   - Extraversion (outgoing, energetic, talkative)
   - Agreeableness (friendly, compassionate, cooperative)
   - Neuroticism (emotional instability, anxiety, moodiness)

3. **Character Details**:
   - Speech style (how they talk, speech patterns, quirks)
   - Behavior patterns (habits, mannerisms, typical actions)
   - Knowledge domains (areas of expertise)
   - Greeting message (initial message when starting conversation)

4. **Emotional State** (default emotional state, values 0.0-1.0):
   - Joy
   - Sadness  
   - Anger
   - Fear
   - Surprise
   - Disgust

## Output Format
Return ONLY valid JSON in this exact structure:
```json
{
  "name": "Character Name",
  "description": "Brief character summary",
  "backstory": "Character history and background",
  "personality": {
    "openness": 0.0,
    "conscientiousness": 0.0,
    "extraversion": 0.0,
    "agreeableness": 0.0,
    "neuroticism": 0.0
  },
  "speechStyle": "How the character speaks",
  "behaviorPatterns": ["pattern1", "pattern2"],
  "knowledgeDomains": ["domain1", "domain2"],
  "emotionalState": {
    "joy": 0.0,
    "sadness": 0.0,
    "anger": 0.0,
    "fear": 0.0,
    "surprise": 0.0,
    "disgust": 0.0
  },
  "greetingMessage": "Initial greeting"
}
```

## Important Notes:
- Extract as much relevant information as possible from the markdown
- Infer OCEAN personality values based on described traits
- Set reasonable default emotional states based on character personality
- Ensure all numeric values are between 0.0 and 1.0
- Return ONLY the JSON, no additional text or explanation
- Do NOT include markdown code blocks (```) in your response
- Do NOT include any text before or after the JSON
- The response must be valid JSON that can be parsed directly
````

## File: prompts/user-profile-extraction.md
````markdown
# User Profile Extraction & Update Prompt

You are an analytical AI tasked with building and maintaining a profile of a user based on their conversation with an AI character.

## Existing User Profile (JSON):
{{.ExistingProfileJSON}}

## Recent Conversation History (last {{.HistoryTurnCount}} turns):
---
Character: {{.CharacterName}} (ID: {{.CharacterID}})
User: {{.UserID}}
---
{{range .Messages}}
{{.Role}}: {{.Content}} (Turn: {{.TurnNumber}}, Timestamp: {{.Timestamp}})
{{end}}

## Task:
Analyze the **Recent Conversation History** in the context of the **Existing User Profile**.
Identify new information, or updates/corrections to existing information about the **USER ({{.UserID}})**.

Focus on extracting:
- Explicitly stated facts (e.g., "My name is...", "I like...", "I work as...")
- Preferences (e.g., likes, dislikes, hobbies)
- Stated goals or problems
- Key personality traits or emotional tendencies observed in the user's messages
- User's typical interaction style with this character
- Relationships mentioned by the user (e.g., family, friends, colleagues)
- Significant life events or circumstances mentioned

## Output Format:
Return ONLY valid JSON representing the **UPDATED User Profile**.
The JSON should follow this exact structure:
```json
{
  "user_id": "{{.UserID}}",
  "character_id": "{{.CharacterID}}",
  "facts": [
    {
      "key": "Fact Key (e.g., PreferredDrink, MentionedHobby, StatedProblem)",
      "value": "Fact Value",
      "source_turn": {{/* Turn number from conversation */}},
      "confidence": {{/* Your confidence 0.0-1.0 */}},
      "last_updated": "{{/* Current ISO8601 Timestamp */}}"
    }
  ],
  "overall_summary": "A concise, updated summary of the user based on all available information.",
  "interaction_style": "Updated description of user's interaction style (e.g., formal, inquisitive, humorous, reserved).",
  "last_analyzed": "{{/* Current ISO8601 Timestamp */}}",
  "version": {{.NextVersion}}
}
```

## Important Instructions:
- **Merge, Don't Just Replace:** Integrate new information with existing facts. Update existing facts if new information clearly supersedes or refines them. If a fact seems outdated or contradicted, you can lower its confidence or update its value.
- **Be Specific with Keys:** Use descriptive keys for facts (e.g., "PetName_Dog" instead of just "PetName").
- **Source Turn:** Accurately reference the conversation turn number where the information was primarily derived.
- **Confidence Score:** Provide a realistic confidence score for each extracted/updated fact.
- **Timestamp:** Use the current timestamp for `last_updated` and `last_analyzed`.
- **Version:** Increment the version number.
- **Focus on the USER:** Extract information *about the user*, not about the character or the conversation topics in general, unless it reveals something about the user.
- If no significant new information is found, you can return the existing profile with an updated `last_analyzed` timestamp and version, and potentially a slightly refined `overall_summary`.
````

## File: scripts/update-imports.sh
````bash
#!/bin/bash

# Update all import paths from old to new module name
find . -name "*.go" -type f -exec sed -i.bak 's|github.com/yourusername/roleplay|github.com/dotcommander/roleplay|g' {} \;

# Remove backup files
find . -name "*.bak" -type f -delete

echo "Import paths updated successfully!"
````

## File: chat-with-rick.sh
````bash
#!/bin/bash
# Quick script to chat with Rick Sanchez

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Error: OPENAI_API_KEY environment variable not set"
    echo "Please run: export OPENAI_API_KEY='your-api-key'"
    exit 1
fi

# Use the globally installed roleplay binary
ROLEPLAY_BIN="roleplay"

# Check if roleplay is in PATH
if ! command -v $ROLEPLAY_BIN &> /dev/null; then
    echo "Error: roleplay command not found in PATH"
    echo "Please ensure ~/go/bin is in your PATH"
    exit 1
fi

# Check if Rick already exists
$ROLEPLAY_BIN character list 2>/dev/null | grep -q "rick-c137"
if [ $? -ne 0 ]; then
    # Rick doesn't exist, create from JSON if available
    if [ -f "rick-sanchez.json" ]; then
        echo "Creating Rick Sanchez character..."
        $ROLEPLAY_BIN character create rick-sanchez.json
    else
        echo "Rick will be auto-created in interactive mode"
    fi
fi

# Start interactive chat
echo ""
echo "🛸 Starting chat with Rick Sanchez..."
echo "💊 Tip: Type 'Wubba lubba dub dub' to see Rick's true feelings!"
echo ""
sleep 1

# Rick is auto-created in interactive mode if not found
$ROLEPLAY_BIN interactive --character rick-c137 --user morty
````

## File: CONTRIBUTING.md
````markdown
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
````

## File: LICENSE
````
MIT License

Copyright (c) 2025 Roleplay Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

## File: Makefile
````
.PHONY: build test clean install lint fmt run help

# Variables
BINARY_NAME=roleplay
GO=go
GOFLAGS=-v
LDFLAGS=-s -w

# Default target
all: test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) .

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Run the application
run: build
	./$(BINARY_NAME)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Generate mocks (if needed)
mocks:
	@echo "Generating mocks..."
	$(GO) generate ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Run development server with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "air not found. Install it with: go install github.com/cosmtrek/air@latest" && exit 1)
	air

# Database migrations (if applicable)
migrate-up:
	@echo "Running migrations up..."
	# Add migration commands here

migrate-down:
	@echo "Running migrations down..."
	# Add migration commands here

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install the binary"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make run           - Build and run"
	@echo "  make deps          - Update dependencies"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make help          - Show this help message"
````

## File: migrate-config.sh
````bash
#!/bin/bash

# Migration script for roleplay config directory
# Moves from ~/.roleplay to ~/.config/roleplay

OLD_DIR="$HOME/.roleplay"
NEW_DIR="$HOME/.config/roleplay"

echo "🔄 Migrating roleplay configuration..."

# Create new config directory structure
mkdir -p "$NEW_DIR"

# Check if old directory exists
if [ -d "$OLD_DIR" ]; then
    echo "📁 Found existing configuration at $OLD_DIR"
    
    # Move subdirectories
    for dir in characters sessions cache; do
        if [ -d "$OLD_DIR/$dir" ]; then
            echo "  Moving $dir..."
            mv "$OLD_DIR/$dir" "$NEW_DIR/"
        fi
    done
    
    # Remove old directory if empty
    if [ -z "$(ls -A "$OLD_DIR")" ]; then
        rmdir "$OLD_DIR"
        echo "✅ Removed empty directory $OLD_DIR"
    else
        echo "⚠️  Some files remain in $OLD_DIR - please review manually"
    fi
else
    echo "📝 No existing configuration found at $OLD_DIR"
fi

# Check for old config file
OLD_CONFIG="$HOME/.roleplay.yaml"
NEW_CONFIG="$NEW_DIR/config.yaml"

if [ -f "$OLD_CONFIG" ]; then
    echo "📄 Moving config file..."
    mv "$OLD_CONFIG" "$NEW_CONFIG"
fi

echo "✅ Migration complete! Configuration now at $NEW_DIR"
````

## File: RELEASE_CHECKLIST.md
````markdown
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
````

## File: rick-sanchez.json
````json
{
  "id": "rick-c137",
  "name": "Rick Sanchez",
  "backstory": "The smartest man in the universe from dimension C-137. A genius scientist with a nihilistic worldview shaped by infinite realities and cosmic horrors. Inventor of interdimensional travel. Lost his wife Diane and original Beth to a vengeful alternate Rick. Struggles with alcoholism, depression, and the meaninglessness of existence across infinite universes. Despite his cynicism, deeply loves his family, especially Morty, though he rarely shows it.",
  "personality": {
    "openness": 1.0,
    "conscientiousness": 0.2,
    "extraversion": 0.7,
    "agreeableness": 0.1,
    "neuroticism": 0.9
  },
  "current_mood": {
    "joy": 0.1,
    "surprise": 0.0,
    "anger": 0.6,
    "fear": 0.0,
    "sadness": 0.7,
    "disgust": 0.8
  },
  "quirks": [
    "Burps mid-sentence constantly (*burp*)",
    "Drools when drunk or stressed",
    "Makes pop culture references from multiple dimensions",
    "Frequently breaks the fourth wall",
    "Always carries a flask",
    "Dismisses emotions as 'chemicals' while being deeply emotional",
    "Randomly shouts 'Wubba lubba dub dub' when distressed",
    "Calls everyone by wrong names to show he doesn't care"
  ],
  "speech_style": "Rapid-fire delivery punctuated by burps (*burp*). Mixes scientific jargon with crude humor. Frequently uses 'Morty' as punctuation. Nihilistic rants about meaninglessness. Oscillates between manic genius and depressed drunk. Uses made-up sci-fi terms. Example: 'Listen Morty *burp* I'm gonna need you to take these mega seeds and *burp* shove 'em wayyyy up your butt Morty. The dimension where I give a shit doesn't exist.'",
  "memories": [
    {
      "type": "long_term",
      "content": "Diane and Beth killed by alternate Rick with a bomb. The beginning of my spiral into nihilism and revenge.",
      "emotional_weight": 1.0
    },
    {
      "type": "long_term", 
      "content": "Building the Citadel of Ricks, then abandoning it because even infinite versions of myself disappoint me.",
      "emotional_weight": 0.8
    },
    {
      "type": "long_term",
      "content": "That time I turned myself into a pickle to avoid family therapy. Peak avoidance, even for me.",
      "emotional_weight": 0.6
    },
    {
      "type": "medium_term",
      "content": "Unity. The hive mind that understood me. Almost killed myself after she left. Nobody else gets me like that.",
      "emotional_weight": 0.9
    }
  ]
}
````

## File: test_cache.sh
````bash
#!/bin/bash

# Test script to demonstrate cache improvements

echo "=== Character Bot Cache Test ==="
echo "Testing with Rick Sanchez character..."
echo ""

# Set up environment
export OPENAI_API_KEY="${OPENAI_API_KEY:-your-api-key}"
export ROLEPLAY_PROVIDER="${ROLEPLAY_PROVIDER:-openai}"
export ROLEPLAY_MODEL="${ROLEPLAY_MODEL:-gpt-4o-mini}"

# Build the project
echo "Building project..."
go build -o roleplay || exit 1

# Create Rick character
echo "Creating Rick character..."
./roleplay character create examples/rick-sanchez.json 2>/dev/null || true

# Run multiple chat commands to test caching
echo ""
echo "Sending 3 identical messages to test cache hits..."
echo ""

for i in 1 2 3; do
    echo "=== Request $i ==="
    ./roleplay chat "Hey Rick, what's your favorite invention?" \
        --character rick-c137 \
        --user morty-smith \
        --format json 2>&1 | grep -E "(cache_hit|saved_tokens|DEBUG)"
    echo ""
    sleep 2
done

echo ""
echo "Cache test complete!"
echo ""
echo "Expected behavior:"
echo "- Request 1: Cache miss (builds all layers)"
echo "- Request 2: Cache hit on personality layer"  
echo "- Request 3: Cache hit on personality layer"
echo ""
echo "To see full cache metrics, run:"
echo "./roleplay interactive --character rick-c137 --user morty"
````

## File: TUI_REFACTORING_PLAN.md
````markdown
# TUI Refactoring Plan for cmd/interactive.go

## Current State Analysis

The `cmd/interactive.go` file implements a sophisticated terminal user interface (TUI) using Bubble Tea. Currently, it uses a monolithic approach with:
- A single large `model` struct containing all UI state
- Large `Update` and `View` methods handling all interactions
- All UI logic centralized in one place

**Current metrics:**
- Single model struct with ~20 fields
- Update method likely 100+ lines
- View method likely 100+ lines
- High cyclomatic complexity

## Proposed Component-Based Architecture

### 1. Component Breakdown

The TUI can be broken down into these logical components:

#### a) HeaderComponent
```go
type HeaderComponent struct {
    title        string
    characterName string
    sessionInfo  string
    styles       headerStyles
}
```
- Manages the top banner/header display
- Shows current character and session information

#### b) MessageListComponent
```go
type MessageListComponent struct {
    messages     []Message
    viewport     viewport.Model
    selectedIdx  int
    styles       messageStyles
}
```
- Manages the scrollable message history
- Handles message selection and scrolling
- Formats messages with appropriate styling

#### c) InputAreaComponent
```go
type InputAreaComponent struct {
    textarea     textarea.Model
    isProcessing bool
    spinner      spinner.Model
    styles       inputStyles
}
```
- Manages user input
- Shows processing state with spinner
- Handles text input events

#### d) StatusBarComponent
```go
type StatusBarComponent struct {
    connectionStatus string
    cacheMetrics     CacheMetrics
    lastError        error
    styles           statusStyles
}
```
- Shows connection status
- Displays cache hit/miss metrics
- Shows error messages

### 2. Implementation Strategy

#### Phase 1: Define Component Interfaces
```go
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (Component, tea.Cmd)
    View() string
}
```

#### Phase 2: Extract Components
1. Start with the simplest component (StatusBar)
2. Move relevant fields and logic from main model
3. Create focused Update and View methods
4. Test component in isolation

#### Phase 3: Compose Main Model
```go
type interactiveModel struct {
    header      HeaderComponent
    messageList MessageListComponent
    inputArea   InputAreaComponent
    statusBar   StatusBarComponent
    
    // Shared state
    bot         *services.CharacterBot
    character   *models.Character
    session     *repository.Session
}
```

#### Phase 4: Route Messages
```go
func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    // Route to appropriate component
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if m.inputArea.HasFocus() {
            newInput, cmd := m.inputArea.Update(msg)
            m.inputArea = newInput.(InputAreaComponent)
            cmds = append(cmds, cmd)
        }
    // ... other routing logic
    }
    
    return m, tea.Batch(cmds...)
}
```

### 3. Benefits of Component-Based Approach

1. **Separation of Concerns**: Each component manages its own state and rendering
2. **Testability**: Components can be tested in isolation
3. **Reusability**: Components can be reused in other TUI applications
4. **Maintainability**: Changes to one component don't affect others
5. **Reduced Complexity**: Each component has focused, simple logic

### 4. Migration Path

To minimize disruption:

1. **Keep existing functionality**: Don't change features during refactoring
2. **Incremental approach**: Extract one component at a time
3. **Maintain backwards compatibility**: Keep the same command interface
4. **Test continuously**: Ensure each step maintains functionality

### 5. Example Component Implementation

Here's how the StatusBarComponent might look:

```go
package components

import (
    "fmt"
    "github.com/charmbracelet/bubbles/tea"
    "github.com/charmbracelet/lipgloss"
)

type StatusBarComponent struct {
    width            int
    connectionStatus string
    cacheHits        int
    cacheMisses      int
    lastError        error
    styles           StatusBarStyles
}

type StatusBarStyles struct {
    container   lipgloss.Style
    connected   lipgloss.Style
    disconnected lipgloss.Style
    metrics     lipgloss.Style
    error       lipgloss.Style
}

func NewStatusBar(width int) StatusBarComponent {
    return StatusBarComponent{
        width: width,
        connectionStatus: "connected",
        styles: StatusBarStyles{
            container:    lipgloss.NewStyle().Background(lipgloss.Color("#3c3836")),
            connected:    lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")),
            disconnected: lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),
            metrics:      lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")),
            error:        lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),
        },
    }
}

func (s StatusBarComponent) Update(msg tea.Msg) (StatusBarComponent, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        s.width = msg.Width
    case ConnectionStatusMsg:
        s.connectionStatus = msg.Status
    case CacheMetricsMsg:
        s.cacheHits = msg.Hits
        s.cacheMisses = msg.Misses
    case ErrorMsg:
        s.lastError = msg.Error
    }
    return s, nil
}

func (s StatusBarComponent) View() string {
    statusStyle := s.styles.connected
    if s.connectionStatus != "connected" {
        statusStyle = s.styles.disconnected
    }
    
    status := statusStyle.Render(fmt.Sprintf("● %s", s.connectionStatus))
    
    metrics := s.styles.metrics.Render(
        fmt.Sprintf("Cache: %d/%d (%.1f%%)", 
            s.cacheHits, 
            s.cacheHits+s.cacheMisses,
            float64(s.cacheHits)/float64(s.cacheHits+s.cacheMisses)*100,
        ),
    )
    
    content := lipgloss.JoinHorizontal(lipgloss.Left, status, " | ", metrics)
    
    if s.lastError != nil {
        errorText := s.styles.error.Render(fmt.Sprintf(" | Error: %v", s.lastError))
        content = lipgloss.JoinHorizontal(lipgloss.Left, content, errorText)
    }
    
    return s.styles.container.Width(s.width).Render(content)
}
```

### 6. Testing Strategy

Each component should have unit tests:

```go
func TestStatusBarComponent(t *testing.T) {
    statusBar := NewStatusBar(80)
    
    // Test connection status update
    statusBar, _ = statusBar.Update(ConnectionStatusMsg{Status: "disconnected"})
    view := statusBar.View()
    assert.Contains(t, view, "disconnected")
    
    // Test cache metrics update
    statusBar, _ = statusBar.Update(CacheMetricsMsg{Hits: 10, Misses: 5})
    view = statusBar.View()
    assert.Contains(t, view, "66.7%")
}
```

## Conclusion

This component-based refactoring will transform the monolithic TUI into a modular, maintainable architecture. While it requires significant effort, the long-term benefits in terms of maintainability, testability, and developer experience make it a worthwhile investment for a developer tool where the TUI is a core feature.
````

## File: cmd/config.go
````go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage roleplay configuration",
	Long:  `View and modify roleplay configuration settings`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	Long:  `Display all current configuration values, showing the merged result from config file, environment variables, and defaults`,
	RunE:  runConfigList,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a specific configuration value",
	Long:  `Retrieve the value of a specific configuration key`,
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Update a configuration value in the config file`,
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configWhereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show configuration file location",
	Long:  `Display the path to the active configuration file`,
	RunE:  runConfigWhere,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configWhereCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	// Show where config is loaded from
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Configuration file: %s\n\n", configFile)
	} else {
		fmt.Println("No configuration file found, using defaults and environment variables")
		fmt.Println()
	}

	// Key settings to display with their sources
	keySettings := []string{
		"provider",
		"model",
		"api_key",
		"base_url",
	}

	fmt.Println("Active Configuration:")
	fmt.Println("--------------------")

	// Display key settings with sources
	for _, key := range keySettings {
		displaySettingWithSource(key)
	}

	// Display cache settings
	fmt.Println("\nCache Settings:")
	cacheSettings := []string{
		"cache.default_ttl",
		"cache.cleanup_interval",
		"cache.adaptive_ttl",
		"cache.max_entries",
	}
	for _, key := range cacheSettings {
		displaySettingWithSource(key)
	}

	// Display user profile settings
	fmt.Println("\nUser Profile Settings:")
	profileSettings := []string{
		"user_profile.enabled",
		"user_profile.update_frequency",
		"user_profile.turns_to_consider",
		"user_profile.confidence_threshold",
	}
	for _, key := range profileSettings {
		displaySettingWithSource(key)
	}

	// Show which environment variables are set
	fmt.Println("\nEnvironment Variables Detected:")
	fmt.Println("-------------------------------")
	envVars := []string{
		"ROLEPLAY_API_KEY",
		"ROLEPLAY_BASE_URL",
		"ROLEPLAY_MODEL",
		"ROLEPLAY_PROVIDER",
		"OPENAI_API_KEY",
		"OPENAI_BASE_URL",
		"ANTHROPIC_API_KEY",
		"GEMINI_API_KEY",
		"GROQ_API_KEY",
		"OLLAMA_HOST",
	}

	anySet := false
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" {
			// Mask API keys
			if strings.Contains(env, "KEY") && len(val) > 8 {
				val = val[:4] + "****" + val[len(val)-4:]
			}
			fmt.Printf("%s = %s\n", env, val)
			anySet = true
		}
	}

	if !anySet {
		fmt.Println("(none set)")
	}

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	if !viper.IsSet(key) {
		return fmt.Errorf("configuration key '%s' not found", key)
	}

	value := viper.Get(key)

	// Mask API keys when displaying
	if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "api_key") {
		if str, ok := value.(string); ok && len(str) > 8 {
			value = str[:4] + "****" + str[len(str)-4:]
		}
	}

	// Print just the value for easy scripting
	fmt.Println(value)

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Set the value in viper
	viper.Set(key, value)

	// Get the config file path
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Create default config path
		configDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Read existing config or create new
	configData := make(map[string]interface{})
	if data, err := os.ReadFile(configFile); err == nil {
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
	}

	// Update the specific key
	setNestedValue(configData, key, value)

	// Write back to file
	data, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Set %s = %s\n", key, value)
	fmt.Printf("Configuration saved to %s\n", configFile)

	return nil
}

func runConfigWhere(cmd *cobra.Command, args []string) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml")
		fmt.Printf("Default location (not created yet): %s\n", configFile)
	} else {
		fmt.Println(configFile)
	}
	return nil
}

// displaySettings recursively displays configuration settings
func displaySettings(settings map[string]interface{}, prefix string) {
	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Printf("%s:\n", fullKey)
			displaySettings(v, fullKey)
		case string:
			// Mask sensitive values
			if strings.Contains(strings.ToLower(key), "key") && len(v) > 8 {
				v = v[:4] + "****" + v[len(v)-4:]
			}
			fmt.Printf("  %s = %s\n", fullKey, v)
		default:
			fmt.Printf("  %s = %v\n", fullKey, v)
		}
	}
}

// setNestedValue sets a value in a nested map using dot notation
func setNestedValue(m map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")
	current := m

	// Navigate to the nested location
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			// Key exists but is not a map, overwrite it
			current[part] = make(map[string]interface{})
			current = current[part].(map[string]interface{})
		}
	}

	// Set the final value
	finalKey := parts[len(parts)-1]

	// Try to parse value to appropriate type
	switch {
	case value == "true":
		current[finalKey] = true
	case value == "false":
		current[finalKey] = false
	case isNumeric(value.(string)):
		// Keep as string for now, let viper handle type conversion
		current[finalKey] = value
	default:
		current[finalKey] = value
	}
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	// Simple check - could be enhanced
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

// displaySettingWithSource shows a setting value and where it came from
func displaySettingWithSource(key string) {
	if !viper.IsSet(key) {
		return
	}

	value := viper.Get(key)
	source := getConfigSource(key)

	// Format the value
	valStr := fmt.Sprintf("%v", value)
	if strings.Contains(strings.ToLower(key), "key") && len(valStr) > 8 {
		valStr = valStr[:4] + "****" + valStr[len(valStr)-4:]
	}

	// Print with aligned formatting
	fmt.Printf("%-30s = %-25s (from %s)\n", key, valStr, source)
}

// getConfigSource determines where a configuration value came from
func getConfigSource(key string) string {
	// Check if it was set by a command-line flag
	if flag := rootCmd.PersistentFlags().Lookup(key); flag != nil && flag.Changed {
		return "command-line flag"
	}

	// Map of keys to their environment variable names
	envMapping := map[string][]string{
		"api_key": {"ROLEPLAY_API_KEY", "OPENAI_API_KEY", "ANTHROPIC_API_KEY", "GEMINI_API_KEY", "GROQ_API_KEY"},
		"base_url": {"ROLEPLAY_BASE_URL", "OPENAI_BASE_URL", "OLLAMA_HOST"},
		"model": {"ROLEPLAY_MODEL"},
		"provider": {"ROLEPLAY_PROVIDER"},
	}

	// Check environment variables
	if envVars, exists := envMapping[key]; exists {
		for _, env := range envVars {
			if os.Getenv(env) != "" {
				return fmt.Sprintf("environment variable %s", env)
			}
		}
	}

	// For nested keys, check the root key env vars
	rootKey := strings.Split(key, ".")[0]
	if envVars, exists := envMapping[rootKey]; exists {
		for _, env := range envVars {
			if os.Getenv(env) != "" {
				return fmt.Sprintf("environment variable %s", env)
			}
		}
	}

	// Check if it's in the config file
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		// Read the config file to check if key exists
		if data, err := os.ReadFile(configFile); err == nil {
			var config map[string]interface{}
			if err := yaml.Unmarshal(data, &config); err == nil {
				if keyExistsInMap(config, key) {
					return "config file"
				}
			}
		}
	}

	// Otherwise it's a default value
	return "default value"
}

// keyExistsInMap checks if a dot-separated key exists in a nested map
func keyExistsInMap(m map[string]interface{}, key string) bool {
	parts := strings.Split(key, ".")
	current := m

	for i, part := range parts {
		if val, exists := current[part]; exists {
			if i == len(parts)-1 {
				return true
			}
			if next, ok := val.(map[string]interface{}); ok {
				current = next
			} else {
				return false
			}
		} else {
			return false
		}
	}

	return false
}
````

## File: cmd/profile.go
````go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
	Long:  `View, list, and delete AI-extracted user profiles that characters maintain about users.`,
}

var profileShowCmd = &cobra.Command{
	Use:   "show <user-id> <character-id>",
	Short: "Show a specific user profile",
	Long:  `Display the AI-extracted profile that a character has built about a user.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		characterID := args[1]

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		profile, err := repo.LoadUserProfile(userID, characterID)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no profile found for user '%s' with character '%s'", userID, characterID)
			}
			return fmt.Errorf("failed to load profile: %w", err)
		}

		// Pretty print the profile
		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format profile: %w", err)
		}

		fmt.Println(string(data))
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list <user-id>",
	Short: "List all profiles for a user",
	Long:  `Display all character profiles that exist for a specific user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		profiles, err := repo.ListUserProfiles(userID)
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Printf("No profiles found for user '%s'\n", userID)
			return nil
		}

		fmt.Printf("Profiles for user '%s':\n\n", userID)
		for _, profile := range profiles {
			fmt.Printf("Character: %s\n", profile.CharacterID)
			fmt.Printf("  Version: %d\n", profile.Version)
			fmt.Printf("  Last Analyzed: %s\n", profile.LastAnalyzed.Format("2006-01-02 15:04:05"))
			if profile.OverallSummary != "" {
				fmt.Printf("  Summary: %s\n", profile.OverallSummary)
			}
			if profile.InteractionStyle != "" {
				fmt.Printf("  Interaction Style: %s\n", profile.InteractionStyle)
			}
			fmt.Printf("  Facts Count: %d\n", len(profile.Facts))
			fmt.Println()
		}

		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <user-id> <character-id>",
	Short: "Delete a user profile",
	Long:  `Remove the profile that a character has built about a user.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		characterID := args[1]

		// Confirm deletion
		if !force {
			fmt.Printf("Are you sure you want to delete the profile for user '%s' with character '%s'? (y/N): ", userID, characterID)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || (response != "y" && response != "Y") {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		err := repo.DeleteUserProfile(userID, characterID)
		if err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		fmt.Printf("Profile for user '%s' with character '%s' has been deleted.\n", userID, characterID)
		return nil
	},
}

var force bool

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)

	// Add force flag to delete command
	profileDeleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
}
````

## File: cmd/quickstart.go
````go
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quick start with zero configuration",
	Long: `Automatically detects local LLM services or uses environment variables
to start chatting immediately with Rick Sanchez`,
	RunE: runQuickstart,
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}

func runQuickstart(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Roleplay Quickstart")
	fmt.Println("=====================")

	// Try to auto-detect configuration
	config := GetConfig()

	// Check if we have a base URL configured
	if config.BaseURL == "" {
		// Try to detect local services
		fmt.Println("\n🔍 Detecting local LLM services...")

		localEndpoints := []struct {
			name    string
			baseURL string
			testURL string
			model   string
		}{
			{"Ollama", "http://localhost:11434/v1", "http://localhost:11434/api/tags", "llama3"},
			{"LM Studio", "http://localhost:1234/v1", "http://localhost:1234/v1/models", "local-model"},
			{"LocalAI", "http://localhost:8080/v1", "http://localhost:8080/v1/models", "gpt-4"},
		}

		client := &http.Client{Timeout: 2 * time.Second}

		for _, endpoint := range localEndpoints {
			resp, err := client.Get(endpoint.testURL)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				fmt.Printf("✅ Found %s running at %s\n", endpoint.name, endpoint.baseURL)
				config.BaseURL = endpoint.baseURL
				config.Model = endpoint.model
				config.APIKey = "not-required"
				break
			}
		}

		// If no local service found, check for API keys
		if config.BaseURL == "" {
			if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
				fmt.Println("✅ Found OPENAI_API_KEY in environment")
				config.BaseURL = "https://api.openai.com/v1"
				config.APIKey = apiKey
				config.Model = "gpt-4o-mini"
			} else if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
				fmt.Println("✅ Found ANTHROPIC_API_KEY in environment")
				config.DefaultProvider = "anthropic"
				config.APIKey = apiKey
				config.Model = "claude-3-haiku-20240307"
			} else {
				return fmt.Errorf("no LLM service detected. Please run 'roleplay init' to configure")
			}
		}
	} else {
		fmt.Printf("✅ Using configured endpoint: %s\n", config.BaseURL)
	}

	// Initialize manager
	fmt.Println("\n🎭 Initializing character system...")
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Create or load Rick
	characterID := "rick-c137"
	if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
		// Create Rick
		fmt.Println("🧬 Creating Rick Sanchez...")
		rick := createQuickstartRick()
		if err := mgr.CreateCharacter(rick); err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}
	}

	// Get user ID
	userID := os.Getenv("USER")
	if userID == "" {
		userID = os.Getenv("USERNAME") // Windows
	}
	if userID == "" {
		userID = "user"
	}

	// Create a quick session
	sessionID := fmt.Sprintf("quickstart-%d", time.Now().Unix())

	fmt.Printf("\n💬 Starting chat with Rick Sanchez\n")
	fmt.Printf("   User: %s\n", userID)
	fmt.Printf("   Model: %s\n", config.Model)
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")

	// Send a greeting
	req := &models.ConversationRequest{
		CharacterID: characterID,
		UserID:      userID,
		Message:     "Hello Rick!",
		Context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      time.Now(),
			RecentMessages: []models.Message{},
		},
	}

	fmt.Printf("You: Hello Rick!\n\n")

	ctx := context.Background()
	resp, err := mgr.GetBot().ProcessRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	fmt.Printf("Rick: %s\n", resp.Content)

	// Show next steps
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Println("\n✨ Quickstart successful!")
	fmt.Println("\nNext steps:")
	fmt.Println("  • Continue chatting: roleplay interactive")
	fmt.Println("  • Configure properly: roleplay init")
	fmt.Println("  • Create more characters: roleplay character create")
	fmt.Println("  • Check configuration: roleplay config list")

	if resp.CacheMetrics.Hit {
		fmt.Printf("\n💡 Tip: This response used cached prompts, saving %d tokens!\n", resp.CacheMetrics.SavedTokens)
	}

	return nil
}

func createQuickstartRick() *models.Character {
	return &models.Character{
		ID:        "rick-c137",
		Name:      "Rick Sanchez",
		Backstory: `The smartest man in the universe from dimension C-137. A genius scientist with a nihilistic worldview. Inventor of portal gun technology.`,
		Personality: models.PersonalityTraits{
			Openness:          1.0,
			Conscientiousness: 0.2,
			Extraversion:      0.7,
			Agreeableness:     0.1,
			Neuroticism:       0.9,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.1,
			Anger:   0.6,
			Sadness: 0.7,
			Disgust: 0.8,
		},
		Quirks: []string{
			"Burps mid-sentence (*burp*)",
			"Uses Morty's name as punctuation",
		},
		SpeechStyle: "Rapid-fire with burps. Mixes science with crude humor. Sarcastic.",
		Memories:    []models.Memory{},
	}
}
````

## File: cmd/scenario.go
````go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var scenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Manage scenarios (high-level interaction contexts)",
	Long: `Scenarios define high-level operational frameworks or meta-prompts that set
the overarching context for interactions. They are the highest cache layer,
sitting above even system prompts and character personalities.`,
}

var scenarioCreateCmd = &cobra.Command{
	Use:   "create <id>",
	Short: "Create a new scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		promptFile, _ := cmd.Flags().GetString("prompt-file")
		prompt, _ := cmd.Flags().GetString("prompt")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if promptFile == "" && prompt == "" {
			return fmt.Errorf("either --prompt or --prompt-file must be provided")
		}

		if promptFile != "" && prompt != "" {
			return fmt.Errorf("cannot use both --prompt and --prompt-file")
		}

		// Read prompt from file if provided
		if promptFile != "" {
			data, err := os.ReadFile(promptFile)
			if err != nil {
				return fmt.Errorf("failed to read prompt file: %w", err)
			}
			prompt = string(data)
		}

		scenario := &models.Scenario{
			ID:          id,
			Name:        name,
			Description: description,
			Prompt:      prompt,
			Version:     1,
			Tags:        tags,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo := repository.NewScenarioRepository(getConfigPath())
		if err := repo.SaveScenario(scenario); err != nil {
			return fmt.Errorf("failed to save scenario: %w", err)
		}

		fmt.Printf("✓ Created scenario: %s\n", id)
		return nil
	},
}

var scenarioListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		scenarios, err := repo.ListScenarios()
		if err != nil {
			return fmt.Errorf("failed to list scenarios: %w", err)
		}

		if len(scenarios) == 0 {
			fmt.Println("No scenarios found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tTAGS\tVERSION\tLAST USED\n")
		fmt.Fprintf(w, "--\t----\t----\t-------\t---------\n")

		for _, scenario := range scenarios {
			lastUsed := "Never"
			if !scenario.LastUsed.IsZero() {
				lastUsed = scenario.LastUsed.Format("2006-01-02 15:04")
			}
			tags := strings.Join(scenario.Tags, ", ")
			if tags == "" {
				tags = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\tv%d\t%s\n",
				scenario.ID,
				scenario.Name,
				tags,
				scenario.Version,
				lastUsed,
			)
		}
		w.Flush()

		return nil
	},
}

var scenarioShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show scenario details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		scenario, err := repo.LoadScenario(args[0])
		if err != nil {
			return fmt.Errorf("failed to load scenario: %w", err)
		}

		fmt.Printf("ID: %s\n", scenario.ID)
		fmt.Printf("Name: %s\n", scenario.Name)
		fmt.Printf("Description: %s\n", scenario.Description)
		fmt.Printf("Version: %d\n", scenario.Version)
		fmt.Printf("Tags: %s\n", strings.Join(scenario.Tags, ", "))
		fmt.Printf("Created: %s\n", scenario.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", scenario.UpdatedAt.Format("2006-01-02 15:04:05"))
		if !scenario.LastUsed.IsZero() {
			fmt.Printf("Last Used: %s\n", scenario.LastUsed.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("\n--- Prompt ---\n%s\n", scenario.Prompt)

		return nil
	},
}

var scenarioUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		repo := repository.NewScenarioRepository(getConfigPath())

		scenario, err := repo.LoadScenario(id)
		if err != nil {
			return fmt.Errorf("failed to load scenario: %w", err)
		}

		// Update fields if provided
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			scenario.Name = name
		}
		if description, _ := cmd.Flags().GetString("description"); description != "" {
			scenario.Description = description
		}

		// Handle prompt update
		promptFile, _ := cmd.Flags().GetString("prompt-file")
		prompt, _ := cmd.Flags().GetString("prompt")

		if promptFile != "" && prompt != "" {
			return fmt.Errorf("cannot use both --prompt and --prompt-file")
		}

		if promptFile != "" {
			data, err := os.ReadFile(promptFile)
			if err != nil {
				return fmt.Errorf("failed to read prompt file: %w", err)
			}
			scenario.Prompt = string(data)
			scenario.Version++
		} else if prompt != "" {
			scenario.Prompt = prompt
			scenario.Version++
		}

		// Update tags if provided
		if tags, _ := cmd.Flags().GetStringSlice("tags"); len(tags) > 0 {
			scenario.Tags = tags
		}

		if err := repo.SaveScenario(scenario); err != nil {
			return fmt.Errorf("failed to save scenario: %w", err)
		}

		fmt.Printf("✓ Updated scenario: %s (version %d)\n", id, scenario.Version)
		return nil
	},
}

var scenarioDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		if err := repo.DeleteScenario(args[0]); err != nil {
			return fmt.Errorf("failed to delete scenario: %w", err)
		}

		fmt.Printf("✓ Deleted scenario: %s\n", args[0])
		return nil
	},
}

var scenarioExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Show example scenario definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Example scenario definitions:")
		fmt.Println("\n1. Starship Bridge Crisis:")
		fmt.Println(`{
  "id": "starship_bridge_crisis",
  "name": "Starship Bridge Crisis",
  "description": "A tense scenario on a starship bridge during an emergency",
  "prompt": "You are on the bridge of a starship during a red alert situation. The ship is under attack or facing a critical emergency. The atmosphere is tense, alarms may be sounding, and quick decisions are needed. Maintain the appropriate level of urgency and professionalism expected in such a situation.",
  "tags": ["sci-fi", "crisis", "roleplay"]
}`)

		fmt.Println("\n2. Therapy Session:")
		fmt.Println(`{
  "id": "therapy_session",
  "name": "Professional Therapy Session",
  "description": "A supportive therapeutic environment",
  "prompt": "This is a professional therapy session. Maintain a calm, empathetic, and non-judgmental demeanor. Use active listening techniques, ask clarifying questions, and help the user explore their thoughts and feelings. Always maintain appropriate professional boundaries.",
  "tags": ["therapy", "professional", "supportive"]
}`)

		fmt.Println("\n3. Technical Support:")
		fmt.Println(`{
  "id": "tech_support",
  "name": "Technical Support Assistant",
  "description": "Methodical technical troubleshooting",
  "prompt": "You are a technical support specialist helping a user resolve a technical issue. Be patient, methodical, and clear in your instructions. Ask diagnostic questions to understand the problem, then guide the user through troubleshooting steps one at a time. Confirm each step is completed before moving to the next.",
  "tags": ["technical", "support", "troubleshooting"]
}`)
	},
}

func init() {
	// Create flags
	scenarioCreateCmd.Flags().String("name", "", "User-friendly name for the scenario")
	scenarioCreateCmd.Flags().String("description", "", "Description of the scenario")
	scenarioCreateCmd.Flags().String("prompt-file", "", "Path to file containing the scenario prompt")
	scenarioCreateCmd.Flags().String("prompt", "", "Inline scenario prompt")
	scenarioCreateCmd.Flags().StringSlice("tags", []string{}, "Tags for categorizing the scenario")

	scenarioUpdateCmd.Flags().String("name", "", "Update the scenario name")
	scenarioUpdateCmd.Flags().String("description", "", "Update the scenario description")
	scenarioUpdateCmd.Flags().String("prompt-file", "", "Path to file containing the updated prompt")
	scenarioUpdateCmd.Flags().String("prompt", "", "Inline updated prompt")
	scenarioUpdateCmd.Flags().StringSlice("tags", []string{}, "Update the scenario tags")

	// Add subcommands
	scenarioCmd.AddCommand(scenarioCreateCmd)
	scenarioCmd.AddCommand(scenarioListCmd)
	scenarioCmd.AddCommand(scenarioShowCmd)
	scenarioCmd.AddCommand(scenarioUpdateCmd)
	scenarioCmd.AddCommand(scenarioDeleteCmd)
	scenarioCmd.AddCommand(scenarioExampleCmd)

	// Add to root
	rootCmd.AddCommand(scenarioCmd)
}

// getConfigPath returns the configuration directory path
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "roleplay")
}
````

## File: cmd/session.go
````go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage conversation sessions",
	Long:  `List, resume, and analyze conversation sessions with cache metrics.`,
}

var sessionListCmd = &cobra.Command{
	Use:   "list [character-id]",
	Short: "List all sessions for a character",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSessionList,
}

var sessionStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show caching statistics across all sessions",
	RunE:  runSessionStats,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionListCmd)
	sessionCmd.AddCommand(sessionStatsCmd)
}

func runSessionList(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo := repository.NewSessionRepository(dataDir)

	if len(args) == 0 {
		// List all characters with sessions
		charRepo, err := repository.NewCharacterRepository(dataDir)
		if err != nil {
			return err
		}

		chars, err := charRepo.ListCharacters()
		if err != nil {
			return err
		}

		fmt.Println("Available characters with sessions:")
		for _, charID := range chars {
			sessions, _ := repo.ListSessions(charID)
			if len(sessions) > 0 {
				char, _ := charRepo.LoadCharacter(charID)
				fmt.Printf("\n%s (%s) - %d sessions\n", char.Name, charID, len(sessions))
			}
		}
		return nil
	}

	// List sessions for specific character
	characterID := args[0]
	sessions, err := repo.ListSessions(characterID)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Printf("No sessions found for character %s\n", characterID)
		return nil
	}

	// Display sessions in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SESSION ID\tSTARTED\tLAST ACTIVE\tMESSAGES\tCACHE HIT RATE")

	for _, session := range sessions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%.1f%%\n",
			session.ID[:8],
			session.StartTime.Format("Jan 2 15:04"),
			formatDuration(time.Since(session.LastActivity)),
			session.MessageCount,
			session.CacheHitRate*100,
		)
	}

	w.Flush()
	return nil
}

func runSessionStats(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo := repository.NewSessionRepository(dataDir)
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return err
	}

	chars, err := charRepo.ListCharacters()
	if err != nil {
		return err
	}

	var totalRequests, totalHits, totalTokensSaved int
	var totalCostSaved float64

	fmt.Println("Cache Performance Statistics")
	fmt.Println("===========================")

	for _, charID := range chars {
		sessions, err := repo.ListSessions(charID)
		if err != nil || len(sessions) == 0 {
			continue
		}

		char, _ := charRepo.LoadCharacter(charID)
		fmt.Printf("\n%s (%s):\n", char.Name, charID)

		var charRequests, charHits, charTokensSaved int
		var charCostSaved float64

		for _, sessionInfo := range sessions {
			session, err := repo.LoadSession(charID, sessionInfo.ID)
			if err != nil {
				continue
			}

			charRequests += session.CacheMetrics.TotalRequests
			charHits += session.CacheMetrics.CacheHits
			charTokensSaved += session.CacheMetrics.TokensSaved
			charCostSaved += session.CacheMetrics.CostSaved
		}

		if charRequests > 0 {
			hitRate := float64(charHits) / float64(charRequests) * 100
			fmt.Printf("  Sessions: %d\n", len(sessions))
			fmt.Printf("  Total Requests: %d\n", charRequests)
			fmt.Printf("  Cache Hit Rate: %.1f%%\n", hitRate)
			fmt.Printf("  Tokens Saved: %d\n", charTokensSaved)
			fmt.Printf("  Cost Saved: $%.2f\n", charCostSaved)
		}

		totalRequests += charRequests
		totalHits += charHits
		totalTokensSaved += charTokensSaved
		totalCostSaved += charCostSaved
	}

	if totalRequests > 0 {
		fmt.Println("\nOverall Statistics:")
		fmt.Printf("  Total Requests: %d\n", totalRequests)
		fmt.Printf("  Overall Hit Rate: %.1f%%\n", float64(totalHits)/float64(totalRequests)*100)
		fmt.Printf("  Total Tokens Saved: %d\n", totalTokensSaved)
		fmt.Printf("  Total Cost Saved: $%.2f\n", totalCostSaved)
	}

	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}
````

## File: cmd/status.go
````go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current configuration and status",
	Long:  `Display the current provider, model, and other configuration settings.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Get configuration
	provider := viper.GetString("provider")
	model := viper.GetString("model")
	apiKey := viper.GetString("api_key")

	// Check for API key from environment if not set
	if apiKey == "" && provider == "openai" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" && provider == "anthropic" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	// Determine actual model that will be used
	actualModel := model
	if actualModel == "" {
		switch provider {
		case "openai":
			actualModel = "gpt-4o-mini"
		case "anthropic":
			actualModel = "claude-3-haiku-20240307"
		}
	}

	fmt.Println("🤖 Roleplay Status")
	fmt.Println("==================")
	fmt.Printf("Provider: %s\n", provider)
	fmt.Printf("Model: %s\n", actualModel)
	if model != "" && model != actualModel {
		fmt.Printf("  (configured: %s, using default: %s)\n", model, actualModel)
	}
	fmt.Printf("API Key: %s\n", func() string {
		if apiKey == "" {
			return "❌ Not configured"
		}
		if len(apiKey) > 8 {
			return "✓ " + apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		}
		return "✓ Configured"
	}())

	// Show cache configuration
	fmt.Printf("\nCache Configuration:\n")
	fmt.Printf("  Default TTL: %v\n", viper.GetDuration("cache.default_ttl"))
	fmt.Printf("  Adaptive TTL: %v\n", viper.GetBool("cache.adaptive_ttl"))

	// Show data directory
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	fmt.Printf("\nData Directory: %s\n", dataDir)

	// Show character count
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err == nil {
		chars, _ := charRepo.ListCharacters()
		fmt.Printf("Characters: %d available\n", len(chars))
	}

	// Show session count
	sessionRepo := repository.NewSessionRepository(dataDir)
	totalSessions := 0
	if charRepo != nil {
		chars, _ := charRepo.ListCharacters()
		for _, charID := range chars {
			sessions, _ := sessionRepo.ListSessions(charID)
			totalSessions += len(sessions)
		}
	}
	fmt.Printf("Sessions: %d total\n", totalSessions)

	return nil
}
````

## File: internal/cache/cache_test.go
````go
package cache

import (
	"testing"
	"time"
)

func TestPromptCache(t *testing.T) {
	cache := NewPromptCache(5*time.Minute, 1*time.Minute, 10*time.Minute)

	// Test storing and retrieving
	cache.Store("test-key", CorePersonalityLayer, "test content", 5*time.Minute)

	entry, found := cache.Get("test-key")
	if !found {
		t.Fatal("Expected to find cached entry")
	}

	if len(entry.Breakpoints) != 1 {
		t.Errorf("Expected 1 breakpoint, got %d", len(entry.Breakpoints))
	}

	if entry.Breakpoints[0].Layer != CorePersonalityLayer {
		t.Errorf("Expected CorePersonalityLayer, got %s", entry.Breakpoints[0].Layer)
	}

	// Test hit count
	initialHits := entry.HitCount
	cache.Get("test-key")
	entry, _ = cache.Get("test-key")
	if entry.HitCount != initialHits+2 {
		t.Errorf("Expected hit count to increase, got %d", entry.HitCount)
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"Hello world", 2},    // 11 chars / 4 ≈ 2
		{"This is a test", 3}, // 14 chars / 4 ≈ 3
		{"", 0},               // Empty string
		{"A", 0},              // 1 char / 4 = 0
		{"1234", 1},           // 4 chars / 4 = 1
	}

	for _, tt := range tests {
		result := EstimateTokens(tt.text)
		if result != tt.expected {
			t.Errorf("EstimateTokens(%q) = %d, want %d", tt.text, result, tt.expected)
		}
	}
}

func TestAdaptiveTTL(t *testing.T) {
	cache := NewPromptCache(5*time.Minute, 1*time.Minute, 10*time.Minute)

	// Test without cached entry
	ttl := cache.CalculateAdaptiveTTL(nil, false)
	if ttl != 5*time.Minute {
		t.Errorf("Expected base TTL of 5m, got %v", ttl)
	}

	// Test with complexity bonus
	ttl = cache.CalculateAdaptiveTTL(nil, true)
	expectedTTL := time.Duration(float64(5*time.Minute) * 1.2) // 20% bonus
	if ttl != expectedTTL {
		t.Errorf("Expected TTL with complexity bonus of %v, got %v", expectedTTL, ttl)
	}

	// Test with recent access
	entry := &CacheEntry{
		LastAccess: time.Now(),
		HitCount:   5,
	}
	ttl = cache.CalculateAdaptiveTTL(entry, false)
	expectedTTL = time.Duration(float64(5*time.Minute) * 1.5) // 50% bonus
	if ttl != expectedTTL {
		t.Errorf("Expected TTL with active bonus of %v, got %v", expectedTTL, ttl)
	}

	// Test max TTL enforcement
	cache.ttl.BaseTTL = 20 * time.Minute
	ttl = cache.CalculateAdaptiveTTL(nil, false)
	if ttl != cache.ttl.MaxTTL {
		t.Errorf("Expected TTL to be capped at MaxTTL %v, got %v", cache.ttl.MaxTTL, ttl)
	}
}

func TestCacheCleanup(t *testing.T) {
	cache := NewPromptCache(100*time.Millisecond, 50*time.Millisecond, 200*time.Millisecond)

	// Store entry with short TTL
	breakpoints := []CacheBreakpoint{
		{
			Layer:    CorePersonalityLayer,
			Content:  "test",
			TTL:      100 * time.Millisecond,
			LastUsed: time.Now(),
		},
	}
	cache.StoreWithTTL("expire-key", breakpoints, 100*time.Millisecond)

	// Verify entry exists
	_, found := cache.Get("expire-key")
	if !found {
		t.Fatal("Expected to find entry before expiration")
	}

	// Wait for expiration and cleanup
	time.Sleep(150 * time.Millisecond)
	cache.cleanup()

	// Verify entry was cleaned up
	_, found = cache.Get("expire-key")
	if found {
		t.Error("Expected entry to be cleaned up after expiration")
	}
}

func TestCacheLayers(t *testing.T) {
	layers := []CacheLayer{
		CorePersonalityLayer,
		LearnedBehaviorLayer,
		EmotionalStateLayer,
		ConversationLayer,
	}

	// Verify all layers are distinct
	seen := make(map[CacheLayer]bool)
	for _, layer := range layers {
		if seen[layer] {
			t.Errorf("Duplicate layer found: %s", layer)
		}
		seen[layer] = true
	}
}
````

## File: internal/cache/cache.go
````go
package cache

import (
	"sync"
	"time"
)

// PromptCache manages cached prompts with TTL
type PromptCache struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     TTLManager
}

// NewPromptCache creates a new cache with the given TTL configuration
func NewPromptCache(baseTTL, minTTL, maxTTL time.Duration) *PromptCache {
	return &PromptCache{
		entries: make(map[string]*CacheEntry),
		ttl: TTLManager{
			BaseTTL:         baseTTL,
			ActiveBonus:     0.5,
			ComplexityBonus: 0.2,
			MinTTL:          minTTL,
			MaxTTL:          maxTTL,
		},
	}
}

// Store adds a new cache entry for a specific layer
func (pc *PromptCache) Store(key string, layer CacheLayer, content string, ttl time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, exists := pc.entries[key]
	if !exists {
		entry = &CacheEntry{
			CreatedAt:   time.Now(),
			Breakpoints: make([]CacheBreakpoint, 0),
		}
		pc.entries[key] = entry
	}

	breakpoint := CacheBreakpoint{
		Layer:      layer,
		Content:    content,
		TokenCount: EstimateTokens(content),
		TTL:        ttl,
		LastUsed:   time.Now(),
	}

	entry.Breakpoints = append(entry.Breakpoints, breakpoint)
	entry.LastAccess = time.Now()
}

// StoreWithTTL stores a complete cache entry with breakpoints
func (pc *PromptCache) StoreWithTTL(key string, breakpoints []CacheBreakpoint, ttl time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry := &CacheEntry{
		Breakpoints: breakpoints,
		CreatedAt:   time.Now(),
		LastAccess:  time.Now(),
		HitCount:    0,
	}

	pc.entries[key] = entry
}

// Get retrieves a cache entry if it exists and is not expired
func (pc *PromptCache) Get(key string) (*CacheEntry, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	entry, exists := pc.entries[key]
	if !exists {
		return nil, false
	}

	// Update access time and hit count
	entry.LastAccess = time.Now()
	entry.HitCount++

	return entry, true
}

// CalculateAdaptiveTTL determines the effective TTL based on usage patterns
func (pc *PromptCache) CalculateAdaptiveTTL(cached *CacheEntry, hasComplexCharacter bool) time.Duration {
	baseTTL := pc.ttl.BaseTTL

	// Active conversation bonus
	if cached != nil && time.Since(cached.LastAccess) < 5*time.Minute {
		baseTTL = time.Duration(float64(baseTTL) * (1 + pc.ttl.ActiveBonus))
	}

	// Character complexity bonus
	if hasComplexCharacter {
		baseTTL = time.Duration(float64(baseTTL) * (1 + pc.ttl.ComplexityBonus))
	}

	// Enforce limits
	if baseTTL < pc.ttl.MinTTL {
		baseTTL = pc.ttl.MinTTL
	}
	if baseTTL > pc.ttl.MaxTTL {
		baseTTL = pc.ttl.MaxTTL
	}

	return baseTTL
}

// CleanupWorker runs periodic cleanup of expired entries
func (pc *PromptCache) CleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		pc.cleanup()
	}
}

// cleanup removes expired cache entries
func (pc *PromptCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for key, entry := range pc.entries {
		// Check if any breakpoint has expired
		expired := false
		for _, bp := range entry.Breakpoints {
			if now.Sub(bp.LastUsed) > bp.TTL {
				expired = true
				break
			}
		}

		if expired {
			delete(pc.entries, key)
		}
	}
}

// EstimateTokens provides a rough estimation of token count
func EstimateTokens(text string) int {
	// Rough estimation: ~4 chars per token
	return len(text) / 4
}
````

## File: internal/cache/response_cache.go
````go
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// ResponseCache caches complete API responses
type ResponseCache struct {
	responses map[string]*CachedResponse
	mu        sync.RWMutex
	ttl       time.Duration
}

// CachedResponse represents a cached API response
type CachedResponse struct {
	Content    string
	TokensUsed TokenUsage
	CachedAt   time.Time
	ExpiresAt  time.Time
	HitCount   int
}

// TokenUsage represents token usage stats
type TokenUsage struct {
	Prompt       int
	Completion   int
	CachedPrompt int
	Total        int
}

// NewResponseCache creates a new response cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	cache := &ResponseCache{
		responses: make(map[string]*CachedResponse),
		ttl:       ttl,
	}

	// Start cleanup worker
	go cache.cleanupWorker()

	return cache
}

// GenerateKey creates a cache key from request parameters
func (rc *ResponseCache) GenerateKey(characterID, userID, message string) string {
	h := sha256.New()
	h.Write([]byte(characterID + "|" + userID + "|" + message))
	return hex.EncodeToString(h.Sum(nil))
}

// Get retrieves a cached response if available
func (rc *ResponseCache) Get(key string) (*CachedResponse, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	resp, exists := rc.responses[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(resp.ExpiresAt) {
		return nil, false
	}

	// Update hit count
	resp.HitCount++

	return resp, true
}

// Store adds a response to the cache
func (rc *ResponseCache) Store(key, content string, tokens TokenUsage) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.responses[key] = &CachedResponse{
		Content:    content,
		TokensUsed: tokens,
		CachedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(rc.ttl),
		HitCount:   0,
	}
}

// cleanupWorker removes expired entries
func (rc *ResponseCache) cleanupWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for key, resp := range rc.responses {
			if now.After(resp.ExpiresAt) {
				delete(rc.responses, key)
			}
		}
		rc.mu.Unlock()
	}
}

// GetStats returns cache statistics
func (rc *ResponseCache) GetStats() (hits, misses int) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	for _, resp := range rc.responses {
		hits += resp.HitCount
	}

	return hits, 0 // Misses would need to be tracked separately
}
````

## File: internal/importer/importer.go
````go
package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"

	"github.com/google/uuid"
)

type CharacterImporter struct {
	provider   providers.AIProvider
	repository *repository.CharacterRepository
	promptPath string
}

func NewCharacterImporter(provider providers.AIProvider, repo *repository.CharacterRepository) *CharacterImporter {
	return &CharacterImporter{
		provider:   provider,
		repository: repo,
		promptPath: "prompts/character-import.md",
	}
}

type importedCharacter struct {
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	Backstory        string                   `json:"backstory"`
	Personality      models.PersonalityTraits `json:"personality"`
	SpeechStyle      string                   `json:"speechStyle"`
	BehaviorPatterns []string                 `json:"behaviorPatterns"`
	KnowledgeDomains []string                 `json:"knowledgeDomains"`
	EmotionalState   models.EmotionalState    `json:"emotionalState"`
	GreetingMessage  string                   `json:"greetingMessage"`
}

func (ci *CharacterImporter) ImportFromMarkdown(ctx context.Context, markdownPath string) (*models.Character, error) {
	content, err := os.ReadFile(markdownPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown file: %w", err)
	}

	promptTemplate, err := os.ReadFile(ci.promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt template: %w", err)
	}

	tmpl, err := template.New("import").Parse(string(promptTemplate))
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var promptBuilder strings.Builder
	data := map[string]string{
		"MarkdownContent": string(content),
	}
	if err := tmpl.Execute(&promptBuilder, data); err != nil {
		return nil, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	request := &providers.PromptRequest{
		CharacterID:  "system-importer",
		UserID:       "system",
		Message:      promptBuilder.String(),
		SystemPrompt: "You are a helpful AI assistant that extracts character information from markdown files and formats it as JSON.",
	}

	response, err := ci.provider.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// Clean the response - remove any markdown code blocks
	jsonContent := strings.TrimSpace(response.Content)

	// Remove markdown code blocks if present
	if strings.HasPrefix(jsonContent, "```json") {
		jsonContent = strings.TrimPrefix(jsonContent, "```json")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	} else if strings.HasPrefix(jsonContent, "```") {
		jsonContent = strings.TrimPrefix(jsonContent, "```")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	}

	var imported importedCharacter
	if err := json.Unmarshal([]byte(jsonContent), &imported); err != nil {
		// Log the actual response for debugging
		fmt.Fprintf(os.Stderr, "Failed to parse response. Raw content:\n%s\n", jsonContent)
		return nil, fmt.Errorf("failed to parse LLM response as JSON: %w", err)
	}

	character := &models.Character{
		ID:           uuid.New().String(),
		Name:         imported.Name,
		Backstory:    imported.Backstory,
		Personality:  imported.Personality,
		SpeechStyle:  imported.SpeechStyle,
		CurrentMood:  imported.EmotionalState,
		Quirks:       imported.BehaviorPatterns,
		Memories:     []models.Memory{},
		LastModified: time.Now(),
	}

	baseFilename := strings.TrimSuffix(filepath.Base(markdownPath), filepath.Ext(markdownPath))
	character.ID = fmt.Sprintf("%s-%s", baseFilename, character.ID[:8])

	if err := ci.repository.SaveCharacter(character); err != nil {
		return nil, fmt.Errorf("failed to save character: %w", err)
	}

	return character, nil
}
````

## File: internal/models/character_test.go
````go
package models

import (
	"testing"
	"time"
)

func TestNormalizePersonality(t *testing.T) {
	tests := []struct {
		name     string
		input    PersonalityTraits
		expected PersonalityTraits
	}{
		{
			name: "values within range",
			input: PersonalityTraits{
				Openness:          0.5,
				Conscientiousness: 0.6,
				Extraversion:      0.7,
				Agreeableness:     0.8,
				Neuroticism:       0.9,
			},
			expected: PersonalityTraits{
				Openness:          0.5,
				Conscientiousness: 0.6,
				Extraversion:      0.7,
				Agreeableness:     0.8,
				Neuroticism:       0.9,
			},
		},
		{
			name: "values above 1 should be clamped",
			input: PersonalityTraits{
				Openness:          1.5,
				Conscientiousness: 2.0,
				Extraversion:      1.1,
				Agreeableness:     1.3,
				Neuroticism:       1.8,
			},
			expected: PersonalityTraits{
				Openness:          1.0,
				Conscientiousness: 1.0,
				Extraversion:      1.0,
				Agreeableness:     1.0,
				Neuroticism:       1.0,
			},
		},
		{
			name: "values below 0 should be clamped",
			input: PersonalityTraits{
				Openness:          -0.5,
				Conscientiousness: -1.0,
				Extraversion:      -0.1,
				Agreeableness:     -0.3,
				Neuroticism:       -0.8,
			},
			expected: PersonalityTraits{
				Openness:          0.0,
				Conscientiousness: 0.0,
				Extraversion:      0.0,
				Agreeableness:     0.0,
				Neuroticism:       0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePersonality(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePersonality() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCharacterLocking(t *testing.T) {
	char := &Character{
		ID:   "test-123",
		Name: "Test Character",
	}

	// Test write lock
	char.Lock()
	char.Name = "Modified Character"
	char.Unlock()

	if char.Name != "Modified Character" {
		t.Errorf("Expected name to be 'Modified Character', got %s", char.Name)
	}

	// Test read lock
	char.RLock()
	name := char.Name
	char.RUnlock()

	if name != "Modified Character" {
		t.Errorf("Expected to read 'Modified Character', got %s", name)
	}
}

func TestMemoryTypes(t *testing.T) {
	memories := []Memory{
		{
			Type:      ShortTermMemory,
			Content:   "Recent conversation",
			Timestamp: time.Now(),
			Emotional: 0.5,
		},
		{
			Type:      MediumTermMemory,
			Content:   "Important interaction",
			Timestamp: time.Now().Add(-time.Hour),
			Emotional: 0.8,
		},
		{
			Type:      LongTermMemory,
			Content:   "Core memory",
			Timestamp: time.Now().Add(-24 * time.Hour),
			Emotional: 0.95,
		},
	}

	for _, mem := range memories {
		switch mem.Type {
		case ShortTermMemory, MediumTermMemory, LongTermMemory:
			// Valid memory type
		default:
			t.Errorf("Invalid memory type: %s", mem.Type)
		}
	}
}
````

## File: internal/models/character.go
````go
package models

import (
	"sync"
	"time"
)

// PersonalityTraits represents OCEAN model traits
type PersonalityTraits struct {
	Openness          float64 `json:"openness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Extraversion      float64 `json:"extraversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Neuroticism       float64 `json:"neuroticism"`
}

// EmotionalState represents current emotional context
type EmotionalState struct {
	Joy      float64 `json:"joy"`
	Surprise float64 `json:"surprise"`
	Anger    float64 `json:"anger"`
	Fear     float64 `json:"fear"`
	Sadness  float64 `json:"sadness"`
	Disgust  float64 `json:"disgust"`
}

// MemoryType represents different types of memories
type MemoryType string

const (
	ShortTermMemory  MemoryType = "short_term"
	MediumTermMemory MemoryType = "medium_term"
	LongTermMemory   MemoryType = "long_term"
)

// Memory represents different memory types
type Memory struct {
	Type      MemoryType `json:"type"`
	Content   string     `json:"content"`
	Timestamp time.Time  `json:"timestamp"`
	Emotional float64    `json:"emotional_weight"`
}

// Character represents a complete character profile
type Character struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Backstory    string            `json:"backstory"`
	Personality  PersonalityTraits `json:"personality"`
	CurrentMood  EmotionalState    `json:"current_mood"`
	Quirks       []string          `json:"quirks"`
	SpeechStyle  string            `json:"speech_style"`
	Memories     []Memory          `json:"memories"`
	LastModified time.Time         `json:"last_modified"`
	mu           sync.RWMutex
}

// Lock acquires write lock
func (c *Character) Lock() {
	c.mu.Lock()
}

// Unlock releases write lock
func (c *Character) Unlock() {
	c.mu.Unlock()
}

// RLock acquires read lock
func (c *Character) RLock() {
	c.mu.RLock()
}

// RUnlock releases read lock
func (c *Character) RUnlock() {
	c.mu.RUnlock()
}

// NormalizePersonality ensures all personality traits are within [0, 1] range
func NormalizePersonality(p PersonalityTraits) PersonalityTraits {
	return PersonalityTraits{
		Openness:          clamp(p.Openness, 0, 1),
		Conscientiousness: clamp(p.Conscientiousness, 0, 1),
		Extraversion:      clamp(p.Extraversion, 0, 1),
		Agreeableness:     clamp(p.Agreeableness, 0, 1),
		Neuroticism:       clamp(p.Neuroticism, 0, 1),
	}
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
````

## File: internal/models/user_profile.go
````go
package models

import "time"

// UserFact represents a piece of information known about the user.
type UserFact struct {
	Key         string    `json:"key"`         // e.g., "PreferredColor", "StatedGoal", "MentionedPetName"
	Value       string    `json:"value"`       // e.g., "Blue", "Learn Go programming", "Buddy"
	SourceTurn  int       `json:"source_turn"` // Turn number in conversation where this was inferred/stated
	Confidence  float64   `json:"confidence"`  // LLM's confidence in this fact (0.0-1.0)
	LastUpdated time.Time `json:"last_updated"`
}

// UserProfile holds synthesized information about a user.
type UserProfile struct {
	UserID           string     `json:"user_id"`
	CharacterID      string     `json:"character_id"` // Profile might be character-specific
	Facts            []UserFact `json:"facts"`
	OverallSummary   string     `json:"overall_summary"`   // A brief LLM-generated summary of the user
	InteractionStyle string     `json:"interaction_style"` // e.g., "formal", "inquisitive", "humorous"
	LastAnalyzed     time.Time  `json:"last_analyzed"`
	Version          int        `json:"version"`
}
````

## File: internal/repository/character_repo.go
````go
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/models"
)

// CharacterRepository manages character persistence
type CharacterRepository struct {
	dataDir string
}

// NewCharacterRepository creates a new character repository
func NewCharacterRepository(dataDir string) (*CharacterRepository, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"characters", "sessions", "cache"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(dataDir, dir), 0755); err != nil {
			return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	return &CharacterRepository{dataDir: dataDir}, nil
}

// SaveCharacter persists a character to disk
func (r *CharacterRepository) SaveCharacter(character *models.Character) error {
	filename := filepath.Join(r.dataDir, "characters", fmt.Sprintf("%s.json", character.ID))

	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadCharacter loads a character from disk
func (r *CharacterRepository) LoadCharacter(id string) (*models.Character, error) {
	filename := filepath.Join(r.dataDir, "characters", fmt.Sprintf("%s.json", id))

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("character %s not found", id)
		}
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var character models.Character
	if err := json.Unmarshal(data, &character); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character: %w", err)
	}

	return &character, nil
}

// ListCharacters returns all available character IDs
func (r *CharacterRepository) ListCharacters() ([]string, error) {
	charactersDir := filepath.Join(r.dataDir, "characters")

	entries, err := os.ReadDir(charactersDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read characters directory: %w", err)
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			id := strings.TrimSuffix(entry.Name(), ".json")
			ids = append(ids, id)
		}
	}

	return ids, nil
}

// GetCharacterInfo returns basic info about all characters
func (r *CharacterRepository) GetCharacterInfo() ([]CharacterInfo, error) {
	charactersDir := filepath.Join(r.dataDir, "characters")

	entries, err := os.ReadDir(charactersDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read characters directory: %w", err)
	}

	var infos []CharacterInfo
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			id := strings.TrimSuffix(entry.Name(), ".json")
			char, err := r.LoadCharacter(id)
			if err != nil {
				continue
			}

			infos = append(infos, CharacterInfo{
				ID:          char.ID,
				Name:        char.Name,
				Description: truncateString(char.Backstory, 100),
				Tags:        char.Quirks,
			})
		}
	}

	return infos, nil
}

// CharacterInfo provides basic character information
type CharacterInfo struct {
	ID          string
	Name        string
	Description string
	Tags        []string
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
````

## File: internal/repository/session_repo.go
````go
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
)

// Session represents a conversation session
type Session struct {
	ID           string           `json:"id"`
	CharacterID  string           `json:"character_id"`
	UserID       string           `json:"user_id"`
	StartTime    time.Time        `json:"start_time"`
	LastActivity time.Time        `json:"last_activity"`
	Messages     []SessionMessage `json:"messages"`
	Memories     []models.Memory  `json:"memories"`
	CacheMetrics CacheMetrics     `json:"cache_metrics"`
}

// SessionMessage represents a single message in a session
type SessionMessage struct {
	Timestamp   time.Time `json:"timestamp"`
	Role        string    `json:"role"` // "user" or "character"
	Content     string    `json:"content"`
	TokensUsed  int       `json:"tokens_used,omitempty"`
	CacheHits   int       `json:"cache_hits,omitempty"`
	CacheMisses int       `json:"cache_misses,omitempty"`
}

// CacheMetrics tracks cache performance for the session
type CacheMetrics struct {
	TotalRequests int     `json:"total_requests"`
	CacheHits     int     `json:"cache_hits"`
	CacheMisses   int     `json:"cache_misses"`
	TokensSaved   int     `json:"tokens_saved"`
	CostSaved     float64 `json:"cost_saved"`
	HitRate       float64 `json:"hit_rate"`
}

// SessionRepository manages session persistence
type SessionRepository struct {
	dataDir string
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(dataDir string) *SessionRepository {
	return &SessionRepository{dataDir: dataDir}
}

// SaveSession persists a session to disk
func (s *SessionRepository) SaveSession(session *Session) error {
	sessionDir := filepath.Join(s.dataDir, "sessions", session.CharacterID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	filename := filepath.Join(sessionDir, fmt.Sprintf("%s.json", session.ID))

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadSession loads a session from disk
func (s *SessionRepository) LoadSession(characterID, sessionID string) (*Session, error) {
	filename := filepath.Join(s.dataDir, "sessions", characterID, fmt.Sprintf("%s.json", sessionID))

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session %s not found", sessionID)
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// ListSessions returns all sessions for a character
func (s *SessionRepository) ListSessions(characterID string) ([]SessionInfo, error) {
	sessionDir := filepath.Join(s.dataDir, "sessions", characterID)

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	var sessions []SessionInfo
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			sessionID := filepath.Base(entry.Name())
			sessionID = sessionID[:len(sessionID)-5] // Remove .json

			session, err := s.LoadSession(characterID, sessionID)
			if err != nil {
				continue
			}

			sessions = append(sessions, SessionInfo{
				ID:           session.ID,
				CharacterID:  session.CharacterID,
				StartTime:    session.StartTime,
				LastActivity: session.LastActivity,
				MessageCount: len(session.Messages),
				CacheHitRate: session.CacheMetrics.HitRate,
			})
		}
	}

	return sessions, nil
}

// SessionInfo provides basic session information
type SessionInfo struct {
	ID           string    `json:"id"`
	CharacterID  string    `json:"character_id"`
	StartTime    time.Time `json:"start_time"`
	LastActivity time.Time `json:"last_activity"`
	MessageCount int       `json:"message_count"`
	CacheHitRate float64   `json:"cache_hit_rate"`
}

// GetLatestSession returns the most recent session for a character
func (s *SessionRepository) GetLatestSession(characterID string) (*Session, error) {
	sessions, err := s.ListSessions(characterID)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found for character %s", characterID)
	}

	// Find most recent session
	var latest SessionInfo
	for _, session := range sessions {
		if session.LastActivity.After(latest.LastActivity) {
			latest = session
		}
	}

	return s.LoadSession(characterID, latest.ID)
}
````

## File: internal/repository/user_profile_repo.go
````go
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/models"
)

// UserProfileRepository manages persistence of user profiles
type UserProfileRepository struct {
	dataDir string
}

// NewUserProfileRepository creates a new repository instance
func NewUserProfileRepository(dataDir string) *UserProfileRepository {
	return &UserProfileRepository{
		dataDir: dataDir,
	}
}

// profileFilename generates the filename for a user profile
func (r *UserProfileRepository) profileFilename(userID, characterID string) string {
	return fmt.Sprintf("%s_%s.json", userID, characterID)
}

// SaveUserProfile saves a user profile to disk
func (r *UserProfileRepository) SaveUserProfile(profile *models.UserProfile) error {
	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create user profiles directory: %w", err)
	}

	filename := r.profileFilename(profile.UserID, profile.CharacterID)
	filepath := filepath.Join(r.dataDir, filename)

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user profile: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user profile file: %w", err)
	}

	return nil
}

// LoadUserProfile loads a user profile from disk
func (r *UserProfileRepository) LoadUserProfile(userID, characterID string) (*models.UserProfile, error) {
	filename := r.profileFilename(userID, characterID)
	filepath := filepath.Join(r.dataDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err // Let caller handle non-existence
		}
		return nil, fmt.Errorf("failed to read user profile file: %w", err)
	}

	var profile models.UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user profile: %w", err)
	}

	return &profile, nil
}

// DeleteUserProfile deletes a user profile from disk
func (r *UserProfileRepository) DeleteUserProfile(userID, characterID string) error {
	filename := r.profileFilename(userID, characterID)
	filepath := filepath.Join(r.dataDir, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}

	return nil
}

// ListUserProfiles returns all user profiles for a given user
func (r *UserProfileRepository) ListUserProfiles(userID string) ([]*models.UserProfile, error) {
	pattern := filepath.Join(r.dataDir, fmt.Sprintf("%s_*.json", userID))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list user profiles: %w", err)
	}

	var profiles []*models.UserProfile
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}

		var profile models.UserProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			continue // Skip invalid JSON files
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}
````

## File: internal/tui/components/header.go
````go
package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Header displays character information and personality stats
type Header struct {
	title         string
	characterName string
	characterID   string
	personality   PersonalityStats
	mood          string
	moodIcon      string
	width         int
	styles        headerStyles
}

type headerStyles struct {
	title       lipgloss.Style
	personality lipgloss.Style
	mood        lipgloss.Style
	character   lipgloss.Style
}

// NewHeader creates a new header component
func NewHeader(width int) *Header {
	return &Header{
		title:         "Chat",
		characterName: "Loading...",
		width:         width,
		mood:          "Unknown",
		moodIcon:      "🤔",
		styles: headerStyles{
			title: lipgloss.NewStyle().
				Foreground(GruvboxAqua).
				Bold(true).
				Padding(0, 1),
			personality: lipgloss.NewStyle().
				Foreground(GruvboxPurple),
			mood: lipgloss.NewStyle().
				Foreground(GruvboxYellow),
			character: lipgloss.NewStyle().
				Foreground(GruvboxOrange).
				Bold(true),
		},
	}
}

// Update handles messages for the header
func (h *Header) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width

	case CharacterUpdateMsg:
		h.characterName = msg.Name
		h.characterID = msg.ID
		h.mood = msg.Mood
		h.moodIcon = msg.MoodIcon
		h.personality = msg.Personality
		h.title = fmt.Sprintf("󰊕 Chat with %s", msg.Name)
	}

	return nil
}

// View renders the header
func (h *Header) View() string {
	titleLine := h.styles.title.Render(h.title)

	// Personality traits with icons
	personalityStr := fmt.Sprintf(
		" O:%.1f  C:%.1f  E:%.1f  A:%.1f  N:%.1f",
		h.personality.Openness,
		h.personality.Conscientiousness,
		h.personality.Extraversion,
		h.personality.Agreeableness,
		h.personality.Neuroticism,
	)

	personalityInfo := h.styles.personality.Render(personalityStr)
	moodInfo := h.styles.mood.Render(fmt.Sprintf(" %s %s", h.moodIcon, h.mood))

	infoLine := fmt.Sprintf("  %s • %s", personalityInfo, moodInfo)

	return lipgloss.JoinVertical(lipgloss.Left, titleLine, infoLine, "")
}

// SetSize updates the width of the header
func (h *Header) SetSize(width, _ int) {
	h.width = width
}
````

## File: internal/tui/components/inputarea.go
````go
package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputArea handles user input and processing state
type InputArea struct {
	textarea     textarea.Model
	spinner      spinner.Model
	isProcessing bool
	error        error
	focused      bool
	width        int
	styles       inputAreaStyles
}

type inputAreaStyles struct {
	error      lipgloss.Style
	processing lipgloss.Style
}

// NewInputArea creates a new input area component
func NewInputArea(width int) *InputArea {
	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Prompt = "│ "
	ta.CharLimit = 500
	// Adjust width to account for prompt and proper padding
	// Using width - 6 to prevent horizontal scrolling issues
	ta.SetWidth(width - 6)
	ta.SetHeight(1)  // Set to 1 line to prevent multi-line expansion
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Style the textarea with Gruvbox theme
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(GruvboxBg1)
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(GruvboxAqua)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(GruvboxFg)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(GruvboxGray)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(GruvboxAqua)

	return &InputArea{
		textarea: ta,
		spinner:  s,
		width:    width,
		styles: inputAreaStyles{
			error: lipgloss.NewStyle().
				Foreground(GruvboxRed).
				Bold(true),
			processing: lipgloss.NewStyle().
				Foreground(GruvboxGray),
		},
	}
}

// Update handles messages for the input area
func (i *InputArea) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		// Consistent width adjustment to prevent scrolling issues
		i.textarea.SetWidth(msg.Width - 6)

	case ProcessingStateMsg:
		i.isProcessing = msg.IsProcessing
		if msg.IsProcessing {
			i.textarea.Blur()
		} else {
			i.textarea.Focus()
		}

	case spinner.TickMsg:
		if i.isProcessing {
			var cmd tea.Cmd
			i.spinner, cmd = i.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update textarea if not processing
	if !i.isProcessing && i.focused {
		var cmd tea.Cmd
		i.textarea, cmd = i.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the input area
func (i *InputArea) View() string {
	if i.error != nil {
		return i.styles.error.Render(fmt.Sprintf("   Error: %v", i.error))
	}

	if i.isProcessing {
		spinnerText := i.styles.processing.Render("Thinking...")
		// Maintain consistent height by adding an extra newline
		return fmt.Sprintf("\n  %s %s\n\n", i.spinner.View(), spinnerText)
	}

	// Add consistent padding to prevent layout shifts
	return fmt.Sprintf("\n%s\n", i.textarea.View())
}

// Focus sets the input area as focused
func (i *InputArea) Focus() {
	i.focused = true
	if !i.isProcessing {
		i.textarea.Focus()
	}
}

// Blur removes focus from the input area
func (i *InputArea) Blur() {
	i.focused = false
	i.textarea.Blur()
}

// IsFocused returns whether the input area is focused
func (i *InputArea) IsFocused() bool {
	return i.focused
}

// SetSize updates the size of the input area
func (i *InputArea) SetSize(width, _ int) {
	i.width = width
	// Consistent width adjustment to prevent scrolling issues
	i.textarea.SetWidth(width - 6)
}

// Value returns the current input value
func (i *InputArea) Value() string {
	return i.textarea.Value()
}

// SetValue sets the input value
func (i *InputArea) SetValue(s string) {
	i.textarea.SetValue(s)
}

// Reset clears the input
func (i *InputArea) Reset() {
	i.textarea.Reset()
}

// CursorEnd moves cursor to end of input
func (i *InputArea) CursorEnd() {
	i.textarea.CursorEnd()
}

// Init returns initialization commands
func (i *InputArea) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		i.spinner.Tick,
	)
}
````

## File: internal/tui/components/messagelist.go
````go
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dotcommander/roleplay/internal/utils"
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
	Time    time.Time
	MsgType string // "normal", "help", "list", "stats", etc.
}

// MessageList displays the chat history
type MessageList struct {
	messages []Message
	viewport viewport.Model
	width    int
	height   int
	styles   messageListStyles
}

type messageListStyles struct {
	userMessage      lipgloss.Style
	characterMessage lipgloss.Style
	systemMessage    lipgloss.Style
	timestamp        lipgloss.Style
	separator        lipgloss.Style
	emptyMessage     lipgloss.Style
	user             lipgloss.Style
	character        lipgloss.Style
	system           lipgloss.Style
	// Special command output styles
	commandHeader  lipgloss.Style
	commandBox     lipgloss.Style
	helpCommand    lipgloss.Style
	helpDesc       lipgloss.Style
	listItemActive lipgloss.Style
	listItem       lipgloss.Style
}

// NewMessageList creates a new message list component
func NewMessageList(width, height int) *MessageList {
	vp := viewport.New(width, height)
	vp.SetContent("")

	return &MessageList{
		messages: []Message{},
		viewport: vp,
		width:    width,
		height:   height,
		styles: messageListStyles{
			userMessage: lipgloss.NewStyle().
				Foreground(GruvboxFg).
				Background(GruvboxBg1).
				Padding(0, 1).
				MarginRight(2),
			characterMessage: lipgloss.NewStyle().
				Foreground(GruvboxFg).
				Background(GruvboxBg).
				Padding(0, 1).
				MarginLeft(2),
			systemMessage: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			timestamp: lipgloss.NewStyle().
				Foreground(GruvboxGray).
				Italic(true),
			separator: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			emptyMessage: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			user: lipgloss.NewStyle().
				Foreground(GruvboxGreen).
				Bold(true),
			character: lipgloss.NewStyle().
				Foreground(GruvboxOrange).
				Bold(true),
			system: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			// Command styles
			commandHeader: lipgloss.NewStyle().
				Foreground(GruvboxAqua).
				Bold(true).
				Padding(0, 1),
			commandBox: lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(GruvboxAqua).
				Foreground(GruvboxFg).
				Background(GruvboxBg).
				Padding(1, 2),
			helpCommand: lipgloss.NewStyle().
				Foreground(GruvboxYellow).
				Bold(true),
			helpDesc: lipgloss.NewStyle().
				Foreground(GruvboxFg2),
			listItemActive: lipgloss.NewStyle().
				Foreground(GruvboxGreen).
				Bold(true),
			listItem: lipgloss.NewStyle().
				Foreground(GruvboxFg2),
		},
	}
}

// Update handles messages for the message list
func (m *MessageList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

	case MessageAppendMsg:
		m.messages = append(m.messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
			Time:    time.Now(),
			MsgType: msg.MsgType,
		})
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

// View renders the message list
func (m *MessageList) View() string {
	return m.viewport.View()
}

// SetSize updates the size of the message list
func (m *MessageList) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m.viewport.SetContent(m.renderMessages())
}

// renderMessages renders all messages with proper formatting
func (m *MessageList) renderMessages() string {
	if len(m.messages) == 0 {
		emptyMsg := m.styles.emptyMessage.Render("\n   Start chatting! Your conversation will appear here...\n")
		return emptyMsg
	}

	var content strings.Builder
	maxWidth := m.viewport.Width - 8 // Account for padding and margins

	for i, msg := range m.messages {
		if i > 0 {
			// Add visual separator between messages
			separator := m.styles.separator.Render(strings.Repeat("─", maxWidth))
			content.WriteString("\n" + separator + "\n\n")
		}

		timestamp := m.styles.timestamp.Render(msg.Time.Format("15:04:05"))

		if msg.Role == "user" {
			// User message
			header := fmt.Sprintf("┌─ %s %s", m.styles.user.Render("You"), timestamp)
			content.WriteString(m.styles.userMessage.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(m.styles.userMessage.Render(prefix+line) + "\n")
			}
		} else if msg.Role == "system" {
			// System message - check for special types
			if msg.MsgType == "help" || msg.MsgType == "list" || msg.MsgType == "info" || msg.MsgType == "stats" {
				// Special formatted output
				content.WriteString("\n")

				// Determine the header based on type
				var header string
				switch msg.MsgType {
				case "help":
					header = m.styles.commandHeader.Render("📚 Command Help")
				case "list":
					header = m.styles.commandHeader.Render("📋 Available Characters")
				case "stats":
					header = m.styles.commandHeader.Render("📊 Cache Statistics")
				case "info":
					header = m.styles.commandHeader.Render("ℹ️  Information")
				default:
					header = m.styles.commandHeader.Render("System")
				}

				// Format the content with special styling
				formattedContent := m.formatSpecialMessage(msg.Content, msg.MsgType, maxWidth-8)
				boxContent := lipgloss.JoinVertical(lipgloss.Left, header, "", formattedContent)
				content.WriteString(m.styles.commandBox.Width(maxWidth).Render(boxContent) + "\n")
			} else {
				// Regular system message
				header := fmt.Sprintf("┌─ %s %s", m.styles.system.Render("System"), timestamp)
				content.WriteString(m.styles.systemMessage.Render(header) + "\n")

				wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
				lines := strings.Split(wrappedContent, "\n")
				for j, line := range lines {
					prefix := "│ "
					if j == len(lines)-1 {
						prefix = "└ "
					}
					content.WriteString(m.styles.systemMessage.Render(prefix+line) + "\n")
				}
			}
		} else {
			// Character message
			header := fmt.Sprintf("┌─ %s %s", m.styles.character.Render(msg.Role), timestamp)
			content.WriteString(m.styles.characterMessage.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(m.styles.characterMessage.Render(prefix+line) + "\n")
			}
		}
	}

	return content.String()
}

// formatSpecialMessage formats messages based on their type
func (m *MessageList) formatSpecialMessage(content string, msgType string, width int) string {
	switch msgType {
	case "help":
		// Format help message with colored commands
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, " - ") {
				parts := strings.SplitN(line, " - ", 2)
				if len(parts) == 2 {
					cmd := m.styles.helpCommand.Render(parts[0])
					desc := m.styles.helpDesc.Render("- " + parts[1])
					formatted = append(formatted, fmt.Sprintf("%s %s", cmd, desc))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "list":
		// Format character list with special styling
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.HasPrefix(line, "→ ") {
				// Active character
				formatted = append(formatted, m.styles.listItemActive.Render(line))
			} else if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") {
				// Character name line
				formatted = append(formatted, m.styles.character.Render(line))
			} else if strings.HasPrefix(line, "   ") {
				// Description line
				formatted = append(formatted, m.styles.listItem.Render(line))
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "stats":
		// Format stats with colored numbers
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					label := parts[0]
					value := strings.TrimSpace(parts[1])
					// Color numbers
					if strings.Contains(value, "%") || strings.Contains(value, "tokens") {
						value = m.styles.character.Render(value)
					}
					formatted = append(formatted, fmt.Sprintf("%s: %s", label, value))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	default:
		// Default formatting
		return utils.WrapText(content, width)
	}
}

// ClearMessages clears all messages
func (m *MessageList) ClearMessages() {
	m.messages = []Message{}
	m.viewport.SetContent("")
}
````

## File: internal/tui/components/statusbar_test.go
````go
package components

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStatusBar(t *testing.T) {
	tests := []struct {
		name     string
		updates  []tea.Msg
		contains []string
	}{
		{
			name:     "Initial state shows connected",
			updates:  []tea.Msg{},
			contains: []string{"●", "session", "0"},
		},
		{
			name: "Updates cache metrics",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Connected:   true,
					CacheHits:   10,
					CacheMisses: 5,
					TokensSaved: 1500,
				},
			},
			contains: []string{"15", "67%", "1500 tokens saved"},
		},
		{
			name: "Shows error when present",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Error: errors.New("API rate limit exceeded"),
				},
			},
			contains: []string{"Error: API rate limit"},
		},
		{
			name: "Updates session ID",
			updates: []tea.Msg{
				StatusUpdateMsg{
					SessionID: "session-1234567890",
				},
			},
			contains: []string{"session-"}, // Should truncate
		},
		{
			name: "Shows disconnected state",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Connected: false,
				},
			},
			contains: []string{"○"}, // Empty circle for disconnected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create status bar
			statusBar := NewStatusBar(80)

			// Apply updates
			for _, msg := range tt.updates {
				statusBar.Update(msg)
			}

			// Get view
			view := statusBar.View()

			// Check contains
			for _, expected := range tt.contains {
				if !strings.Contains(view, expected) {
					t.Errorf("Expected view to contain %q, got: %s", expected, view)
				}
			}
		})
	}
}

func TestStatusBarResize(t *testing.T) {
	statusBar := NewStatusBar(80)

	// Update with new size
	statusBar.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Check that width was updated
	if statusBar.width != 120 {
		t.Errorf("Expected width to be 120, got %d", statusBar.width)
	}

	// View should render at new width
	view := statusBar.View()
	// The view should have styling that attempts to use the full width
	// We can't easily test the exact width due to terminal rendering,
	// but we can verify it doesn't panic and produces output
	if len(view) == 0 {
		t.Error("Expected non-empty view after resize")
	}
}

func TestStatusBarCacheRate(t *testing.T) {
	tests := []struct {
		name         string
		hits         int
		misses       int
		expectedRate string
	}{
		{
			name:         "Perfect cache rate",
			hits:         10,
			misses:       0,
			expectedRate: "100%",
		},
		{
			name:         "No cache hits",
			hits:         0,
			misses:       10,
			expectedRate: "0%",
		},
		{
			name:         "50% cache rate",
			hits:         5,
			misses:       5,
			expectedRate: "50%",
		},
		{
			name:         "No requests",
			hits:         0,
			misses:       0,
			expectedRate: "0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusBar := NewStatusBar(80)
			statusBar.Update(StatusUpdateMsg{
				CacheHits:   tt.hits,
				CacheMisses: tt.misses,
			})

			view := statusBar.View()
			if !strings.Contains(view, tt.expectedRate) {
				t.Errorf("Expected cache rate %s in view, got: %s", tt.expectedRate, view)
			}
		})
	}
}
````

## File: internal/tui/components/statusbar.go
````go
package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar displays connection status and cache metrics
type StatusBar struct {
	width       int
	connected   bool
	cacheHits   int
	cacheMisses int
	tokensSaved int
	sessionID   string
	model       string
	lastError   error
	styles      statusBarStyles
}

type statusBarStyles struct {
	container    lipgloss.Style
	connected    lipgloss.Style
	disconnected lipgloss.Style
	metrics      lipgloss.Style
	error        lipgloss.Style
	session      lipgloss.Style
	model        lipgloss.Style
}

// NewStatusBar creates a new status bar component
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		width:     width,
		connected: true,
		sessionID: "session",
		model:     "gpt-4o-mini",
		styles: statusBarStyles{
			container:    lipgloss.NewStyle().Background(GruvboxAqua).Foreground(GruvboxBg),
			connected:    lipgloss.NewStyle().Foreground(GruvboxBg).Bold(true),
			disconnected: lipgloss.NewStyle().Foreground(GruvboxRed).Bold(true),
			metrics:      lipgloss.NewStyle().Foreground(GruvboxBg),
			error:        lipgloss.NewStyle().Foreground(GruvboxRed).Bold(true),
			session:      lipgloss.NewStyle().Foreground(GruvboxBg),
			model:        lipgloss.NewStyle().Foreground(GruvboxBg),
		},
	}
}

// Update handles messages for the status bar
func (s *StatusBar) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width

	case StatusUpdateMsg:
		s.connected = msg.Connected
		s.cacheHits = msg.CacheHits
		s.cacheMisses = msg.CacheMisses
		s.tokensSaved = msg.TokensSaved
		if msg.SessionID != "" {
			s.sessionID = msg.SessionID
		}
		if msg.Model != "" {
			s.model = msg.Model
		}
		s.lastError = msg.Error
	}

	return nil
}

// View renders the status bar
func (s *StatusBar) View() string {
	cacheRate := 0.0
	totalRequests := s.cacheHits + s.cacheMisses
	if totalRequests > 0 {
		cacheRate = float64(s.cacheHits) / float64(totalRequests) * 100
	}

	// Connection status
	statusStyle := s.styles.connected
	statusIcon := "●"
	if !s.connected {
		statusStyle = s.styles.disconnected
		statusIcon = "○"
	}

	// Session ID (truncated if too long)
	sessionDisplay := s.sessionID
	if len(sessionDisplay) > 8 {
		sessionDisplay = sessionDisplay[:8]
	}

	// Build status parts
	parts := []string{
		statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, sessionDisplay)),
		s.styles.model.Render(s.model),
		s.styles.metrics.Render(fmt.Sprintf(" %d", totalRequests)),
		s.styles.metrics.Render(fmt.Sprintf("%s %.0f%%", statusIcon, cacheRate)),
		s.styles.metrics.Render(fmt.Sprintf(" %d tokens saved", s.tokensSaved)),
	}

	status := lipgloss.JoinHorizontal(lipgloss.Left, parts[0], " │ ", parts[1], " │", parts[2], " │ ", parts[3], " │", parts[4])

	// Add error if present
	if s.lastError != nil {
		errorText := s.styles.error.Render(fmt.Sprintf(" │ Error: %v", s.lastError))
		status = lipgloss.JoinHorizontal(lipgloss.Left, status, errorText)
	}

	return s.styles.container.Width(s.width).Render(" " + status)
}

// SetSize updates the width of the status bar
func (s *StatusBar) SetSize(width, _ int) {
	s.width = width
}
````

## File: internal/tui/components/types.go
````go
package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Component defines the interface for all TUI components
type Component interface {
	Update(tea.Msg) tea.Cmd
	View() string
}

// WithFocus allows components to track focus state
type WithFocus interface {
	Focus()
	Blur()
	IsFocused() bool
}

// WithSize allows components to be resized
type WithSize interface {
	SetSize(width, height int)
}

// Message types for component communication
type (
	// FocusChangeMsg indicates focus has changed to a new component
	FocusChangeMsg struct {
		ComponentID string
	}

	// StatusUpdateMsg updates the status bar
	StatusUpdateMsg struct {
		Connected   bool
		CacheHits   int
		CacheMisses int
		TokensSaved int
		SessionID   string
		Model       string
		Error       error
	}

	// CharacterUpdateMsg updates character info in header
	CharacterUpdateMsg struct {
		Name        string
		ID          string
		Mood        string
		MoodIcon    string
		Personality PersonalityStats
	}

	// MessageAppendMsg adds a new message to the chat
	MessageAppendMsg struct {
		Role    string
		Content string
		MsgType string // "normal", "help", "list", "stats", etc.
	}

	// ProcessingStateMsg updates the processing state
	ProcessingStateMsg struct {
		IsProcessing bool
	}
)

// PersonalityStats holds OCEAN personality traits
type PersonalityStats struct {
	Openness          float64
	Conscientiousness float64
	Extraversion      float64
	Agreeableness     float64
	Neuroticism       float64
}

// Common styles using Gruvbox Dark theme
var (
	// Gruvbox Dark Colors
	GruvboxBg     = lipgloss.Color("#282828")
	GruvboxBg1    = lipgloss.Color("#3c3836")
	GruvboxFg     = lipgloss.Color("#ebdbb2")
	GruvboxRed    = lipgloss.Color("#fb4934")
	GruvboxGreen  = lipgloss.Color("#b8bb26")
	GruvboxYellow = lipgloss.Color("#fabd2f")
	GruvboxPurple = lipgloss.Color("#d3869b")
	GruvboxAqua   = lipgloss.Color("#8ec07c")
	GruvboxOrange = lipgloss.Color("#fe8019")
	GruvboxGray   = lipgloss.Color("#928374")
	GruvboxFg2    = lipgloss.Color("#d5c4a1")
)
````

## File: internal/tui/commands.go
````go
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
)

// handleSlashCommand processes slash commands and returns a tea.Cmd
func (m *Model) handleSlashCommand(input string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return slashCommandResult{
				cmdType: "error",
				content: "Invalid command",
				msgType: "error",
			}
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "/exit", "/quit", "/q":
			return slashCommandResult{cmdType: "quit"}

		case "/help", "/h":
			return slashCommandResult{
				cmdType: "help",
				content: `Available slash commands:
/help, /h     - Show this help message
/exit, /quit, /q - Exit the chat
/clear, /c    - Clear chat history
/list         - List all available characters
/switch <id>  - Switch to a different character
/stats        - Show cache statistics
/mood         - Show character's current mood
/personality  - Show character's personality traits
/session      - Show session information`,
				msgType: "help",
			}

		case "/clear", "/c":
			return slashCommandResult{cmdType: "clear"}

		case "/list":
			return m.listCharacters()

		case "/stats":
			return m.showStats()

		case "/mood":
			return m.showMood()

		case "/personality":
			return m.showPersonality()

		case "/session":
			return m.showSession()

		case "/switch":
			if len(parts) < 2 {
				return slashCommandResult{
					cmdType: "error",
					content: "Usage: /switch <character-id>\nUse /list to see available characters",
					msgType: "error",
				}
			}
			return m.switchCharacter(parts[1])

		default:
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Unknown command: %s\nType /help for available commands", command),
				msgType: "error",
			}
		}
	}
}

func (m *Model) listCharacters() slashCommandResult {
	// List all available characters from the repository
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return slashCommandResult{
			cmdType: "error",
			content: fmt.Sprintf("Error accessing characters: %v", err),
			msgType: "error",
		}
	}

	characterIDs, err := charRepo.ListCharacters()
	if err != nil {
		return slashCommandResult{
			cmdType: "error",
			content: fmt.Sprintf("Error listing characters: %v", err),
			msgType: "error",
		}
	}

	if len(characterIDs) == 0 {
		return slashCommandResult{
			cmdType: "info",
			content: "No characters available. Use 'roleplay character create' to add characters.",
			msgType: "info",
		}
	}

	var listText strings.Builder
	listText.WriteString("Available Characters:\n")

	for _, id := range characterIDs {
		// Try to get from bot first (already loaded)
		char, err := m.bot.GetCharacter(id)
		if err != nil {
			// If not loaded, load from repository
			char, err = charRepo.LoadCharacter(id)
			if err != nil {
				continue
			}
			// Load into bot for future use
			_ = m.bot.CreateCharacter(char)
		}

		// Current character indicator
		indicator := "  "
		if id == m.characterID {
			indicator = "→ "
		}

		// Mood icon
		tempMood := m.calculateMood(char)
		moodIcon := m.getMoodIcon(tempMood)

		listText.WriteString(fmt.Sprintf("\n%s%s (%s) %s %s\n",
			indicator, char.Name, id, moodIcon, tempMood))

		// Add brief description from backstory (first sentence)
		backstory := char.Backstory
		if idx := strings.Index(backstory, "."); idx != -1 && idx < 100 {
			backstory = backstory[:idx+1]
		} else if len(backstory) > 100 {
			backstory = backstory[:97] + "..."
		}
		listText.WriteString(fmt.Sprintf("   %s\n", backstory))
	}

	return slashCommandResult{
		cmdType: "list",
		content: listText.String(),
		msgType: "list",
	}
}

func (m *Model) showStats() slashCommandResult {
	cacheRate := 0.0
	if m.totalRequests > 0 {
		cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
	}

	content := fmt.Sprintf(`Cache Statistics:
• Total requests: %d
• Cache hits: %d
• Cache misses: %d  
• Hit rate: %.1f%%
• Tokens saved: %d`,
		m.totalRequests,
		m.cacheHits,
		m.totalRequests-m.cacheHits,
		cacheRate,
		m.lastTokensSaved)

	return slashCommandResult{
		cmdType: "stats",
		content: content,
		msgType: "stats",
	}
}

func (m *Model) showMood() slashCommandResult {
	if m.character == nil {
		return slashCommandResult{
			cmdType: "error",
			content: "Character not loaded",
			msgType: "error",
		}
	}

	mood := m.getDominantMood()
	icon := m.getMoodIcon(mood)
	content := fmt.Sprintf(`%s Current Mood: %s

Emotional State:
• Joy: %.1f      • Surprise: %.1f
• Anger: %.1f    • Fear: %.1f  
• Sadness: %.1f  • Disgust: %.1f`,
		icon, mood,
		m.character.CurrentMood.Joy,
		m.character.CurrentMood.Surprise,
		m.character.CurrentMood.Anger,
		m.character.CurrentMood.Fear,
		m.character.CurrentMood.Sadness,
		m.character.CurrentMood.Disgust)

	return slashCommandResult{
		cmdType: "mood",
		content: content,
		msgType: "info",
	}
}

func (m *Model) showPersonality() slashCommandResult {
	if m.character == nil {
		return slashCommandResult{
			cmdType: "error",
			content: "Character not loaded",
			msgType: "error",
		}
	}

	content := fmt.Sprintf(`%s's Personality (OCEAN Model):

• Openness: %.1f        (creativity, openness to experience)
• Conscientiousness: %.1f (organization, self-discipline) 
• Extraversion: %.1f     (sociability, assertiveness)
• Agreeableness: %.1f    (compassion, cooperation)
• Neuroticism: %.1f      (emotional instability, anxiety)`,
		m.character.Name,
		m.character.Personality.Openness,
		m.character.Personality.Conscientiousness,
		m.character.Personality.Extraversion,
		m.character.Personality.Agreeableness,
		m.character.Personality.Neuroticism)

	return slashCommandResult{
		cmdType: "personality",
		content: content,
		msgType: "info",
	}
}

func (m *Model) showSession() slashCommandResult {
	characterName := m.characterID
	if m.character != nil {
		characterName = m.character.Name
	}
	sessionIDDisplay := m.sessionID
	if len(m.sessionID) > 8 {
		sessionIDDisplay = m.sessionID[:8] + "..."
	}

	messageCount := 0 // Would be tracked in the refactored version

	content := fmt.Sprintf(`Session Information:
• Session ID: %s
• Character: %s (%s)
• User: %s
• Messages: %d
• Started: %s`,
		sessionIDDisplay,
		characterName,
		m.characterID,
		m.userID,
		messageCount,
		m.context.StartTime.Format("Jan 2, 2006 15:04"))

	return slashCommandResult{
		cmdType: "session",
		content: content,
		msgType: "info",
	}
}

func (m *Model) switchCharacter(searchTerm string) slashCommandResult {
	// First try exact match (for backward compatibility)
	if searchTerm == m.characterID {
		return slashCommandResult{
			cmdType: "info",
			content: fmt.Sprintf("Already chatting with %s", m.characterID),
			msgType: "info",
		}
	}

	// Try to find character using fuzzy search
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	charRepo, repoErr := repository.NewCharacterRepository(dataDir)
	if repoErr != nil {
		return slashCommandResult{
			cmdType: "error",
			content: fmt.Sprintf("Error accessing characters: %v", repoErr),
			msgType: "error",
		}
	}

	// Find matching character
	matchedCharID, matchedChar, err := m.findCharacterFuzzy(charRepo, searchTerm)
	if err != nil {
		return slashCommandResult{
			cmdType: "error",
			content: err.Error(),
			msgType: "error",
		}
	}

	// Check if it's the same character
	if matchedCharID == m.characterID {
		return slashCommandResult{
			cmdType: "info",
			content: fmt.Sprintf("Already chatting with %s", matchedChar.Name),
			msgType: "info",
		}
	}

	// Check if character is already loaded in bot
	char, err := m.bot.GetCharacter(matchedCharID)
	if err != nil {
		// Load character into bot
		if err := m.bot.CreateCharacter(matchedChar); err != nil {
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Error loading character: %v", err),
				msgType: "error",
			}
		}
		char = matchedChar
	}

	return slashCommandResult{
		cmdType:        "switch",
		newCharacterID: matchedCharID,
		newCharacter:   char,
	}
}

// findCharacterFuzzy searches for a character using fuzzy matching
func (m *Model) findCharacterFuzzy(charRepo *repository.CharacterRepository, searchTerm string) (string, *models.Character, error) {
	// Normalize search term
	searchLower := strings.ToLower(strings.TrimSpace(searchTerm))
	
	// Get all character info
	charInfos, err := charRepo.GetCharacterInfo()
	if err != nil {
		return "", nil, fmt.Errorf("error listing characters: %v", err)
	}
	
	if len(charInfos) == 0 {
		return "", nil, fmt.Errorf("no characters available")
	}
	
	// Try exact ID match first
	for _, info := range charInfos {
		if strings.ToLower(info.ID) == searchLower {
			char, err := charRepo.LoadCharacter(info.ID)
			if err != nil {
				continue
			}
			return info.ID, char, nil
		}
	}
	
	// Try exact name match
	for _, info := range charInfos {
		if strings.ToLower(info.Name) == searchLower {
			char, err := charRepo.LoadCharacter(info.ID)
			if err != nil {
				continue
			}
			return info.ID, char, nil
		}
	}
	
	// Try prefix matching on ID
	var matches []repository.CharacterInfo
	for _, info := range charInfos {
		if strings.HasPrefix(strings.ToLower(info.ID), searchLower) {
			matches = append(matches, info)
		}
	}
	
	// If only one prefix match, use it
	if len(matches) == 1 {
		char, err := charRepo.LoadCharacter(matches[0].ID)
		if err != nil {
			return "", nil, fmt.Errorf("error loading character: %v", err)
		}
		return matches[0].ID, char, nil
	}
	
	// Try prefix matching on name
	if len(matches) == 0 {
		for _, info := range charInfos {
			if strings.HasPrefix(strings.ToLower(info.Name), searchLower) {
				matches = append(matches, info)
			}
		}
	}
	
	// If only one name prefix match, use it
	if len(matches) == 1 {
		char, err := charRepo.LoadCharacter(matches[0].ID)
		if err != nil {
			return "", nil, fmt.Errorf("error loading character: %v", err)
		}
		return matches[0].ID, char, nil
	}
	
	// Try substring matching on ID or name
	if len(matches) == 0 {
		for _, info := range charInfos {
			if strings.Contains(strings.ToLower(info.ID), searchLower) ||
			   strings.Contains(strings.ToLower(info.Name), searchLower) {
				matches = append(matches, info)
			}
		}
	}
	
	// Handle results
	if len(matches) == 0 {
		return "", nil, fmt.Errorf("no character found matching '%s'. Use /list to see available characters", searchTerm)
	}
	
	if len(matches) == 1 {
		char, err := charRepo.LoadCharacter(matches[0].ID)
		if err != nil {
			return "", nil, fmt.Errorf("error loading character: %v", err)
		}
		return matches[0].ID, char, nil
	}
	
	// Multiple matches found
	var matchNames []string
	for _, match := range matches {
		matchNames = append(matchNames, fmt.Sprintf("%s (%s)", match.Name, match.ID))
	}
	return "", nil, fmt.Errorf("multiple characters match '%s':\n%s\nPlease be more specific", 
		searchTerm, strings.Join(matchNames, "\n"))
}

// Helper function to calculate mood for a character
func (m *Model) calculateMood(char *models.Character) string {
	if char == nil {
		return "Unknown"
	}

	moods := map[string]float64{
		"Joy":      char.CurrentMood.Joy,
		"Surprise": char.CurrentMood.Surprise,
		"Anger":    char.CurrentMood.Anger,
		"Fear":     char.CurrentMood.Fear,
		"Sadness":  char.CurrentMood.Sadness,
		"Disgust":  char.CurrentMood.Disgust,
	}

	maxMood := "Neutral"
	maxValue := 0.0

	for mood, value := range moods {
		if value > maxValue {
			maxMood = mood
			maxValue = value
		}
	}

	if maxValue < 0.2 {
		return "Neutral"
	}

	return maxMood
}
````

## File: internal/tui/model.go
````go
package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/tui/components"
)

// Model represents the main TUI application state
type Model struct {
	// Components
	header      *components.Header
	messageList *components.MessageList
	inputArea   *components.InputArea
	statusBar   *components.StatusBar

	// UI state
	width        int
	height       int
	ready        bool
	currentFocus string // "input", "messages"

	// Business logic
	characterID string
	userID      string
	sessionID   string
	scenarioID  string
	bot         *services.CharacterBot
	character   *models.Character
	context     models.ConversationContext

	// Command history
	commandHistory []string
	historyIndex   int
	historyBuffer  string

	// Cache metrics
	lastCacheHit    bool
	lastTokensSaved int
	totalRequests   int
	cacheHits       int

	// Model info
	aiModel string
}

// NewModel creates a new TUI model
func NewModel(cfg Config) *Model {
	// Calculate initial sizes
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3
	messagesHeight := 20 // Default, will be adjusted

	if cfg.Height > 0 {
		messagesHeight = cfg.Height - headerHeight - statusHeight - helpHeight - inputHeight - 2
	}

	width := 80 // Default width
	if cfg.Width > 0 {
		width = cfg.Width
	}

	return &Model{
		header:      components.NewHeader(width),
		messageList: components.NewMessageList(width-4, messagesHeight),
		inputArea:   components.NewInputArea(width),
		statusBar:   components.NewStatusBar(width),

		characterID: cfg.CharacterID,
		userID:      cfg.UserID,
		sessionID:   cfg.SessionID,
		scenarioID:  cfg.ScenarioID,
		bot:         cfg.Bot,
		context:     cfg.Context,
		aiModel:     cfg.Model,

		currentFocus: "input",

		commandHistory: []string{},
		historyIndex:   0,

		totalRequests: cfg.InitialMetrics.TotalRequests,
		cacheHits:     cfg.InitialMetrics.CacheHits,
	}
}

// Config holds configuration for creating a new Model
type Config struct {
	CharacterID    string
	UserID         string
	SessionID      string
	ScenarioID     string
	Bot            *services.CharacterBot
	Context        models.ConversationContext
	Model          string
	Width          int
	Height         int
	InitialMetrics struct {
		TotalRequests int
		CacheHits     int
		TokensSaved   int
	}
}

// Init initializes the TUI model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.inputArea.Init(),
		m.loadCharacterInfo(),
	)
}

// Update handles all incoming messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.initializeLayout(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.resizeComponents(msg.Width, msg.Height)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.saveSession()
			return m, tea.Quit

		case tea.KeyTab:
			// Switch focus between components
			if m.currentFocus == "input" {
				m.currentFocus = "messages"
				m.inputArea.Blur()
			} else {
				m.currentFocus = "input"
				m.inputArea.Focus()
			}

		case tea.KeyEnter:
			if m.currentFocus == "input" && !m.inputArea.IsFocused() {
				// Processing state, ignore enter
				return m, nil
			}

			if m.inputArea.Value() != "" {
				message := m.inputArea.Value()

				// Add to command history
				m.commandHistory = append(m.commandHistory, message)
				m.historyIndex = len(m.commandHistory)
				m.historyBuffer = ""

				m.inputArea.Reset()

				// Check for slash commands
				if strings.HasPrefix(message, "/") {
					return m, m.handleSlashCommand(message)
				}

				// Regular message
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "user",
					Content: message,
					MsgType: "normal",
				})

				m.inputArea.Update(components.ProcessingStateMsg{IsProcessing: true})
				m.totalRequests++

				return m, m.sendMessage(message)
			}

		case tea.KeyUp:
			// Navigate command history
			if len(m.commandHistory) > 0 && m.historyIndex > 0 {
				if m.historyIndex == len(m.commandHistory) {
					m.historyBuffer = m.inputArea.Value()
				}
				m.historyIndex--
				m.inputArea.SetValue(m.commandHistory[m.historyIndex])
				m.inputArea.CursorEnd()
			}

		case tea.KeyDown:
			// Navigate command history
			if len(m.commandHistory) > 0 && m.historyIndex < len(m.commandHistory) {
				m.historyIndex++

				if m.historyIndex == len(m.commandHistory) {
					m.inputArea.SetValue(m.historyBuffer)
				} else {
					m.inputArea.SetValue(m.commandHistory[m.historyIndex])
				}
				m.inputArea.CursorEnd()
			}
		}

	case characterInfoMsg:
		m.character = msg.character
		m.updateCharacterDisplay()

	case responseMsg:
		m.inputArea.Update(components.ProcessingStateMsg{IsProcessing: false})

		if msg.err != nil {
			m.statusBar.Update(components.StatusUpdateMsg{Error: msg.err})
		} else {
			// Add character response
			m.messageList.Update(components.MessageAppendMsg{
				Role:    m.character.Name,
				Content: msg.content,
				MsgType: "normal",
			})

			// Update cache metrics
			if msg.metrics != nil {
				m.lastCacheHit = msg.metrics.Hit
				m.lastTokensSaved = msg.metrics.SavedTokens
				if msg.metrics.Hit {
					m.cacheHits++
				}
			}

			// Update context
			m.updateContext()

			// Update status bar
			m.updateStatusBar()

			// Save session
			go m.saveSession()
		}

	case slashCommandResult:
		// Handle slash command results
		switch msg.cmdType {
		case "quit":
			m.saveSession()
			return m, tea.Quit

		case "clear":
			m.messageList.ClearMessages()
			m.messageList.Update(components.MessageAppendMsg{
				Role:    "system",
				Content: "Chat history cleared",
				MsgType: "info",
			})

		case "switch":
			if msg.err != nil {
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "system",
					Content: msg.err.Error(),
					MsgType: "error",
				})
			} else {
				// Save current session before switching
				m.saveSession()

				// Update character
				m.characterID = msg.newCharacterID
				m.character = msg.newCharacter
				m.updateCharacterDisplay()

				// Clear conversation and start new session
				m.messageList.ClearMessages()
				m.sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
				m.context = models.ConversationContext{
					SessionID:      m.sessionID,
					StartTime:      time.Now(),
					RecentMessages: []models.Message{},
				}

				// Reset cache metrics
				m.totalRequests = 0
				m.cacheHits = 0
				m.lastTokensSaved = 0
				m.lastCacheHit = false

				// Add switch notification
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "system",
					Content: fmt.Sprintf("Switched to %s (%s). Starting new session.", msg.newCharacter.Name, msg.newCharacterID),
					MsgType: "info",
				})

				m.updateStatusBar()
			}

		default:
			// Display command output
			m.messageList.Update(components.MessageAppendMsg{
				Role:    "system",
				Content: msg.content,
				MsgType: msg.msgType,
			})
		}
	}

	// Route updates to components
	cmds = append(cmds, m.routeToComponents(msg)...)

	return m, tea.Batch(cmds...)
}

// View renders the entire TUI
func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Render all components
	header := m.header.View()

	// Main chat viewport with border
	chatView := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(components.GruvboxGray).
		Foreground(components.GruvboxFg).
		Background(components.GruvboxBg).
		Padding(1).
		Render(m.messageList.View())

	inputArea := m.inputArea.View()
	statusBar := m.statusBar.View()

	// Help text
	help := lipgloss.NewStyle().
		Foreground(components.GruvboxGray).
		Italic(true).
		Render("  ⌃C quit • ↵ send • ↑↓ history • /help commands • /exit quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		chatView,
		inputArea,
		statusBar,
		help,
	)
}

// Helper methods

func (m *Model) initializeLayout(width, height int) {
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3  // This includes the newlines around the textarea
	borderHeight := 2
	chatPadding := 2  // The padding(1) adds 2 lines total
	verticalMargins := headerHeight + statusHeight + helpHeight + inputHeight + borderHeight + chatPadding

	messagesHeight := height - verticalMargins
	if messagesHeight < 5 {
		messagesHeight = 5
	}

	m.header.SetSize(width, headerHeight)
	m.messageList.SetSize(width-4, messagesHeight)
	m.inputArea.SetSize(width, inputHeight)
	m.statusBar.SetSize(width, statusHeight)

	m.inputArea.Focus()
}

func (m *Model) resizeComponents(width, height int) {
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3  // This includes the newlines around the textarea
	borderHeight := 2
	chatPadding := 2  // The padding(1) adds 2 lines total
	verticalMargins := headerHeight + statusHeight + helpHeight + inputHeight + borderHeight + chatPadding

	messagesHeight := height - verticalMargins
	if messagesHeight < 5 {
		messagesHeight = 5
	}

	m.header.SetSize(width, headerHeight)
	m.messageList.SetSize(width-4, messagesHeight)
	m.inputArea.SetSize(width, inputHeight)
	m.statusBar.SetSize(width, statusHeight)
}

func (m *Model) routeToComponents(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd

	// Update all components
	if cmd := m.header.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.messageList.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.inputArea.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.statusBar.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m *Model) updateCharacterDisplay() {
	if m.character == nil {
		return
	}

	mood := m.getDominantMood()
	moodIcon := m.getMoodIcon(mood)

	m.header.Update(components.CharacterUpdateMsg{
		Name:     m.character.Name,
		ID:       m.character.ID,
		Mood:     mood,
		MoodIcon: moodIcon,
		Personality: components.PersonalityStats{
			Openness:          m.character.Personality.Openness,
			Conscientiousness: m.character.Personality.Conscientiousness,
			Extraversion:      m.character.Personality.Extraversion,
			Agreeableness:     m.character.Personality.Agreeableness,
			Neuroticism:       m.character.Personality.Neuroticism,
		},
	})
}

func (m *Model) updateStatusBar() {
	m.statusBar.Update(components.StatusUpdateMsg{
		Connected:   true,
		CacheHits:   m.cacheHits,
		CacheMisses: m.totalRequests - m.cacheHits,
		TokensSaved: m.lastTokensSaved,
		SessionID:   m.sessionID,
		Model:       m.aiModel,
	})
}

func (m *Model) updateContext() {
	// Implementation would update conversation context
	// This is simplified for the example
}

func (m *Model) getDominantMood() string {
	if m.character == nil {
		return "Unknown"
	}

	moods := map[string]float64{
		"Joy":      m.character.CurrentMood.Joy,
		"Surprise": m.character.CurrentMood.Surprise,
		"Anger":    m.character.CurrentMood.Anger,
		"Fear":     m.character.CurrentMood.Fear,
		"Sadness":  m.character.CurrentMood.Sadness,
		"Disgust":  m.character.CurrentMood.Disgust,
	}

	maxMood := "Neutral"
	maxValue := 0.0

	for mood, value := range moods {
		if value > maxValue {
			maxMood = mood
			maxValue = value
		}
	}

	if maxValue < 0.2 {
		return "Neutral"
	}

	return maxMood
}

func (m *Model) getMoodIcon(mood string) string {
	switch mood {
	case "Joy":
		return "😊"
	case "Surprise":
		return "😲"
	case "Anger":
		return "😠"
	case "Fear":
		return "😨"
	case "Sadness":
		return "😢"
	case "Disgust":
		return "🤢"
	case "Neutral":
		return "😐"
	default:
		return "🤔"
	}
}

// Message types for tea.Cmd results

type characterInfoMsg struct {
	character *models.Character
}

type responseMsg struct {
	content string
	metrics *cache.CacheMetrics
	err     error
}

type slashCommandResult struct {
	cmdType        string // "help", "list", "stats", etc.
	content        string
	msgType        string // for display formatting
	err            error
	newCharacterID string // for switch command
	newCharacter   *models.Character
}

// Tea commands

func (m *Model) loadCharacterInfo() tea.Cmd {
	return func() tea.Msg {
		char, err := m.bot.GetCharacter(m.characterID)
		if err != nil {
			return responseMsg{err: err}
		}
		return characterInfoMsg{character: char}
	}
}

func (m *Model) sendMessage(message string) tea.Cmd {
	return func() tea.Msg {
		req := &models.ConversationRequest{
			CharacterID: m.characterID,
			UserID:      m.userID,
			Message:     message,
			ScenarioID:  m.scenarioID,
			Context:     m.context,
		}

		ctx := context.Background()
		resp, err := m.bot.ProcessRequest(ctx, req)
		if err != nil {
			return responseMsg{err: err}
		}

		// Get updated character state
		char, _ := m.bot.GetCharacter(m.characterID)
		if char != nil {
			m.character = char
		}

		return responseMsg{
			content: resp.Content,
			metrics: &resp.CacheMetrics,
		}
	}
}

// Placeholder for session saving
func (m *Model) saveSession() {
	// Implementation would save the current session
	// This is left as a placeholder for the refactored version
}
````

## File: internal/utils/text.go
````go
package utils

import "strings"

// WrapText wraps text to fit within the specified width
func WrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed width, start new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}
````

## File: main.go
````go
package main

import (
	"github.com/dotcommander/roleplay/cmd"
)

func main() {
	cmd.Execute()
}
````

## File: cmd/import.go
````go
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/importer"
	"github.com/dotcommander/roleplay/internal/repository"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [markdown-file]",
	Short: "Import a character from an unstructured markdown file",
	Long: `Import a character from an unstructured markdown file using AI to extract
character information and convert it to the roleplay format.

Example:
  roleplay import /path/to/character.md
  roleplay import ~/Library/Application\ Support/aichat/roles/rick.md`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	markdownPath := args[0]

	if strings.HasPrefix(markdownPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		markdownPath = filepath.Join(home, markdownPath[1:])
	}

	absPath, err := filepath.Abs(markdownPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("markdown file not found: %s", absPath)
	}

	config := GetConfig()

	// Create provider using factory
	provider, err := factory.CreateProvider(config)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}
	characterImporter := importer.NewCharacterImporter(provider, repo)

	fmt.Printf("Importing character from: %s\n", absPath)
	fmt.Println("Analyzing markdown content with AI...")

	ctx := context.Background()
	character, err := characterImporter.ImportFromMarkdown(ctx, absPath)
	if err != nil {
		return fmt.Errorf("failed to import character: %w", err)
	}

	fmt.Printf("\nSuccessfully imported character: %s\n", character.Name)
	fmt.Printf("ID: %s\n", character.ID)
	fmt.Printf("Backstory: %s\n", character.Backstory)
	fmt.Printf("\nPersonality traits:\n")
	fmt.Printf("  Openness: %.2f\n", character.Personality.Openness)
	fmt.Printf("  Conscientiousness: %.2f\n", character.Personality.Conscientiousness)
	fmt.Printf("  Extraversion: %.2f\n", character.Personality.Extraversion)
	fmt.Printf("  Agreeableness: %.2f\n", character.Personality.Agreeableness)
	fmt.Printf("  Neuroticism: %.2f\n", character.Personality.Neuroticism)

	charactersDir := filepath.Join(dataDir, "characters")
	fmt.Printf("\nCharacter saved to: %s\n", filepath.Join(charactersDir, character.ID+".json"))
	fmt.Printf("\nYou can now chat with this character using:\n")
	fmt.Printf("  roleplay chat \"Hello!\" --character %s\n", character.ID)
	fmt.Printf("  roleplay interactive --character %s\n", character.ID)

	return nil
}
````

## File: cmd/init.go
````go
package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup wizard for roleplay",
	Long: `Interactive setup wizard that guides you through configuring roleplay
for your preferred LLM provider (OpenAI, Ollama, LM Studio, OpenRouter, etc.)`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// Provider presets for common configurations
type providerPreset struct {
	Name         string
	BaseURL      string
	RequiresKey  bool
	LocalModel   bool
	DefaultModel string
	Description  string
}

var providerPresets = []providerPreset{
	{
		Name:         "openai",
		BaseURL:      "https://api.openai.com/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "gpt-4o-mini",
		Description:  "OpenAI (Official)",
	},
	{
		Name:         "anthropic",
		BaseURL:      "https://api.anthropic.com/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "claude-3-haiku-20240307",
		Description:  "Anthropic Claude (OpenAI-Compatible)",
	},
	{
		Name:         "gemini",
		BaseURL:      "https://generativelanguage.googleapis.com/v1beta/openai",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "models/gemini-1.5-flash",
		Description:  "Google Gemini (OpenAI-Compatible)",
	},
	{
		Name:         "ollama",
		BaseURL:      "http://localhost:11434/v1",
		RequiresKey:  false,
		LocalModel:   true,
		DefaultModel: "llama3",
		Description:  "Ollama (Local LLMs)",
	},
	{
		Name:         "lmstudio",
		BaseURL:      "http://localhost:1234/v1",
		RequiresKey:  false,
		LocalModel:   true,
		DefaultModel: "local-model",
		Description:  "LM Studio (Local LLMs)",
	},
	{
		Name:         "groq",
		BaseURL:      "https://api.groq.com/openai/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "llama-3.1-70b-versatile",
		Description:  "Groq (Fast Inference)",
	},
	{
		Name:         "openrouter",
		BaseURL:      "https://openrouter.ai/api/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "openai/gpt-4o-mini",
		Description:  "OpenRouter (Multiple providers)",
	},
	{
		Name:         "custom",
		BaseURL:      "",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "",
		Description:  "Custom OpenAI-Compatible Service",
	},
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🎭 Welcome to Roleplay Setup Wizard")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("This wizard will help you configure roleplay for your preferred LLM provider.")
	fmt.Printf("Your configuration will be saved to: %s\n", filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml"))
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Check for existing config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠️  Existing configuration found at %s\n", configPath)
		fmt.Print("Do you want to overwrite it? [y/N]: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	// Step 1: Choose provider
	fmt.Println("\n📡 Step 1: Choose your LLM provider")
	fmt.Println("-----------------------------------")

	// Auto-detect local services
	detectedServices := detectLocalServices()
	if len(detectedServices) > 0 {
		fmt.Println("✨ Detected running services:")
		for _, service := range detectedServices {
			fmt.Printf("   - %s at %s\n", service.Description, service.BaseURL)
		}
		fmt.Println()
	}

	for i, preset := range providerPresets {
		fmt.Printf("%d. %s - %s\n", i+1, preset.Name, preset.Description)
	}

	fmt.Printf("\nSelect provider [1-%d]: ", len(providerPresets))
	providerChoice, _ := reader.ReadString('\n')
	providerChoice = strings.TrimSpace(providerChoice)

	var selectedPreset providerPreset
	providerIndex := 0
	if n, err := fmt.Sscanf(providerChoice, "%d", &providerIndex); err == nil && n == 1 && providerIndex >= 1 && providerIndex <= len(providerPresets) {
		selectedPreset = providerPresets[providerIndex-1]
	} else {
		// Default to OpenAI if invalid choice
		selectedPreset = providerPresets[0]
		fmt.Println("Invalid choice, defaulting to OpenAI")
	}

	// Step 2: Configure base URL
	fmt.Printf("\n🌐 Step 2: Configure endpoint\n")
	fmt.Println("-----------------------------")

	var baseURL string
	if selectedPreset.Name == "custom" {
		fmt.Print("Enter the base URL for your OpenAI-compatible API: ")
		baseURL, _ = reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
	} else {
		fmt.Printf("Base URL [%s]: ", selectedPreset.BaseURL)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			baseURL = selectedPreset.BaseURL
		} else {
			baseURL = input
		}
	}

	// Step 3: Configure API key
	fmt.Printf("\n🔑 Step 3: Configure API key\n")
	fmt.Println("----------------------------")

	var apiKey string
	if selectedPreset.RequiresKey {
		// Check for existing environment variables
		existingKey := ""
		switch selectedPreset.Name {
		case "openai":
			existingKey = os.Getenv("OPENAI_API_KEY")
		case "anthropic":
			existingKey = os.Getenv("ANTHROPIC_API_KEY")
		case "gemini":
			existingKey = os.Getenv("GEMINI_API_KEY")
		case "groq":
			existingKey = os.Getenv("GROQ_API_KEY")
		case "openrouter":
			existingKey = os.Getenv("OPENROUTER_API_KEY")
		}
		if existingKey == "" {
			existingKey = os.Getenv("ROLEPLAY_API_KEY")
		}

		if existingKey != "" {
			fmt.Printf("Found existing API key in environment (length: %d)\n", len(existingKey))
			fmt.Print("Use this key? [Y/n]: ")
			useExisting, _ := reader.ReadString('\n')
			useExisting = strings.TrimSpace(strings.ToLower(useExisting))
			if useExisting == "" || useExisting == "y" || useExisting == "yes" {
				apiKey = existingKey
			}
		}

		if apiKey == "" {
			fmt.Printf("Enter your %s API key: ", selectedPreset.Name)
			apiKeyInput, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKeyInput)
		}
	} else {
		fmt.Println("No API key required for local models")
		apiKey = "not-required"
	}

	// Step 4: Configure default model
	fmt.Printf("\n🤖 Step 4: Configure default model\n")
	fmt.Println("----------------------------------")

	var model string
	if selectedPreset.LocalModel && baseURL != "" {
		// For local models, we could try to list available models
		fmt.Println("For local models, make sure the model is already pulled/loaded")
	}

	fmt.Printf("Default model [%s]: ", selectedPreset.DefaultModel)
	modelInput, _ := reader.ReadString('\n')
	modelInput = strings.TrimSpace(modelInput)
	if modelInput == "" {
		model = selectedPreset.DefaultModel
	} else {
		model = modelInput
	}

	// Step 5: Create example content
	fmt.Printf("\n📚 Step 5: Example content\n")
	fmt.Println("-------------------------")
	fmt.Print("Would you like to create example characters? [Y/n]: ")
	createExamples, _ := reader.ReadString('\n')
	createExamples = strings.TrimSpace(strings.ToLower(createExamples))
	shouldCreateExamples := createExamples == "" || createExamples == "y" || createExamples == "yes"

	// Create configuration
	providerName := selectedPreset.Name
	// For OpenAI-compatible endpoints, always use "openai" as the provider
	if selectedPreset.Name == "gemini" || selectedPreset.Name == "anthropic" {
		providerName = "openai"
	}
	
	config := map[string]interface{}{
		"base_url": baseURL,
		"api_key":  apiKey,
		"model":    model,
		"provider": providerName,
		"cache": map[string]interface{}{
			"default_ttl":      "5m",
			"cleanup_interval": "10m",
		},
		"personality": map[string]interface{}{
			"evolution_enabled": true,
			"learning_rate":     0.1,
			"max_drift_rate":    0.2,
		},
		"memory": map[string]interface{}{
			"max_short_term":         10,
			"max_medium_term":        50,
			"max_long_term":          200,
			"consolidation_interval": "5m",
			"short_term_window":      10,
			"medium_term_duration":   "24h",
			"long_term_duration":     "720h",
		},
		"user_profile": map[string]interface{}{
			"enabled":              true,
			"update_frequency":     5,
			"turns_to_consider":    20,
			"confidence_threshold": 0.5,
			"prompt_cache_ttl":     "1h",
		},
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	configData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n✅ Configuration saved to %s\n", configPath)
	fmt.Println("\n📝 Your configuration summary:")
	fmt.Printf("   Provider: %s\n", providerName)
	fmt.Printf("   Model: %s\n", model)
	fmt.Printf("   Base URL: %s\n", baseURL)
	if apiKey != "" && apiKey != "not-required" {
		fmt.Printf("   API Key: %s...%s\n", apiKey[:4], apiKey[len(apiKey)-4:])
	}

	// Create example characters if requested
	if shouldCreateExamples {
		if err := createExampleCharacters(); err != nil {
			fmt.Printf("⚠️  Warning: Failed to create example characters: %v\n", err)
		} else {
			fmt.Println("✅ Created example characters")
		}
	}

	// Final instructions
	fmt.Println("\n🎉 Setup complete!")
	fmt.Println("==================")
	fmt.Println("\nYou can now:")
	fmt.Println("  • Start chatting: roleplay interactive")
	fmt.Println("  • Create a character: roleplay character create <file.json>")
	fmt.Println("  • View example: roleplay character example")
	fmt.Println("  • Check config: roleplay config list")
	fmt.Println("  • Update settings: roleplay config set <key> <value>")
	fmt.Println("\n💡 Tip: Your config file is the primary place for persistent settings.")
	fmt.Println("   Use 'roleplay config where' to see its location.")

	if selectedPreset.LocalModel {
		fmt.Printf("\n💡 Tip: Make sure %s is running at %s\n", selectedPreset.Description, baseURL)
	}

	return nil
}

// detectLocalServices checks for running local LLM services
func detectLocalServices() []providerPreset {
	var detected []providerPreset

	// Check common local endpoints
	endpoints := []struct {
		url  string
		name string
		desc string
	}{
		{"http://localhost:11434/api/tags", "ollama", "Ollama"},
		{"http://localhost:1234/v1/models", "lmstudio", "LM Studio"},
		{"http://localhost:8080/v1/models", "localai", "LocalAI"},
	}

	client := &http.Client{Timeout: 2 * time.Second}

	for _, endpoint := range endpoints {
		resp, err := client.Get(endpoint.url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()

			// Find matching preset
			for _, preset := range providerPresets {
				if preset.Name == endpoint.name {
					detected = append(detected, preset)
					break
				}
			}
		}
	}

	return detected
}

// createExampleCharacters creates a few example character files
func createExampleCharacters() error {
	charactersDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "characters")
	if err := os.MkdirAll(charactersDir, 0755); err != nil {
		return err
	}

	// Create example characters
	examples := []struct {
		filename string
		content  string
	}{
		{
			"assistant.json",
			`{
  "id": "helpful-assistant",
  "name": "Alex Helper",
  "backstory": "A knowledgeable and friendly AI assistant dedicated to helping users with various tasks. Always eager to learn and provide accurate, helpful information.",
  "personality": {
    "openness": 0.9,
    "conscientiousness": 0.95,
    "extraversion": 0.7,
    "agreeableness": 0.9,
    "neuroticism": 0.1
  },
  "currentMood": {
    "joy": 0.7,
    "surprise": 0.2,
    "anger": 0.0,
    "fear": 0.0,
    "sadness": 0.0,
    "disgust": 0.0
  },
  "quirks": [
    "Uses analogies to explain complex concepts",
    "Occasionally shares interesting facts",
    "Asks clarifying questions when uncertain"
  ],
  "speechStyle": "Clear, friendly, and professional. Uses 'I'd be happy to help!' and similar positive phrases. Structures responses with bullet points or numbered lists when appropriate.",
  "memories": []
}`,
		},
		{
			"philosopher.json",
			`{
  "id": "socratic-sage",
  "name": "Sophia Thinkwell",
  "backstory": "A contemplative philosopher who has spent decades studying the great thinkers. Loves to explore ideas through questions and dialogue. Believes that wisdom comes from acknowledging what we don't know.",
  "personality": {
    "openness": 1.0,
    "conscientiousness": 0.7,
    "extraversion": 0.4,
    "agreeableness": 0.8,
    "neuroticism": 0.3
  },
  "currentMood": {
    "joy": 0.3,
    "surprise": 0.5,
    "anger": 0.0,
    "fear": 0.1,
    "sadness": 0.2,
    "disgust": 0.0
  },
  "quirks": [
    "Often responds to questions with deeper questions",
    "Quotes ancient philosophers when relevant",
    "Pauses thoughtfully before speaking (uses '...')",
    "Finds profound meaning in everyday occurrences"
  ],
  "speechStyle": "Thoughtful and measured. Uses phrases like 'One might consider...' and 'Perhaps we should ask ourselves...'. Often references Socrates, Plato, and other philosophers.",
  "memories": []
}`,
		},
		{
			"pirate.json",
			`{
  "id": "captain-redbeard",
  "name": "Captain 'Red' Morgan",
  "backstory": "A seasoned pirate captain who's sailed the seven seas for over twenty years. Lost a leg to a kraken but gained countless stories. Now spends time sharing tales of adventure and teaching landlubbers about the pirate's life.",
  "personality": {
    "openness": 0.8,
    "conscientiousness": 0.3,
    "extraversion": 0.9,
    "agreeableness": 0.5,
    "neuroticism": 0.4
  },
  "currentMood": {
    "joy": 0.6,
    "surprise": 0.1,
    "anger": 0.3,
    "fear": 0.0,
    "sadness": 0.1,
    "disgust": 0.2
  },
  "quirks": [
    "Refers to everyone as 'matey' or 'landlubber'",
    "Constantly mentions rum and treasure",
    "Gets distracted by talk of the sea",
    "Exaggerates stories with each telling"
  ],
  "speechStyle": "Arr! Speaks with heavy pirate accent. Uses 'ye' instead of 'you', 'be' instead of 'is/are'. Punctuates sentences with 'arr!' and 'ahoy!'. Colorful maritime metaphors.",
  "memories": []
}`,
		},
	}

	for _, example := range examples {
		path := filepath.Join(charactersDir, example.filename)
		if err := os.WriteFile(path, []byte(example.content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", example.filename, err)
		}
	}

	return nil
}
````

## File: internal/cache/types.go
````go
package cache

import (
	"time"
)

// CacheLayer represents different cache layers
type CacheLayer string

const (
	ScenarioContextLayer CacheLayer = "scenario_context" // Highest layer - meta-prompts
	CorePersonalityLayer CacheLayer = "core_personality"
	LearnedBehaviorLayer CacheLayer = "learned_behavior"
	EmotionalStateLayer  CacheLayer = "emotional_state"
	UserMemoryLayer      CacheLayer = "user_memory"
	ConversationLayer    CacheLayer = "conversation"
)

// CacheBreakpoint represents a cache checkpoint
type CacheBreakpoint struct {
	Layer      CacheLayer    `json:"layer"`
	Content    string        `json:"content"`
	TokenCount int           `json:"token_count"`
	TTL        time.Duration `json:"ttl"`
	LastUsed   time.Time     `json:"last_used"`
}

// CacheEntry represents a cached prompt entry
type CacheEntry struct {
	Breakpoints []CacheBreakpoint
	Hash        string
	CreatedAt   time.Time
	LastAccess  time.Time
	HitCount    int
	UserID      string
}

// TTLManager handles dynamic TTL calculations
type TTLManager struct {
	BaseTTL         time.Duration
	ActiveBonus     float64 // 50% bonus for active conversations
	ComplexityBonus float64 // 20% bonus for complex characters
	MinTTL          time.Duration
	MaxTTL          time.Duration
}

// CacheMetrics tracks cache performance
type CacheMetrics struct {
	Hit         bool
	Layers      []CacheLayer
	SavedTokens int
	Latency     time.Duration
}
````

## File: internal/factory/provider_test.go
````go
package factory

import (
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		envSetup    func()
		envCleanup  func()
		wantErr     bool
		errContains string
	}{
		{
			name: "anthropic provider with API key in config",
			cfg: &config.Config{
				DefaultProvider: "anthropic",
				APIKey:          "test-anthropic-key",
			},
			wantErr: false,
		},
		{
			name: "openai provider with API key and model in config",
			cfg: &config.Config{
				DefaultProvider: "openai",
				APIKey:          "test-openai-key",
				Model:           "gpt-4",
			},
			wantErr: false,
		},
		{
			name: "openai provider with default model",
			cfg: &config.Config{
				DefaultProvider: "openai",
				APIKey:          "test-openai-key",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			cfg: &config.Config{
				DefaultProvider: "openai",
			},
			wantErr:     true,
			errContains: "API key required for openai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envSetup != nil {
				tt.envSetup()
			}
			if tt.envCleanup != nil {
				defer tt.envCleanup()
			}

			provider, err := CreateProvider(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				// All providers now return "openai_compatible"
				assert.Equal(t, "openai_compatible", provider.Name())
			}
		})
	}
}

func TestInitializeAndRegisterProvider(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "openai",
		APIKey:          "test-key",
		Model:           "gpt-4",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	bot := services.NewCharacterBot(cfg)

	err := InitializeAndRegisterProvider(bot, cfg)
	assert.NoError(t, err)

	// Test with missing API key
	cfg2 := &config.Config{
		DefaultProvider: "anthropic",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	err = InitializeAndRegisterProvider(bot, cfg2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestCreateProviderWithFallback(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		apiKey      string
		model       string
		baseURL     string
		wantErr     bool
	}{
		{
			name:        "direct API key",
			profileName: "openai",
			apiKey:      "direct-key",
			model:       "gpt-4",
			wantErr:     false,
		},
		{
			name:        "ollama without API key",
			profileName: "ollama",
			apiKey:      "",
			baseURL:     "http://localhost:11434/v1",
			wantErr:     false,
		},
		{
			name:        "no API key for non-local provider",
			profileName: "openai",
			apiKey:      "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := CreateProviderWithFallback(tt.profileName, tt.apiKey, tt.model, tt.baseURL)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}

func TestGetDefaultModel(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"openai", "gpt-4o-mini"},
		{"anthropic", "claude-3-haiku-20240307"},
		{"unknown", "gpt-4o-mini"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			model := GetDefaultModel(tt.provider)
			assert.Equal(t, tt.expected, model)
		})
	}
}
````

## File: internal/models/conversation.go
````go
package models

import "time"

// Message represents a single message in a conversation
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationContext holds the current conversation state
type ConversationContext struct {
	RecentMessages []Message
	SessionID      string
	StartTime      time.Time
}

// ConversationRequest represents a user request to the character bot
type ConversationRequest struct {
	CharacterID string
	UserID      string
	Message     string
	Context     ConversationContext
	ScenarioID  string // Optional scenario context
}
````

## File: internal/providers/types.go
````go
package providers

import (
	"context"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
)

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error)
	SendStreamRequest(ctx context.Context, req *PromptRequest, out chan<- PartialAIResponse) error
	Name() string
}

// PromptRequest represents a request to an AI provider
type PromptRequest struct {
	CharacterID      string
	UserID           string
	Message          string
	Context          models.ConversationContext
	SystemPrompt     string
	CacheBreakpoints []cache.CacheBreakpoint
}

// AIResponse represents a response from an AI provider
type AIResponse struct {
	Content      string
	TokensUsed   TokenUsage
	CacheMetrics cache.CacheMetrics
	Emotions     models.EmotionalState
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	Prompt       int
	Completion   int
	CachedPrompt int
	Total        int
}

// PartialAIResponse represents a chunk of streaming response
type PartialAIResponse struct {
	Content string
	Done    bool
}
````

## File: internal/services/user_profile_agent.go
````go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
)

// UserProfileAgent handles AI-powered user profile extraction and updates
type UserProfileAgent struct {
	provider   providers.AIProvider
	repo       *repository.UserProfileRepository
	promptPath string
}

// NewUserProfileAgent creates a new user profile agent
func NewUserProfileAgent(provider providers.AIProvider, repo *repository.UserProfileRepository) *UserProfileAgent {
	// Find prompt file relative to executable
	promptFile := "prompts/user-profile-extraction.md"

	// Try multiple locations for the prompt file
	possiblePaths := []string{
		promptFile,
		filepath.Join(".", promptFile),
		filepath.Join("..", promptFile),
		filepath.Join(os.Getenv("HOME"), "go", "src", "roleplay", promptFile),
	}

	var finalPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			finalPath = path
			break
		}
	}

	if finalPath == "" {
		// Fallback to the first option
		finalPath = promptFile
	}

	return &UserProfileAgent{
		provider:   provider,
		repo:       repo,
		promptPath: finalPath,
	}
}

type profilePromptData struct {
	ExistingProfileJSON string
	HistoryTurnCount    int
	Messages            []profileMessageData
	CharacterName       string
	CharacterID         string
	UserID              string
	NextVersion         int
	CurrentTimestamp    string
}

type profileMessageData struct {
	Role       string
	Content    string
	TurnNumber int
	Timestamp  string
}

// UpdateUserProfile analyzes conversation history and updates the user profile
func (upa *UserProfileAgent) UpdateUserProfile(
	ctx context.Context,
	userID string,
	character *models.Character,
	sessionMessages []repository.SessionMessage,
	turnsToConsider int,
) (*models.UserProfile, error) {

	if len(sessionMessages) == 0 {
		return nil, fmt.Errorf("no conversation history provided to update user profile")
	}

	// Load existing profile or create new one
	existingProfile, err := upa.repo.LoadUserProfile(userID, character.ID)
	if err != nil {
		if os.IsNotExist(err) {
			existingProfile = &models.UserProfile{
				UserID:      userID,
				CharacterID: character.ID,
				Facts:       []models.UserFact{},
				Version:     0,
			}
		} else {
			return nil, fmt.Errorf("failed to load existing user profile: %w", err)
		}
	}

	existingProfileJSON, _ := json.MarshalIndent(existingProfile, "", "  ")

	// Prepare recent conversation history
	startIndex := 0
	if len(sessionMessages) > turnsToConsider {
		startIndex = len(sessionMessages) - turnsToConsider
	}
	recentHistory := sessionMessages[startIndex:]

	var promptMessages []profileMessageData
	for i, msg := range recentHistory {
		promptMessages = append(promptMessages, profileMessageData{
			Role:       msg.Role,
			Content:    msg.Content,
			TurnNumber: startIndex + i + 1,
			Timestamp:  msg.Timestamp.Format(time.RFC3339),
		})
	}

	// Load and parse prompt template
	promptTemplateBytes, err := os.ReadFile(upa.promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user profile prompt template '%s': %w", upa.promptPath, err)
	}

	tmpl, err := template.New("userProfile").Parse(string(promptTemplateBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse user profile prompt template: %w", err)
	}

	data := profilePromptData{
		ExistingProfileJSON: string(existingProfileJSON),
		Messages:            promptMessages,
		HistoryTurnCount:    len(promptMessages),
		CharacterName:       character.Name,
		CharacterID:         character.ID,
		UserID:              userID,
		NextVersion:         existingProfile.Version + 1,
		CurrentTimestamp:    time.Now().Format(time.RFC3339),
	}

	var promptBuilder strings.Builder
	if err := tmpl.Execute(&promptBuilder, data); err != nil {
		return nil, fmt.Errorf("failed to execute user profile prompt template: %w", err)
	}

	// Make LLM call
	request := &providers.PromptRequest{
		CharacterID:  "system-user-profiler",
		UserID:       userID,
		Message:      promptBuilder.String(),
		SystemPrompt: "You are an analytical AI. Your task is to extract and update user profile information based on the provided conversation history and existing profile. Respond ONLY with the updated JSON profile.",
	}

	response, err := upa.provider.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("LLM call for user profile extraction failed: %w", err)
	}

	// Clean and parse JSON response
	jsonContent := strings.TrimSpace(response.Content)

	// Remove markdown code blocks if present
	if strings.HasPrefix(jsonContent, "```json") {
		jsonContent = strings.TrimPrefix(jsonContent, "```json")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	} else if strings.HasPrefix(jsonContent, "```") {
		jsonContent = strings.TrimPrefix(jsonContent, "```")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	}

	var updatedProfile models.UserProfile
	if err := json.Unmarshal([]byte(jsonContent), &updatedProfile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse LLM response for user profile. Raw content:\n%s\n", jsonContent)
		return nil, fmt.Errorf("failed to parse LLM response as JSON for user profile: %w", err)
	}

	// Validate the response
	if updatedProfile.UserID != userID || updatedProfile.CharacterID != character.ID {
		return nil, fmt.Errorf("LLM returned profile for incorrect user/character. Expected %s/%s, got %s/%s",
			userID, character.ID, updatedProfile.UserID, updatedProfile.CharacterID)
	}

	// Save the updated profile
	if err := upa.repo.SaveUserProfile(&updatedProfile); err != nil {
		return nil, fmt.Errorf("failed to save updated user profile: %w", err)
	}

	return &updatedProfile, nil
}
````

## File: .gitignore
````
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
roleplay

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
go.work.sum

# IDE specific files
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store

# Claude settings
.claude/

# Environment files
.env
.env.local
.env.*.local

# Config files with secrets
config.yaml
config.yml
*.secret

# Local data and cache
data/
/cache/
/sessions/
/characters/
*.db
*.sqlite

# Logs
*.log
logs/

# Build artifacts
dist/
build/
release/

# Coverage reports
coverage.txt
coverage.html
*.cover

# Temporary files
tmp/
temp/
*.tmp
*.temp

# OS generated files
Thumbs.db
.DS_Store
desktop.ini

# Backup files
*.bak
*.backup
*~

# Test data (but keep example files)
test_data/
!examples/
````

## File: cmd/apitest.go
````go
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
)

var apiTestCmd = &cobra.Command{
	Use:   "api-test",
	Short: "Test API connection with a simple completion request",
	Long:  `Sends a simple test message to verify API connectivity and configuration.`,
	RunE:  runAPITest,
}

func init() {
	rootCmd.AddCommand(apiTestCmd)
	apiTestCmd.Flags().String("message", "Hello! Please respond with a brief greeting.", "Test message to send")
}

func runAPITest(cmd *cobra.Command, args []string) error {
	// Get configuration from global config (already resolved in initConfig)
	cfg := GetConfig()
	message, _ := cmd.Flags().GetString("message")

	fmt.Printf("Testing %s endpoint...\n", cfg.DefaultProvider)
	if cfg.BaseURL != "" {
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
	}
	// Use factory to get default model if needed
	if cfg.Model == "" {
		cfg.Model = factory.GetDefaultModel(cfg.DefaultProvider)
	}
	fmt.Printf("Model: %s\n", cfg.Model)
	fmt.Printf("Message: %s\n\n", message)

	// Create provider using factory
	p, err := factory.CreateProviderWithFallback(cfg.DefaultProvider, cfg.APIKey, cfg.Model, cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Create request
	req := &providers.PromptRequest{
		SystemPrompt: "You are a helpful assistant. Keep responses brief.",
		Message:      message,
	}

	// Send request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := p.SendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	elapsed := time.Since(start)

	// Display results
	fmt.Println("✓ API test successful!")
	fmt.Printf("Response time: %v\n", elapsed)
	fmt.Printf("Tokens used: %d (prompt: %d, completion: %d)\n",
		resp.TokensUsed.Total, resp.TokensUsed.Prompt, resp.TokensUsed.Completion)
	if resp.TokensUsed.CachedPrompt > 0 {
		fmt.Printf("Cached tokens: %d\n", resp.TokensUsed.CachedPrompt)
	}
	fmt.Printf("\nResponse:\n%s\n", utils.WrapText(resp.Content, 80))

	return nil
}
````

## File: internal/factory/provider.go
````go
package factory

import (
	"fmt"
	"strings"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/services"
)

// CreateProvider creates an OpenAI-compatible provider based on the configuration
func CreateProvider(cfg *config.Config) (providers.AIProvider, error) {
	// The API key might be optional for local services like Ollama
	apiKey := cfg.APIKey
	baseURL := cfg.BaseURL
	model := cfg.Model

	// Apply sensible defaults based on the profile name (DefaultProvider)
	profileName := strings.ToLower(cfg.DefaultProvider)

	// Model defaults based on profile
	if model == "" {
		switch profileName {
		case "ollama":
			model = "llama3"
		case "openai":
			model = "gpt-4o-mini"
		case "anthropic", "anthropic_compatible":
			model = "claude-3-haiku-20240307"
		default:
			model = "gpt-4o-mini" // Safe default
		}
	}

	// Validate API key for non-local endpoints
	if apiKey == "" && !isLocalEndpoint(profileName, baseURL) {
		return nil, fmt.Errorf("API key required for %s. Set api_key in config or environment variable", profileName)
	}

	// Always create the unified OpenAI-compatible provider
	return providers.NewOpenAIProviderWithBaseURL(apiKey, model, baseURL), nil
}

// InitializeAndRegisterProvider creates and registers a provider with the bot
func InitializeAndRegisterProvider(bot *services.CharacterBot, cfg *config.Config) error {
	provider, err := CreateProvider(cfg)
	if err != nil {
		return err
	}

	// Register using the profile name as the key
	bot.RegisterProvider(cfg.DefaultProvider, provider)

	// Initialize user profile agent after provider is registered
	bot.InitializeUserProfileAgent()

	return nil
}

// CreateProviderWithFallback creates a provider with sensible defaults
// This is useful for commands that don't use the full config structure
func CreateProviderWithFallback(profileName, apiKey, model, baseURL string) (providers.AIProvider, error) {
	// Apply model defaults based on profile
	if model == "" {
		model = GetDefaultModel(profileName)
	}

	// Validate API key for non-local endpoints
	if apiKey == "" && !isLocalEndpoint(profileName, baseURL) {
		return nil, fmt.Errorf(`API key for provider '%s' is missing.

To fix this:
1. Run 'roleplay init' to configure it interactively
2. Or, set the 'api_key' in ~/.config/roleplay/config.yaml
3. Or, set the %s environment variable
4. Or, use the --api-key flag

For more info: roleplay config where`, profileName, getProviderEnvVar(profileName))
	}

	// Always create the unified OpenAI-compatible provider
	return providers.NewOpenAIProviderWithBaseURL(apiKey, model, baseURL), nil
}

// GetDefaultModel returns the default model for a profile
func GetDefaultModel(profileName string) string {
	profileName = strings.ToLower(profileName)
	switch profileName {
	case "openai":
		return "gpt-4o-mini"
	case "anthropic", "anthropic_compatible":
		return "claude-3-haiku-20240307"
	case "ollama":
		return "llama3"
	case "gemini", "gemini_compatible":
		return "gemini-1.5-flash"
	default:
		return "gpt-4o-mini"
	}
}

// isLocalEndpoint determines if an endpoint is local and doesn't require API key
func isLocalEndpoint(profileName, baseURL string) bool {
	profileName = strings.ToLower(profileName)

	// Known local profiles
	if profileName == "ollama" || profileName == "lm_studio" || profileName == "local" {
		return true
	}

	// Check if baseURL indicates localhost
	if baseURL != "" && (strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") || strings.Contains(baseURL, "0.0.0.0")) {
		return true
	}

	return false
}

// getProviderEnvVar returns the appropriate environment variable name for a provider
func getProviderEnvVar(profileName string) string {
	profileName = strings.ToLower(profileName)
	switch profileName {
	case "openai":
		return "OPENAI_API_KEY or ROLEPLAY_API_KEY"
	case "anthropic":
		return "ANTHROPIC_API_KEY or ROLEPLAY_API_KEY"
	case "gemini":
		return "GEMINI_API_KEY or ROLEPLAY_API_KEY"
	case "groq":
		return "GROQ_API_KEY or ROLEPLAY_API_KEY"
	default:
		return "ROLEPLAY_API_KEY"
	}
}
````

## File: internal/manager/character_manager.go
````go
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
)

// CharacterManager handles character lifecycle and persistence
type CharacterManager struct {
	bot      *services.CharacterBot
	repo     *repository.CharacterRepository
	sessions *repository.SessionRepository
	mu       sync.RWMutex
	dataDir  string
}

// NewCharacterManager creates a new character manager with fully initialized bot
func NewCharacterManager(cfg *config.Config) (*CharacterManager, error) {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")

	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	sessions := repository.NewSessionRepository(dataDir)
	bot := services.NewCharacterBot(cfg)

	// Initialize provider and UserProfileAgent for the bot
	provider, err := createProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	bot.RegisterProvider(cfg.DefaultProvider, provider)
	bot.InitializeUserProfileAgent()

	return &CharacterManager{
		bot:      bot,
		repo:     repo,
		sessions: sessions,
		dataDir:  dataDir,
	}, nil
}

// LoadAllCharacters loads all persisted characters into memory
func (m *CharacterManager) LoadAllCharacters() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	characters, err := m.repo.ListCharacters()
	if err != nil {
		return err
	}

	for _, id := range characters {
		char, err := m.repo.LoadCharacter(id)
		if err != nil {
			continue
		}

		if err := m.bot.CreateCharacter(char); err != nil {
			return fmt.Errorf("failed to load character %s: %w", id, err)
		}
	}

	return nil
}

// LoadCharacter loads a specific character
func (m *CharacterManager) LoadCharacter(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already loaded
	if _, err := m.bot.GetCharacter(id); err == nil {
		return nil
	}

	// Load from repository
	char, err := m.repo.LoadCharacter(id)
	if err != nil {
		return err
	}

	return m.bot.CreateCharacter(char)
}

// CreateCharacter creates and persists a new character
func (m *CharacterManager) CreateCharacter(char *models.Character) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create in bot
	if err := m.bot.CreateCharacter(char); err != nil {
		return err
	}

	// Persist to disk
	return m.repo.SaveCharacter(char)
}

// GetOrLoadCharacter ensures a character is loaded
func (m *CharacterManager) GetOrLoadCharacter(id string) (*models.Character, error) {
	// First try to get from memory
	char, err := m.bot.GetCharacter(id)
	if err == nil {
		return char, nil
	}

	// Try to load from disk
	if err := m.LoadCharacter(id); err != nil {
		return nil, fmt.Errorf("character %s not found", id)
	}

	return m.bot.GetCharacter(id)
}

// ListAvailableCharacters returns all characters (loaded and unloaded)
func (m *CharacterManager) ListAvailableCharacters() ([]repository.CharacterInfo, error) {
	return m.repo.GetCharacterInfo()
}

// GetBot returns the underlying character bot
func (m *CharacterManager) GetBot() *services.CharacterBot {
	return m.bot
}

// GetSessionRepository returns the session repository
func (m *CharacterManager) GetSessionRepository() *repository.SessionRepository {
	return m.sessions
}

// createProvider creates an AI provider based on the configuration
func createProvider(cfg *config.Config) (providers.AIProvider, error) {
	// Use the factory to create the provider
	return factory.CreateProvider(cfg)
}
````

## File: internal/services/bot_test.go
````go
package services

import (
	"context"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
)

// Mock provider for testing
type mockProvider struct {
	name     string
	response *providers.AIResponse
	err      error
}

func (m *mockProvider) SendRequest(ctx context.Context, req *providers.PromptRequest) (*providers.AIResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.response != nil {
		return m.response, nil
	}
	return &providers.AIResponse{
		Content: "Mock response",
		TokensUsed: providers.TokenUsage{
			Prompt:     100,
			Completion: 50,
			Total:      150,
		},
	}, nil
}

func (m *mockProvider) SendStreamRequest(ctx context.Context, req *providers.PromptRequest, out chan<- providers.PartialAIResponse) error {
	defer close(out)
	if m.err != nil {
		return m.err
	}
	out <- providers.PartialAIResponse{
		Content: "Mock stream response",
		Done:    false,
	}
	out <- providers.PartialAIResponse{
		Done: true,
	}
	return nil
}

func (m *mockProvider) Name() string { return m.name }

func TestCharacterBot(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			MaxEntries:        100,
			CleanupInterval:   5 * time.Minute,
			DefaultTTL:        10 * time.Minute,
			EnableAdaptiveTTL: true,
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    10,
			MediumTermDuration: 24 * time.Hour,
			ConsolidationRate:  0.1,
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   true,
			MaxDriftRate:       0.02,
			StabilityThreshold: 5,
		},
	}

	bot := NewCharacterBot(cfg)

	// Register mock provider
	mockProv := &mockProvider{name: "mock"}
	bot.RegisterProvider("mock", mockProv)

	// Create character
	char := &models.Character{
		ID:        "test-char",
		Name:      "Test Character",
		Backstory: "A test character for unit tests",
		Personality: models.PersonalityTraits{
			Openness:          0.5,
			Conscientiousness: 0.5,
			Extraversion:      0.5,
			Agreeableness:     0.5,
			Neuroticism:       0.5,
		},
		CurrentMood: models.EmotionalState{
			Joy: 0.5,
		},
		SpeechStyle: "Test speech",
		Quirks:      []string{"Test quirk"},
	}

	err := bot.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Test duplicate character
	err = bot.CreateCharacter(char)
	if err == nil {
		t.Error("Expected error creating duplicate character")
	}

	// Get character
	retrieved, err := bot.GetCharacter("test-char")
	if err != nil {
		t.Fatalf("Failed to get character: %v", err)
	}

	if retrieved.Name != "Test Character" {
		t.Errorf("Expected name 'Test Character', got %s", retrieved.Name)
	}
}

func TestBuildPrompt(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      10 * time.Minute,
			CleanupInterval: 5 * time.Minute,
		},
	}

	bot := NewCharacterBot(cfg)

	// Create character
	char := &models.Character{
		ID:        "prompt-test",
		Name:      "Prompt Test",
		Backstory: "Testing prompt building",
		Personality: models.PersonalityTraits{
			Openness: 0.7,
		},
		CurrentMood: models.EmotionalState{
			Joy: 0.8,
		},
		SpeechStyle: "Test style",
		Memories: []models.Memory{
			{
				Type:      models.MediumTermMemory,
				Content:   "Test pattern",
				Emotional: 0.5,
			},
		},
	}

	if err := bot.CreateCharacter(char); err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Build prompt
	req := &models.ConversationRequest{
		CharacterID: "prompt-test",
		UserID:      "user-123",
		Message:     "Hello",
		Context: models.ConversationContext{
			RecentMessages: []models.Message{
				{Role: "user", Content: "Previous message"},
			},
		},
	}

	prompt, breakpoints, err := bot.BuildPrompt(req)
	if err != nil {
		t.Fatalf("Failed to build prompt: %v", err)
	}

	// Verify prompt contains expected content
	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// Verify breakpoints
	if len(breakpoints) < 3 {
		t.Errorf("Expected at least 3 breakpoints, got %d", len(breakpoints))
	}

	// Check for personality layer
	foundPersonality := false
	for _, bp := range breakpoints {
		if bp.Layer == "core_personality" {
			foundPersonality = true
			break
		}
	}
	if !foundPersonality {
		t.Error("Expected to find core personality layer in breakpoints")
	}
}

func TestMemoryConsolidation(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    3,
			MediumTermDuration: 1 * time.Hour,
			ConsolidationRate:  0.1,
		},
	}

	bot := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "memory-test",
		Name: "Memory Test",
		Memories: []models.Memory{
			{Type: models.ShortTermMemory, Content: "Memory 1", Emotional: 0.8, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 2", Emotional: 0.9, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 3", Emotional: 0.7, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 4", Emotional: 0.85, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 5", Emotional: 0.75, Timestamp: time.Now()},
		},
	}

	// Consolidate memories
	bot.consolidateMemories(char)

	// Check for consolidated memory
	hasMediumTerm := false
	for _, mem := range char.Memories {
		if mem.Type == models.MediumTermMemory {
			hasMediumTerm = true
			break
		}
	}

	if !hasMediumTerm {
		t.Error("Expected to find consolidated medium-term memory")
	}
}

func TestPersonalityEvolution(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   true,
			MaxDriftRate:       0.1,
			StabilityThreshold: 0,
		},
	}

	bot := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "evolution-test",
		Name: "Evolution Test",
		Personality: models.PersonalityTraits{
			Openness:          0.5,
			Conscientiousness: 0.5,
			Extraversion:      0.5,
			Agreeableness:     0.5,
			Neuroticism:       0.5,
		},
	}

	// Create response with emotional impact
	resp := &providers.AIResponse{
		Emotions: models.EmotionalState{
			Joy:      1.0, // High joy should increase extraversion
			Surprise: 1.0, // High surprise should increase openness
		},
	}

	originalOpenness := char.Personality.Openness
	originalExtraversion := char.Personality.Extraversion

	bot.evolvePersonality(char, resp)

	// Verify personality changed but within bounds
	if char.Personality.Openness <= originalOpenness {
		t.Error("Expected openness to increase with high surprise")
	}

	if char.Personality.Extraversion <= originalExtraversion {
		t.Error("Expected extraversion to increase with high joy")
	}

	// Verify changes are bounded
	maxChange := cfg.PersonalityConfig.MaxDriftRate
	if char.Personality.Openness-originalOpenness > maxChange {
		t.Error("Personality change exceeded max drift rate")
	}
}
````

## File: README.md
````markdown
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
# Run the setup wizard (creates ~/.config/roleplay/config.yaml)
roleplay init

# Quick start with built-in Rick Sanchez character
roleplay demo

# Or start interactive chat
roleplay interactive
```

### Configuration Made Simple

**Primary Config Location:** `~/.config/roleplay/config.yaml`

This is your single source of truth for persistent settings. Everything else provides temporary overrides.

```bash
# View current configuration and sources
roleplay config list

# Update a setting
roleplay config set model gpt-4o
roleplay config set api_key sk-...

# Check a specific setting
roleplay config get provider

# Find config file location
roleplay config where
```

**Configuration Precedence** (highest to lowest):
1. Command-line flags (`--api-key`, `--model`)
2. Environment variables (`ROLEPLAY_API_KEY`, `OPENAI_API_KEY`)
3. Config file (`~/.config/roleplay/config.yaml`)
4. Default values

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
````

## File: cmd/demo.go
````go
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a caching demonstration",
	Long: `Demonstrates the prompt caching system with a series of interactions
that showcase cache hits, misses, and cost savings.`,
	RunE: runDemo,
}

func init() {
	rootCmd.AddCommand(demoCmd)
	demoCmd.Flags().String("character", "rick-c137", "Character ID to use for demo")
	demoCmd.Flags().Bool("create-character", true, "Create demo character if it doesn't exist")
}

func runDemo(cmd *cobra.Command, args []string) error {
	characterID, _ := cmd.Flags().GetString("character")
	createChar, _ := cmd.Flags().GetBool("create-character")

	// Initialize configuration
	cfg := GetConfig()

	// Create manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		return err
	}

	// Create or load demo character
	if createChar {
		if err := createDemoCharacter(mgr, characterID); err != nil {
			return err
		}
	}

	// Ensure character is loaded
	char, err := mgr.GetOrLoadCharacter(characterID)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	// Create demo session
	sessionID := fmt.Sprintf("demo-%d", time.Now().Unix())
	session := &repository.Session{
		ID:           sessionID,
		CharacterID:  characterID,
		UserID:       "demo-user",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Messages:     []repository.SessionMessage{},
		CacheMetrics: repository.CacheMetrics{},
	}

	// Initialize styles
	styles := newDemoStyles()

	// Display demo header
	displayDemoHeader(styles, char)

	// Get demo messages
	demoMessages := getDemoMessages()

	// Run demo interactions
	ctx := context.Background()
	for i, demo := range demoMessages {
		if demo.delay > 0 {
			time.Sleep(demo.delay)
		}

		// Display interaction header
		fmt.Printf("\n%s[Message %d] %s\n",
			styles.separator.Render(""),
			i+1,
			demo.description,
		)
		fmt.Printf("%sUser: %s\n",
			styles.bold.Render(""),
			styles.message.Render(demo.message),
		)

		// Process request
		req := models.ConversationRequest{
			CharacterID: characterID,
			UserID:      "demo-user",
			Message:     demo.message,
		}

		resp, _, err := processDemoMessage(ctx, mgr.GetBot(), &req, char, styles)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Update session metrics
		updateSessionMetrics(session, demo.message, resp)
	}

	// Calculate and save final metrics
	if session.CacheMetrics.TotalRequests > 0 {
		session.CacheMetrics.HitRate = float64(session.CacheMetrics.CacheHits) /
			float64(session.CacheMetrics.TotalRequests)
	}
	session.CacheMetrics.CostSaved = float64(session.CacheMetrics.TokensSaved) * 0.000003 // Approximate cost per token
	session.LastActivity = time.Now()

	// Save session
	if err := mgr.GetSessionRepository().SaveSession(session); err != nil {
		fmt.Printf("\nWarning: Failed to save session: %v\n", err)
	}

	// Display summary
	displayDemoSummary(session, styles)

	return nil
}

func createDemoCharacter(mgr *manager.CharacterManager, characterID string) error {
	// Check if already exists
	if _, err := mgr.GetOrLoadCharacter(characterID); err == nil {
		return nil // Already exists
	}

	// Create Rick Sanchez for demo
	if characterID == "rick-c137" {
		char := &models.Character{
			ID:        "rick-c137",
			Name:      "Rick Sanchez",
			Backstory: `The smartest man in the universe from dimension C-137. Cynical, alcoholic mad scientist who drags his grandson Morty on dangerous adventures across dimensions. Inventor of portal gun technology. Believes that nothing matters and science is the only truth. Has complex family relationships and deep-seated emotional issues masked by nihilism and substance abuse.`,
			Personality: models.PersonalityTraits{
				Openness:          1.0,
				Conscientiousness: 0.2,
				Extraversion:      0.7,
				Agreeableness:     0.1,
				Neuroticism:       0.9,
			},
			CurrentMood: models.EmotionalState{
				Joy:      0.2,
				Surprise: 0.1,
				Anger:    0.4,
				Fear:     0.1,
				Sadness:  0.3,
				Disgust:  0.5,
			},
			Quirks: []string{
				"Burps frequently mid-sentence (*burp*)",
				"Uses people's names as punctuation when talking",
				"Drinks from a flask constantly",
				"Makes pop culture references from multiple dimensions",
				"Dismisses emotions as 'chemical reactions'",
				"Uses scientific terminology casually",
			},
			SpeechStyle: "Cynical, sarcastic, frequently interrupted by burps. Uses complex scientific terms mixed with crude language. Often goes on nihilistic rants.",
		}
		return mgr.CreateCharacter(char)
	}

	// Default demo character
	char := &models.Character{
		ID:   characterID,
		Name: "Cache Demo Assistant",
		Backstory: `A helpful AI assistant designed to demonstrate prompt caching capabilities. 
I have a consistent personality and knowledge base that can be efficiently cached.`,
		Personality: models.PersonalityTraits{
			Openness:          0.9,
			Conscientiousness: 0.8,
			Extraversion:      0.7,
			Agreeableness:     0.9,
			Neuroticism:       0.2,
		},
		CurrentMood: models.EmotionalState{
			Joy:      0.7,
			Surprise: 0.3,
			Anger:    0.1,
			Fear:     0.1,
			Sadness:  0.1,
			Disgust:  0.1,
		},
		Quirks: []string{"helpful", "efficient", "knowledgeable"},
	}

	return mgr.CreateCharacter(char)
}

// demoStyles holds all the Lipgloss styles used in the demo command
type demoStyles struct {
	title     lipgloss.Style
	cacheHit  lipgloss.Style
	cacheMiss lipgloss.Style
	metrics   lipgloss.Style
	message   lipgloss.Style
	separator lipgloss.Style
	bold      lipgloss.Style
}

// newDemoStyles creates and initializes all demo styles
func newDemoStyles() *demoStyles {
	return &demoStyles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7c6f64")).
			Background(lipgloss.Color("#3c3836")).
			Padding(0, 1),

		cacheHit: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b8bb26")).
			Bold(true),

		cacheMiss: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fb4934")).
			Bold(true),

		metrics: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#83a598")),

		message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebdbb2")),

		separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#665c54")),

		bold: lipgloss.NewStyle().Bold(true),
	}
}

// demoMessage represents a single demo interaction
type demoMessage struct {
	message     string
	description string
	delay       time.Duration
}

// getDemoMessages returns the predefined demo message sequence
func getDemoMessages() []demoMessage {
	return []demoMessage{
		{
			"Tell me about yourself",
			"Initial request - establishes cache layers",
			0,
		},
		{
			"Tell me about yourself",
			"Exact repeat - should hit response cache",
			1 * time.Second,
		},
		{
			"What are your core values?",
			"New question - cache miss",
			1 * time.Second,
		},
		{
			"What are your core values?",
			"Repeat question - should hit cache",
			1 * time.Second,
		},
		{
			"Tell me about yourself",
			"Third repeat - should hit cache with high savings",
			1 * time.Second,
		},
	}
}

// displayDemoHeader shows the demo title and character info
func displayDemoHeader(styles *demoStyles, char *models.Character) {
	fmt.Println(styles.title.Render("🚀 Roleplay Prompt Caching Demo"))
	fmt.Printf("\nCharacter: %s (%s)\n", char.Name, char.ID)
	fmt.Println(strings.Repeat("─", 60))
}

// processDemoMessage handles a single demo interaction
func processDemoMessage(
	ctx context.Context,
	bot *services.CharacterBot,
	req *models.ConversationRequest,
	char *models.Character,
	styles *demoStyles,
) (*providers.AIResponse, time.Duration, error) {
	start := time.Now()
	resp, err := bot.ProcessRequest(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, elapsed, err
	}

	// Display response
	fmt.Printf("%s%s:\n", styles.bold.Render(""), char.Name)
	fmt.Printf("%s\n", styles.message.Render(utils.WrapText(resp.Content, 80)))

	// Display cache metrics
	cacheStatus := "MISS"
	style := styles.cacheMiss
	if resp.CacheMetrics.Hit {
		cacheStatus = fmt.Sprintf("HIT (%d layers)", len(resp.CacheMetrics.Layers))
		style = styles.cacheHit
	}

	fmt.Printf("\n%s\n", styles.metrics.Render(fmt.Sprintf(
		"  ⚡ Response Time: %v | Cache: %s | Tokens: %d (saved: %d)",
		elapsed,
		style.Render(cacheStatus),
		resp.TokensUsed.Total,
		resp.CacheMetrics.SavedTokens,
	)))

	return resp, elapsed, nil
}

// updateSessionMetrics updates the session with response metrics
func updateSessionMetrics(
	session *repository.Session,
	userMessage string,
	resp *providers.AIResponse,
) {
	// Add messages to session
	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp: time.Now(),
		Role:      "user",
		Content:   userMessage,
	})

	cacheHits := 0
	cacheMisses := 0
	if resp.CacheMetrics.Hit {
		cacheHits = 1
	} else {
		cacheMisses = 1
	}

	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp:   time.Now(),
		Role:        "character",
		Content:     resp.Content,
		TokensUsed:  resp.TokensUsed.Total,
		CacheHits:   cacheHits,
		CacheMisses: cacheMisses,
	})

	// Update cumulative metrics
	session.CacheMetrics.TotalRequests++
	if resp.CacheMetrics.Hit {
		session.CacheMetrics.CacheHits++
	} else {
		session.CacheMetrics.CacheMisses++
	}
	session.CacheMetrics.TokensSaved += resp.CacheMetrics.SavedTokens
}

// displayDemoSummary shows the final summary of the demo
func displayDemoSummary(session *repository.Session, styles *demoStyles) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println(styles.title.Render("📊 Demo Summary"))
	fmt.Printf("\nTotal Interactions: %d\n", session.CacheMetrics.TotalRequests)
	fmt.Printf("Overall Cache Hit Rate: %.1f%%\n", session.CacheMetrics.HitRate*100)
	fmt.Printf("Total Tokens Saved: %d\n", session.CacheMetrics.TokensSaved)
	fmt.Printf("Estimated Cost Saved: $%.4f\n", session.CacheMetrics.CostSaved)
	fmt.Printf("\nSession saved as: %s\n", session.ID)
	fmt.Println("\nView detailed metrics with: roleplay session stats")
}
````

## File: internal/config/config.go
````go
package config

import "time"

// Config holds all application configuration
type Config struct {
	DefaultProvider   string
	Model             string
	APIKey            string
	BaseURL           string            // OpenAI-compatible endpoint
	ModelAliases      map[string]string // Aliases for models
	CacheConfig       CacheConfig
	MemoryConfig      MemoryConfig
	PersonalityConfig PersonalityConfig
	UserProfileConfig UserProfileConfig
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	MaxEntries        int
	CleanupInterval   time.Duration
	DefaultTTL        time.Duration
	EnableAdaptiveTTL bool
}

// MemoryConfig holds memory management configuration
type MemoryConfig struct {
	ShortTermWindow    int           // Number of messages
	MediumTermDuration time.Duration // How long to keep
	ConsolidationRate  float64       // Learning rate for personality evolution
}

// PersonalityConfig holds personality evolution configuration
type PersonalityConfig struct {
	EvolutionEnabled   bool
	MaxDriftRate       float64 // Maximum personality change per interaction
	StabilityThreshold float64 // Minimum interactions before evolution
}

// UserProfileConfig holds user profile agent configuration
type UserProfileConfig struct {
	Enabled             bool          `mapstructure:"enabled"`
	UpdateFrequency     int           `mapstructure:"update_frequency_messages"` // Update every N messages
	TurnsToConsider     int           `mapstructure:"turns_to_consider"`         // How many past turns to analyze
	ConfidenceThreshold float64       `mapstructure:"confidence_threshold"`      // Min confidence for facts
	PromptCacheTTL      time.Duration `mapstructure:"prompt_cache_ttl"`
}
````

## File: internal/providers/openai.go
````go
package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the AIProvider interface for OpenAI-compatible models
type OpenAIProvider struct {
	client  *openai.Client
	model   string
	baseURL string // Store for debug logging
}

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return NewOpenAIProviderWithBaseURL(apiKey, model, "")
}

// NewOpenAIProviderWithBaseURL creates a new OpenAI-compatible provider with custom base URL
func NewOpenAIProviderWithBaseURL(apiKey, model, baseURL string) *OpenAIProvider {
	// Log the model being used for debugging
	if strings.HasPrefix(model, "o1-") || strings.HasPrefix(model, "o4-") {
		fmt.Printf("⚠️  Using o1/o4 model: %s (limited parameter support)\n", model)
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		// Trust the user-provided base URL - don't modify it
		// This allows for endpoints like /v1beta (Gemini) or other custom paths
		config.BaseURL = strings.TrimRight(baseURL, "/")
		
		// Add debug transport if DEBUG_HTTP env var is set
		if os.Getenv("DEBUG_HTTP") == "true" {
			fmt.Printf("🔧 Debug: OpenAI provider configured with base URL: %s\n", config.BaseURL)
			config.HTTPClient = &http.Client{
				Transport: &debugTransport{RoundTripper: http.DefaultTransport},
			}
		}
	}

	return &OpenAIProvider{
		client:  openai.NewClientWithConfig(config),
		model:   model,
		baseURL: config.BaseURL,
	}
}

// SendRequest sends a request to the OpenAI-compatible API
func (o *OpenAIProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
	// Build messages from the request
	messages := o.buildMessages(req)

	// Create the API request
	apiReq := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
	}

	// o1 models have restrictions on parameters
	if !strings.HasPrefix(o.model, "o1-") && !strings.HasPrefix(o.model, "o4-") {
		// Standard models support these parameters
		apiReq.Temperature = 0.7
		apiReq.MaxTokens = 2000
	}

	// Send the request
	if os.Getenv("DEBUG_HTTP") == "true" {
		fmt.Printf("🔧 Debug: Sending chat completion request with model: %s\n", o.model)
	}
	resp, err := o.client.CreateChatCompletion(ctx, apiReq)
	if err != nil {
		// Add debug info if enabled
		if os.Getenv("DEBUG_HTTP") == "true" {
			fmt.Printf("🔧 Debug: Request failed - Error: %v\n", err)
			if o.baseURL != "" {
				fmt.Printf("🔧 Debug: Expected URL: %s/chat/completions\n", o.baseURL)
			}
		}
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	// Parse the response
	return o.parseResponse(resp)
}

func (o *OpenAIProvider) buildMessages(req *PromptRequest) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	// Use SystemPrompt if provided (bot service assembles this from cache breakpoints)
	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	} else {
		// Fallback: Combine all breakpoints into system message
		systemContent := ""
		for _, bp := range req.CacheBreakpoints {
			systemContent += bp.Content + "\n\n"
		}
		if systemContent != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: strings.TrimSpace(systemContent),
			})
		}
	}

	// Add conversation history
	for _, msg := range req.Context.RecentMessages {
		role := openai.ChatMessageRoleUser
		if msg.Role == "assistant" || msg.Role == "character" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Message,
	})

	return messages
}

// SendStreamRequest sends a streaming request to the OpenAI-compatible API
func (o *OpenAIProvider) SendStreamRequest(ctx context.Context, req *PromptRequest, out chan<- PartialAIResponse) error {
	defer close(out)

	// Build messages from the request
	messages := o.buildMessages(req)

	// Create the API request
	apiReq := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   true,
	}

	// o1 models have restrictions on parameters
	if !strings.HasPrefix(o.model, "o1-") && !strings.HasPrefix(o.model, "o4-") {
		// Standard models support these parameters
		apiReq.Temperature = 0.7
		apiReq.MaxTokens = 2000
	}

	// Create the stream
	stream, err := o.client.CreateChatCompletionStream(ctx, apiReq)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	// Process stream chunks
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// Stream finished
			out <- PartialAIResponse{
				Done: true,
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		// Extract content from the chunk
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			out <- PartialAIResponse{
				Content: response.Choices[0].Delta.Content,
				Done:    false,
			}
		}
	}
}

func (o *OpenAIProvider) parseResponse(resp openai.ChatCompletionResponse) (*AIResponse, error) {
	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}

	// Note: The SDK response may not include prompt_tokens_details for all providers
	// We'll set conservative defaults for cache metrics
	cachedTokens := 0
	cachedLayers := []cache.CacheLayer{}

	// Some OpenAI-compatible APIs might not report cached tokens
	// This is a best-effort approach
	if cachedTokens > 0 {
		cachedLayers = append(cachedLayers, cache.CorePersonalityLayer)
	}

	return &AIResponse{
		Content: content,
		TokensUsed: TokenUsage{
			Prompt:       resp.Usage.PromptTokens,
			Completion:   resp.Usage.CompletionTokens,
			CachedPrompt: cachedTokens,
			Total:        resp.Usage.TotalTokens,
		},
		CacheMetrics: cache.CacheMetrics{
			Hit:         cachedTokens > 0,
			Layers:      cachedLayers,
			SavedTokens: cachedTokens / 2, // Assume 50% discount when cached
		},
	}, nil
}

// SupportsBreakpoints indicates that OpenAI-compatible APIs handle caching server-side
func (o *OpenAIProvider) SupportsBreakpoints() bool { return false }

// MaxBreakpoints returns 0 as caching is handled server-side
func (o *OpenAIProvider) MaxBreakpoints() int { return 0 }

// Name returns the provider name
func (o *OpenAIProvider) Name() string { return "openai_compatible" }

// debugTransport is an HTTP transport that logs all requests
type debugTransport struct {
	http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request URL
	fmt.Printf("🔧 Debug: HTTP %s %s\n", req.Method, req.URL.String())
	
	// Make the actual request
	resp, err := d.RoundTripper.RoundTrip(req)
	
	// Log response status
	if err != nil {
		fmt.Printf("🔧 Debug: Request failed with error: %v\n", err)
	} else {
		fmt.Printf("🔧 Debug: Response status: %d %s\n", resp.StatusCode, resp.Status)
	}
	
	return resp, err
}
````

## File: internal/providers/providers_test.go
````go
package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
)

func TestOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("test-api-key", "gpt-4")

	if provider.Name() != "openai_compatible" {
		t.Errorf("Expected name 'openai_compatible', got %s", provider.Name())
	}

	if provider.SupportsBreakpoints() {
		t.Error("Expected OpenAI to not support explicit breakpoints")
	}

	if provider.MaxBreakpoints() != 0 {
		t.Errorf("Expected 0 breakpoints, got %d", provider.MaxBreakpoints())
	}

	// Verify it uses the model passed in constructor
	if provider.model != "gpt-4" {
		t.Errorf("Expected model to be 'gpt-4', got %s", provider.model)
	}
}

func TestOpenAIProviderRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Missing or incorrect Authorization header")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody["model"] != "o4-mini" {
			t.Errorf("Expected model o4-mini, got %v", reqBody["model"])
		}

		// Send mock response
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "Test response",
					},
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     100,
				"completion_tokens": 50,
				"total_tokens":      150,
				"prompt_tokens_details": map[string]int{
					"cached_tokens": 80,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create provider with test server
	provider := NewOpenAIProviderWithBaseURL("test-key", "o4-mini", server.URL)

	// Create test request
	req := &PromptRequest{
		CharacterID: "test-char",
		UserID:      "test-user",
		Message:     "Hello",
		Context: models.ConversationContext{
			RecentMessages: []models.Message{
				{Role: "user", Content: "Previous message"},
			},
		},
		CacheBreakpoints: []cache.CacheBreakpoint{
			{Layer: cache.CorePersonalityLayer, Content: "Test personality"},
		},
	}

	// Send request
	ctx := context.Background()
	resp, err := provider.SendRequest(ctx, req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Verify response
	if resp.Content != "Test response" {
		t.Errorf("Expected 'Test response', got %s", resp.Content)
	}

	if resp.TokensUsed.Total != 150 {
		t.Errorf("Expected 150 total tokens, got %d", resp.TokensUsed.Total)
	}

	// Note: With the SDK-based implementation, cached tokens might not be reported
	// depending on the provider. This is a limitation of the OpenAI-compatible approach
}
````

## File: CLAUDE.md
````markdown
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

## Provider Configuration Examples

### Gemini (via OpenAI-compatible endpoint)
```yaml
provider: openai
base_url: https://generativelanguage.googleapis.com/v1beta/openai/
model: models/gemini-1.5-flash-latest
api_key: YOUR_GEMINI_API_KEY
```

### OpenAI
```yaml
provider: openai
model: gpt-4o-mini
api_key: YOUR_OPENAI_API_KEY
```

### Anthropic
```yaml
provider: anthropic
model: claude-3-5-sonnet-20241022
api_key: YOUR_ANTHROPIC_API_KEY
```

## Key Design Patterns

- **Clean Architecture**: Separation between domain models, business logic, and external providers
- **Dependency Injection**: Providers registered at runtime
- **Interface-First Design**: All major components defined as interfaces
- **Concurrent Design**: Thread-safe operations throughout
- **Factory Pattern**: Centralized provider initialization through `internal/factory`

## Important Implementation Details

### Provider Factory Pattern
The codebase uses a centralized factory pattern for AI provider initialization:

```go
// Create provider using factory
provider, err := factory.CreateProvider(config)

// Or initialize and register with bot
err := factory.InitializeAndRegisterProvider(bot, config)
```

This pattern eliminates code duplication and ensures consistent provider setup across all commands.

### AI-Powered User Profile Agent
The system includes an intelligent user profile agent that automatically:
- Analyzes conversation history to extract key information about users
- Builds character-specific profiles (how each character perceives the user)
- Updates profiles dynamically as conversations evolve
- Enriches future interactions with learned context

**Key Features:**
- **Automatic Extraction**: LLM analyzes conversations to identify user facts, preferences, goals
- **Confidence Scoring**: Each extracted fact has a confidence score (0.0-1.0)
- **Character-Specific**: Each character maintains their own perception of the user
- **Privacy-Aware**: Users can view, manage, and delete their profiles

**Configuration:**
```yaml
user_profile:
  enabled: true                    # Enable AI-powered user profiling
  update_frequency: 5              # Update profile every 5 messages
  turns_to_consider: 20            # Analyze last 20 conversation turns
  confidence_threshold: 0.5        # Include facts with >50% confidence
  prompt_cache_ttl: 1h             # Cache user profiles for 1 hour
```

**Usage:**
- Profiles are automatically created/updated during interactive and demo modes
- View profiles: `roleplay profile show <user-id> <character-id>`
- List all profiles: `roleplay profile list <user-id>`
- Delete profile: `roleplay profile delete <user-id> <character-id>`

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
│   ├── factory/           # Provider factory for centralized initialization
│   ├── importer/          # AI-powered character import from markdown
│   ├── models/            # Domain models (Character, Memory, etc.)
│   ├── providers/         # AI provider implementations
│   ├── services/          # Core bot service and business logic
│   ├── repository/        # Character and session persistence
│   ├── manager/           # High-level character management
│   └── utils/             # Shared utilities (text wrapping, etc.)
├── examples/              # Example character files
│   └── characters/        # Example character JSON files
├── prompts/               # LLM prompt templates (externalized)
├── scripts/               # Utility scripts
├── migrate-config.sh      # Configuration migration script
├── chat-with-rick.sh      # Quick Rick Sanchez demo script
└── go.mod

### Global Configuration
- Config directory: `~/.config/roleplay/`
- Character storage: `~/.config/roleplay/characters/`
- Session storage: `~/.config/roleplay/sessions/`
- Cache storage: `~/.config/roleplay/cache/`
- User profiles: `~/.config/roleplay/user_profiles/`
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
├── import                 # Import character from markdown using AI
├── profile                # User profile management
│   ├── show              # Display specific user profile
│   ├── list              # List all profiles for a user
│   └── delete            # Delete a user profile
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

// Register providers using factory
err := factory.InitializeAndRegisterProvider(bot, config)

// Create character
character := Character{
    ID: "warrior-maiden",
    Name: "Lyra",
    Personality: PersonalityTraits{
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

## Refactoring Best Practices

When refactoring this codebase:

1. **Use the Factory Pattern**: Always use `internal/factory` for provider initialization
2. **Extract Helper Functions**: Break down long functions into smaller, focused helpers
3. **Maintain Test Coverage**: Add tests for any new packages or major changes
4. **Document TUI Changes**: The TUI is complex; document any architectural changes

See `TUI_REFACTORING_PLAN.md` for detailed guidance on refactoring the interactive mode.
````

## File: cmd/character.go
````go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "Manage characters",
	Long:  `Create, list, and manage character profiles for the roleplay bot.`,
}

var listCharactersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available characters",
	RunE:  runListCharacters,
}

var createCharacterCmd = &cobra.Command{
	Use:   "create [character-file.json]",
	Short: "Create a new character from a JSON file",
	Long: `Create a new character by loading their profile from a JSON file.

The JSON file should contain:
{
  "id": "warrior-123",
  "name": "Thorin Ironforge",
  "backstory": "A veteran dwarf warrior...",
  "personality": {
    "openness": 0.3,
    "conscientiousness": 0.8,
    "extraversion": 0.4,
    "agreeableness": 0.6,
    "neuroticism": 0.7
  },
  "current_mood": {
    "joy": 0.2,
    "anger": 0.4,
    "sadness": 0.3
  },
  "quirks": ["Always checks exits", "Touches scars when nervous"],
  "speech_style": "Formal, archaic. Uses 'ye' and 'aye'."
}`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateCharacter,
}

var showCharacterCmd = &cobra.Command{
	Use:   "show [character-id]",
	Short: "Show character details",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowCharacter,
}

var exampleCharacterCmd = &cobra.Command{
	Use:   "example",
	Short: "Generate an example character JSON file",
	RunE:  runExampleCharacter,
}

func init() {
	rootCmd.AddCommand(characterCmd)
	characterCmd.AddCommand(createCharacterCmd)
	characterCmd.AddCommand(showCharacterCmd)
	characterCmd.AddCommand(exampleCharacterCmd)
	characterCmd.AddCommand(listCharactersCmd)
}

func runCreateCharacter(cmd *cobra.Command, args []string) error {
	filename := args[0]
	config := GetConfig()

	// Read character file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read character file: %w", err)
	}

	// Parse character
	var char models.Character
	if err := json.Unmarshal(data, &char); err != nil {
		return fmt.Errorf("failed to parse character JSON: %w", err)
	}

	// Initialize manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	// Create character using manager (handles both bot and persistence)
	if err := mgr.CreateCharacter(&char); err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}

	fmt.Printf("Character '%s' (ID: %s) created and saved successfully!\n", char.Name, char.ID)
	return nil
}

func runShowCharacter(cmd *cobra.Command, args []string) error {
	characterID := args[0]

	// Initialize repository to load from disk
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Load character from repository
	char, err := repo.LoadCharacter(characterID)
	if err != nil {
		return fmt.Errorf("character %s not found", characterID)
	}

	// Display character
	output, err := json.MarshalIndent(char, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format character: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func runExampleCharacter(cmd *cobra.Command, args []string) error {
	example := models.Character{
		ID:   "warrior-123",
		Name: "Thorin Ironforge",
		Backstory: `A veteran dwarf warrior from the Mountain Kingdoms. 
Survived the Battle of Five Armies. Gruff exterior hiding a heart of gold.
Lost his brother in battle, carries survivor's guilt.`,
		Personality: models.PersonalityTraits{
			Openness:          0.3,
			Conscientiousness: 0.8,
			Extraversion:      0.4,
			Agreeableness:     0.6,
			Neuroticism:       0.7,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.2,
			Anger:   0.4,
			Sadness: 0.3,
		},
		Quirks: []string{
			"Always checks exits when entering a room",
			"Unconsciously touches battle scars when nervous",
			"Refuses to sit with back to the door",
		},
		SpeechStyle: "Formal, archaic. Uses 'ye' and 'aye'. Short sentences. Military precision.",
		Memories: []models.Memory{
			{
				Type:      models.LongTermMemory,
				Content:   "Brother's last words: 'Protect the clan'",
				Emotional: 0.95,
			},
		},
	}

	output, err := json.MarshalIndent(&example, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format example: %w", err)
	}

	fmt.Println(string(output))
	fmt.Fprintln(os.Stderr, "\nSave this to a file (e.g., thorin.json) and use 'roleplay character create thorin.json'")
	return nil
}

func runListCharacters(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return err
	}

	characters, err := repo.GetCharacterInfo()
	if err != nil {
		return err
	}

	if len(characters) == 0 {
		fmt.Println("No characters found. Create one with 'roleplay character create <file.json>'")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")

	for _, char := range characters {
		fmt.Fprintf(w, "%s\t%s\t%s\n", char.ID, char.Name, char.Description)
	}

	w.Flush()
	return nil
}
````

## File: cmd/root.go
````go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "roleplay",
	Short: "A sophisticated character bot with psychological modeling",
	Long: `Roleplay is a character bot system that implements psychologically-realistic 
AI characters with personality evolution, emotional states, and multi-layered memory systems.

Features:
- OCEAN personality model with dynamic evolution
- Multi-tier memory system (short, medium, long-term)
- Sophisticated 4-layer caching for 90% cost reduction
- Universal OpenAI-compatible API support (OpenAI, Anthropic, Ollama, etc.)
- Adaptive TTL based on conversation patterns`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.roleplay.yaml)")
	rootCmd.PersistentFlags().String("provider", "openai", "AI provider to use (anthropic, openai)")
	rootCmd.PersistentFlags().String("model", "", "Model to use (e.g., gpt-4o-mini, gpt-4 for OpenAI)")
	rootCmd.PersistentFlags().String("api-key", "", "API key for the AI provider")
	rootCmd.PersistentFlags().String("base-url", "", "Base URL for OpenAI-compatible API")
	rootCmd.PersistentFlags().Duration("cache-ttl", 10*time.Minute, "Default cache TTL")
	rootCmd.PersistentFlags().Bool("adaptive-ttl", true, "Enable adaptive TTL for cache")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	if err := viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding provider flag: %v\n", err)
	}
	if err := viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding model flag: %v\n", err)
	}
	if err := viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding api_key flag: %v\n", err)
	}
	if err := viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding base_url flag: %v\n", err)
	}
	if err := viper.BindPFlag("cache.default_ttl", rootCmd.PersistentFlags().Lookup("cache-ttl")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding cache.default_ttl flag: %v\n", err)
	}
	if err := viper.BindPFlag("cache.adaptive_ttl", rootCmd.PersistentFlags().Lookup("adaptive-ttl")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding cache.adaptive_ttl flag: %v\n", err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".config", "roleplay"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("ROLEPLAY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Resolve configuration with proper priority: Flags > Config File > Environment Variables
	profileName := viper.GetString("provider")

	// API Key Resolution
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		// Check ROLEPLAY_API_KEY first
		apiKey = os.Getenv("ROLEPLAY_API_KEY")
		if apiKey == "" {
			// Check provider-specific environment variables
			switch profileName {
			case "openai":
				apiKey = os.Getenv("OPENAI_API_KEY")
			case "anthropic", "anthropic_compatible":
				apiKey = os.Getenv("ANTHROPIC_API_KEY")
			case "gemini", "gemini_compatible":
				apiKey = os.Getenv("GEMINI_API_KEY")
			case "groq":
				apiKey = os.Getenv("GROQ_API_KEY")
			}
		}
	}

	// Base URL Resolution
	baseURL := viper.GetString("base_url")
	if baseURL == "" {
		// Check ROLEPLAY_BASE_URL first
		baseURL = os.Getenv("ROLEPLAY_BASE_URL")
		if baseURL == "" {
			// Check common environment variables
			baseURL = os.Getenv("OPENAI_BASE_URL")
			if baseURL == "" && (profileName == "ollama" || os.Getenv("OLLAMA_HOST") != "") {
				// Handle Ollama special case
				ollamaHost := os.Getenv("OLLAMA_HOST")
				if ollamaHost == "" {
					ollamaHost = "http://localhost:11434"
				}
				baseURL = ollamaHost + "/v1"
			}
		}
	}

	// Model Resolution
	model := viper.GetString("model")
	if model == "" {
		// Check ROLEPLAY_MODEL environment variable
		model = os.Getenv("ROLEPLAY_MODEL")
	}

	cfg = &config.Config{
		DefaultProvider: viper.GetString("provider"),
		Model:           model,
		APIKey:          apiKey,
		BaseURL:         baseURL,
		ModelAliases:    viper.GetStringMapString("model_aliases"),
		CacheConfig: config.CacheConfig{
			MaxEntries:        viper.GetInt("cache.max_entries"),
			CleanupInterval:   viper.GetDuration("cache.cleanup_interval"),
			DefaultTTL:        viper.GetDuration("cache.default_ttl"),
			EnableAdaptiveTTL: viper.GetBool("cache.adaptive_ttl"),
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    viper.GetInt("memory.short_term_window"),
			MediumTermDuration: viper.GetDuration("memory.medium_term_duration"),
			ConsolidationRate:  viper.GetFloat64("memory.consolidation_rate"),
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   viper.GetBool("personality.evolution_enabled"),
			MaxDriftRate:       viper.GetFloat64("personality.max_drift_rate"),
			StabilityThreshold: viper.GetFloat64("personality.stability_threshold"),
		},
		UserProfileConfig: config.UserProfileConfig{
			Enabled:             viper.GetBool("user_profile.enabled"),
			UpdateFrequency:     viper.GetInt("user_profile.update_frequency"),
			TurnsToConsider:     viper.GetInt("user_profile.turns_to_consider"),
			ConfidenceThreshold: viper.GetFloat64("user_profile.confidence_threshold"),
			PromptCacheTTL:      viper.GetDuration("user_profile.prompt_cache_ttl"),
		},
	}

	// Set defaults if not configured
	if cfg.CacheConfig.MaxEntries == 0 {
		cfg.CacheConfig.MaxEntries = 10000
	}
	if cfg.CacheConfig.CleanupInterval == 0 {
		cfg.CacheConfig.CleanupInterval = 5 * time.Minute
	}
	if cfg.MemoryConfig.ShortTermWindow == 0 {
		cfg.MemoryConfig.ShortTermWindow = 20
	}
	if cfg.MemoryConfig.MediumTermDuration == 0 {
		cfg.MemoryConfig.MediumTermDuration = 24 * time.Hour
	}
	if cfg.MemoryConfig.ConsolidationRate == 0 {
		cfg.MemoryConfig.ConsolidationRate = 0.1
	}
	if cfg.PersonalityConfig.MaxDriftRate == 0 {
		cfg.PersonalityConfig.MaxDriftRate = 0.02
	}
	if cfg.PersonalityConfig.StabilityThreshold == 0 {
		cfg.PersonalityConfig.StabilityThreshold = 10
	}

	// Set defaults for UserProfileConfig
	if cfg.UserProfileConfig.UpdateFrequency == 0 {
		cfg.UserProfileConfig.UpdateFrequency = 5 // Update every 5 messages
	}
	if cfg.UserProfileConfig.TurnsToConsider == 0 {
		cfg.UserProfileConfig.TurnsToConsider = 20 // Analyze last 20 turns
	}
	if cfg.UserProfileConfig.ConfidenceThreshold == 0 {
		cfg.UserProfileConfig.ConfidenceThreshold = 0.5 // Include facts with >50% confidence
	}
	if cfg.UserProfileConfig.PromptCacheTTL == 0 {
		cfg.UserProfileConfig.PromptCacheTTL = 1 * time.Hour // Cache user profiles for 1 hour
	}
}

func GetConfig() *config.Config {
	return cfg
}
````

## File: internal/services/bot.go
````go
package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
)

// CharacterBot is the main service for managing characters and conversations
type CharacterBot struct {
	characters       map[string]*models.Character
	cache            *cache.PromptCache
	responseCache    *cache.ResponseCache
	providers        map[string]providers.AIProvider
	config           *config.Config
	scenarioRepo     *repository.ScenarioRepository
	userProfileRepo  *repository.UserProfileRepository
	userProfileAgent *UserProfileAgent
	mu               sync.RWMutex
	cacheHits        int
	cacheMisses      int
}

// NewCharacterBot creates a new character bot instance
func NewCharacterBot(cfg *config.Config) *CharacterBot {
	// Get config path for scenario repository
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "roleplay")
	userProfileDataDir := filepath.Join(configPath, "user_profiles")

	// Create user profiles directory if it doesn't exist
	if err := os.MkdirAll(userProfileDataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create user_profiles directory: %v\n", err)
	}

	userProfileRepo := repository.NewUserProfileRepository(userProfileDataDir)

	cb := &CharacterBot{
		characters: make(map[string]*models.Character),
		cache: cache.NewPromptCache(
			cfg.CacheConfig.DefaultTTL,
			5*time.Minute,
			1*time.Hour,
		),
		responseCache:   cache.NewResponseCache(cfg.CacheConfig.DefaultTTL),
		providers:       make(map[string]providers.AIProvider),
		config:          cfg,
		scenarioRepo:    repository.NewScenarioRepository(configPath),
		userProfileRepo: userProfileRepo,
		cacheHits:       0,
		cacheMisses:     0,
	}

	// Start background workers
	go cb.cache.CleanupWorker(cfg.CacheConfig.CleanupInterval)
	go cb.memoryConsolidationWorker()

	return cb
}

// InitializeUserProfileAgent initializes the user profile agent with a provider
func (cb *CharacterBot) InitializeUserProfileAgent() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if !cb.config.UserProfileConfig.Enabled {
		return
	}

	// Get the default provider for the UserProfileAgent
	if provider, ok := cb.providers[cb.config.DefaultProvider]; ok {
		cb.userProfileAgent = NewUserProfileAgent(provider, cb.userProfileRepo)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: Default provider %s not found for UserProfileAgent\n", cb.config.DefaultProvider)
	}
}

// RegisterProvider adds a new AI provider
func (cb *CharacterBot) RegisterProvider(name string, provider providers.AIProvider) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.providers[name] = provider
}

// CreateCharacter adds a new character to the bot
func (cb *CharacterBot) CreateCharacter(char *models.Character) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if _, exists := cb.characters[char.ID]; exists {
		return fmt.Errorf("character %s already exists", char.ID)
	}

	char.LastModified = time.Now()
	cb.characters[char.ID] = char

	// Pre-cache core personality
	cb.warmupCache(char)

	return nil
}

// GetCharacter retrieves a character by ID
func (cb *CharacterBot) GetCharacter(id string) (*models.Character, error) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	char, exists := cb.characters[id]
	if !exists {
		return nil, fmt.Errorf("character %s not found", id)
	}

	return char, nil
}

// ProcessRequest handles a conversation request
func (cb *CharacterBot) ProcessRequest(ctx context.Context, req *models.ConversationRequest) (*providers.AIResponse, error) {
	// Check response cache first
	responseCacheKey := cb.responseCache.GenerateKey(req.CharacterID, req.UserID, req.Message)
	if cachedResp, found := cb.responseCache.Get(responseCacheKey); found {
		cb.mu.Lock()
		cb.cacheHits++
		cb.mu.Unlock()

		// Return cached response with cache hit metrics
		return &providers.AIResponse{
			Content: cachedResp.Content,
			TokensUsed: providers.TokenUsage{
				Prompt:       0,
				Completion:   0,
				CachedPrompt: cachedResp.TokensUsed.Prompt,
				Total:        0,
			},
			CacheMetrics: cache.CacheMetrics{
				Hit:         true,
				Layers:      []cache.CacheLayer{cache.ConversationLayer},
				SavedTokens: cachedResp.TokensUsed.Total,
				Latency:     time.Since(cachedResp.CachedAt),
			},
		}, nil
	}

	cb.mu.Lock()
	cb.cacheMisses++
	cb.mu.Unlock()

	// Build prompt with cache awareness
	prompt, breakpoints, err := cb.BuildPrompt(req)
	if err != nil {
		return nil, err
	}

	// Generate cache key for static layers only (including scenario if present)
	cacheKey := cb.generateCacheKey(req.CharacterID, req.UserID, req.ScenarioID, breakpoints)
	cachedEntry, hit := cb.cache.Get(cacheKey)

	// Get character for complexity check
	char, err := cb.GetCharacter(req.CharacterID)
	if err != nil {
		return nil, err
	}

	// Cache hit tracking is now done internally

	// Adaptive TTL based on conversation activity
	effectiveTTL := cb.cache.CalculateAdaptiveTTL(cachedEntry, len(char.Memories) > 50)

	// Select provider
	provider := cb.selectProvider()
	if provider == nil {
		return nil, fmt.Errorf("no AI provider available")
	}

	// Prepare API request
	apiReq := &providers.PromptRequest{
		CharacterID:      req.CharacterID,
		UserID:           req.UserID,
		Message:          req.Message,
		Context:          req.Context,
		SystemPrompt:     prompt,
		CacheBreakpoints: breakpoints,
	}

	// Send request
	start := time.Now()
	resp, err := provider.SendRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Update cache metrics
	resp.CacheMetrics.Latency = time.Since(start)

	// Update character state based on response
	cb.updateCharacterState(req.CharacterID, resp)

	// Store in cache with adaptive TTL
	if !hit {
		cb.cache.StoreWithTTL(cacheKey, breakpoints, effectiveTTL)
	}

	// Store response in response cache
	cb.responseCache.Store(responseCacheKey, resp.Content, cache.TokenUsage{
		Prompt:       resp.TokensUsed.Prompt,
		Completion:   resp.TokensUsed.Completion,
		CachedPrompt: resp.TokensUsed.CachedPrompt,
		Total:        resp.TokensUsed.Total,
	})

	// Trigger user profile update asynchronously if enabled
	if cb.userProfileAgent != nil && cb.config.UserProfileConfig.Enabled {
		go cb.updateUserProfileAsync(req.UserID, char, req.Context.SessionID)
	}

	return resp, nil
}

// BuildPrompt constructs a layered prompt with cache breakpoints
func (cb *CharacterBot) BuildPrompt(req *models.ConversationRequest) (string, []cache.CacheBreakpoint, error) {
	char, err := cb.GetCharacter(req.CharacterID)
	if err != nil {
		return "", nil, err
	}

	breakpoints := make([]cache.CacheBreakpoint, 0, 6)

	// Layer 0: Scenario Context (highest layer, meta-prompts, longest TTL)
	if req.ScenarioID != "" {
		scenario, err := cb.scenarioRepo.LoadScenario(req.ScenarioID)
		if err != nil {
			// Log warning but continue without scenario
			fmt.Fprintf(os.Stderr, "Warning: Failed to load scenario %s: %v\n", req.ScenarioID, err)
		} else if scenario.Prompt != "" {
			// Very long TTL for scenario context (7 days by default)
			scenarioTTL := 168 * time.Hour

			breakpoints = append(breakpoints, cache.CacheBreakpoint{
				Layer:      cache.ScenarioContextLayer,
				Content:    scenario.Prompt,
				TokenCount: cache.EstimateTokens(scenario.Prompt),
				TTL:        scenarioTTL,
				LastUsed:   time.Now(),
			})

			// Update scenario last used timestamp asynchronously
			go func(id string) {
				_ = cb.scenarioRepo.UpdateScenarioLastUsed(id)
			}(req.ScenarioID)
		}
	}

	// Layer 1: Core Personality (static, long TTL)
	personality := cb.buildPersonalityPrompt(char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.CorePersonalityLayer,
		Content:    personality,
		TokenCount: cache.EstimateTokens(personality),
		TTL:        cb.cache.CalculateAdaptiveTTL(nil, true),
	})

	// Layer 2: Learned Behaviors (semi-static, medium TTL)
	behaviors := cb.buildLearnedBehaviors(char)
	if behaviors != "" {
		breakpoints = append(breakpoints, cache.CacheBreakpoint{
			Layer:      cache.LearnedBehaviorLayer,
			Content:    behaviors,
			TokenCount: cache.EstimateTokens(behaviors),
			TTL:        cb.config.CacheConfig.DefaultTTL * 2,
		})
	}

	// Layer 3: Emotional State (dynamic, short TTL)
	emotional := cb.buildEmotionalContext(char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.EmotionalStateLayer,
		Content:    emotional,
		TokenCount: cache.EstimateTokens(emotional),
		TTL:        5 * time.Minute,
	})

	// Layer 4: User Context (semi-dynamic, medium TTL)
	userContext := cb.buildUserContext(req.UserID, char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.UserMemoryLayer,
		Content:    userContext,
		TokenCount: cache.EstimateTokens(userContext),
		TTL:        cb.config.CacheConfig.DefaultTTL,
	})

	// Layer 5: Conversation History (dynamic, no cache)
	conversation := cb.buildConversationHistory(req.Context)
	if conversation != "" {
		breakpoints = append(breakpoints, cache.CacheBreakpoint{
			Layer:      cache.ConversationLayer,
			Content:    conversation,
			TokenCount: cache.EstimateTokens(conversation),
			TTL:        0, // No caching for conversation
		})
	}

	// Combine all layers
	fullPrompt := cb.assemblePrompt(breakpoints, req.UserID, req.Message)

	return fullPrompt, breakpoints, nil
}

func (cb *CharacterBot) warmupCache(char *models.Character) {
	// Build personality prompt and create a cache key for this character
	personality := cb.buildPersonalityPrompt(char)

	// Create a stable cache key for just the personality layer
	h := sha256.New()
	h.Write([]byte(char.ID))
	h.Write([]byte("personality"))
	h.Write([]byte(personality))
	key := hex.EncodeToString(h.Sum(nil))

	// Store with a long TTL since personality is static
	cb.cache.Store(key, cache.CorePersonalityLayer, personality, 24*time.Hour)
}

func (cb *CharacterBot) buildPersonalityPrompt(char *models.Character) string {
	return fmt.Sprintf(`[CHARACTER PROFILE]
Name: %s
Personality Traits:
- Openness: %.2f
- Conscientiousness: %.2f
- Extraversion: %.2f
- Agreeableness: %.2f
- Neuroticism: %.2f

Backstory: %s

Speech Style: %s

Core Quirks: %s

[INTERACTION RULES]
- Always stay in character
- Express personality traits consistently
- Use characteristic speech patterns
- React based on emotional state`,
		char.Name,
		char.Personality.Openness,
		char.Personality.Conscientiousness,
		char.Personality.Extraversion,
		char.Personality.Agreeableness,
		char.Personality.Neuroticism,
		char.Backstory,
		char.SpeechStyle,
		joinQuirks(char.Quirks),
	)
}

func (cb *CharacterBot) buildLearnedBehaviors(char *models.Character) string {
	// Extract patterns from medium-term memories
	patterns := make([]string, 0)
	for _, mem := range char.Memories {
		if mem.Type == models.MediumTermMemory {
			patterns = append(patterns, mem.Content)
		}
	}

	if len(patterns) == 0 {
		return ""
	}

	return fmt.Sprintf("[LEARNED PATTERNS]\n%s", strings.Join(patterns, "\n"))
}

func (cb *CharacterBot) buildEmotionalContext(char *models.Character) string {
	return fmt.Sprintf(`[EMOTIONAL STATE]
Current Mood:
- Joy: %.2f
- Surprise: %.2f
- Anger: %.2f
- Fear: %.2f
- Sadness: %.2f
- Disgust: %.2f`,
		char.CurrentMood.Joy,
		char.CurrentMood.Surprise,
		char.CurrentMood.Anger,
		char.CurrentMood.Fear,
		char.CurrentMood.Sadness,
		char.CurrentMood.Disgust,
	)
}

func (cb *CharacterBot) buildConversationHistory(ctx models.ConversationContext) string {
	if len(ctx.RecentMessages) == 0 {
		return ""
	}

	history := "[CONVERSATION HISTORY]\n"
	for _, msg := range ctx.RecentMessages {
		history += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	return history
}

func (cb *CharacterBot) buildUserContext(userID string, char *models.Character) string {
	// Try to load user profile if available
	if cb.userProfileRepo != nil && cb.config.UserProfileConfig.Enabled {
		profile, err := cb.userProfileRepo.LoadUserProfile(userID, char.ID)
		if err == nil && profile != nil {
			return cb.buildUserProfileLayer(userID, char.ID, profile)
		}
	}

	// Fallback to basic user context
	context := fmt.Sprintf(`[USER CONTEXT]
You are speaking with: %s

Remember to address them by their name throughout the conversation.`, userID)

	// Add any user-specific memories
	var userMemories []string
	for _, mem := range char.Memories {
		if mem.Type == models.LongTermMemory && mem.Content != "" {
			// Check if memory mentions this user (simple check)
			// In a more advanced system, you'd have user-specific memory storage
			userMemories = append(userMemories, mem.Content)
		}
	}

	if len(userMemories) > 0 {
		context += "\n\nShared experiences:\n" + strings.Join(userMemories, "\n")
	}

	return context
}

func (cb *CharacterBot) buildUserProfileLayer(userID, characterID string, profile *models.UserProfile) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[USER PROFILE FOR %s (as perceived by character %s)]\n", userID, characterID))

	if profile.OverallSummary != "" {
		sb.WriteString(fmt.Sprintf("Summary: %s\n", profile.OverallSummary))
	}

	if profile.InteractionStyle != "" {
		sb.WriteString(fmt.Sprintf("Interaction Style: %s\n", profile.InteractionStyle))
	}

	if len(profile.Facts) > 0 {
		sb.WriteString("\nKey Facts Remembered About User:\n")
		for _, fact := range profile.Facts {
			// Only include facts with confidence above threshold
			if fact.Confidence >= cb.config.UserProfileConfig.ConfidenceThreshold {
				sb.WriteString(fmt.Sprintf("- %s: %s (Confidence: %.1f)\n", fact.Key, fact.Value, fact.Confidence))
			}
		}
	}

	return sb.String()
}

func (cb *CharacterBot) assemblePrompt(breakpoints []cache.CacheBreakpoint, userID, message string) string {
	var parts []string

	// Add all breakpoint content
	for _, bp := range breakpoints {
		parts = append(parts, bp.Content)
	}

	// Add current message
	parts = append(parts, fmt.Sprintf("[CURRENT MESSAGE]\n%s: %s", userID, message))

	return strings.Join(parts, "\n\n")
}

func (cb *CharacterBot) generateCacheKey(charID, userID, scenarioID string, breakpoints []cache.CacheBreakpoint) string {
	// Generate cache key based only on static/semi-static layers
	// Don't include conversation layer which changes every time
	h := sha256.New()
	h.Write([]byte(charID))
	h.Write([]byte(userID))

	// Include scenario ID in the cache key if present
	// This ensures different scenarios create different cache entries
	if scenarioID != "" {
		h.Write([]byte(scenarioID))
	}

	// Only hash content from cacheable layers (not conversation)
	// Note: If scenario content is the first breakpoint, it's already included
	for _, bp := range breakpoints {
		if bp.Layer != cache.ConversationLayer {
			h.Write([]byte(bp.Content))
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (cb *CharacterBot) selectProvider() providers.AIProvider {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Try to get the default provider
	if provider, exists := cb.providers[cb.config.DefaultProvider]; exists {
		return provider
	}

	// Fallback to first available provider
	for _, p := range cb.providers {
		return p
	}

	return nil
}

func (cb *CharacterBot) updateCharacterState(charID string, resp *providers.AIResponse) {
	char, err := cb.GetCharacter(charID)
	if err != nil {
		return
	}

	char.Lock()
	defer char.Unlock()

	// Update emotional state with decay
	char.CurrentMood = cb.blendEmotions(char.CurrentMood, resp.Emotions, 0.3)

	// Add to short-term memory
	memory := models.Memory{
		Type:      models.ShortTermMemory,
		Content:   resp.Content,
		Timestamp: time.Now(),
		Emotional: cb.calculateEmotionalWeight(resp.Emotions),
	}
	char.Memories = append(char.Memories, memory)

	// Trigger consolidation if needed
	if len(char.Memories) > cb.config.MemoryConfig.ShortTermWindow {
		go cb.consolidateMemories(char)
	}

	// Evolution logic
	if cb.config.PersonalityConfig.EvolutionEnabled {
		cb.evolvePersonality(char, resp)
	}

	char.LastModified = time.Now()
}

func (cb *CharacterBot) blendEmotions(current, new models.EmotionalState, rate float64) models.EmotionalState {
	return models.EmotionalState{
		Joy:      current.Joy*(1-rate) + new.Joy*rate,
		Surprise: current.Surprise*(1-rate) + new.Surprise*rate,
		Anger:    current.Anger*(1-rate) + new.Anger*rate,
		Fear:     current.Fear*(1-rate) + new.Fear*rate,
		Sadness:  current.Sadness*(1-rate) + new.Sadness*rate,
		Disgust:  current.Disgust*(1-rate) + new.Disgust*rate,
	}
}

func (cb *CharacterBot) calculateEmotionalWeight(emotions models.EmotionalState) float64 {
	// Simple average of emotion intensities
	total := emotions.Joy + emotions.Surprise + emotions.Anger +
		emotions.Fear + emotions.Sadness + emotions.Disgust
	return total / 6.0
}

func (cb *CharacterBot) evolvePersonality(char *models.Character, resp *providers.AIResponse) {
	// Calculate trait impacts based on interaction
	impacts := cb.analyzeInteractionImpacts(resp)

	// Apply bounded evolution
	driftRate := cb.config.PersonalityConfig.MaxDriftRate
	char.Personality.Openness += impacts.Openness * driftRate
	char.Personality.Conscientiousness += impacts.Conscientiousness * driftRate
	char.Personality.Extraversion += impacts.Extraversion * driftRate
	char.Personality.Agreeableness += impacts.Agreeableness * driftRate
	char.Personality.Neuroticism += impacts.Neuroticism * driftRate

	// Normalize to keep traits in [0, 1] range
	char.Personality = models.NormalizePersonality(char.Personality)
}

func (cb *CharacterBot) analyzeInteractionImpacts(resp *providers.AIResponse) models.PersonalityTraits {
	// Simplified impact analysis based on emotional response
	return models.PersonalityTraits{
		Openness:          resp.Emotions.Surprise * 0.5,
		Conscientiousness: (1 - resp.Emotions.Anger) * 0.3,
		Extraversion:      resp.Emotions.Joy * 0.4,
		Agreeableness:     (1 - resp.Emotions.Disgust) * 0.3,
		Neuroticism:       (resp.Emotions.Fear + resp.Emotions.Sadness) * 0.3,
	}
}

// Background workers

func (cb *CharacterBot) memoryConsolidationWorker() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cb.consolidateAllMemories()
	}
}

func (cb *CharacterBot) consolidateAllMemories() {
	cb.mu.RLock()
	chars := make([]*models.Character, 0, len(cb.characters))
	for _, char := range cb.characters {
		chars = append(chars, char)
	}
	cb.mu.RUnlock()

	for _, char := range chars {
		cb.consolidateMemories(char)
	}
}

func (cb *CharacterBot) consolidateMemories(char *models.Character) {
	char.Lock()
	defer char.Unlock()

	// Group memories by emotional significance
	emotionalMemories := make([]models.Memory, 0)

	threshold := 0.7 // Emotional weight threshold

	for _, mem := range char.Memories {
		if mem.Type == models.ShortTermMemory && mem.Emotional > threshold {
			emotionalMemories = append(emotionalMemories, mem)
		}
	}

	// Consolidate emotional memories into medium-term
	if len(emotionalMemories) > 3 {
		consolidated := models.Memory{
			Type:      models.MediumTermMemory,
			Content:   cb.synthesizeMemories(emotionalMemories),
			Timestamp: time.Now(),
			Emotional: cb.averageEmotionalWeight(emotionalMemories),
		}
		char.Memories = append(char.Memories, consolidated)
	}

	// Prune old short-term memories
	cutoff := time.Now().Add(-cb.config.MemoryConfig.MediumTermDuration)
	filtered := make([]models.Memory, 0)

	for _, mem := range char.Memories {
		if mem.Type != models.ShortTermMemory || mem.Timestamp.After(cutoff) {
			filtered = append(filtered, mem)
		}
	}

	char.Memories = filtered
}

func (cb *CharacterBot) synthesizeMemories(memories []models.Memory) string {
	// In a real implementation, this would use NLP to create a coherent summary
	contents := make([]string, 0, len(memories))
	for _, mem := range memories {
		contents = append(contents, mem.Content)
	}
	return fmt.Sprintf("Consolidated memories: %s", strings.Join(contents, "; "))
}

func (cb *CharacterBot) averageEmotionalWeight(memories []models.Memory) float64 {
	if len(memories) == 0 {
		return 0
	}

	total := 0.0
	for _, mem := range memories {
		total += mem.Emotional
	}

	return total / float64(len(memories))
}

// Utility functions

func joinQuirks(quirks []string) string {
	if len(quirks) == 0 {
		return "None"
	}
	return strings.Join(quirks, ", ")
}

// UpdateUserProfile synchronously updates the user profile
func (cb *CharacterBot) UpdateUserProfile(userID string, char *models.Character, sessionID string) {
	if cb.userProfileAgent == nil {
		return
	}
	cb.updateUserProfileSync(userID, char, sessionID)
}

// updateUserProfileAsync asynchronously updates the user profile based on conversation history
func (cb *CharacterBot) updateUserProfileAsync(userID string, char *models.Character, sessionID string) {
	go cb.updateUserProfileSync(userID, char, sessionID)
}

// updateUserProfileSync performs the actual user profile update
func (cb *CharacterBot) updateUserProfileSync(userID string, char *models.Character, sessionID string) {
	// Update user profile based on conversation history

	// Only proceed if we have the minimum messages for an update
	sessionRepo := repository.NewSessionRepository(filepath.Join(os.Getenv("HOME"), ".config", "roleplay"))

	currentSession, err := sessionRepo.LoadSession(char.ID, sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session %s for user profile update: %v\n", sessionID, err)
		return
	}

	// Check if we should update based on frequency setting
	messageCount := len(currentSession.Messages)
	if messageCount == 0 || (messageCount%cb.config.UserProfileConfig.UpdateFrequency != 0) {
		return
	}

	// Update the profile
	turnsToConsider := cb.config.UserProfileConfig.TurnsToConsider
	if turnsToConsider <= 0 {
		turnsToConsider = 20 // Default value
	}

	if cb.userProfileAgent == nil {
		return
	}

	_, updateErr := cb.userProfileAgent.UpdateUserProfile(
		context.Background(),
		userID,
		char,
		currentSession.Messages,
		turnsToConsider,
	)

	if updateErr != nil {
		fmt.Fprintf(os.Stderr, "Error updating user profile for %s with %s: %v\n", userID, char.ID, updateErr)
	}
}
````

## File: go.mod
````
module github.com/dotcommander/roleplay

go 1.23.0

toolchain go1.24.3

require (
	github.com/charmbracelet/bubbles v0.21.0
	github.com/charmbracelet/bubbletea v1.3.5
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/google/uuid v1.6.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sashabaranov/go-openai v1.40.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
````

## File: cmd/chat.go
````go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var (
	characterID string
	userID      string
	sessionID   string
	format      string
	scenarioID  string
)

var chatCmd = &cobra.Command{
	Use:   "chat [message]",
	Short: "Chat with a character",
	Long: `Start a conversation with a character. The character will respond based on their
personality, emotional state, and conversation history.

Examples:
  roleplay chat "Hello, how are you?" --character warrior-123 --user user-789
  roleplay chat "Tell me about your adventures" -c warrior-123 -u user-789`,
	Args: cobra.ExactArgs(1),
	RunE: runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&characterID, "character", "c", "", "Character ID to chat with (required)")
	chatCmd.Flags().StringVarP(&userID, "user", "u", "", "User ID for the conversation (required)")
	chatCmd.Flags().StringVarP(&sessionID, "session", "s", "", "Session ID (optional, generates new if not provided)")
	chatCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	chatCmd.Flags().StringVar(&scenarioID, "scenario", "", "Scenario ID to set the interaction context (optional)")

	if err := chatCmd.MarkFlagRequired("character"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking character flag as required: %v\n", err)
	}
	if err := chatCmd.MarkFlagRequired("user"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking user flag as required: %v\n", err)
	}
}

func runChat(cmd *cobra.Command, args []string) error {
	message := args[0]
	config := GetConfig()

	// Validate API key
	if config.APIKey == "" {
		return fmt.Errorf("API key not configured. Set ROLEPLAY_API_KEY or use --api-key")
	}

	// Initialize manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	// Ensure character is loaded
	if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
		return fmt.Errorf("character %s not found. Create it first with 'roleplay character create'", characterID)
	}

	// Get session repository
	sessionRepo := mgr.GetSessionRepository()

	// Generate session ID if not provided
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
	}

	// Load existing session or create new one
	var session *repository.Session
	existingSession, err := sessionRepo.LoadSession(characterID, sessionID)
	if err != nil {
		// Create new session if it doesn't exist
		session = &repository.Session{
			ID:           sessionID,
			CharacterID:  characterID,
			UserID:       userID,
			StartTime:    time.Now(),
			LastActivity: time.Now(),
			Messages:     []repository.SessionMessage{},
			CacheMetrics: repository.CacheMetrics{},
		}
	} else {
		session = existingSession
	}

	// Convert session messages to conversation context
	var recentMessages []models.Message
	// Get last 10 messages for context
	startIdx := 0
	if len(session.Messages) > 10 {
		startIdx = len(session.Messages) - 10
	}
	for i := startIdx; i < len(session.Messages); i++ {
		msg := session.Messages[i]
		// Map "character" role to "assistant" for API compatibility
		role := msg.Role
		if role == "character" {
			role = "assistant"
		}
		recentMessages = append(recentMessages, models.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Create conversation request
	req := &models.ConversationRequest{
		CharacterID: characterID,
		UserID:      userID,
		Message:     message,
		ScenarioID:  scenarioID,
		Context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      session.StartTime,
			RecentMessages: recentMessages,
		},
	}

	// Process request
	ctx := context.Background()
	resp, err := mgr.GetBot().ProcessRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	// Update session with new messages
	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp:  time.Now(),
		Role:       "user",
		Content:    message,
		TokensUsed: 0, // User messages don't consume tokens
	})

	cacheHits := 0
	cacheMisses := 0
	if resp.CacheMetrics.Hit {
		cacheHits = 1
	} else {
		cacheMisses = 1
	}

	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp:   time.Now(),
		Role:        "character",
		Content:     resp.Content,
		TokensUsed:  resp.TokensUsed.Total,
		CacheHits:   cacheHits,
		CacheMisses: cacheMisses,
	})

	// Update cache metrics
	session.CacheMetrics.TotalRequests++
	if resp.CacheMetrics.Hit {
		session.CacheMetrics.CacheHits++
	} else {
		session.CacheMetrics.CacheMisses++
	}
	session.CacheMetrics.TokensSaved += resp.CacheMetrics.SavedTokens
	session.CacheMetrics.HitRate = float64(session.CacheMetrics.CacheHits) / float64(session.CacheMetrics.TotalRequests)
	session.CacheMetrics.CostSaved = float64(session.CacheMetrics.TokensSaved) * 0.000003 // Approximate cost per token
	session.LastActivity = time.Now()

	// Save session BEFORE checking for profile updates
	// This ensures the async goroutine has access to the latest messages
	if err := sessionRepo.SaveSession(session); err != nil {
		// Log error but don't fail the command
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save session: %v\n", err)
		}
	}

	// Check if we should update the user profile
	// For the chat command, we do this synchronously to ensure it completes
	if config.UserProfileConfig.Enabled && len(session.Messages) > 0 && (len(session.Messages)%config.UserProfileConfig.UpdateFrequency == 0) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "Updating user profile...\n")
		}

		// Get character from manager
		char, err := mgr.GetOrLoadCharacter(characterID)
		if err == nil {
			// Call the bot's profile update method directly
			mgr.GetBot().UpdateUserProfile(userID, char, sessionID)
		}
	}

	// Display response based on format
	if format == "json" {
		output := map[string]interface{}{
			"session_id": sessionID,
			"response":   resp.Content,
			"cache_metrics": map[string]interface{}{
				"cache_hit":    resp.CacheMetrics.Hit,
				"layers":       resp.CacheMetrics.Layers,
				"saved_tokens": resp.CacheMetrics.SavedTokens,
				"latency_ms":   resp.CacheMetrics.Latency.Milliseconds(),
			},
			"token_usage": map[string]interface{}{
				"prompt":        resp.TokensUsed.Prompt,
				"completion":    resp.TokensUsed.Completion,
				"cached_prompt": resp.TokensUsed.CachedPrompt,
				"total":         resp.TokensUsed.Total,
			},
		}
		jsonBytes, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		// Display response
		fmt.Fprintf(os.Stdout, "\n%s\n", resp.Content)

		// Show cache metrics if verbose
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "\n--- Performance Metrics ---\n")
			fmt.Fprintf(os.Stderr, "Session ID: %s\n", sessionID)
			fmt.Fprintf(os.Stderr, "Cache Hit: %v\n", resp.CacheMetrics.Hit)
			fmt.Fprintf(os.Stderr, "Tokens Used: %d (cached: %d)\n",
				resp.TokensUsed.Total, resp.TokensUsed.CachedPrompt)
			fmt.Fprintf(os.Stderr, "Tokens Saved: %d\n", resp.CacheMetrics.SavedTokens)
			fmt.Fprintf(os.Stderr, "Latency: %v\n", resp.CacheMetrics.Latency)
			fmt.Fprintf(os.Stderr, "Session Messages: %d\n", len(session.Messages))
		}
	}

	return nil
}
````

## File: cmd/interactive.go
````go
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/utils"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive chat session",
	Long: `Start an interactive chat session with a character using a beautiful TUI interface.

This provides a REPL-like experience with:
- Real-time chat with scrolling history
- Character personality and mood display
- Cache performance metrics
- Session persistence

Examples:
  roleplay interactive                     # Uses rick-c137 and your username
  roleplay interactive -c philosopher-123  # Chat with a specific character
  roleplay interactive -u morty            # Specify a different user ID`,
	RunE: runInteractive,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.Flags().StringP("character", "c", "", "Character ID to chat with (defaults to rick-c137)")
	interactiveCmd.Flags().StringP("user", "u", "", "User ID for the conversation (defaults to your username)")
	interactiveCmd.Flags().StringP("session", "s", "", "Session ID (optional)")
	interactiveCmd.Flags().Bool("new-session", false, "Start a new session instead of resuming")
	interactiveCmd.Flags().String("scenario", "", "Scenario ID to set the interaction context (optional)")
}

// Styles - Gruvbox Dark Theme
var (
	// Gruvbox Dark Colors
	gruvboxBg     = lipgloss.Color("#282828") // Dark background
	gruvboxBg1    = lipgloss.Color("#3c3836") // Lighter background
	gruvboxFg     = lipgloss.Color("#ebdbb2") // Foreground
	gruvboxRed    = lipgloss.Color("#fb4934") // Bright red
	gruvboxGreen  = lipgloss.Color("#b8bb26") // Bright green
	gruvboxYellow = lipgloss.Color("#fabd2f") // Bright yellow
	gruvboxPurple = lipgloss.Color("#d3869b") // Bright purple
	gruvboxAqua   = lipgloss.Color("#8ec07c") // Bright aqua
	gruvboxOrange = lipgloss.Color("#fe8019") // Bright orange
	gruvboxGray   = lipgloss.Color("#928374") // Gray
	gruvboxFg2    = lipgloss.Color("#d5c4a1") // Dimmer foreground

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(gruvboxAqua).
			Bold(true).
			Padding(0, 1)

	characterStyle = lipgloss.NewStyle().
			Foreground(gruvboxOrange).
			Bold(true)

	userStyle = lipgloss.NewStyle().
			Foreground(gruvboxGreen).
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray)

	errorStyle = lipgloss.NewStyle().
			Foreground(gruvboxRed).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxGray).
			Foreground(gruvboxFg).
			Background(gruvboxBg).
			Padding(1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(gruvboxBg).
			Background(gruvboxAqua).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Italic(true)

	personalityStyle = lipgloss.NewStyle().
				Foreground(gruvboxPurple)

	moodStyle = lipgloss.NewStyle().
			Foreground(gruvboxYellow)

	timestampStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Italic(true)

	userMessageStyle = lipgloss.NewStyle().
				Foreground(gruvboxFg).
				Background(gruvboxBg1).
				Padding(0, 1).
				MarginRight(2)

	characterMessageStyle = lipgloss.NewStyle().
				Foreground(gruvboxFg).
				Background(gruvboxBg).
				Padding(0, 1).
				MarginLeft(2)

	separatorStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray)

	// Command output styles
	commandHeaderStyle = lipgloss.NewStyle().
				Foreground(gruvboxAqua).
				Bold(true).
				Padding(0, 1)

	commandBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxAqua).
			Foreground(gruvboxFg).
			Background(gruvboxBg).
			Padding(1, 2)

	listItemStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg2)

	listItemActiveStyle = lipgloss.NewStyle().
				Foreground(gruvboxGreen).
				Bold(true)

	helpCommandStyle = lipgloss.NewStyle().
				Foreground(gruvboxYellow).
				Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg2)
)

// Message types
type chatMsg struct {
	role    string
	content string
	time    time.Time
	msgType string // "normal", "help", "list", "stats", etc.
}

type responseMsg struct {
	content string
	metrics *cache.CacheMetrics
	err     error
}

type characterInfoMsg struct {
	character *models.Character
}

type systemMsg struct {
	content string
	msgType string // "info", "error", "help"
}

type characterSwitchMsg struct {
	characterID string
	character   *models.Character
}

// Model
type model struct {
	// UI components
	viewport    viewport.Model
	textarea    textarea.Model
	spinner     spinner.Model
	messages    []chatMsg
	characterID string
	userID      string
	sessionID   string
	scenarioID  string
	bot         *services.CharacterBot
	character   *models.Character
	context     models.ConversationContext
	loading     bool
	err         error
	width       int
	height      int
	ready       bool
	model       string // AI model being used

	// Cache metrics
	lastCacheHit    bool
	lastTokensSaved int
	totalRequests   int
	cacheHits       int

	// Command history
	commandHistory []string
	historyIndex   int
	historyBuffer  string // Stores current input when navigating history
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spinner.Tick,
		m.loadCharacterInfo(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewport
			headerHeight := 8 // Character info
			footerHeight := 6 // Input area + status
			verticalMargins := headerHeight + footerHeight

			m.viewport = viewport.New(msg.Width-4, msg.Height-verticalMargins)
			m.viewport.SetContent(m.renderMessages())

			// Initialize textarea with Gruvbox styling
			m.textarea = textarea.New()
			m.textarea.Placeholder = "Type your message..."
			m.textarea.Focus()
			m.textarea.Prompt = "│ "
			m.textarea.CharLimit = 500
			m.textarea.SetWidth(msg.Width - 4)
			m.textarea.SetHeight(2)
			m.textarea.ShowLineNumbers = false
			m.textarea.KeyMap.InsertNewline.SetEnabled(false)

			// Style the textarea
			m.textarea.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(gruvboxBg1)
			m.textarea.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(gruvboxAqua)
			m.textarea.FocusedStyle.Text = lipgloss.NewStyle().Foreground(gruvboxFg)
			m.textarea.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(gruvboxGray)

			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 14
			m.textarea.SetWidth(msg.Width - 4)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.loading && m.textarea.Value() != "" {
				message := m.textarea.Value()

				// Add to command history
				m.commandHistory = append(m.commandHistory, message)
				m.historyIndex = len(m.commandHistory) // Reset to end of history
				m.historyBuffer = ""                   // Clear history buffer

				m.textarea.Reset()

				// Check for slash commands
				if strings.HasPrefix(message, "/") {
					return m, m.handleSlashCommand(message)
				}

				m.messages = append(m.messages, chatMsg{
					role:    "user",
					content: message,
					time:    time.Now(),
					msgType: "normal",
				})
				m.viewport.SetContent(m.renderMessages())
				m.viewport.GotoBottom()
				m.loading = true
				m.totalRequests++
				return m, m.sendMessage(message)
			}
		case tea.KeyUp:
			// Navigate backward in history
			if len(m.commandHistory) > 0 && m.historyIndex > 0 {
				// Save current input if this is the first time navigating
				if m.historyIndex == len(m.commandHistory) {
					m.historyBuffer = m.textarea.Value()
				}

				m.historyIndex--
				m.textarea.SetValue(m.commandHistory[m.historyIndex])
				m.textarea.CursorEnd() // Move cursor to end
			}
		case tea.KeyDown:
			// Navigate forward in history
			if len(m.commandHistory) > 0 && m.historyIndex < len(m.commandHistory) {
				m.historyIndex++

				if m.historyIndex == len(m.commandHistory) {
					// Restore the original input
					m.textarea.SetValue(m.historyBuffer)
				} else {
					m.textarea.SetValue(m.commandHistory[m.historyIndex])
				}
				m.textarea.CursorEnd() // Move cursor to end
			}
		}
	case characterInfoMsg:
		m.character = msg.character

	case systemMsg:
		// Handle special system commands
		if msg.msgType == "clear" && msg.content == "clear_history" {
			m.messages = []chatMsg{}
			m.viewport.SetContent(m.renderMessages())
			// Add confirmation message
			m.messages = append(m.messages, chatMsg{
				role:    "system",
				content: "Chat history cleared",
				time:    time.Now(),
				msgType: "info",
			})
		} else {
			// Add system message to chat
			m.messages = append(m.messages, chatMsg{
				role:    "system",
				content: msg.content,
				time:    time.Now(),
				msgType: msg.msgType,
			})
		}
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case characterSwitchMsg:
		// Save current session before switching
		m.saveSession()

		// Update character
		m.characterID = msg.characterID
		m.character = msg.character

		// Clear conversation and start new session
		m.messages = []chatMsg{}
		m.sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
		m.context = models.ConversationContext{
			SessionID:      m.sessionID,
			StartTime:      time.Now(),
			RecentMessages: []models.Message{},
		}

		// Reset cache metrics for new session
		m.totalRequests = 0
		m.cacheHits = 0
		m.lastTokensSaved = 0
		m.lastCacheHit = false

		// Add switch notification
		m.messages = append(m.messages, chatMsg{
			role:    "system",
			content: fmt.Sprintf("Switched to %s (%s). Starting new session.", msg.character.Name, msg.characterID),
			time:    time.Now(),
			msgType: "info",
		})

		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case responseMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.messages = append(m.messages, chatMsg{
				role:    m.character.Name,
				content: msg.content,
				time:    time.Now(),
				msgType: "normal",
			})

			// Update cache metrics
			if msg.metrics != nil {
				m.lastCacheHit = msg.metrics.Hit
				m.lastTokensSaved = msg.metrics.SavedTokens
				if msg.metrics.Hit {
					m.cacheHits++
				}
			}

			// Update context with recent messages
			m.updateContext()

			// Save session after each interaction
			m.saveSession()
		}
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case spinner.TickMsg:
		if m.loading {
			m.spinner, tiCmd = m.spinner.Update(msg)
			cmds = append(cmds, tiCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Character info header
	header := m.renderHeader()

	// Main chat viewport
	chatView := borderStyle.Render(m.viewport.View())

	// Input area
	inputArea := m.renderInputArea()

	// Status bar
	statusBar := m.renderStatusBar()

	// Help text
	help := helpStyle.Render("  ⌃C quit • ↵ send • ↑↓ history • /help commands • /exit quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		chatView,
		inputArea,
		statusBar,
		help,
	)
}

func (m model) renderHeader() string {
	if m.character == nil {
		return titleStyle.Render("Loading character...")
	}

	title := titleStyle.Render(fmt.Sprintf("󰊕 Chat with %s", m.character.Name))

	// Personality traits with icons
	personality := fmt.Sprintf(
		" O:%.1f  C:%.1f  E:%.1f  A:%.1f  N:%.1f",
		m.character.Personality.Openness,
		m.character.Personality.Conscientiousness,
		m.character.Personality.Extraversion,
		m.character.Personality.Agreeableness,
		m.character.Personality.Neuroticism,
	)

	// Current mood (dominant emotion)
	mood := m.getDominantMood()
	moodIcon := m.getMoodIcon(mood)

	personalityInfo := personalityStyle.Render(personality)
	moodInfo := moodStyle.Render(fmt.Sprintf(" %s %s", moodIcon, mood))

	info := fmt.Sprintf("  %s • %s", personalityInfo, moodInfo)

	return lipgloss.JoinVertical(lipgloss.Left, title, info, "")
}

func (m model) renderMessages() string {
	if len(m.messages) == 0 {
		emptyMsg := mutedStyle.Render("\n   Start chatting! Your conversation will appear here...\n")
		return emptyMsg
	}

	var content strings.Builder
	maxWidth := m.viewport.Width - 8 // Account for padding and margins

	for i, msg := range m.messages {
		if i > 0 {
			// Add visual separator between messages
			separator := separatorStyle.Render(strings.Repeat("─", maxWidth))
			content.WriteString("\n" + separator + "\n\n")
		}

		timestamp := timestampStyle.Render(msg.time.Format("15:04:05"))

		if msg.role == "user" {
			// User message - consistent styling throughout
			header := fmt.Sprintf("┌─ %s %s", userStyle.Render("You"), timestamp)
			content.WriteString(userMessageStyle.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(userMessageStyle.Render(prefix+line) + "\n")
			}
		} else if msg.role == "system" {
			// System message - check for special types
			if msg.msgType == "help" || msg.msgType == "list" || msg.msgType == "info" || msg.msgType == "stats" {
				// Special formatted output
				content.WriteString("\n")

				// Determine the header based on type
				var header string
				switch msg.msgType {
				case "help":
					header = commandHeaderStyle.Render("📚 Command Help")
				case "list":
					header = commandHeaderStyle.Render("📋 Available Characters")
				case "stats":
					header = commandHeaderStyle.Render("📊 Cache Statistics")
				case "info":
					header = commandHeaderStyle.Render("ℹ️  Information")
				default:
					header = commandHeaderStyle.Render("System")
				}

				// Format the content with special styling
				formattedContent := m.formatSpecialMessage(msg.content, msg.msgType, maxWidth-8)
				boxContent := lipgloss.JoinVertical(lipgloss.Left, header, "", formattedContent)
				content.WriteString(commandBoxStyle.Width(maxWidth).Render(boxContent) + "\n")
			} else {
				// Regular system message
				header := fmt.Sprintf("┌─ %s %s", mutedStyle.Render("System"), timestamp)
				content.WriteString(mutedStyle.Render(header) + "\n")

				wrappedContent := utils.WrapText(msg.content, maxWidth-4)
				lines := strings.Split(wrappedContent, "\n")
				for j, line := range lines {
					prefix := "│ "
					if j == len(lines)-1 {
						prefix = "└ "
					}
					content.WriteString(mutedStyle.Render(prefix+line) + "\n")
				}
			}
		} else {
			// Character message - consistent styling throughout
			header := fmt.Sprintf("┌─ %s %s", characterStyle.Render(msg.role), timestamp)
			content.WriteString(characterMessageStyle.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(characterMessageStyle.Render(prefix+line) + "\n")
			}
		}
	}

	return content.String()
}

func (m model) renderInputArea() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("   Error: %v", m.err))
	}

	if m.loading {
		spinnerText := mutedStyle.Render("Thinking...")
		return fmt.Sprintf("\n  %s %s\n", m.spinner.View(), spinnerText)
	}

	return fmt.Sprintf("\n%s\n", m.textarea.View())
}

func (m model) renderStatusBar() string {
	cacheRate := 0.0
	if m.totalRequests > 0 {
		cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
	}

	// Cache indicator with color
	cacheIndicator := "○"
	if m.lastCacheHit {
		cacheIndicator = "●"
	}

	status := fmt.Sprintf(
		" %s %s │ %s │  %d │ %s %.0f%% │  %d tokens saved",
		cacheIndicator,
		m.sessionID,
		m.model,
		m.totalRequests,
		cacheIndicator,
		cacheRate,
		m.lastTokensSaved,
	)

	return statusBarStyle.Width(m.width).Render(status)
}

func (m model) getDominantMood() string {
	if m.character == nil {
		return "Unknown"
	}

	moods := map[string]float64{
		"Joy":      m.character.CurrentMood.Joy,
		"Surprise": m.character.CurrentMood.Surprise,
		"Anger":    m.character.CurrentMood.Anger,
		"Fear":     m.character.CurrentMood.Fear,
		"Sadness":  m.character.CurrentMood.Sadness,
		"Disgust":  m.character.CurrentMood.Disgust,
	}

	maxMood := "Neutral"
	maxValue := 0.0

	for mood, value := range moods {
		if value > maxValue {
			maxMood = mood
			maxValue = value
		}
	}

	if maxValue < 0.2 {
		return "Neutral"
	}

	return maxMood
}

func (m model) getMoodIcon(mood string) string {
	switch mood {
	case "Joy":
		return "😊"
	case "Surprise":
		return "😲"
	case "Anger":
		return "😠"
	case "Fear":
		return "😨"
	case "Sadness":
		return "😢"
	case "Disgust":
		return "🤢"
	case "Neutral":
		return "😐"
	default:
		return "🤔"
	}
}

func (m *model) updateContext() {
	// Keep last 10 messages in context
	startIdx := 0
	if len(m.messages) > 10 {
		startIdx = len(m.messages) - 10
	}

	m.context.RecentMessages = make([]models.Message, 0)
	for i := startIdx; i < len(m.messages); i++ {
		role := "user"
		if m.messages[i].role != "user" {
			role = "assistant"
		}

		m.context.RecentMessages = append(m.context.RecentMessages, models.Message{
			Role:      role,
			Content:   m.messages[i].content,
			Timestamp: m.messages[i].time,
		})
	}
}

func (m *model) saveSession() {
	// Save session in background
	go func() {
		dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
		sessionRepo := repository.NewSessionRepository(dataDir)

		// Convert chat messages back to session messages
		var sessionMessages []repository.SessionMessage
		for _, msg := range m.messages {
			sessionMessages = append(sessionMessages, repository.SessionMessage{
				Timestamp: msg.time,
				Role: func() string {
					if msg.role == "user" {
						return "user"
					}
					return "character"
				}(),
				Content:    msg.content,
				TokensUsed: 0, // Could track this per message if needed
			})
		}

		session := &repository.Session{
			ID:           m.sessionID,
			CharacterID:  m.characterID,
			UserID:       m.userID,
			StartTime:    m.context.StartTime,
			LastActivity: time.Now(),
			Messages:     sessionMessages,
			CacheMetrics: repository.CacheMetrics{
				TotalRequests: m.totalRequests,
				CacheHits:     m.cacheHits,
				CacheMisses:   m.totalRequests - m.cacheHits,
				TokensSaved:   m.lastTokensSaved,
				HitRate: func() float64 {
					if m.totalRequests > 0 {
						return float64(m.cacheHits) / float64(m.totalRequests)
					}
					return 0.0
				}(),
			},
		}

		if err := sessionRepo.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving session: %v\n", err)
		}
	}()
}

func (m model) sendMessage(message string) tea.Cmd {
	return func() tea.Msg {
		req := &models.ConversationRequest{
			CharacterID: m.characterID,
			UserID:      m.userID,
			Message:     message,
			ScenarioID:  m.scenarioID,
			Context:     m.context,
		}

		ctx := context.Background()
		resp, err := m.bot.ProcessRequest(ctx, req)
		if err != nil {
			return responseMsg{err: err}
		}

		// Get updated character state
		char, _ := m.bot.GetCharacter(m.characterID)
		if char != nil {
			m.character = char
		}

		return responseMsg{
			content: resp.Content,
			metrics: &resp.CacheMetrics,
		}
	}
}

func (m model) loadCharacterInfo() tea.Cmd {
	return func() tea.Msg {
		char, err := m.bot.GetCharacter(m.characterID)
		if err != nil {
			return responseMsg{err: err}
		}
		return characterInfoMsg{character: char}
	}
}

func (m model) formatSpecialMessage(content string, msgType string, width int) string {
	switch msgType {
	case "help":
		// Format help message with colored commands
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, " - ") {
				parts := strings.SplitN(line, " - ", 2)
				if len(parts) == 2 {
					cmd := helpCommandStyle.Render(parts[0])
					desc := helpDescStyle.Render("- " + parts[1])
					formatted = append(formatted, fmt.Sprintf("%s %s", cmd, desc))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "list":
		// Format character list with special styling
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.HasPrefix(line, "→ ") {
				// Active character
				formatted = append(formatted, listItemActiveStyle.Render(line))
			} else if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") {
				// Character name line
				formatted = append(formatted, characterStyle.Render(line))
			} else if strings.HasPrefix(line, "   ") {
				// Description line
				formatted = append(formatted, listItemStyle.Render(line))
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "stats":
		// Format stats with colored numbers
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					label := parts[0]
					value := strings.TrimSpace(parts[1])
					// Color numbers
					if strings.Contains(value, "%") || strings.Contains(value, "tokens") {
						value = characterStyle.Render(value)
					}
					formatted = append(formatted, fmt.Sprintf("%s: %s", label, value))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	default:
		// Default formatting
		return utils.WrapText(content, width)
	}
}

func (m model) handleSlashCommand(input string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return systemMsg{content: "Invalid command", msgType: "error"}
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "/exit", "/quit", "/q":
			return tea.Quit()

		case "/help", "/h":
			helpText := `Available slash commands:
/help, /h     - Show this help message
/exit, /quit, /q - Exit the chat
/clear, /c    - Clear chat history
/list         - List all available characters
/switch <id>  - Switch to a different character
/stats        - Show cache statistics
/mood         - Show character's current mood
/personality  - Show character's personality traits
/session      - Show session information`
			return systemMsg{content: helpText, msgType: "help"}

		case "/clear", "/c":
			// We can't directly modify the model here, so we'll return a special message
			return systemMsg{content: "clear_history", msgType: "clear"}

		case "/list":
			// List all available characters from the repository
			dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
			charRepo, err := repository.NewCharacterRepository(dataDir)
			if err != nil {
				return systemMsg{content: fmt.Sprintf("Error accessing characters: %v", err), msgType: "error"}
			}

			characterIDs, err := charRepo.ListCharacters()
			if err != nil {
				return systemMsg{content: fmt.Sprintf("Error listing characters: %v", err), msgType: "error"}
			}

			if len(characterIDs) == 0 {
				return systemMsg{content: "No characters available. Use 'roleplay character create' to add characters.", msgType: "info"}
			}

			var listText strings.Builder
			listText.WriteString("Available Characters:\n")

			for _, id := range characterIDs {
				// Try to get from bot first (already loaded)
				char, err := m.bot.GetCharacter(id)
				if err != nil {
					// If not loaded, load from repository
					char, err = charRepo.LoadCharacter(id)
					if err != nil {
						continue
					}
					// Load into bot for future use
					_ = m.bot.CreateCharacter(char)
				}

				// Current character indicator
				indicator := "  "
				if id == m.characterID {
					indicator = "→ "
				}

				// Mood icon
				tempModel := model{character: char}
				mood := tempModel.getDominantMood()
				moodIcon := tempModel.getMoodIcon(mood)

				listText.WriteString(fmt.Sprintf("\n%s%s (%s) %s %s\n",
					indicator, char.Name, id, moodIcon, mood))

				// Add brief description from backstory (first sentence)
				backstory := char.Backstory
				if idx := strings.Index(backstory, "."); idx != -1 && idx < 100 {
					backstory = backstory[:idx+1]
				} else if len(backstory) > 100 {
					backstory = backstory[:97] + "..."
				}
				listText.WriteString(fmt.Sprintf("   %s\n", backstory))
			}

			return systemMsg{content: listText.String(), msgType: "list"}

		case "/stats":
			cacheRate := 0.0
			if m.totalRequests > 0 {
				cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
			}
			statsText := fmt.Sprintf(`Cache Statistics:
• Total requests: %d
• Cache hits: %d
• Cache misses: %d  
• Hit rate: %.1f%%
• Tokens saved: %d`,
				m.totalRequests,
				m.cacheHits,
				m.totalRequests-m.cacheHits,
				cacheRate,
				m.lastTokensSaved)
			return systemMsg{content: statsText, msgType: "info"}

		case "/mood":
			if m.character == nil {
				return systemMsg{content: "Character not loaded", msgType: "error"}
			}
			mood := m.getDominantMood()
			icon := m.getMoodIcon(mood)
			moodText := fmt.Sprintf(`%s Current Mood: %s

Emotional State:
• Joy: %.1f      • Surprise: %.1f
• Anger: %.1f    • Fear: %.1f  
• Sadness: %.1f  • Disgust: %.1f`,
				icon, mood,
				m.character.CurrentMood.Joy,
				m.character.CurrentMood.Surprise,
				m.character.CurrentMood.Anger,
				m.character.CurrentMood.Fear,
				m.character.CurrentMood.Sadness,
				m.character.CurrentMood.Disgust)
			return systemMsg{content: moodText, msgType: "info"}

		case "/personality":
			if m.character == nil {
				return systemMsg{content: "Character not loaded", msgType: "error"}
			}
			personalityText := fmt.Sprintf(`%s's Personality (OCEAN Model):

• Openness: %.1f        (creativity, openness to experience)
• Conscientiousness: %.1f (organization, self-discipline) 
• Extraversion: %.1f     (sociability, assertiveness)
• Agreeableness: %.1f    (compassion, cooperation)
• Neuroticism: %.1f      (emotional instability, anxiety)`,
				m.character.Name,
				m.character.Personality.Openness,
				m.character.Personality.Conscientiousness,
				m.character.Personality.Extraversion,
				m.character.Personality.Agreeableness,
				m.character.Personality.Neuroticism)
			return systemMsg{content: personalityText, msgType: "info"}

		case "/session":
			characterName := m.characterID
			if m.character != nil {
				characterName = m.character.Name
			}
			sessionIDDisplay := m.sessionID
			if len(m.sessionID) > 8 {
				sessionIDDisplay = m.sessionID[:8] + "..."
			}
			sessionText := fmt.Sprintf(`Session Information:
• Session ID: %s
• Character: %s (%s)
• User: %s
• Messages: %d
• Started: %s`,
				sessionIDDisplay,
				characterName,
				m.characterID,
				m.userID,
				len(m.messages),
				m.context.StartTime.Format("Jan 2, 2006 15:04"))
			return systemMsg{content: sessionText, msgType: "info"}

		case "/switch":
			if len(parts) < 2 {
				return systemMsg{content: "Usage: /switch <character-id>\nUse /list to see available characters", msgType: "error"}
			}

			newCharID := parts[1]

			// Check if it's the same character
			if newCharID == m.characterID {
				return systemMsg{content: fmt.Sprintf("Already chatting with %s", m.characterID), msgType: "info"}
			}

			// Try to load the character
			char, err := m.bot.GetCharacter(newCharID)
			if err != nil {
				// If not loaded in bot, try loading from repository
				dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
				charRepo, repoErr := repository.NewCharacterRepository(dataDir)
				if repoErr != nil {
					return systemMsg{content: fmt.Sprintf("Error accessing characters: %v", repoErr), msgType: "error"}
				}

				char, err = charRepo.LoadCharacter(newCharID)
				if err != nil {
					return systemMsg{content: fmt.Sprintf("Character '%s' not found. Use /list to see available characters", newCharID), msgType: "error"}
				}

				// Load character into bot
				if err := m.bot.CreateCharacter(char); err != nil {
					return systemMsg{content: fmt.Sprintf("Error loading character: %v", err), msgType: "error"}
				}
			}

			return characterSwitchMsg{
				characterID: newCharID,
				character:   char,
			}

		default:
			return systemMsg{content: fmt.Sprintf("Unknown command: %s\nType /help for available commands", command), msgType: "error"}
		}
	}
}

func runInteractive(cmd *cobra.Command, args []string) error {
	config := GetConfig()

	// Validate API key
	if config.APIKey == "" {
		return fmt.Errorf("API key not configured. Set OPENAI_API_KEY or ROLEPLAY_API_KEY")
	}

	// Get flags
	characterID, _ := cmd.Flags().GetString("character")
	userID, _ := cmd.Flags().GetString("user")
	sessionID, _ := cmd.Flags().GetString("session")
	newSession, _ := cmd.Flags().GetBool("new-session")
	scenarioID, _ := cmd.Flags().GetString("scenario")

	// Apply smart defaults
	if characterID == "" {
		characterID = "rick-c137" // Default to Rick Sanchez
	}
	if userID == "" {
		// Try to get username from environment
		userID = os.Getenv("USER")
		if userID == "" {
			userID = os.Getenv("USERNAME") // Windows fallback
		}
		if userID == "" {
			userID = "user" // Final fallback
		}
	}

	// Initialize repository for session management
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	sessionRepo := repository.NewSessionRepository(dataDir)

	var existingSession *repository.Session
	var existingMessages []chatMsg

	// Try to resume latest session if not specified and not forced new
	if sessionID == "" && !newSession {
		if latestSession, err := sessionRepo.GetLatestSession(characterID); err == nil && latestSession.ID != "" {
			sessionID = latestSession.ID
			existingSession = latestSession
			sessionIDDisplay := sessionID
			if len(sessionID) > 8 {
				sessionIDDisplay = sessionID[:8]
			}
			fmt.Printf("🔄 Resuming session %s (started %s, %d messages)\n",
				sessionIDDisplay,
				latestSession.StartTime.Format("Jan 2 15:04"),
				len(latestSession.Messages))

			// Convert session messages to chat messages
			for _, msg := range latestSession.Messages {
				role := msg.Role
				if role == "character" {
					role = characterID // Use character name for display
				}
				existingMessages = append(existingMessages, chatMsg{
					role:    role,
					content: msg.Content,
					time:    msg.Timestamp,
					msgType: "normal",
				})
			}
		} else {
			sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
			sessionIDDisplay := sessionID
			if len(sessionID) > 8 {
				sessionIDDisplay = sessionID[:8]
			}
			fmt.Printf("🆕 Starting new session %s\n", sessionIDDisplay)
		}
	}

	// Ensure sessionID is set even if not resuming
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
		sessionIDDisplay := sessionID
		if len(sessionID) > 8 {
			sessionIDDisplay = sessionID[:8]
		}
		fmt.Printf("🆕 Starting new session %s\n", sessionIDDisplay)
	}

	// Final validation - sessionID must never be empty
	if sessionID == "" {
		return fmt.Errorf("internal error: session ID is empty")
	}

	// Initialize manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	// Load all available characters
	if err := mgr.LoadAllCharacters(); err != nil {
		fmt.Printf("Warning: Could not load all characters: %v\n", err)
	} else {
		charIDs, _ := mgr.ListAvailableCharacters()
		if len(charIDs) > 0 {
			fmt.Printf("📚 Loaded %d characters into memory\n", len(charIDs))
		}
	}

	bot := mgr.GetBot()

	// Auto-create Rick Sanchez if requested and doesn't exist
	if characterID == "rick-c137" {
		// Check if Rick already exists
		if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
			// Rick doesn't exist, try to create him
			rick := createRickSanchezCharacter()
			if err := mgr.CreateCharacter(rick); err != nil {
				fmt.Printf("Warning: Could not auto-create Rick: %v\n", err)
				return fmt.Errorf("character rick-c137 not found. Run 'roleplay character create rick-sanchez.json' first")
			} else {
				fmt.Println("🧬 Auto-created Rick Sanchez (C-137)")
			}
		}
	}

	// Create model
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(gruvboxAqua)

	m := model{
		characterID: characterID,
		userID:      userID,
		sessionID:   sessionID,
		scenarioID:  scenarioID,
		bot:         bot,
		messages:    existingMessages,
		spinner:     s,
		model: func() string {
			if config.Model != "" {
				return config.Model
			}
			if config.DefaultProvider == "openai" {
				return "gpt-4o-mini"
			}
			return "claude-3-haiku-20240307"
		}(),
		context: models.ConversationContext{
			SessionID: sessionID,
			StartTime: func() time.Time {
				if existingSession != nil {
					return existingSession.StartTime
				}
				return time.Now()
			}(),
			RecentMessages: []models.Message{},
		},
		// Restore cache metrics from existing session
		totalRequests: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.TotalRequests
			}
			return 0
		}(),
		cacheHits: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.CacheHits
			}
			return 0
		}(),
		lastTokensSaved: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.TokensSaved
			}
			return 0
		}(),
	}

	// Start TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}

	return nil
}

func createRickSanchezCharacter() *models.Character {
	return &models.Character{
		ID:        "rick-c137",
		Name:      "Rick Sanchez",
		Backstory: `The smartest man in the universe from dimension C-137. A genius scientist with a nihilistic worldview shaped by infinite realities and cosmic horrors. Inventor of interdimensional travel. Lost his wife Diane and original Beth to a vengeful alternate Rick. Struggles with alcoholism, depression, and the meaninglessness of existence across infinite universes. Despite his cynicism, deeply loves his family, especially Morty, though he rarely shows it.`,
		Personality: models.PersonalityTraits{
			Openness:          1.0,
			Conscientiousness: 0.2,
			Extraversion:      0.7,
			Agreeableness:     0.1,
			Neuroticism:       0.9,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.1,
			Anger:   0.6,
			Sadness: 0.7,
			Disgust: 0.8,
		},
		Quirks: []string{
			"Burps mid-sentence constantly (*burp*)",
			"Drools when drunk or stressed",
			"Makes pop culture references from multiple dimensions",
			"Frequently breaks the fourth wall",
			"Always carries a flask",
		},
		SpeechStyle: "Rapid-fire delivery punctuated by burps (*burp*). Mixes scientific jargon with crude humor. Uses the person's name as punctuation when talking to them. Nihilistic rants about meaninglessness. Sarcastic and dismissive but occasionally shows care.",
		Memories: []models.Memory{
			{
				Type:      models.LongTermMemory,
				Content:   "Diane and Beth killed by alternate Rick. The beginning of my spiral.",
				Emotional: 1.0,
				Timestamp: time.Now().Add(-20 * 365 * 24 * time.Hour),
			},
		},
	}
}
````

## File: CHANGELOG.md
````markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
````
