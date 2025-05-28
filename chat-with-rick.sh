#!/bin/bash
# Quick script to chat with Rick Sanchez

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "Error: OPENAI_API_KEY environment variable not set"
    echo "Please run: export OPENAI_API_KEY='your-api-key'"
    exit 1
fi

# Use the globally installed roleplay binary
ROLEPLAY_BIN="roleplay"

# Check if roleplay is in PATH
if ! command -v $ROLEPLAY_BIN &> /dev/null; then
    echo "Error: roleplay command not found in PATH"
    echo "Please ensure ~/go/bin is in your PATH"
    exit 1
fi

# Check if Rick already exists
$ROLEPLAY_BIN character list 2>/dev/null | grep -q "rick-c137"
if [ $? -ne 0 ]; then
    # Rick doesn't exist, create from JSON if available
    if [ -f "rick-sanchez.json" ]; then
        echo "Creating Rick Sanchez character..."
        $ROLEPLAY_BIN character create rick-sanchez.json
    else
        echo "Rick will be auto-created in interactive mode"
    fi
fi

# Start interactive chat
echo ""
echo "🛸 Starting chat with Rick Sanchez..."
echo "💊 Tip: Type 'Wubba lubba dub dub' to see Rick's true feelings!"
echo ""
sleep 1

# Rick is auto-created in interactive mode if not found
$ROLEPLAY_BIN interactive --character rick-c137 --user morty