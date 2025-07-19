package bridge_test

import (
	"context"
	"fmt"
	"log"

	"github.com/dotcommander/roleplay/pkg/bridge"
)

func ExampleCharactersConverter() {
	// Create a converter and registry
	converter := bridge.NewCharactersConverter()
	registry := bridge.NewConverterRegistry()
	_ = registry.Register(converter)

	// Example character from the Characters system
	charactersChar := &bridge.CharactersCharacter{
		ID:        "warrior-001",
		Name:      "Valeria the Bold",
		Age:       28,
		Gender:    "female",
		Archetype: "warrior",
		Traits:    []string{"brave", "loyal", "determined", "honorable"},
		Experiences: []string{
			"veteran of the Northern Wars",
			"captain of the Royal Guard",
			"dragon slayer",
		},
		Narrative: "Born into a family of warriors, Valeria proved herself on countless battlefields. Her unwavering loyalty and tactical brilliance earned her the position of Captain of the Royal Guard at the young age of 25.",
		Attributes: map[string]interface{}{
			"origin": "Mountain Clans of the North",
			"skills": []interface{}{"swordsmanship", "tactics", "leadership"},
			"weapon": "Ancestral greatsword 'Oathkeeper'",
		},
		Persona: &bridge.CharactersPersona{
			CoreIdentity: bridge.CharactersCoreIdentity{
				Worldview:      "Honor and duty above all else",
				CoreMotivation: "Protect the innocent and uphold justice",
				CoreFear:       "Failing those who depend on her",
				Secret:         "Doubts her worthiness to wield Oathkeeper",
			},
			Communication: bridge.CharactersCommunication{
				VoicePacing:       "measured and commanding",
				VocabularyTier:    "formal military",
				SentenceStructure: "direct and concise",
				VerbalTics:        []string{"By my oath...", "Steel yourself"},
				ForbiddenTopics:   []string{"her father's death", "the Betrayal at Ironhold"},
			},
			Behavior: bridge.CharactersBehavior{
				DecisionHeuristic: "What would bring the most honor?",
				ConflictStyle:     "direct confrontation with honor",
				Quirks:            []string{"polishes armor when thinking", "never breaks eye contact"},
			},
		},
	}

	// Convert to universal format
	ctx := context.Background()
	universal, err := converter.ToUniversal(ctx, charactersChar)
	if err != nil {
		log.Fatal(err)
	}

	// Display the converted character
	fmt.Printf("Name: %s\n", universal.Name)
	fmt.Printf("Description: %s\n", universal.Description)
	fmt.Printf("Background: %s\n", universal.Background)
	fmt.Printf("Speech Style: %s\n", universal.SpeechStyle)
	fmt.Printf("Personality (OCEAN):\n")
	fmt.Printf("  Openness: %.2f\n", universal.Personality.Openness)
	fmt.Printf("  Conscientiousness: %.2f\n", universal.Personality.Conscientiousness)
	fmt.Printf("  Extraversion: %.2f\n", universal.Personality.Extraversion)
	fmt.Printf("  Agreeableness: %.2f\n", universal.Personality.Agreeableness)
	fmt.Printf("  Neuroticism: %.2f\n", universal.Personality.Neuroticism)

	// Output:
	// Name: Valeria the Bold
	// Description: 28 year old female warrior
	// Background: Born into a family of warriors, Valeria proved herself on countless battlefields. Her unwavering loyalty and tactical brilliance earned her the position of Captain of the Royal Guard at the young age of 25.
	// Speech Style: measured and commanding pacing, formal military vocabulary, direct and concise sentences
	// Personality (OCEAN):
	//   Openness: 0.50
	//   Conscientiousness: 0.67
	//   Extraversion: 0.71
	//   Agreeableness: 0.71
	//   Neuroticism: 0.50
}

func ExampleCharactersConverter_bidirectional() {
	converter := bridge.NewCharactersConverter()
	ctx := context.Background()

	// Start with a UniversalCharacter
	universal := &bridge.UniversalCharacter{
		ID:          "mage-001",
		Name:        "Eldara the Wise",
		Description: "Ancient elven archmage",
		Background:  "Having lived for over 500 years, Eldara has witnessed the rise and fall of empires.",
		Personality: bridge.PersonalityTraits{
			Openness:          0.95, // Very open to new ideas and experiences
			Conscientiousness: 0.80, // Disciplined in magical studies
			Extraversion:      0.30, // Introverted scholar
			Agreeableness:     0.70, // Kind but can be stern
			Neuroticism:       0.20, // Emotionally stable
		},
		Traits:       []string{"wise", "patient", "mysterious", "powerful"},
		Behaviors:    []string{"speaks in riddles", "observes before acting"},
		SpeechStyle:  "archaic and poetic",
		Quirks:       []string{"eyes glow when casting", "levitates instead of sitting"},
		Catchphrases: []string{"Time reveals all truths", "Magic flows where will goes"},
		Topics:       []string{"ancient history", "arcane theory", "prophecies"},
		Motivations:  []string{"preserve magical knowledge", "guide young mages"},
		Fears:        []string{"the death of magic", "forgotten lore being lost forever"},
	}

	// Convert to Roleplay format
	roleplayChar, err := converter.FromUniversal(ctx, universal)
	if err != nil {
		log.Fatal(err)
	}

	// The converter returns a roleplay Character
	fmt.Printf("Converted successfully to roleplay.Character format\n")
	_ = roleplayChar // Would be used in the roleplay system

	// Output:
	// Converted successfully to roleplay.Character format
}