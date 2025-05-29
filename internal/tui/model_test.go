package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/services"
)

func TestNewModel(t *testing.T) {
	// Create a test bot
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}
	bot := services.NewCharacterBot(cfg)

	// Create TUI config
	tuiConfig := Config{
		CharacterID: "test-char",
		UserID:      "test-user",
		SessionID:   "test-session",
		ScenarioID:  "test-scenario",
		Bot:         bot,
		Context: models.ConversationContext{
			SessionID: "test-session",
			StartTime: time.Now(),
		},
		Model:  "test-model",
		Width:  100,
		Height: 30,
		InitialMetrics: struct {
			TotalRequests int
			CacheHits     int
			TokensSaved   int
		}{
			TotalRequests: 10,
			CacheHits:     5,
			TokensSaved:   1000,
		},
	}

	model := NewModel(tuiConfig)

	// Verify model creation
	if model == nil {
		t.Fatal("NewModel returned nil")
	}

	// Verify fields
	if model.characterID != "test-char" {
		t.Errorf("CharacterID mismatch: got %s, want test-char", model.characterID)
	}
	if model.userID != "test-user" {
		t.Errorf("UserID mismatch: got %s, want test-user", model.userID)
	}
	if model.sessionID != "test-session" {
		t.Errorf("SessionID mismatch: got %s, want test-session", model.sessionID)
	}
	if model.scenarioID != "test-scenario" {
		t.Errorf("ScenarioID mismatch: got %s, want test-scenario", model.scenarioID)
	}
	if model.bot == nil {
		t.Error("Bot is nil")
	}
	if model.aiModel != "test-model" {
		t.Errorf("Model mismatch: got %s, want test-model", model.aiModel)
	}

	// Verify components
	if model.header == nil {
		t.Error("Header component is nil")
	}
	if model.messageList == nil {
		t.Error("MessageList component is nil")
	}
	if model.inputArea == nil {
		t.Error("InputArea component is nil")
	}
	if model.statusBar == nil {
		t.Error("StatusBar component is nil")
	}

	// Verify initial state
	if model.currentFocus != "input" {
		t.Errorf("Initial focus mismatch: got %s, want input", model.currentFocus)
	}
	if model.totalRequests != 10 {
		t.Errorf("TotalRequests mismatch: got %d, want 10", model.totalRequests)
	}
	if model.cacheHits != 5 {
		t.Errorf("CacheHits mismatch: got %d, want 5", model.cacheHits)
	}
}

func TestModelInit(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}
	bot := services.NewCharacterBot(cfg)

	tuiConfig := Config{
		CharacterID: "test-char",
		UserID:      "test-user",
		Bot:         bot,
		Width:       100,
		Height:      30,
	}

	model := NewModel(tuiConfig)
	cmd := model.Init()

	// Init should return a batch command
	if cmd == nil {
		t.Error("Init() returned nil command")
	}
}

func TestModelUpdate(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}
	bot := services.NewCharacterBot(cfg)

	tuiConfig := Config{
		CharacterID: "test-char",
		UserID:      "test-user",
		Bot:         bot,
		Width:       100,
		Height:      30,
	}

	model := NewModel(tuiConfig)

	// Test window size message
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, cmd := model.Update(sizeMsg)

	// Verify model is returned
	if updatedModel == nil {
		t.Error("Update returned nil model")
	}

	// Cast back to our model type
	m, ok := updatedModel.(*Model)
	if !ok {
		t.Fatal("Update returned wrong model type")
	}

	// Verify size was updated
	if m.width != 120 {
		t.Errorf("Width not updated: got %d, want 120", m.width)
	}
	if m.height != 40 {
		t.Errorf("Height not updated: got %d, want 40", m.height)
	}

	// Test ready state
	if !m.ready {
		t.Error("Model should be ready after window size message")
	}

	_ = cmd // cmd might be nil or a batch command
}

func TestCommandHistory(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}
	bot := services.NewCharacterBot(cfg)

	tuiConfig := Config{
		CharacterID: "test-char",
		UserID:      "test-user",
		Bot:         bot,
		Width:       100,
		Height:      30,
	}

	model := NewModel(tuiConfig)

	// Add some commands to history
	model.commandHistory = []string{"command 1", "command 2", "command 3"}
	model.historyIndex = len(model.commandHistory)

	// Verify initial state
	if len(model.commandHistory) != 3 {
		t.Errorf("Command history length mismatch: got %d, want 3", len(model.commandHistory))
	}
	if model.historyIndex != 3 {
		t.Errorf("History index mismatch: got %d, want 3", model.historyIndex)
	}
}

func TestCacheMetrics(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}
	bot := services.NewCharacterBot(cfg)

	tuiConfig := Config{
		CharacterID: "test-char",
		UserID:      "test-user",
		Bot:         bot,
		Width:       100,
		Height:      30,
		InitialMetrics: struct {
			TotalRequests int
			CacheHits     int
			TokensSaved   int
		}{
			TotalRequests: 100,
			CacheHits:     75,
			TokensSaved:   5000,
		},
	}

	model := NewModel(tuiConfig)

	// Verify metrics
	if model.totalRequests != 100 {
		t.Errorf("TotalRequests mismatch: got %d, want 100", model.totalRequests)
	}
	if model.cacheHits != 75 {
		t.Errorf("CacheHits mismatch: got %d, want 75", model.cacheHits)
	}

	// Calculate hit rate
	hitRate := float64(model.cacheHits) / float64(model.totalRequests) * 100
	expectedRate := 75.0
	if hitRate != expectedRate {
		t.Errorf("Hit rate mismatch: got %.1f%%, want %.1f%%", hitRate, expectedRate)
	}
}