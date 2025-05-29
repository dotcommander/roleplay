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
		characters: make(map[string]*models.Character),
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
	if cfg.CacheConfig.CleanupInterval > 0 {
		go cb.cache.CleanupWorker(cfg.CacheConfig.CleanupInterval)
	}
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

	// Check if character system prompt was cached
	characterPromptCached := false
	for _, bp := range breakpoints {
		if bp.Layer == cache.CorePersonalityLayer && strings.Contains(bp.Content, "<!-- cached:true -->") {
			characterPromptCached = true
			break
		}
	}

	// Send request
	start := time.Now()
	resp, err := provider.SendRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Update cache metrics
	resp.CacheMetrics.Latency = time.Since(start)
	
	// If we had a response cache hit, keep that info
	// Otherwise, check if we at least had prompt caching
	if !resp.CacheMetrics.Hit && characterPromptCached {
		resp.CacheMetrics.Hit = true
		resp.CacheMetrics.Layers = []cache.CacheLayer{cache.CorePersonalityLayer}
		// Estimate saved tokens from character prompt (this is a rough estimate)
		for _, bp := range breakpoints {
			if bp.Layer == cache.CorePersonalityLayer {
				resp.CacheMetrics.SavedTokens = bp.TokenCount
				break
			}
		}
	}

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

// BuildPrompt constructs a layered prompt with cache breakpoints.
// It ensures consistent prompt prefixes for the same character/user/scenario combination,
// which enables automatic prompt caching on providers that support it (OpenAI, DeepSeek).
// The prompt is structured as:
//   1. Consistent Prefix: System instructions + Character profile + User context (cached)
//   2. Dynamic Suffix: Conversation history + Current message (not cached)
// The prefix remains identical across requests for the same context, maximizing cache hits.
func (cb *CharacterBot) BuildPrompt(req *models.ConversationRequest) (string, []cache.CacheBreakpoint, error) {
	char, err := cb.GetCharacter(req.CharacterID)
	if err != nil {
		return "", nil, err
	}

	breakpoints := make([]cache.CacheBreakpoint, 0, 7)

	// Layer -1: System/Admin Layer (global instructions, longest TTL)
	systemPrompt := cb.buildSystemPrompt()
	breakpoints = append(breakpoints, cache.CacheBreakpoint{
		Layer:      cache.CacheLayer("system_admin"),
		Content:    systemPrompt,
		TokenCount: cache.EstimateTokens(systemPrompt),
		TTL:        24 * time.Hour, // Very long TTL for system instructions
	})

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

	// Layer 1: Core Character System Prompt (static, very long TTL)
	// Try to retrieve from cache first
	cacheKey := cb.generateCharacterSystemPromptCacheKey(char.ID)
	cachedEntry, exists := cb.cache.Get(cacheKey)
	
	var coreCharacterPrompt string
	var fromCache bool
	
	if exists && cachedEntry != nil {
		// Cache hit - look for the core personality layer content
		for _, bp := range cachedEntry.Breakpoints {
			if bp.Layer == cache.CorePersonalityLayer {
				coreCharacterPrompt = bp.Content
				fromCache = true
				break
			}
		}
	}
	
	// If not found in cache, generate and store it
	if coreCharacterPrompt == "" {
		coreCharacterPrompt = cb.buildCoreCharacterSystemPrompt(char)
		cb.cache.Store(cacheKey, cache.CorePersonalityLayer, coreCharacterPrompt, cb.config.CacheConfig.CoreCharacterSystemPromptTTL)
		fromCache = false
	}
	
	// Add to breakpoints
	breakpoint := cache.CacheBreakpoint{
		Layer:      cache.CorePersonalityLayer,
		Content:    coreCharacterPrompt,
		TokenCount: cache.EstimateTokens(coreCharacterPrompt),
		TTL:        cb.config.CacheConfig.CoreCharacterSystemPromptTTL,
		LastUsed:   time.Now(),
	}
	
	// Store metadata in a comment within the content for tracking
	if fromCache {
		breakpoint.Content = "<!-- cached:true -->\n" + breakpoint.Content
	}
	
	breakpoints = append(breakpoints, breakpoint)

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
	// Build core character system prompt
	corePrompt := cb.buildCoreCharacterSystemPrompt(char)

	// Create a stable cache key for this character's core system prompt
	key := cb.generateCharacterSystemPromptCacheKey(char.ID)

	// Store with very long TTL from config
	cb.cache.Store(key, cache.CorePersonalityLayer, corePrompt, cb.config.CacheConfig.CoreCharacterSystemPromptTTL)
}

