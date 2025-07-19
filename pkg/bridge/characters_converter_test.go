package bridge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	roleplayModels "github.com/dotcommander/roleplay/internal/models"
)

func TestCharactersConverter_CanConvert(t *testing.T) {
	converter := NewCharactersConverter()

	tests := []struct {
		name     string
		data     interface{}
		expected bool
	}{
		{
			name: "CharactersCharacter pointer",
			data: &CharactersCharacter{
				ID:   "test-1",
				Name: "Test Character",
			},
			expected: true,
		},
		{
			name: "CharactersCharacter value",
			data: CharactersCharacter{
				ID:   "test-2",
				Name: "Test Character 2",
			},
			expected: true,
		},
		{
			name: "Map with Characters fields",
			data: map[string]interface{}{
				"name":       "Test",
				"traits":     []string{"brave", "loyal"},
				"attributes": map[string]interface{}{},
			},
			expected: true,
		},
		{
			name: "Map with archetype",
			data: map[string]interface{}{
				"name":      "Test",
				"archetype": "warrior",
			},
			expected: true,
		},
		{
			name:     "Wrong type",
			data:     "not a character",
			expected: false,
		},
		{
			name:     "Empty map",
			data:     map[string]interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.CanConvert(tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCharactersConverter_ToUniversal(t *testing.T) {
	converter := NewCharactersConverter()
	ctx := context.Background()

	t.Run("Basic conversion", func(t *testing.T) {
		char := &CharactersCharacter{
			ID:        "warrior-123",
			Name:      "Lyra",
			Age:       25,
			Gender:    "female",
			Archetype: "warrior",
			Traits:    []string{"brave", "loyal", "determined"},
			Experiences: []string{"combat training", "leadership"},
			Narrative: "A seasoned warrior with years of combat experience.",
			Attributes: map[string]interface{}{
				"origin": "Northern Kingdom",
				"skills": []interface{}{"swordsmanship", "tactics"},
			},
		}

		result, err := converter.ToUniversal(ctx, char)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "warrior-123", result.ID)
		assert.Equal(t, "Lyra", result.Name)
		assert.Equal(t, "25 year old female warrior", result.Description)
		assert.Equal(t, "A seasoned warrior with years of combat experience.", result.Background)
		assert.Equal(t, []string{"brave", "loyal", "determined"}, result.Traits)
		assert.Contains(t, result.Behaviors, "combat training")
		assert.Contains(t, result.Behaviors, "leadership")
		assert.Equal(t, "characters", result.Source)
		assert.Equal(t, "1.0", result.Version)

		// Check personality was analyzed
		assert.GreaterOrEqual(t, result.Personality.Conscientiousness, 0.5) // brave, loyal, determined suggest high conscientiousness
		assert.GreaterOrEqual(t, result.Personality.Extraversion, 0.5)     // leadership suggests extraversion

		// Check source data preservation
		assert.Equal(t, "warrior", result.SourceData["archetype"])
		assert.Equal(t, 25, result.SourceData["age"])
		assert.Equal(t, "female", result.SourceData["gender"])
	})

	t.Run("With persona", func(t *testing.T) {
		char := &CharactersCharacter{
			ID:        "mystic-456",
			Name:      "Eldrin",
			Age:       150,
			Gender:    "male",
			Archetype: "mystic",
			Traits:    []string{"wise", "mysterious", "patient"},
			Persona: &CharactersPersona{
				CoreIdentity: CharactersCoreIdentity{
					Worldview:      "All things are connected through the cosmic web",
					CoreMotivation: "Seeking universal truth",
					CoreFear:       "The corruption of magic",
					Secret:         "Once misused power for personal gain",
				},
				Communication: CharactersCommunication{
					VoicePacing:       "slow and deliberate",
					VocabularyTier:    "archaic",
					SentenceStructure: "complex",
					VerbalTics:        []string{"Indeed...", "As the ancients say..."},
					ForbiddenTopics:   []string{"dark magic", "necromancy"},
				},
				Behavior: CharactersBehavior{
					DecisionHeuristic: "Consult the cosmic patterns",
					ConflictStyle:     "avoidance through wisdom",
					Quirks:            []string{"speaks in riddles", "never gives direct answers"},
				},
			},
		}

		result, err := converter.ToUniversal(ctx, char)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "slow and deliberate pacing, archaic vocabulary, complex sentences", result.SpeechStyle)
		assert.Equal(t, []string{"speaks in riddles", "never gives direct answers"}, result.Quirks)
		assert.Equal(t, []string{"Indeed...", "As the ancients say..."}, result.Catchphrases)
		assert.Equal(t, []string{"dark magic", "necromancy"}, result.Boundaries)
		assert.Contains(t, result.Motivations, "Seeking universal truth")
		assert.Contains(t, result.Fears, "The corruption of magic")
		assert.Contains(t, result.SystemPrompt, "All things are connected through the cosmic web")
	})

	t.Run("Map conversion", func(t *testing.T) {
		data := map[string]interface{}{
			"id":        "rogue-789",
			"name":      "Shadow",
			"age":       30.0, // JSON numbers come as float64
			"gender":    "non-binary",
			"archetype": "rogue",
			"traits":    []interface{}{"stealthy", "cunning", "independent"},
			"experiences": []interface{}{"thievery", "espionage"},
			"attributes": map[string]interface{}{
				"skills": []interface{}{"lockpicking", "disguise"},
			},
		}

		result, err := converter.ToUniversal(ctx, data)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "rogue-789", result.ID)
		assert.Equal(t, "Shadow", result.Name)
		assert.Contains(t, result.Topics, "lockpicking")
		assert.Contains(t, result.Topics, "disguise")
	})
}

func TestCharactersConverter_FromUniversal(t *testing.T) {
	converter := NewCharactersConverter()
	ctx := context.Background()

	t.Run("Basic conversion", func(t *testing.T) {
		uc := &UniversalCharacter{
			ID:          "test-123",
			Name:        "Elena",
			Description: "A skilled diplomat",
			Background:  "Raised in the royal court, trained in diplomacy and statecraft.",
			Personality: PersonalityTraits{
				Openness:          0.8,
				Conscientiousness: 0.7,
				Extraversion:      0.9,
				Agreeableness:     0.8,
				Neuroticism:       0.3,
			},
			Traits:      []string{"charismatic", "diplomatic", "perceptive"},
			Behaviors:   []string{"reads body language", "chooses words carefully"},
			SpeechStyle: "formal and eloquent",
			Quirks:      []string{"adjusts jewelry when nervous"},
			Catchphrases: []string{"Let us find common ground"},
			Topics:      []string{"politics", "etiquette", "history"},
			Motivations: []string{"peace between nations", "justice for all"},
			Fears:       []string{"war", "betrayal"},
			Relationships: map[string]string{
				"King Aldric": "Former mentor",
				"Lady Vera":   "Trusted advisor",
			},
			UpdatedAt: time.Now(),
			SourceData: map[string]interface{}{
				"age":       28,
				"gender":    "female",
				"archetype": "diplomat",
			},
		}

		result, err := converter.FromUniversal(ctx, uc)
		require.NoError(t, err)
		require.NotNil(t, result)

		char, ok := result.(*roleplayModels.Character)
		require.True(t, ok)

		assert.Equal(t, "test-123", char.ID)
		assert.Equal(t, "Elena", char.Name)
		assert.Equal(t, "Raised in the royal court, trained in diplomacy and statecraft.", char.Backstory)
		assert.Equal(t, "formal and eloquent", char.SpeechStyle)
		assert.Equal(t, []string{"adjusts jewelry when nervous"}, char.Quirks)
		assert.Equal(t, []string{"Let us find common ground"}, char.CatchPhrases)
		assert.Equal(t, "28", char.Age)
		assert.Equal(t, "female", char.Gender)
		assert.Equal(t, "diplomat", char.Occupation)

		// Check personality conversion
		assert.Equal(t, 0.8, char.Personality.Openness)
		assert.Equal(t, 0.7, char.Personality.Conscientiousness)
		assert.Equal(t, 0.9, char.Personality.Extraversion)
		assert.Equal(t, 0.8, char.Personality.Agreeableness)
		assert.Equal(t, 0.3, char.Personality.Neuroticism)

		// Check fields extraction
		assert.Contains(t, char.Goals, "peace between nations")
		assert.Contains(t, char.Goals, "justice for all")
		assert.Contains(t, char.Fears, "war")
		assert.Contains(t, char.Fears, "betrayal")
		assert.Equal(t, "Former mentor", char.Relationships["King Aldric"])
		assert.Equal(t, "Trusted advisor", char.Relationships["Lady Vera"])

		// Check behaviors
		assert.Contains(t, char.BehaviorPatterns, "reads body language")
		assert.Contains(t, char.BehaviorPatterns, "chooses words carefully")
	})

	t.Run("With persona source data", func(t *testing.T) {
		uc := &UniversalCharacter{
			ID:   "warrior-456",
			Name: "Marcus",
			Personality: PersonalityTraits{
				Openness:          0.4,
				Conscientiousness: 0.9,
				Extraversion:      0.6,
				Agreeableness:     0.5,
				Neuroticism:       0.3,
			},
			UpdatedAt: time.Now(),
			SourceData: map[string]interface{}{
				"persona": &CharactersPersona{
					CoreIdentity: CharactersCoreIdentity{
						Worldview: "Honor above all else",
					},
					Behavior: CharactersBehavior{
						DecisionHeuristic: "Follow the warrior code",
						ConflictStyle:     "Direct confrontation",
					},
				},
			},
		}

		result, err := converter.FromUniversal(ctx, uc)
		require.NoError(t, err)

		char, ok := result.(*roleplayModels.Character)
		require.True(t, ok)

		assert.Equal(t, "Follow the warrior code", char.DecisionMaking)
		assert.Equal(t, "Direct confrontation", char.ConflictStyle)
		assert.Equal(t, "Honor above all else", char.WorldView)
	})

	t.Run("With dialogue examples", func(t *testing.T) {
		uc := &UniversalCharacter{
			ID:   "test-789",
			Name: "Sage",
			Personality: PersonalityTraits{
				Openness:          0.9,
				Conscientiousness: 0.7,
				Extraversion:      0.4,
				Agreeableness:     0.8,
				Neuroticism:       0.2,
			},
			Examples: []ConversationExample{
				{
					User:      "What is the meaning of life?",
					Character: "Ah, the eternal question. Perhaps the meaning is in the seeking itself.",
				},
				{
					User:      "Can you teach me magic?",
					Character: "Magic cannot be taught, only discovered within oneself.",
				},
			},
			UpdatedAt: time.Now(),
		}

		result, err := converter.FromUniversal(ctx, uc)
		require.NoError(t, err)

		char, ok := result.(*roleplayModels.Character)
		require.True(t, ok)

		assert.Len(t, char.DialogueExamples, 2)
		assert.Contains(t, char.DialogueExamples[0], "What is the meaning of life?")
		assert.Contains(t, char.DialogueExamples[0], "Ah, the eternal question")
		assert.Contains(t, char.DialogueExamples[1], "Can you teach me magic?")
		assert.Contains(t, char.DialogueExamples[1], "Magic cannot be taught")
	})
}

func TestCharactersConverter_TraitExtraction(t *testing.T) {
	converter := NewCharactersConverter()

	t.Run("Extract beliefs and morals", func(t *testing.T) {
		traits := []string{
			"believes in justice",
			"faithful servant",
			"honorable",
			"strong moral compass",
			"ethical leader",
			"principled",
			"brave",
			"clever",
		}

		beliefs, morals := converter.extractBeliefsAndMorals(traits)

		assert.Contains(t, beliefs, "believes in justice")
		assert.Contains(t, beliefs, "faithful servant")
		assert.Contains(t, morals, "honorable")
		assert.Contains(t, morals, "strong moral compass")
		assert.Contains(t, morals, "ethical leader")
		assert.Contains(t, morals, "principled")
		assert.NotContains(t, beliefs, "brave")
		assert.NotContains(t, morals, "clever")
	})

	t.Run("Extract flaws and strengths", func(t *testing.T) {
		traits := []string{
			"brave warrior",
			"loyal friend",
			"arrogant noble",
			"stubborn as a mule",
			"creative thinker",
			"selfish at times",
			"compassionate healer",
			"reckless adventurer",
		}

		flaws, strengths := converter.extractFlawsAndStrengths(traits)

		assert.Contains(t, strengths, "brave warrior")
		assert.Contains(t, strengths, "loyal friend")
		assert.Contains(t, strengths, "creative thinker")
		assert.Contains(t, strengths, "compassionate healer")
		assert.Contains(t, flaws, "arrogant noble")
		assert.Contains(t, flaws, "stubborn as a mule")
		assert.Contains(t, flaws, "selfish at times")
		assert.Contains(t, flaws, "reckless adventurer")
	})
}