# Characters Converter

The `CharactersConverter` provides bidirectional conversion between the Characters system format and the UniversalCharacter format used by the roleplay bridge.

## Features

### From Characters to Universal

The converter handles the following transformations:

1. **Basic Information**: Maps ID, name, age, gender, and archetype
2. **Personality Analysis**: Uses the TraitAnalyzer to convert traits into OCEAN scores
3. **Narrative Content**: Preserves backstory, narrative, and dialogue examples
4. **Differential System**: Merges base templates with modifiers
5. **Persona Mapping**: Converts detailed persona configurations including:
   - Core identity (worldview, motivations, fears)
   - Communication style (voice pacing, vocabulary, verbal tics)
   - Behavioral patterns (decision-making, conflict style, quirks)
   - Boundaries and forbidden topics

### From Universal to Roleplay

The converter creates a fully-featured roleplay Character with:

1. **OCEAN Personality**: Direct mapping of personality traits
2. **Extended Fields**: Maps to roleplay's rich character fields
3. **Emotional State**: Initializes with neutral mood
4. **Memory System**: Sets up empty memory arrays
5. **Behavior Patterns**: Preserves all behavioral information
6. **Dialogue Examples**: Formats conversation examples appropriately

## Character Structure Mapping

### Characters System Fields
- `ID`, `Name`, `Age`, `Gender` → Basic info
- `Archetype` → Used as occupation/role
- `Traits` → Analyzed for OCEAN personality
- `Experiences` → Mapped to behavior patterns
- `Attributes` → Flexible key-value storage
- `Differentials` → Template modifications
- `Persona` → Detailed AI behavior configuration

### Universal Character Fields
- `Personality` → OCEAN model scores
- `Background` → Narrative/backstory
- `SpeechStyle` → Derived from persona communication
- `Quirks`, `Catchphrases` → From persona behavior
- `Topics`, `Motivations`, `Fears` → Extracted from various sources
- `SourceData` → Preserves original format data

## Usage Example

```go
// Create converter
converter := NewCharactersConverter()

// Convert from Characters format
charactersChar := &CharactersCharacter{
    ID:        "warrior-001",
    Name:      "Valeria",
    Archetype: "warrior",
    Traits:    []string{"brave", "loyal"},
    // ... other fields
}

universal, err := converter.ToUniversal(ctx, charactersChar)

// Convert to Roleplay format
roleplayChar, err := converter.FromUniversal(ctx, universal)
```

## Trait Analysis

The converter uses sophisticated trait analysis to map personality descriptions to OCEAN scores:

- **Openness**: Creative, imaginative, curious traits
- **Conscientiousness**: Organized, disciplined, responsible traits  
- **Extraversion**: Social, outgoing, leadership traits
- **Agreeableness**: Cooperative, trusting, helpful traits
- **Neuroticism**: Anxious, stressed, emotional traits

The analysis considers:
1. Direct trait mappings (e.g., "brave" → higher conscientiousness)
2. Behavioral patterns (e.g., "leadership" → higher extraversion)
3. Narrative context (background story analysis)
4. Persona configuration (worldview, communication style)