# Gemini URL Debugging Summary

## Issue
When using the Google Gemini API through the OpenAI-compatible endpoint, we needed to understand how the go-openai SDK constructs URLs.

## Solution
The go-openai SDK correctly appends `/chat/completions` to the base URL. For Gemini:

- Base URL: `https://generativelanguage.googleapis.com/v1beta/openai`
- Constructed URL: `https://generativelanguage.googleapis.com/v1beta/openai/chat/completions`

This matches the expected Gemini endpoint format.

## Debug Features Added
Added optional HTTP debugging to the OpenAI provider:

1. Set environment variable `DEBUG_HTTP=true` to enable debug logging
2. Shows the exact URLs being requested
3. Shows response status codes
4. Minimal overhead when disabled

## Usage Example
```bash
# Enable debug logging
export DEBUG_HTTP=true

# Run with Gemini
roleplay chat "Hello" --provider gemini --model gemini-2.0-flash-exp

# Disable debug logging
unset DEBUG_HTTP
```

## Key Findings
1. The go-openai SDK properly constructs URLs by appending standard OpenAI endpoints to the base URL
2. The Gemini OpenAI-compatible endpoint works correctly when the base URL is set to `https://generativelanguage.googleapis.com/v1beta/openai`
3. The 404 errors were likely due to incorrect base URL configuration, not SDK issues