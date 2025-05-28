package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
)

// CharacterBot is the main service for managing characters and conversations
type CharacterBot struct {
	characters       map[string]*models.Character
	cache            *cache.PromptCache
	responseCache    *cache.ResponseCache
	providers        map[string]providers.AIProvider
	config           *config.Config
	scenarioRepo     *repository.ScenarioRepository
	userProfileRepo  *repository.UserProfileRepository
	userProfileAgent *UserProfileAgent
	mu               sync.RWMutex
	cacheHits        int
	cacheMisses      int
}

// NewCharacterBot creates a new character bot instance
func NewCharacterBot(cfg *config.Config) *CharacterBot {
	// Get config path for scenario repository
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "roleplay")
	userProfileDataDir := filepath.Join(configPath, "user_profiles")

	// Create user profiles directory if it doesn't exist
	if err := os.MkdirAll(userProfileDataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create user_profiles directory: %v\n", err)
	}

	userProfileRepo := repository.NewUserProfileRepository(userProfileDataDir)

	cb := &CharacterBot{
		characters:      make(map[string]*models.Character),
		cache: cache.NewPromptCache(
			cfg.CacheConfig.DefaultTTL,
			5*time.Minute,
			1*time.Hour,
		),
		responseCache:   cache.NewResponseCache(cfg.CacheConfig.DefaultTTL),
		providers:       make(map[string]providers.AIProvider),
		config:          cfg,
		scenarioRepo:    repository.NewScenarioRepository(configPath),
		userProfileRepo: userProfileRepo,
		cacheHits:       0,
		cacheMisses:     0,
	}

	// Start background workers
	go cb.cache.CleanupWorker(cfg.CacheConfig.CleanupInterval)
	go cb.memoryConsolidationWorker()

	return cb
}

// InitializeUserProfileAgent initializes the user profile agent with a provider
func (cb *CharacterBot) InitializeUserProfileAgent() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if !cb.config.UserProfileConfig.Enabled {
		return
	}

	// Get the default provider for the UserProfileAgent
	if provider, ok := cb.providers[cb.config.DefaultProvider]; ok {
		cb.userProfileAgent = NewUserProfileAgent(provider, cb.userProfileRepo)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: Default provider %s not found for UserProfileAgent\n", cb.config.DefaultProvider)
	}
}

// RegisterProvider adds a new AI provider
func (cb *CharacterBot) RegisterProvider(name string, provider providers.AIProvider) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.providers[name] = provider
}

// CreateCharacter adds a new character to the bot
func (cb *CharacterBot) CreateCharacter(char *models.Character) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if _, exists := cb.characters[char.ID]; exists {
		return fmt.Errorf("character %s already exists", char.ID)
	}

	char.LastModified = time.Now()
	cb.characters[char.ID] = char

	// Pre-cache core personality
	cb.warmupCache(char)

	return nil
}

// GetCharacter retrieves a character by ID
func (cb *CharacterBot) GetCharacter(id string) (*models.Character, error) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	char, exists := cb.characters[id]
	if !exists {
		return nil, fmt.Errorf("character %s not found", id)
	}

	return char, nil
}

