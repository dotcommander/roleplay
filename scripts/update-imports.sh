#!/bin/bash

# Update all import paths from old to new module name
find . -name "*.go" -type f -exec sed -i.bak 's|github.com/yourusername/roleplay|github.com/dotcommander/roleplay|g' {} \;

# Remove backup files
find . -name "*.bak" -type f -delete

echo "Import paths updated successfully!"