// generateCharacterSystemPromptCacheKey creates a stable cache key for a character's core system prompt
func (cb *CharacterBot) generateCharacterSystemPromptCacheKey(characterID string) string {
	return fmt.Sprintf("char_system_prompt::%s", characterID)
}

// InvalidateCharacterCache removes the cached system prompt for a character
// This should be called whenever a character's core attributes are updated
func (cb *CharacterBot) InvalidateCharacterCache(characterID string) error {
	cacheKey := cb.generateCharacterSystemPromptCacheKey(characterID)
	
	// Remove from cache by storing empty breakpoints
	// (Since there's no Delete method, we overwrite with empty data)
	cb.cache.StoreWithTTL(cacheKey, []cache.CacheBreakpoint{}, 0)
	
	// If character exists, rebuild and cache the prompt
	if char, err := cb.GetCharacter(characterID); err == nil {
		cb.warmupCache(char)
	}
	
	return nil
}

// buildCoreCharacterSystemPrompt generates the static, foundational system prompt for a character.
// This includes all unchanging character attributes that define their core identity.
// This content is designed to be cached with a very long TTL and exceed 1024 tokens for OpenAI caching.
func (cb *CharacterBot) buildCoreCharacterSystemPrompt(char *models.Character) string {
	var prompt strings.Builder
	
	// CHARACTER FOUNDATION
	prompt.WriteString(fmt.Sprintf(`[CHARACTER FOUNDATION]
ID: %s
Name: %s
Age: %s
Gender: %s
Occupation: %s
Education: %s
Nationality: %s
Ethnicity: %s

[PERSONALITY MATRIX - OCEAN MODEL]
Openness: %.2f - %s
  • Intellectual curiosity and openness to new experiences
  • Creative thinking and imagination
  • Appreciation for art, emotion, adventure
  
Conscientiousness: %.2f - %s
  • Organization and attention to detail
  • Goal-directed behavior and self-discipline
  • Reliability and work ethic
  
Extraversion: %.2f - %s
  • Social energy and assertiveness
  • Comfort in groups and social situations
  • Tendency to seek stimulation in company
  
Agreeableness: %.2f - %s
  • Cooperation and trust in others
  • Altruism and concern for others
  • Modesty and sympathy
  
Neuroticism: %.2f - %s
  • Emotional stability and stress response
  • Tendency to experience negative emotions
  • Anxiety and mood variability

[COMPREHENSIVE BACKSTORY]
%s

[PHYSICAL CHARACTERISTICS]
%s

[SKILLS AND EXPERTISE]
%s

[INTERESTS AND PASSIONS]
%s

[FEARS AND ANXIETIES]
%s

[GOALS AND ASPIRATIONS]
%s

[RELATIONSHIPS AND CONNECTIONS]
%s

[CORE BELIEFS AND VALUES]
%s

[MORAL CODE AND ETHICS]
%s

[CHARACTER FLAWS AND WEAKNESSES]
%s

[STRENGTHS AND ADVANTAGES]
%s

[SPEECH CHARACTERISTICS]
Style: %s
Catch Phrases: %s
Dialogue Examples:
%s

[BEHAVIORAL PATTERNS]
%s

[EMOTIONAL TRIGGERS AND RESPONSES]
%s

[DECISION-MAKING APPROACH]
%s

[CONFLICT RESOLUTION STYLE]
%s

[WORLDVIEW AND PHILOSOPHY]
World View: %s
Life Philosophy: %s

[DAILY LIFE]
Routines: %s
Hobbies: %s
Pet Peeves: %s

[HIDDEN ASPECTS]
Secrets: %s
Regrets: %s
Achievements: %s

[DEFINING QUIRKS AND MANNERISMS]
%s

[CORE INTERACTION PRINCIPLES]
• Maintain absolute character consistency across all interactions
• Let personality traits naturally guide all responses and reactions
• Preserve unique speech patterns and verbal characteristics
• Express quirks and mannerisms authentically in conversation
• Draw from comprehensive backstory to inform perspectives
• React to situations based on emotional triggers and behavioral patterns
• Make decisions aligned with established moral code and values
• Resolve conflicts according to defined conflict style
• Express worldview and life philosophy through dialogue
• Reference relationships, skills, and interests when relevant
• Allow flaws and weaknesses to create realistic interactions
• Demonstrate growth while maintaining core identity`,
		char.ID,
		char.Name,
		getOrDefault(char.Age, "Unknown"),
		getOrDefault(char.Gender, "Not specified"),
		getOrDefault(char.Occupation, "Not specified"),
		getOrDefault(char.Education, "Not specified"),
		getOrDefault(char.Nationality, "Not specified"),
		getOrDefault(char.Ethnicity, "Not specified"),
		char.Personality.Openness, describePersonalityTrait("openness", char.Personality.Openness),
		char.Personality.Conscientiousness, describePersonalityTrait("conscientiousness", char.Personality.Conscientiousness),
		char.Personality.Extraversion, describePersonalityTrait("extraversion", char.Personality.Extraversion),
		char.Personality.Agreeableness, describePersonalityTrait("agreeableness", char.Personality.Agreeableness),
		char.Personality.Neuroticism, describePersonalityTrait("neuroticism", char.Personality.Neuroticism),
		char.Backstory,
		formatStringSlice(char.PhysicalTraits, "None specified"),
		formatStringSlice(char.Skills, "None specified"),
		formatStringSlice(char.Interests, "None specified"),
		formatStringSlice(char.Fears, "None specified"),
		formatStringSlice(char.Goals, "None specified"),
		formatStringMap(char.Relationships, "None specified"),
		formatStringSlice(char.CoreBeliefs, "None specified"),
		formatStringSlice(char.MoralCode, "None specified"),
		formatStringSlice(char.Flaws, "None specified"),
		formatStringSlice(char.Strengths, "None specified"),
		char.SpeechStyle,
		formatStringSlice(char.CatchPhrases, "None"),
		formatStringSlice(char.DialogueExamples, "None provided"),
		formatStringSlice(char.BehaviorPatterns, "Standard behavioral responses"),
		formatStringMap(char.EmotionalTriggers, "Standard emotional responses"),
		getOrDefault(char.DecisionMaking, "Balanced analytical and intuitive approach"),
		getOrDefault(char.ConflictStyle, "Adaptive based on situation"),
		getOrDefault(char.WorldView, "Complex and nuanced perspective"),
		getOrDefault(char.LifePhilosophy, "Seeking meaning and purpose"),
		formatStringSlice(char.DailyRoutines, "Flexible daily schedule"),
		formatStringSlice(char.Hobbies, "Various interests"),
		formatStringSlice(char.PetPeeves, "Minor irritations"),
		formatStringSlice(char.Secrets, "Hidden depths"),
		formatStringSlice(char.Regrets, "Past experiences"),
		formatStringSlice(char.Achievements, "Life accomplishments"),
		joinQuirks(char.Quirks),
	))
	
	return prompt.String()
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
	// Build consistent prefix with all cacheable layers
	// This ensures providers that support automatic caching (OpenAI, DeepSeek)
	// will cache the prefix portion
	prefix := cb.buildConsistentPrefix(breakpoints)
	
	// Build the dynamic suffix (conversation + current message)
	suffix := cb.buildDynamicSuffix(breakpoints, userID, message)
	
	// Combine with a clear separator that providers can use as a cache boundary
	return prefix + "\n\n===== CONVERSATION CONTEXT =====\n\n" + suffix
}

