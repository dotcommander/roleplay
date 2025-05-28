#!/bin/bash

# Migration script for roleplay config directory
# Moves from ~/.roleplay to ~/.config/roleplay

OLD_DIR="$HOME/.roleplay"
NEW_DIR="$HOME/.config/roleplay"

echo "🔄 Migrating roleplay configuration..."

# Create new config directory structure
mkdir -p "$NEW_DIR"

# Check if old directory exists
if [ -d "$OLD_DIR" ]; then
    echo "📁 Found existing configuration at $OLD_DIR"
    
    # Move subdirectories
    for dir in characters sessions cache; do
        if [ -d "$OLD_DIR/$dir" ]; then
            echo "  Moving $dir..."
            mv "$OLD_DIR/$dir" "$NEW_DIR/"
        fi
    done
    
    # Remove old directory if empty
    if [ -z "$(ls -A "$OLD_DIR")" ]; then
        rmdir "$OLD_DIR"
        echo "✅ Removed empty directory $OLD_DIR"
    else
        echo "⚠️  Some files remain in $OLD_DIR - please review manually"
    fi
else
    echo "📝 No existing configuration found at $OLD_DIR"
fi

# Check for old config file
OLD_CONFIG="$HOME/.roleplay.yaml"
NEW_CONFIG="$NEW_DIR/config.yaml"

if [ -f "$OLD_CONFIG" ]; then
    echo "📄 Moving config file..."
    mv "$OLD_CONFIG" "$NEW_CONFIG"
fi

echo "✅ Migration complete! Configuration now at $NEW_DIR"