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