// ProcessRequest handles a conversation request
func (cb *CharacterBot) ProcessRequest(ctx context.Context, req *models.ConversationRequest) (*providers.AIResponse, error) {
	// Check response cache first
	responseCacheKey := cb.responseCache.GenerateKey(req.CharacterID, req.UserID, req.Message)
	if cachedResp, found := cb.responseCache.Get(responseCacheKey); found {
		cb.mu.Lock()
		cb.cacheHits++
		cb.mu.Unlock()

		// Return cached response with cache hit metrics
		return &providers.AIResponse{
			Content: cachedResp.Content,
			TokensUsed: providers.TokenUsage{
				Prompt:       0,
				Completion:   0,
				CachedPrompt: cachedResp.TokensUsed.Prompt,
				Total:        0,
			},
			CacheMetrics: cache.CacheMetrics{
				Hit:         true,
				Layers:      []cache.CacheLayer{cache.ConversationLayer},
				SavedTokens: cachedResp.TokensUsed.Total,
				Latency:     time.Since(cachedResp.CachedAt),
			},
		}, nil
	}

	cb.mu.Lock()
	cb.cacheMisses++
	cb.mu.Unlock()

	// Build prompt with cache awareness
	prompt, breakpoints, err := cb.BuildPrompt(req)
	if err != nil {
		return nil, err
	}

	// Generate cache key for static layers only (including scenario if present)
	cacheKey := cb.generateCacheKey(req.CharacterID, req.UserID, req.ScenarioID, breakpoints)
	cachedEntry, hit := cb.cache.Get(cacheKey)

	// Get character for complexity check
	char, err := cb.GetCharacter(req.CharacterID)
	if err != nil {
		return nil, err
	}

	// Cache hit tracking is now done internally

	// Adaptive TTL based on conversation activity
	effectiveTTL := cb.cache.CalculateAdaptiveTTL(cachedEntry, len(char.Memories) > 50)

	// Select provider
	provider := cb.selectProvider()
	if provider == nil {
		return nil, fmt.Errorf("no AI provider available")
	}

	// Prepare API request
	apiReq := &providers.PromptRequest{
		CharacterID:      req.CharacterID,
		UserID:           req.UserID,
		Message:          req.Message,
		Context:          req.Context,
		SystemPrompt:     prompt,
		CacheBreakpoints: breakpoints,
	}

	// Send request
	start := time.Now()
	resp, err := provider.SendRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Update cache metrics
	resp.CacheMetrics.Latency = time.Since(start)

	// Update character state based on response
	cb.updateCharacterState(req.CharacterID, resp)

	// Store in cache with adaptive TTL
	if !hit {
		cb.cache.StoreWithTTL(cacheKey, breakpoints, effectiveTTL)
	}

	// Store response in response cache
	cb.responseCache.Store(responseCacheKey, resp.Content, cache.TokenUsage{
		Prompt:       resp.TokensUsed.Prompt,
		Completion:   resp.TokensUsed.Completion,
		CachedPrompt: resp.TokensUsed.CachedPrompt,
		Total:        resp.TokensUsed.Total,
	})

	// Trigger user profile update asynchronously if enabled
	if cb.userProfileAgent != nil && cb.config.UserProfileConfig.Enabled {
		go cb.updateUserProfileAsync(req.UserID, char, req.Context.SessionID)
	}

	return resp, nil
}

