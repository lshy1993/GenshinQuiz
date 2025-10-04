#!/bin/bash

# Generate OpenAPI code using oapi-codegen
# This script generates Go types and server stubs from the OpenAPI specification

set -e

echo "Generating OpenAPI code..."

# Create output directory
mkdir -p internal/generated

# Generate types (models)
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
  -package generated \
  -generate types \
  -o internal/generated/types.go \
  api/openapi.yaml

# Generate server interface
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
  -package generated \
  -generate chi-server \
  -o internal/generated/server.go \
  api/openapi.yaml

# Generate client (optional, for testing)
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
  -package generated \
  -generate client \
  -o internal/generated/client.go \
  api/openapi.yaml

echo "OpenAPI code generation completed!"
echo "Generated files:"
echo "  - internal/generated/types.go (data types)"
echo "  - internal/generated/server.go (server interface)"
echo "  - internal/generated/client.go (client for testing)"