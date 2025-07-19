#!/bin/bash

# Script to apply role descriptions to the role database

# Default database path
DB_PATH="$HOME/.config/role/history.db"

# Check if custom path provided
if [ ! -z "$1" ]; then
    DB_PATH="$1"
fi

# Check if database exists
if [ ! -f "$DB_PATH" ]; then
    echo "Database not found at: $DB_PATH"
    echo "Please provide the correct path to the role history database"
    echo "Usage: $0 [database_path]"
    exit 1
fi

# Apply the SQL file
echo "Applying role descriptions to database: $DB_PATH"
sqlite3 "$DB_PATH" < create-descriptions-table.sql

if [ $? -eq 0 ]; then
    echo "Successfully applied role descriptions!"
    echo ""
    echo "You can now use: role list --list-roles"
    echo "to see the Hemingway-style descriptions for each role."
else
    echo "Failed to apply role descriptions."
    exit 1
fi