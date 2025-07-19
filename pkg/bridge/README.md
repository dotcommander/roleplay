# Bridge Package - Character Format Converter

The bridge package provides intelligent conversion between different character formats, enabling seamless interoperability between character creation and roleplay systems.

## Overview

This package implements a sophisticated character conversion system that:
- **Preserves Character Essence**: Maintains personality, backstory, quirks, and speech patterns
- **Intelligent Mapping**: Converts personality traits to OCEAN psychological model
- **Multi-Format Support**: Extensible architecture for adding new character formats
- **Context Analysis**: Analyzes speech patterns and behaviors for enhanced accuracy

## Core Components

### Universal Character Format
A comprehensive interchange format that captures:
- **Personality**: OCEAN traits + original trait descriptions
- **Identity**: Name, background, appearance, age
- **Behavior**: Quirks, speech patterns, catchphrases
- **Content**: Backstory, interests, relationships
- **Metadata**: Source format, conversion details

### Trait Analysis Engine
Intelligent personality mapping with:
- **100+ Trait Mappings**: Comprehensive vocabulary covering personality spectrum
- **Context Awareness**: Analyzes speech patterns and behaviors
- **OCEAN Conversion**: Maps descriptive traits to psychological scores
- **Validation**: Ensures realistic and balanced personality profiles

### Converter Interface
Clean, extensible architecture:
```go
type CharacterConverter interface {
    CanConvert(data interface{}) bool
    ToUniversal(data interface{}) (*UniversalCharacter, error)
    FromUniversal(char *UniversalCharacter) (interface{}, error)
}
```

## Usage

### Import a Character
```go
import "github.com/dotcommander/roleplay/pkg/bridge"

// Auto-detect format and convert
converter := bridge.NewCharactersConverter()
if converter.CanConvert(characterData) {
    universal, err := converter.ToUniversal(characterData)
    if err != nil {
        return err
    }
    
    // Convert to target format
    roleplayChar, err := converter.FromUniversal(universal)
    return err
}
```

### CLI Integration
The bridge seamlessly integrates with command-line tools:

```bash
# Roleplay side - auto-detects Characters format
roleplay character import character.json --verbose

# Characters side - export to roleplay format  
characters export warrior-123 --format roleplay
```

## Supported Formats

### Characters Format
- **Traits Array**: Descriptive personality traits
- **Differential System**: Base templates + modifiers
- **Rich Personas**: Worldview, voice pacing, catchphrases
- **Complex Attributes**: Nested personality, physical, background data

### Roleplay Format
- **OCEAN Personality**: Scientific personality model
- **Memory System**: Multi-tier memory management
- **Evolution Tracking**: Personality change over time
- **Conversation Context**: Speech style and behavioral patterns

## Architecture

```
pkg/bridge/
├── format.go              # UniversalCharacter definition
├── converter.go           # Converter interface and registry
├── trait_analyzer.go      # Personality trait analysis
├── mappings.go            # Trait-to-OCEAN mapping database
├── characters_converter.go # Characters format converter
└── README.md              # This file
```

## Key Features

### Intelligent Trait Mapping
Converts descriptive traits to OCEAN scores:
```
"brave" + "disciplined" + "protective" + "determined"
  ↓
Conscientiousness: 0.74 (High - disciplined, protective, determined)
Neuroticism: 0.50 (Neutral - brave balances other factors)
```

### Speech Pattern Analysis
Extracts personality from communication style:
```
"Short, clipped sentences. Direct communication."
  ↓
Traits: ["direct", "efficient", "serious"]
  ↓
Conscientiousness: +0.1, Extraversion: +0.1
```

### Multi-Source Data Extraction
Comprehensively extracts character data:
- **Quirks**: From both `attributes.personality.quirks` and `persona.quirks`
- **Speech**: From `voice_pacing`, `catchphrases`, and `dialogue_examples`
- **Personality**: From trait arrays, behaviors, and narrative analysis

### Validation and Feedback
- **Conversion Warnings**: Alerts for missing or incomplete data
- **Score Validation**: Ensures OCEAN scores remain within realistic bounds
- **Verbose Mode**: Detailed conversion feedback for debugging

## Examples

### Basic Conversion
```go
// Characters format character
character := CharactersFormat{
    Name: "Lyra Dragonbane",
    Traits: []string{"brave", "disciplined", "protective"},
    Backstory: "Dragon hunter seeking redemption...",
}

// Convert to universal format
converter := NewCharactersConverter()
universal, err := converter.ToUniversal(character)

// Result: UniversalCharacter with OCEAN personality,
// preserved backstory, and extracted speech patterns
```

### Production Example
The bridge successfully converted a test character with:
- ✅ **Traits**: `["brave", "disciplined", "protective", "haunted", "determined"]`
- ✅ **OCEAN Result**: Conscientiousness 0.74 (correctly high from multiple relevant traits)
- ✅ **Content**: Full backstory, 4 quirks, speech style, catchphrases
- ✅ **Authenticity**: Character responds with exact original catchphrases

## Testing

Run the conversion tests:
```bash
go test ./pkg/bridge/...
```

Integration test with real character:
```bash
# Import Characters format
roleplay character import test-character.json --verbose

# Verify authentic conversation
roleplay chat "Hello warrior" --character test-warrior-001 --user test
```

## Extension Points

### Adding New Formats
1. Implement the `CharacterConverter` interface
2. Register with the `ConverterRegistry`
3. Add format detection logic
4. Implement bidirectional conversion methods

### Custom Trait Mappings
```go
analyzer := NewTraitAnalyzer()
analyzer.AddMapping("custom_trait", PersonalityTraits{
    Openness: 0.3,
    Conscientiousness: 0.2,
})
```

### Enhanced Context Analysis
Extend the trait analyzer to handle additional context sources:
- Dialogue examples
- Relationship descriptions  
- Goal statements
- Background narratives

## Performance

- **Conversion Speed**: <1 second per character
- **Memory Efficient**: Minimal overhead for large character sets
- **Scalable**: Designed for batch operations
- **Robust**: Graceful handling of missing or malformed data

## Contributing

When adding new converters or mappings:
1. Follow the existing interface patterns
2. Add comprehensive tests
3. Update documentation with examples
4. Ensure backward compatibility