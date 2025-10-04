#!/bin/bash

# Development setup script for Genshin Quiz Go Backend
# This script sets up the development environment

set -e

echo "🚀 Setting up Genshin Quiz Go Backend Development Environment"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ Created .env file. Please update it with your settings."
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
go mod download

# Install development tools
echo "🔧 Installing development tools..."

# Install goose for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install oapi-codegen for OpenAPI code generation
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Install go-jet for model generation
go install github.com/go-jet/jet/v2/cmd/jet@latest

# Install air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

echo "✅ Development tools installed!"

# Make scripts executable
echo "🔐 Making scripts executable..."
chmod +x scripts/*.sh

# Check if Docker is available
if command -v docker &> /dev/null; then
    echo "🐳 Docker is available. You can use 'docker-compose up' to start the stack."
    
    # Check if PostgreSQL is running locally
    if pg_isready -h localhost -p 5432 -U postgres 2>/dev/null; then
        echo "🐘 Local PostgreSQL detected."
    else
        echo "💡 To start PostgreSQL with Docker: docker-compose up postgres"
    fi
else
    echo "⚠️  Docker not found. Please install Docker to use the containerized setup."
fi

echo ""
echo "🎉 Setup completed! Next steps:"
echo ""
echo "1. Update your .env file with correct database credentials"
echo "2. Start PostgreSQL:"
echo "   - Using Docker: docker-compose up postgres"
echo "   - Or use your local PostgreSQL instance"
echo ""
echo "3. Run database migrations:"
echo "   ./scripts/migrate.sh up"
echo ""
echo "4. Generate Go-Jet models (after migrations):"
echo "   ./scripts/generate_models.sh"
echo ""
echo "5. Generate OpenAPI code:"
echo "   ./scripts/generate_api.sh"
echo ""
echo "6. Start the development server:"
echo "   go run main.go"
echo "   # Or with hot reloading: air"
echo ""
echo "7. Test the API:"
echo "   curl http://localhost:8080/health"
echo ""
echo "📚 Documentation:"
echo "   - API Docs: http://localhost:8080/docs (when implemented)"
echo "   - OpenAPI Spec: api/openapi.yaml"
echo ""