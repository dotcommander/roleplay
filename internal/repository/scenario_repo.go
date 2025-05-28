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
