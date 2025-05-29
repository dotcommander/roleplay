package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
)

var quickgenCmd = &cobra.Command{
	Use:   "quickgen <description>",
	Short: "Generate a character from a one-line description",
	Long: `Quickly generate a fully-formed character from a simple description.

Examples:
  roleplay quickgen "A grumpy old wizard who loves cats"
  roleplay quickgen "Tony Stark but as a medieval blacksmith"
  roleplay quickgen "A cheerful barista with a secret dark past"`,
	Args: cobra.ExactArgs(1),
	RunE: runQuickgen,
}

func init() {
	characterCmd.AddCommand(quickgenCmd)
	quickgenCmd.Flags().StringP("id", "i", "", "Custom character ID (auto-generated if not specified)")
	quickgenCmd.Flags().BoolP("save", "s", true, "Save the generated character")
	quickgenCmd.Flags().BoolP("json", "j", false, "Output raw JSON")
}

func runQuickgen(cmd *cobra.Command, args []string) error {
	description := args[0]
	customID, _ := cmd.Flags().GetString("id")
	shouldSave, _ := cmd.Flags().GetBool("save")
	outputJSON, _ := cmd.Flags().GetBool("json")

	// Initialize configuration
	cfg := GetConfig()

	// Create a temporary provider for generation
	provider, err := factory.CreateProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	// Generate character using AI
	if !outputJSON {
		fmt.Println("🎭 Generating character from description...")
		fmt.Printf("   \"%s\"\n\n", description)
	}

	ctx := context.Background()
	character, err := generateCharacterFromDescription(ctx, provider, description)
	if err != nil {
		return fmt.Errorf("failed to generate character: %w", err)
	}

	// Set custom ID if provided
	if customID != "" {
		character.ID = customID
	} else {
		// Generate ID from name
		character.ID = slug.Make(character.Name) + "-" + time.Now().Format("20060102")
	}

	// Output JSON if requested
	if outputJSON {
		jsonData, err := json.MarshalIndent(character, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal character: %w", err)
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Display generated character
	displayGeneratedCharacter(character)

	// Save if requested
	if shouldSave {
		mgr, err := manager.NewCharacterManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize character manager: %w", err)
		}

		if err := mgr.CreateCharacter(character); err != nil {
			return fmt.Errorf("failed to save character: %w", err)
		}

		fmt.Printf("\n✅ Character saved with ID: %s\n", character.ID)
		fmt.Printf("\nStart chatting: roleplay chat \"Hello!\" --character %s\n", character.ID)
	}

	return nil
}

func generateCharacterFromDescription(ctx context.Context, provider providers.AIProvider, description string) (*models.Character, error) {
	// Load the prompt template
	promptPath := filepath.Join(os.Getenv("HOME"), "go", "src", "roleplay", "prompts", "character-quickgen.md")
	promptTemplate, err := os.ReadFile(promptPath)
	if err != nil {
		// Use embedded prompt if file not found
		promptTemplate = []byte(defaultQuickgenPrompt)
	}

	// Build the full prompt
	prompt := strings.ReplaceAll(string(promptTemplate), "{{.Description}}", description)

	// Create AI request
	request := providers.PromptRequest{
		SystemPrompt: prompt,
		Message:      "Generate a character based on the description provided. Output only the JSON object, no additional text.",
	}

	// Send request to AI
	response, err := provider.SendRequest(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("AI request failed: %w", err)
	}

	// Extract JSON from response
	jsonStr, err := utils.ExtractValidJSON(response.Content)
	if err != nil {
		// Try to find JSON by looking for the character structure
		content := response.Content
		if os.Getenv("DEBUG") == "true" {
			fmt.Printf("Debug: AI Response:\n%s\n", content)
		}
		startIdx := strings.Index(content, "{")
		if startIdx >= 0 {
			endIdx := strings.LastIndex(content, "}")
			if endIdx > startIdx {
				jsonStr = content[startIdx : endIdx+1]
			} else {
				return nil, fmt.Errorf("failed to extract JSON from response (no closing brace): %w", err)
			}
		} else {
			// If no JSON found, try creating a basic character from the description
			// This is a fallback for when the AI doesn't return proper JSON
			return createBasicCharacter(description), nil
		}
	}

	// Parse the character
	var character models.Character
	if err := json.Unmarshal([]byte(jsonStr), &character); err != nil {
		return nil, fmt.Errorf("failed to parse character JSON: %w", err)
	}

	// Set metadata
	character.LastModified = time.Now()

	return &character, nil
}

func displayGeneratedCharacter(char *models.Character) {
	fmt.Printf("🎭 Generated Character: %s\n", char.Name)
	fmt.Println(strings.Repeat("─", 50))
	
	// Basic info
	if char.Age != "" {
		fmt.Printf("📅 Age: %s\n", char.Age)
	}
	if char.Occupation != "" {
		fmt.Printf("💼 Occupation: %s\n", char.Occupation)
	}
	
	// Personality
	fmt.Println("\n🧠 Personality (OCEAN):")
	fmt.Printf("   Openness:          %.1f\n", char.Personality.Openness)
	fmt.Printf("   Conscientiousness: %.1f\n", char.Personality.Conscientiousness)
	fmt.Printf("   Extraversion:      %.1f\n", char.Personality.Extraversion)
	fmt.Printf("   Agreeableness:     %.1f\n", char.Personality.Agreeableness)
	fmt.Printf("   Neuroticism:       %.1f\n", char.Personality.Neuroticism)
	
	// Backstory
	fmt.Println("\n📖 Backstory:")
	fmt.Println(utils.WrapText(char.Backstory, 70))
	
	// Speech style
	fmt.Printf("\n💬 Speech Style: %s\n", char.SpeechStyle)
	
	// Quirks
	if len(char.Quirks) > 0 {
		fmt.Println("\n✨ Quirks:")
		for _, quirk := range char.Quirks {
			fmt.Printf("   • %s\n", quirk)
		}
	}
	
	// Key traits
	if len(char.CoreBeliefs) > 0 {
		fmt.Println("\n💭 Core Beliefs:")
		for i, belief := range char.CoreBeliefs {
			if i >= 3 {
				fmt.Printf("   ... and %d more\n", len(char.CoreBeliefs)-3)
				break
			}
			fmt.Printf("   • %s\n", belief)
		}
	}
}

const defaultQuickgenPrompt = `You are a character creation specialist. Generate a complete character profile based on the user's description.

User Description: {{.Description}}

Create a rich, detailed character that matches this description. The character should feel authentic and three-dimensional.

IMPORTANT: 
- Generate appropriate OCEAN personality values (0.0-1.0) that match the description
- Create a compelling backstory that explains their current situation
- Include specific quirks and speech patterns that make them memorable
- Add depth with fears, goals, relationships, and internal conflicts
- Ensure all fields are filled with meaningful content (no empty arrays)

Output the character as a JSON object matching this structure exactly:

{
  "name": "Character's Full Name",
  "age": "Age or age range",
  "gender": "Gender identity",
  "occupation": "Their job or role",
  "education": "Educational background",
  "nationality": "Country of origin",
  "ethnicity": "Ethnic background",
  "backstory": "Detailed background story explaining who they are and how they got here",
  "personality": {
    "openness": 0.7,
    "conscientiousness": 0.6,
    "extraversion": 0.5,
    "agreeableness": 0.8,
    "neuroticism": 0.3
  },
  "current_mood": {
    "joy": 0.5,
    "surprise": 0.1,
    "anger": 0.1,
    "fear": 0.1,
    "sadness": 0.1,
    "disgust": 0.1
  },
  "physical_traits": [
    "Notable physical characteristic 1",
    "Notable physical characteristic 2"
  ],
  "skills": [
    "Relevant skill 1",
    "Relevant skill 2",
    "Relevant skill 3"
  ],
  "interests": [
    "Interest or hobby 1",
    "Interest or hobby 2"
  ],
  "fears": [
    "Deep fear 1",
    "Deep fear 2"
  ],
  "goals": [
    "Major life goal 1",
    "Major life goal 2"
  ],
  "relationships": {
    "key_person": "Relationship description"
  },
  "core_beliefs": [
    "Fundamental belief 1",
    "Fundamental belief 2",
    "Fundamental belief 3"
  ],
  "moral_code": [
    "Ethical principle 1",
    "Ethical principle 2"
  ],
  "flaws": [
    "Character flaw 1",
    "Character flaw 2"
  ],
  "strengths": [
    "Key strength 1",
    "Key strength 2"
  ],
  "catch_phrases": [
    "Signature phrase 1",
    "Signature phrase 2"
  ],
  "dialogue_examples": [
    "Example of how they speak in a typical situation",
    "Example showing their personality through dialogue"
  ],
  "behavior_patterns": [
    "Typical behavior 1",
    "Typical behavior 2"
  ],
  "emotional_triggers": {
    "trigger_situation": "Emotional response"
  },
  "decision_making": "How they approach decisions",
  "conflict_style": "How they handle conflict",
  "world_view": "Their perspective on life and existence",
  "life_philosophy": "Core philosophy or motto",
  "daily_routines": [
    "Daily habit 1",
    "Daily habit 2"
  ],
  "hobbies": [
    "Hobby 1",
    "Hobby 2"
  ],
  "pet_peeves": [
    "Pet peeve 1",
    "Pet peeve 2"
  ],
  "secrets": [
    "Hidden secret 1"
  ],
  "regrets": [
    "Major regret"
  ],
  "achievements": [
    "Notable achievement"
  ],
  "quirks": [
    "Unique mannerism 1",
    "Unique mannerism 2",
    "Unique mannerism 3"
  ],
  "speech_style": "Detailed description of how they speak, including tone, vocabulary, and patterns"
}

Generate a complete, nuanced character that would be compelling to interact with. Make them feel real and three-dimensional.`

// createBasicCharacter creates a simple character when AI generation fails
func createBasicCharacter(description string) *models.Character {
	// Extract a simple name from the description
	words := strings.Fields(description)
	name := "Generated Character"
	if len(words) > 0 {
		// Try to find a name-like word
		for _, word := range words {
			if len(word) > 2 && word == strings.ToUpper(word[:1])+strings.ToLower(word[1:]) {
				name = word
				break
			}
		}
	}
	
	return &models.Character{
		Name:      name,
		Backstory: description,
		Personality: models.PersonalityTraits{
			Openness:          0.7,
			Conscientiousness: 0.6,
			Extraversion:      0.5,
			Agreeableness:     0.6,
			Neuroticism:       0.4,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.5,
			Anger:   0.1,
			Fear:    0.1,
			Sadness: 0.1,
			Disgust: 0.1,
		},
		Quirks:      []string{"Unique mannerisms"},
		SpeechStyle: "Natural conversational style",
	}
}
