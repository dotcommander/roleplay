#!/bin/bash

# Test script to demonstrate cache improvements

echo "=== Character Bot Cache Test ==="
echo "Testing with Rick Sanchez character..."
echo ""

# Set up environment
export OPENAI_API_KEY="${OPENAI_API_KEY:-your-api-key}"
export ROLEPLAY_PROVIDER="${ROLEPLAY_PROVIDER:-openai}"
export ROLEPLAY_MODEL="${ROLEPLAY_MODEL:-gpt-4o-mini}"

# Build the project
echo "Building project..."
go build -o roleplay || exit 1

# Create Rick character
echo "Creating Rick character..."
./roleplay character create examples/rick-sanchez.json 2>/dev/null || true

# Run multiple chat commands to test caching
echo ""
echo "Sending 3 identical messages to test cache hits..."
echo ""

for i in 1 2 3; do
    echo "=== Request $i ==="
    ./roleplay chat "Hey Rick, what's your favorite invention?" \
        --character rick-c137 \
        --user morty-smith \
        --format json 2>&1 | grep -E "(cache_hit|saved_tokens|DEBUG)"
    echo ""
    sleep 2
done

echo ""
echo "Cache test complete!"
echo ""
echo "Expected behavior:"
echo "- Request 1: Cache miss (builds all layers)"
echo "- Request 2: Cache hit on personality layer"  
echo "- Request 3: Cache hit on personality layer"
echo ""
echo "To see full cache metrics, run:"
echo "./roleplay interactive --character rick-c137 --user morty"