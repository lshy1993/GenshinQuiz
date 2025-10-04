#!/bin/bash

# Generate Go-Jet models from database schema
# This script should be run after database migrations

set -e

# Configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"genshin_quiz"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"password"}
DB_SCHEMA=${DB_SCHEMA:-"public"}

# Build connection string
CONNECTION_STRING="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Generating Go-Jet models from database..."
echo "Database: $CONNECTION_STRING"

# Create models directory
mkdir -p internal/models/generated

# Generate models using Jet generator
go run github.com/go-jet/jet/v2/cmd/jet \
  -dsn="$CONNECTION_STRING" \
  -schema="$DB_SCHEMA" \
  -path="./internal/models/generated"

echo "Go-Jet models generated successfully!"

# Generate type-safe table definitions
echo "Generating table definitions..."

cat > internal/models/tables.go << 'EOF'
package models

import (
	"genshin-quiz-backend/internal/models/generated/genshin_quiz/public/table"
)

// Table definitions for type-safe queries
var (
	Users         = table.Users
	Quizzes       = table.Quizzes
	Questions     = table.Questions
	QuizAttempts  = table.QuizAttempts
	UserAnswers   = table.UserAnswers
)
EOF

echo "Table definitions generated!"
echo "Models are available in internal/models/generated/"