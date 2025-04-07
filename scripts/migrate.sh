#!/bin/bash

# Load environment variables
if [ -f "../.env" ]; then
    source ../.env
elif [ -f ".env" ]; then
    source .env
else
    echo "Error: .env file not found"
    exit 1
fi

# Get the absolute path of the script's directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Run migrations
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$PROJECT_ROOT/internal/db/migrations/000001_init_schema.up.sql"

echo "Migrations completed" 