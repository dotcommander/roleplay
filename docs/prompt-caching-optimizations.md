# Prompt Caching Optimizations

## Implemented Features (v0.8.3)

### 1. User Parameter for Cache Routing ✅
- Added `User` parameter to all OpenAI API requests
- Improves cache hit rates by routing requests with same user-character pair to same servers
- Automatically included in both streaming and non-streaming requests

### 2. Request Rate Limiting ✅
- Implemented per user-character pair rate limiting (14 requests/minute)
- Prevents cache overflow (OpenAI limit is 15/min)
- Graceful error messages with current rate information
- Background cleanup of stale rate limit buckets
- Rate limiter statistics available via `GetRateLimiterStats()`

### 3. Enhanced Cache Metrics ✅
- Added OpenAI `cached_tokens` tracking from API responses
- Session statistics now show:
  - Response Cache Hit Rate (application-level)
  - OpenAI Prompt Cache Hit Rate (provider-level)
  - Total cached tokens served by OpenAI
- Updated demo mode to display both cache types
- Session stats command shows comprehensive metrics

### 4. Provider-Level Cache Detection ✅
- OpenAI provider now extracts `cached_tokens` from `prompt_tokens_details`
- Automatically detects which cache layers were hit based on token count
- Supports 128-token increment detection (1024, 1152, 1280, etc.)

## Architecture

### Rate Limiter
```go
// Located in: internal/services/rate_limiter.go
type RateLimiter struct {
    buckets   map[string]*bucket  // key: "userID:characterID"
    maxRate   int                 // 14 requests/minute
    window    time.Duration       // 1 minute
}
```

### Enhanced Session Metrics
```go
// Updated in: internal/repository/session_repo.go
type CacheMetrics struct {
    TotalRequests       int
    CacheHits           int     // Response cache hits
    CacheMisses         int
    TokensSaved         int
    CachedTokensTotal   int     // OpenAI cached tokens
    PromptCacheHitRate  float64 // OpenAI cache rate
    CostSaved           float64
    HitRate             float64
}
```

### Provider Updates
```go
// Updated in: internal/providers/openai.go
apiReq := openai.ChatCompletionRequest{
    Model:    o.model,
    Messages: messages,
    User:     req.UserID, // Now included for cache routing
}
```

## Usage

### Rate Limiting
The rate limiter automatically prevents excessive requests:
```
Error: rate limit exceeded: 14/14 requests per minute for this user-character pair. Please wait before sending more messages
```

### Cache Metrics
View comprehensive cache statistics:
```bash
roleplay session stats

Cache Performance Statistics
===========================

Rick Sanchez (rick-sanchez):
  Sessions: 3
  Total Requests: 45
  Response Cache Hit Rate: 78.2%
  OpenAI Prompt Cache Rate: 92.5%
  Cached Tokens (OpenAI): 84,320
  Tokens Saved: 12,450
  Cost Saved: $1.24
```

## Benefits

1. **Cost Reduction**: Up to 75% reduction in API costs through optimal caching
2. **Latency Improvement**: Up to 80% faster responses for cached prompts
3. **Resilience**: Rate limiting prevents service disruption
4. **Visibility**: Comprehensive metrics for optimization

## Future Optimizations

### High Priority
- Prompt prefix validation to ensure exact matches
- Prompt structure versioning for cache invalidation

### Medium Priority  
- Align prompt sections to 128-token boundaries
- Prompt warming strategy for popular characters
- Cache overflow handling configuration

### Low Priority
- Document best practices in CLAUDE.md
- Advanced cache analytics dashboard