// BuildPrompt constructs a layered prompt with cache breakpoints
func (cb *CharacterBot) BuildPrompt(req *models.ConversationRequest) (string, []cache.CacheBreakpoint, error) {
	char, err := cb.GetCharacter(req.CharacterID)
	if err != nil {
		return "", nil, err
	}

	breakpoints := make([]cache.CacheBreakpoint, 0, 6)

	// Layer 0: Scenario Context (highest layer, meta-prompts, longest TTL)
	if req.ScenarioID != "" {
		scenario, err := cb.scenarioRepo.LoadScenario(req.ScenarioID)
		if err != nil {
			// Log warning but continue without scenario
			fmt.Fprintf(os.Stderr, "Warning: Failed to load scenario %s: %v\n", req.ScenarioID, err)
		} else if scenario.Prompt != "" {
			// Very long TTL for scenario context (7 days by default)
			scenarioTTL := 168 * time.Hour

			breakpoints = append(breakpoints, cache.CacheBreakpoint{
				Layer:      cache.ScenarioContextLayer,
				Content:    scenario.Prompt,
				TokenCount: cache.EstimateTokens(scenario.Prompt),
				TTL:        scenarioTTL,
				LastUsed:   time.Now(),
			})

			// Update scenario last used timestamp asynchronously
			go func(id string) {
				_ = cb.scenarioRepo.UpdateScenarioLastUsed(id)
			}(req.ScenarioID)
		}
	}

	// Layer 1: Core Personality (static, long TTL)
	personality := cb.buildPersonalityPrompt(char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.CorePersonalityLayer,
		Content:    personality,
		TokenCount: cache.EstimateTokens(personality),
		TTL:        cb.cache.CalculateAdaptiveTTL(nil, true),
	})

	// Layer 2: Learned Behaviors (semi-static, medium TTL)
	behaviors := cb.buildLearnedBehaviors(char)
	if behaviors != "" {
		breakpoints = append(breakpoints, cache.CacheBreakpoint{
			Layer:      cache.LearnedBehaviorLayer,
			Content:    behaviors,
			TokenCount: cache.EstimateTokens(behaviors),
			TTL:        cb.config.CacheConfig.DefaultTTL * 2,
		})
	}

	// Layer 3: Emotional State (dynamic, short TTL)
	emotional := cb.buildEmotionalContext(char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.EmotionalStateLayer,
		Content:    emotional,
		TokenCount: cache.EstimateTokens(emotional),
		TTL:        5 * time.Minute,
	})

	// Layer 4: User Context (semi-dynamic, medium TTL)
	userContext := cb.buildUserContext(req.UserID, char)
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.UserMemoryLayer,
		Content:    userContext,
		TokenCount: cache.EstimateTokens(userContext),
		TTL:        cb.config.CacheConfig.DefaultTTL,
	})

	// Layer 5: Conversation History (dynamic, no cache)
	conversation := cb.buildConversationHistory(req.Context)
	if conversation != "" {
		breakpoints = append(breakpoints, cache.CacheBreakpoint{
			Layer:      cache.ConversationLayer,
			Content:    conversation,
			TokenCount: cache.EstimateTokens(conversation),
			TTL:        0, // No caching for conversation
		})
	}

	// Combine all layers
	fullPrompt := cb.assemblePrompt(breakpoints, req.UserID, req.Message)

	return fullPrompt, breakpoints, nil
}

func (cb *CharacterBot) warmupCache(char *models.Character) {
	// Build personality prompt and create a cache key for this character
	personality := cb.buildPersonalityPrompt(char)

	// Create a stable cache key for just the personality layer
	h := sha256.New()
	h.Write([]byte(char.ID))
	h.Write([]byte("personality"))
	h.Write([]byte(personality))
	key := hex.EncodeToString(h.Sum(nil))

	// Store with a long TTL since personality is static
	cb.cache.Store(key, cache.CorePersonalityLayer, personality, 24*time.Hour)
}

func (cb *CharacterBot) buildPersonalityPrompt(char *models.Character) string {
	return fmt.Sprintf(`[CHARACTER PROFILE]
Name: %s
Personality Traits:
- Openness: %.2f
- Conscientiousness: %.2f
- Extraversion: %.2f
- Agreeableness: %.2f
- Neuroticism: %.2f

Backstory: %s

Speech Style: %s

Core Quirks: %s

[INTERACTION RULES]
- Always stay in character
- Express personality traits consistently
- Use characteristic speech patterns
- React based on emotional state`,
		char.Name,
		char.Personality.Openness,
		char.Personality.Conscientiousness,
		char.Personality.Extraversion,
		char.Personality.Agreeableness,
		char.Personality.Neuroticism,
		char.Backstory,
		char.SpeechStyle,
		joinQuirks(char.Quirks),
	)
}

func (cb *CharacterBot) buildLearnedBehaviors(char *models.Character) string {
	// Extract patterns from medium-term memories
	patterns := make([]string, 0)
	for _, mem := range char.Memories {
		if mem.Type == models.MediumTermMemory {
			patterns = append(patterns, mem.Content)
		}
	}

	if len(patterns) == 0 {
		return ""
	}

	return fmt.Sprintf("[LEARNED PATTERNS]\n%s", strings.Join(patterns, "\n"))
}

