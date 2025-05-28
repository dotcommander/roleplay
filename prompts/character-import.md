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