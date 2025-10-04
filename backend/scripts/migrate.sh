#!/bin/bash

# Database migration script using Goose
# Usage: ./migrate.sh [command]
# Commands: up, down, status, version, create [name]

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"genshin_quiz"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"password"}
MIGRATION_DIR=${MIGRATION_DIR:-"./migrations"}

# Build connection string
CONNECTION_STRING="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "❌ goose is not installed. Installing..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Function to display usage
usage() {
    echo "Usage: $0 [command] [args...]"
    echo ""
    echo "Commands:"
    echo "  up                 Apply all pending migrations"
    echo "  down               Rollback the last migration"
    echo "  status             Show migration status"
    echo "  version            Show current migration version"
    echo "  create [name]      Create a new migration file"
    echo "  reset              Reset database (down to 0, then up)"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 down"
    echo "  $0 create add_user_preferences"
    echo "  $0 status"
}

# Parse command
COMMAND=${1:-"status"}

case $COMMAND in
    "up")
        echo "🔄 Applying migrations..."
        goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" up
        echo "✅ Migrations applied successfully!"
        ;;
    
    "down")
        echo "🔄 Rolling back last migration..."
        goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" down
        echo "✅ Migration rolled back successfully!"
        ;;
    
    "status")
        echo "📊 Migration status:"
        goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" status
        ;;
    
    "version")
        echo "📋 Current migration version:"
        goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" version
        ;;
    
    "create")
        if [ -z "$2" ]; then
            echo "❌ Migration name is required"
            echo "Usage: $0 create [migration_name]"
            exit 1
        fi
        
        echo "📝 Creating new migration: $2"
        goose -dir "$MIGRATION_DIR" create "$2" sql
        echo "✅ Migration file created!"
        ;;
    
    "reset")
        echo "⚠️  This will reset the entire database!"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "🔄 Resetting database..."
            goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" reset
            echo "🔄 Applying all migrations..."
            goose -dir "$MIGRATION_DIR" postgres "$CONNECTION_STRING" up
            echo "✅ Database reset completed!"
        else
            echo "❌ Reset cancelled."
        fi
        ;;
    
    "help"|"-h"|"--help")
        usage
        ;;
    
    *)
        echo "❌ Unknown command: $COMMAND"
        echo ""
        usage
        exit 1
        ;;
esac