func (cb *CharacterBot) buildEmotionalContext(char *models.Character) string {
	return fmt.Sprintf(`[EMOTIONAL STATE]
Current Mood:
- Joy: %.2f
- Surprise: %.2f
- Anger: %.2f
- Fear: %.2f
- Sadness: %.2f
- Disgust: %.2f`,
		char.CurrentMood.Joy,
		char.CurrentMood.Surprise,
		char.CurrentMood.Anger,
		char.CurrentMood.Fear,
		char.CurrentMood.Sadness,
		char.CurrentMood.Disgust,
	)
}

func (cb *CharacterBot) buildConversationHistory(ctx models.ConversationContext) string {
	if len(ctx.RecentMessages) == 0 {
		return ""
	}

	history := "[CONVERSATION HISTORY]\n"
	for _, msg := range ctx.RecentMessages {
		history += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	return history
}

func (cb *CharacterBot) buildUserContext(userID string, char *models.Character) string {
	// Try to load user profile if available
	if cb.userProfileRepo != nil && cb.config.UserProfileConfig.Enabled {
		profile, err := cb.userProfileRepo.LoadUserProfile(userID, char.ID)
		if err == nil && profile != nil {
			return cb.buildUserProfileLayer(userID, char.ID, profile)
		}
	}

	// Fallback to basic user context
	context := fmt.Sprintf(`[USER CONTEXT]
You are speaking with: %s

Remember to address them by their name throughout the conversation.`, userID)

	// Add any user-specific memories
	var userMemories []string
	for _, mem := range char.Memories {
		if mem.Type == models.LongTermMemory && mem.Content != "" {
			// Check if memory mentions this user (simple check)
			// In a more advanced system, you'd have user-specific memory storage
			userMemories = append(userMemories, mem.Content)
		}
	}

	if len(userMemories) > 0 {
		context += "\n\nShared experiences:\n" + strings.Join(userMemories, "\n")
	}

	return context
}

func (cb *CharacterBot) buildUserProfileLayer(userID, characterID string, profile *models.UserProfile) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[USER PROFILE FOR %s (as perceived by character %s)]\n", userID, characterID))
	
	if profile.OverallSummary != "" {
		sb.WriteString(fmt.Sprintf("Summary: %s\n", profile.OverallSummary))
	}
	
	if profile.InteractionStyle != "" {
		sb.WriteString(fmt.Sprintf("Interaction Style: %s\n", profile.InteractionStyle))
	}

	if len(profile.Facts) > 0 {
		sb.WriteString("\nKey Facts Remembered About User:\n")
		for _, fact := range profile.Facts {
			// Only include facts with confidence above threshold
			if fact.Confidence >= cb.config.UserProfileConfig.ConfidenceThreshold {
				sb.WriteString(fmt.Sprintf("- %s: %s (Confidence: %.1f)\n", fact.Key, fact.Value, fact.Confidence))
			}
		}
	}
	
	return sb.String()
}

func (cb *CharacterBot) assemblePrompt(breakpoints []cache.CacheBreakpoint, userID, message string) string {
	var parts []string

	// Add all breakpoint content
	for _, bp := range breakpoints {
		parts = append(parts, bp.Content)
	}

	// Add current message
	parts = append(parts, fmt.Sprintf("[CURRENT MESSAGE]\n%s: %s", userID, message))

	return strings.Join(parts, "\n\n")
}