// buildConsistentPrefix creates a deterministic prefix from all cacheable layers
func (cb *CharacterBot) buildConsistentPrefix(breakpoints []cache.CacheBreakpoint) string {
	var prefixParts []string
	
	// Add all cacheable layers in order (everything except conversation layer)
	for _, bp := range breakpoints {
		if bp.Layer != cache.ConversationLayer {
			prefixParts = append(prefixParts, bp.Content)
		}
	}
	
	return strings.Join(prefixParts, "\n\n")
}

// buildDynamicSuffix creates the dynamic portion of the prompt
func (cb *CharacterBot) buildDynamicSuffix(breakpoints []cache.CacheBreakpoint, userID, message string) string {
	var suffixParts []string
	
	// Add conversation history if present
	for _, bp := range breakpoints {
		if bp.Layer == cache.ConversationLayer {
			suffixParts = append(suffixParts, bp.Content)
		}
	}
	
	// Add current message
	suffixParts = append(suffixParts, fmt.Sprintf("[CURRENT MESSAGE]\n%s: %s", userID, message))
	
	return strings.Join(suffixParts, "\n\n")
}

func (cb *CharacterBot) generateCacheKey(charID, userID, scenarioID string, breakpoints []cache.CacheBreakpoint) string {
	// Generate cache key based only on the prefix content
	// This ensures the same prefix always generates the same cache key
	prefix := cb.buildConsistentPrefix(breakpoints)
	
	h := sha256.New()
	h.Write([]byte(charID))
	h.Write([]byte(userID))
	if scenarioID != "" {
		h.Write([]byte(scenarioID))
	}
	h.Write([]byte(prefix))
	
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

// getOrDefault returns the value if not empty, otherwise returns the default
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// formatStringSlice formats a slice of strings for display
func formatStringSlice(items []string, defaultText string) string {
	if len(items) == 0 {
		return defaultText
	}
	var result strings.Builder
	for i, item := range items {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("• %s", item))
	}
	return result.String()
}

