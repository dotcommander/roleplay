# Bridge Troubleshooting Guide

Common issues and solutions when using the Characters ↔ Roleplay bridge.

## Import Issues

### "Auto-detected format: characters" but conversion fails

**Problem**: File is detected as Characters format but conversion fails.

**Solutions**:
1. Check JSON syntax:
   ```bash
   jq . character.json  # Validates JSON syntax
   ```

2. Verify required fields:
   ```json
   {
     "name": "Character Name",
     "traits": ["trait1", "trait2"]  // At minimum
   }
   ```

3. Use verbose mode for details:
   ```bash
   roleplay character import character.json --verbose
   ```

### Format auto-detection incorrect

**Problem**: Character file detected as wrong format.

**Solutions**:
1. Explicitly specify format:
   ```bash
   roleplay character import character.json --source characters
   ```

2. Check file structure - Characters format should have `traits`, `attributes`, or `persona` fields

3. Rename file extension if needed:
   - `.json` for Characters format
   - `.md` for markdown format

### Missing personality traits

**Problem**: OCEAN scores all show 0.5 (neutral).

**Solutions**:
1. Add more descriptive traits:
   ```json
   {
     "traits": ["brave", "disciplined", "kind", "creative"]
   }
   ```

2. Include persona information:
   ```json
   {
     "persona": {
       "voice_pacing": "Quick, energetic speech",
       "quirks": ["Always smiling", "Gestures while talking"]
     }
   }
   ```

3. Check trait mappings in verbose output

### Character responses seem off

**Problem**: Imported character doesn't respond authentically.

**Solutions**:
1. Verify speech patterns imported:
   ```bash
   cat ~/.config/roleplay/characters/CHARACTER_ID.json | jq .speech_style
   ```

2. Add more specific persona data:
   ```json
   {
     "persona": {
       "catchphrases": ["Signature phrase", "Common saying"],
       "forbidden_topics": ["sensitive subject"],
       "voice_pacing": "Detailed communication style"
     }
   }
   ```

3. Include rich backstory for context

## Export Issues

### Export command not found

**Problem**: `characters export` command doesn't exist.

**Solutions**:
1. Rebuild Characters project:
   ```bash
   cd /path/to/characters
   go build -o characters .
   ```

2. Check if export command is registered in main.go

3. Verify roleplay dependency in go.mod

### Conversion warnings

**Problem**: "No backstory found" or similar warnings.

**Solutions**:
1. Add missing fields to source character:
   ```json
   {
     "backstory": "Character's history and motivation",
     "persona": {
       "quirks": ["Behavioral trait"],
       "voice_pacing": "How they speak"
     }
   }
   ```

2. Use verbose mode to see what's missing:
   ```bash
   characters export CHARACTER_ID --format roleplay --verbose
   ```

## Personality Mapping Issues

### Unexpected OCEAN scores

**Problem**: Personality scores don't match expected character.

**Solutions**:
1. Check trait mappings:
   ```bash
   # Look for specific trait mappings in bridge code
   grep -A 5 -B 5 "trait_name" pkg/bridge/mappings.go
   ```

2. Add context for better analysis:
   ```json
   {
     "traits": ["brave", "leader"],
     "persona": {
       "voice_pacing": "Commanding, direct speech",
       "behaviors": ["Takes charge in crisis"]
     }
   }
   ```

3. Use more specific traits instead of generic ones

### Character seems too neutral

**Problem**: All OCEAN scores around 0.5.

**Solutions**:
1. Add more distinctive traits:
   ```json
   {
     "traits": ["extremely organized", "deeply creative", "very shy"]
   }
   ```

2. Include behavioral descriptions:
   ```json
   {
     "attributes": {
       "personality": {
         "quirks": ["Always early", "Perfectionist", "Avoids crowds"]
       }
     }
   }
   ```

## File and Path Issues

### Character file not found

**Problem**: "File not found" error during import.

**Solutions**:
1. Use absolute path:
   ```bash
   roleplay character import /full/path/to/character.json
   ```

2. Check current directory:
   ```bash
   ls -la *.json
   ```

3. Verify file permissions:
   ```bash
   chmod 644 character.json
   ```

### Configuration directory issues

**Problem**: Characters not saving to expected location.

**Solutions**:
1. Check config directory:
   ```bash
   ls -la ~/.config/roleplay/characters/
   ```

2. Create directory if missing:
   ```bash
   mkdir -p ~/.config/roleplay/characters/
   ```

3. Verify permissions:
   ```bash
   chmod 755 ~/.config/roleplay/
   ```

## Verbose Mode Debugging

Always use verbose mode when troubleshooting:

```bash
# For imports
roleplay character import character.json --verbose

# For exports  
characters export CHARACTER_ID --format roleplay --verbose
```

### Verbose Output Interpretation

**Good conversion**:
```
Auto-detected format: characters
Successfully imported character: Name
Personality traits:
  Conscientiousness: 0.74
Quirks:
  - specific behavioral trait
Speech Style: distinctive communication pattern
```

**Problematic conversion**:
```
Conversion warnings:
  ⚠️  No backstory found
  ⚠️  No quirks found
  ⚠️  No speech style defined
```

## Common Character Format Issues

### Nested attribute structure

**Problem**: Attributes not extracted properly.

**Solution**: Flatten structure or use proper nesting:
```json
{
  "attributes": {
    "personality": {
      "quirks": ["trait1", "trait2"]
    }
  }
}
```

### Inconsistent field names

**Problem**: Fields not mapping correctly.

**Solution**: Use standard field names:
- `traits` (not `characteristics`)
- `backstory` (not `background` or `history`)
- `quirks` (not `habits` or `behaviors`)

## Performance Issues

### Slow conversion

**Problem**: Import/export takes too long.

**Solutions**:
1. Check file size - very large character files may be slow
2. Verify system resources
3. Use specific format instead of auto-detection

### Memory usage

**Problem**: High memory usage during conversion.

**Solutions**:
1. Process characters individually instead of batch
2. Check for memory leaks in custom converters
3. Use streaming for large datasets

## Getting Help

### Debug Information

Collect this information when reporting issues:

```bash
# System info
go version
uname -a

# File info
file character.json
head -20 character.json

# Verbose output
roleplay character import character.json --verbose 2>&1

# Character data (if comfortable sharing)
cat ~/.config/roleplay/characters/CHARACTER_ID.json
```

### Common Solutions Summary

1. **Always use verbose mode** for detailed feedback
2. **Check JSON syntax** with `jq` or similar tool
3. **Verify required fields** exist in source
4. **Use specific trait names** instead of generic ones
5. **Include rich context** (backstory, persona, quirks)
6. **Check file permissions** and paths
7. **Rebuild projects** if commands missing

### Bridge Limitations

Current known limitations:
- Only supports Characters → Roleplay direction (not reverse)
- Some complex nested attributes may not map perfectly
- Cultural/contextual trait interpretations may vary
- Large character sets may require individual processing