func (cb *CharacterBot) generateCacheKey(charID, userID, scenarioID string, breakpoints []cache.CacheBreakpoint) string {
	// Generate cache key based only on static/semi-static layers
	// Don't include conversation layer which changes every time
	h := sha256.New()
	h.Write([]byte(charID))
	h.Write([]byte(userID))

	// Include scenario ID in the cache key if present
	// This ensures different scenarios create different cache entries
	if scenarioID != "" {
		h.Write([]byte(scenarioID))
	}

	// Only hash content from cacheable layers (not conversation)
	// Note: If scenario content is the first breakpoint, it's already included
	for _, bp := range breakpoints {
		if bp.Layer != cache.ConversationLayer {
			h.Write([]byte(bp.Content))
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (cb *CharacterBot) selectProvider() providers.AIProvider {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Try to get the default provider
	if provider, exists := cb.providers[cb.config.DefaultProvider]; exists {
		return provider
	}

	// Fallback to first available provider
	for _, p := range cb.providers {
		return p
	}

	return nil
}

func (cb *CharacterBot) updateCharacterState(charID string, resp *providers.AIResponse) {
	char, err := cb.GetCharacter(charID)
	if err != nil {
		return
	}

	char.Lock()
	defer char.Unlock()

	// Update emotional state with decay
	char.CurrentMood = cb.blendEmotions(char.CurrentMood, resp.Emotions, 0.3)

	// Add to short-term memory
	memory := models.Memory{
		Type:      models.ShortTermMemory,
		Content:   resp.Content,
		Timestamp: time.Now(),
		Emotional: cb.calculateEmotionalWeight(resp.Emotions),
	}
	char.Memories = append(char.Memories, memory)

	// Trigger consolidation if needed
	if len(char.Memories) > cb.config.MemoryConfig.ShortTermWindow {
		go cb.consolidateMemories(char)
	}

	// Evolution logic
	if cb.config.PersonalityConfig.EvolutionEnabled {
		cb.evolvePersonality(char, resp)
	}

	char.LastModified = time.Now()
}

func (cb *CharacterBot) blendEmotions(current, new models.EmotionalState, rate float64) models.EmotionalState {
	return models.EmotionalState{
		Joy:      current.Joy*(1-rate) + new.Joy*rate,
		Surprise: current.Surprise*(1-rate) + new.Surprise*rate,
		Anger:    current.Anger*(1-rate) + new.Anger*rate,
		Fear:     current.Fear*(1-rate) + new.Fear*rate,
		Sadness:  current.Sadness*(1-rate) + new.Sadness*rate,
		Disgust:  current.Disgust*(1-rate) + new.Disgust*rate,
	}
}

func (cb *CharacterBot) calculateEmotionalWeight(emotions models.EmotionalState) float64 {
	// Simple average of emotion intensities
	total := emotions.Joy + emotions.Surprise + emotions.Anger +
		emotions.Fear + emotions.Sadness + emotions.Disgust
	return total / 6.0
}

func (cb *CharacterBot) evolvePersonality(char *models.Character, resp *providers.AIResponse) {
	// Calculate trait impacts based on interaction
	impacts := cb.analyzeInteractionImpacts(resp)

	// Apply bounded evolution
	driftRate := cb.config.PersonalityConfig.MaxDriftRate
	char.Personality.Openness += impacts.Openness * driftRate
	char.Personality.Conscientiousness += impacts.Conscientiousness * driftRate
	char.Personality.Extraversion += impacts.Extraversion * driftRate
	char.Personality.Agreeableness += impacts.Agreeableness * driftRate
	char.Personality.Neuroticism += impacts.Neuroticism * driftRate

	// Normalize to keep traits in [0, 1] range
	char.Personality = models.NormalizePersonality(char.Personality)
}

func (cb *CharacterBot) analyzeInteractionImpacts(resp *providers.AIResponse) models.PersonalityTraits {
	// Simplified impact analysis based on emotional response
	return models.PersonalityTraits{
		Openness:          resp.Emotions.Surprise * 0.5,
		Conscientiousness: (1 - resp.Emotions.Anger) * 0.3,
		Extraversion:      resp.Emotions.Joy * 0.4,
		Agreeableness:     (1 - resp.Emotions.Disgust) * 0.3,
		Neuroticism:       (resp.Emotions.Fear + resp.Emotions.Sadness) * 0.3,
	}
}

// Background workers

func (cb *CharacterBot) memoryConsolidationWorker() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cb.consolidateAllMemories()
	}
}

func (cb *CharacterBot) consolidateAllMemories() {
	cb.mu.RLock()
	chars := make([]*models.Character, 0, len(cb.characters))
	for _, char := range cb.characters {
		chars = append(chars, char)
	}
	cb.mu.RUnlock()

	for _, char := range chars {
		cb.consolidateMemories(char)
	}
}

func (cb *CharacterBot) consolidateMemories(char *models.Character) {
	char.Lock()
	defer char.Unlock()

	// Group memories by emotional significance
	emotionalMemories := make([]models.Memory, 0)

	threshold := 0.7 // Emotional weight threshold

	for _, mem := range char.Memories {
		if mem.Type == models.ShortTermMemory && mem.Emotional > threshold {
			emotionalMemories = append(emotionalMemories, mem)
		}
	}

	// Consolidate emotional memories into medium-term
	if len(emotionalMemories) > 3 {
		consolidated := models.Memory{
			Type:      models.MediumTermMemory,
			Content:   cb.synthesizeMemories(emotionalMemories),
			Timestamp: time.Now(),
			Emotional: cb.averageEmotionalWeight(emotionalMemories),
		}
		char.Memories = append(char.Memories, consolidated)
	}

	// Prune old short-term memories
	cutoff := time.Now().Add(-cb.config.MemoryConfig.MediumTermDuration)
	filtered := make([]models.Memory, 0)

	for _, mem := range char.Memories {
		if mem.Type != models.ShortTermMemory || mem.Timestamp.After(cutoff) {
			filtered = append(filtered, mem)
		}
	}

	char.Memories = filtered
}

func (cb *CharacterBot) synthesizeMemories(memories []models.Memory) string {
	// In a real implementation, this would use NLP to create a coherent summary
	contents := make([]string, 0, len(memories))
	for _, mem := range memories {
		contents = append(contents, mem.Content)
	}
	return fmt.Sprintf("Consolidated memories: %s", strings.Join(contents, "; "))
}

func (cb *CharacterBot) averageEmotionalWeight(memories []models.Memory) float64 {
	if len(memories) == 0 {
		return 0
	}

	total := 0.0
	for _, mem := range memories {
		total += mem.Emotional
	}

	return total / float64(len(memories))
}

// Utility functions

func joinQuirks(quirks []string) string {
	if len(quirks) == 0 {
		return "None"
	}
	return strings.Join(quirks, ", ")
}

// UpdateUserProfile synchronously updates the user profile
func (cb *CharacterBot) UpdateUserProfile(userID string, char *models.Character, sessionID string) {
	if cb.userProfileAgent == nil {
		return
	}
	cb.updateUserProfileSync(userID, char, sessionID)
}

// updateUserProfileAsync asynchronously updates the user profile based on conversation history
func (cb *CharacterBot) updateUserProfileAsync(userID string, char *models.Character, sessionID string) {
	go cb.updateUserProfileSync(userID, char, sessionID)
}

// updateUserProfileSync performs the actual user profile update
func (cb *CharacterBot) updateUserProfileSync(userID string, char *models.Character, sessionID string) {
	// Update user profile based on conversation history
	
	// Only proceed if we have the minimum messages for an update
	sessionRepo := repository.NewSessionRepository(filepath.Join(os.Getenv("HOME"), ".config", "roleplay"))
	
	currentSession, err := sessionRepo.LoadSession(char.ID, sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session %s for user profile update: %v\n", sessionID, err)
		return
	}
	
	// Check if we should update based on frequency setting
	messageCount := len(currentSession.Messages)
	if messageCount == 0 || (messageCount % cb.config.UserProfileConfig.UpdateFrequency != 0) {
		return
	}
	
	// Update the profile
	turnsToConsider := cb.config.UserProfileConfig.TurnsToConsider
	if turnsToConsider <= 0 {
		turnsToConsider = 20 // Default value
	}
	
	if cb.userProfileAgent == nil {
		return
	}
	
	_, updateErr := cb.userProfileAgent.UpdateUserProfile(
		context.Background(),
		userID,
		char,
		currentSession.Messages,
		turnsToConsider,
	)
	
	if updateErr != nil {
		fmt.Fprintf(os.Stderr, "Error updating user profile for %s with %s: %v\n", userID, char.ID, updateErr)
	}
}