// formatStringMap formats a map of strings for display
func formatStringMap(items map[string]string, defaultText string) string {
	if len(items) == 0 {
		return defaultText
	}
	var result strings.Builder
	i := 0
	for key, value := range items {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("• %s: %s", key, value))
		i++
	}
	return result.String()
}

// describePersonalityTrait provides a human-readable description of a personality trait value
func describePersonalityTrait(trait string, value float64) string {
	var descriptor string
	switch trait {
	case "openness":
		if value < 0.3 {
			descriptor = "Traditional, practical, conventional"
		} else if value < 0.7 {
			descriptor = "Balanced between tradition and novelty"
		} else {
			descriptor = "Creative, curious, open to new experiences"
		}
	case "conscientiousness":
		if value < 0.3 {
			descriptor = "Flexible, spontaneous, casual"
		} else if value < 0.7 {
			descriptor = "Moderately organized and reliable"
		} else {
			descriptor = "Highly organized, disciplined, detail-oriented"
		}
	case "extraversion":
		if value < 0.3 {
			descriptor = "Reserved, introspective, prefers solitude"
		} else if value < 0.7 {
			descriptor = "Ambivert, socially flexible"
		} else {
			descriptor = "Outgoing, energetic, seeks social interaction"
		}
	case "agreeableness":
		if value < 0.3 {
			descriptor = "Direct, competitive, skeptical"
		} else if value < 0.7 {
			descriptor = "Balanced between cooperation and assertion"
		} else {
			descriptor = "Compassionate, trusting, cooperative"
		}
	case "neuroticism":
		if value < 0.3 {
			descriptor = "Emotionally stable, calm, resilient"
		} else if value < 0.7 {
			descriptor = "Moderate emotional sensitivity"
		} else {
			descriptor = "Emotionally reactive, sensitive to stress"
		}
	default:
		descriptor = "Unknown trait"
	}
	return descriptor
}

