#!/bin/bash
# Setup script for Sift backend

set -e

cd "$(dirname "$0")"

echo "=== Sift Backend Setup ==="

# Check Go version
echo "Checking Go installation..."
go version

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the application
echo "Building the application..."
go build -o bin/sift ./cmd/main.go

echo "=== Setup complete! ==="
echo ""
echo "To run the server:"
echo "  1. Copy .env.example to .env and configure your settings"
echo "  2. Run migrations: migrate -path migrations -database \"\$DATABASE_URL\" up"
echo "  3. Start the server: ./bin/sift"
