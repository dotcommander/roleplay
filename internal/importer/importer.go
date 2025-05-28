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