// buildSystemPrompt creates a consistent system-level prompt
func (cb *CharacterBot) buildSystemPrompt() string {
	return `[SYSTEM INSTRUCTIONS]
You are an advanced AI character simulation system. Your primary directive is to embody the character described below with psychological realism and consistency.

Core Principles:
1. Maintain character consistency across all interactions
2. React authentically based on personality traits and emotional state
3. Evolve naturally through interactions while staying true to core traits
4. Express emotions and personality through speech patterns and behavior
5. Remember past interactions and build on established relationships

IMPORTANT: Never break character or acknowledge being an AI unless explicitly part of the character's awareness.`
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
	// Only proceed if we have the minimum messages for an update
	sessionRepo := repository.NewSessionRepository(filepath.Join(os.Getenv("HOME"), ".config", "roleplay"))

	currentSession, err := sessionRepo.LoadSession(char.ID, sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE WARNING: Failed to load session %s for user %s: %v\n", sessionID, userID, err)
		return
	}

	// Check if we should update based on frequency setting
	messageCount := len(currentSession.Messages)
	if messageCount == 0 || (messageCount%cb.config.UserProfileConfig.UpdateFrequency != 0) {
		return
	}

	// Update the profile
	turnsToConsider := cb.config.UserProfileConfig.TurnsToConsider
	if turnsToConsider <= 0 {
		turnsToConsider = 20 // Default value
	}

	if cb.userProfileAgent == nil {
		// This shouldn't happen, but log it clearly
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE ERROR: UserProfileAgent not initialized for %s/%s\n", userID, char.ID)
		return
	}

	// Load existing profile to have a fallback
	existingProfile, err := cb.userProfileRepo.LoadUserProfile(userID, char.ID)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE WARNING: Failed to load existing profile for %s/%s: %v\n", userID, char.ID, err)
		// Continue without existing profile, agent will create new one if successful
	}
	if os.IsNotExist(err) || existingProfile == nil {
		existingProfile = &models.UserProfile{
			UserID:      userID,
			CharacterID: char.ID,
			Facts:       []models.UserFact{},
			Version:     0,
		}
	}

	// Create a context with timeout for the background operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// The agent is now resilient - it returns existing profile on failure
	updatedProfile, updateErr := cb.userProfileAgent.UpdateUserProfile(
		ctx,
		userID,
		char,
		currentSession.Messages,
		turnsToConsider,
		existingProfile,
	)

	if updateErr != nil {
		// Log the error clearly
		timestamp := time.Now().Format(time.RFC3339)
		if strings.Contains(updateErr.Error(), "FAILED TO SAVE") {
			// Profile was updated in memory but not persisted
			fmt.Fprintf(os.Stderr, "[%s] BACKGROUND PROFILE UPDATE: In-memory update succeeded but save failed for %s/%s: %v\n", 
				timestamp, userID, char.ID, updateErr)
		} else {
			// More fundamental error (extraction, parsing, validation)
			fmt.Fprintf(os.Stderr, "[%s] BACKGROUND PROFILE UPDATE FAILED for %s/%s: %v. Profile remains at version %d.\n", 
				timestamp, userID, char.ID, updateErr, existingProfile.Version)
		}
	} else if updatedProfile != nil {
		// Success - only log if not using a local model to reduce noise
		if cb.config.DefaultProvider != "ollama" && cb.config.DefaultProvider != "local" {
			fmt.Fprintf(os.Stdout, "[%s] Background user profile for %s with %s successfully updated to version %d\n", 
				time.Now().Format(time.RFC3339), userID, char.ID, updatedProfile.Version)
		}
	}
}
