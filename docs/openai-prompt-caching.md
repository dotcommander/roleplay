# OpenAI Prompt Caching Documentation

## Overview

Prompt caching reduces latency and cost by routing API requests to servers that recently processed the same prompt. This can reduce latency by up to 80% and cost by up to 75%. Prompt Caching works automatically on all API requests (no code changes required) and has no additional fees. It's enabled for all recent models (gpt-4o and newer).

## How It Works

Caching is enabled automatically for prompts that are 1024 tokens or longer. The process follows these steps:

### 1. Cache Routing
- Requests are routed to a machine based on a hash of the initial prefix of the prompt
- The hash typically uses the first 256 tokens (varies by model)
- The `user` parameter can be combined with the prefix hash to influence routing and improve cache hit rates
- If requests for the same prefix/user combination exceed ~15 requests per minute, some may overflow to additional machines, reducing cache effectiveness

### 2. Cache Lookup
- The system checks if the initial portion (prefix) of your prompt exists in the cache on the selected machine

### 3. Cache Hit/Miss
- **Cache Hit**: If a matching prefix is found, the system uses the cached result (significantly decreases latency and reduces costs)
- **Cache Miss**: If no matching prefix is found, the system processes your full prompt and caches the prefix for future requests

### Cache Persistence
- Cached prefixes generally remain active for 5-10 minutes of inactivity
- During off-peak periods, caches may persist for up to one hour

## Requirements

### Token Requirements
- Minimum prompt length: **1024 tokens**
- Cache hits occur in increments of **128 tokens**
- Cached token sequence: 1024, 1152, 1280, 1408, etc.

### API Response
All requests include a `cached_tokens` field in the response:

```json
"usage": {
  "prompt_tokens": 2006,
  "completion_tokens": 300,
  "total_tokens": 2306,
  "prompt_tokens_details": {
    "cached_tokens": 1920
  },
  "completion_tokens_details": {
    "reasoning_tokens": 0,
    "accepted_prediction_tokens": 0,
    "rejected_prediction_tokens": 0
  }
}
```

## What Can Be Cached

1. **Messages**: The complete messages array (system, user, and assistant interactions)
2. **Images**: Images in user messages (as links or base64-encoded data)
   - The `detail` parameter must be set identically as it impacts tokenization
3. **Tool use**: Both the messages array and the list of available `tools`
4. **Structured outputs**: The structured output schema (serves as a prefix to the system message)

## Structuring Prompts for Optimal Caching

Cache hits are only possible for **exact prefix matches**. To maximize caching benefits:

1. Place **static content** (instructions, examples) at the beginning of your prompt
2. Put **variable content** (user-specific information) at the end
3. Keep images and tools identical between requests

## Best Practices

1. **Structure prompts** with static/repeated content at the beginning and dynamic content at the end
2. **Use the `user` parameter** consistently across requests with common prefixes
   - Choose a granularity that keeps each unique prefix-user combination below 15 requests per minute
3. **Monitor cache performance metrics**:
   - Cache hit rates
   - Latency
   - Proportion of tokens cached
4. **Maintain steady request streams** with identical prompt prefixes to minimize cache evictions

## Key Points

### Privacy
- Prompt caches are not shared between organizations
- Only members of the same organization can access caches of identical prompts

### Output Generation
- Prompt Caching does not affect output token generation or the final response
- Only the prompt is cached; the response is computed anew each time

### Cost
- No extra cost for using prompt caching
- Automatic feature with no explicit action needed

### Limitations
- Manual cache clearing is not available
- Cached prompts contribute to TPM rate limits
- Discounting available on Scale Tier but not on Batch API
- Compatible with Zero Data Retention policies

## Optimization Tips for Roleplay Project

Based on this documentation, here are specific optimizations for the roleplay project:

1. **Ensure prompts exceed 1024 tokens**: The expanded character model (1940+ tokens) already meets this requirement
2. **Use consistent `user` parameter**: Pass the user ID to improve cache routing for user-specific conversations
3. **Structure prompts correctly**:
   - System instructions (static) → Character definition (static) → User profile (semi-static) → Conversation history (dynamic)
4. **Monitor request patterns**: Keep requests below 15/minute per unique prefix-user combination
5. **Leverage the 128-token increment**: Structure prompt sections to align with these boundaries for optimal caching