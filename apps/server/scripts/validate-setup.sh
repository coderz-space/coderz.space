#!/bin/bash

# Validation script for CI/CD setup
set -e

echo "🔍 Validating CI/CD Setup..."
echo ""

# Check required files
echo "✓ Checking required files..."
required_files=(
    "dockerfile"
    ".dockerignore"
    ".golangci.yml"
    ".env.example"
    "docker-compose.yml"
    "go.mod"
    "go.sum"
    "Makefile"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✓ $file exists"
    else
        echo "  ✗ $file missing"
        exit 1
    fi
done

echo ""
echo "✓ Checking Go module..."
go mod verify
echo "  ✓ Go modules verified"

echo ""
echo "✓ Checking Go formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "  ✗ Code needs formatting. Run: gofmt -w ."
    exit 1
else
    echo "  ✓ Code is properly formatted"
fi

echo ""
echo "✓ Checking Go vet..."
go vet ./...
echo "  ✓ Go vet passed"

echo ""
echo "✓ Checking Dockerfile syntax..."
if docker build -f dockerfile -t coderz-test:latest . > /dev/null 2>&1; then
    echo "  ✓ Dockerfile builds successfully"
    docker rmi coderz-test:latest > /dev/null 2>&1
else
    echo "  ✗ Dockerfile build failed"
    exit 1
fi

echo ""
echo "✓ Checking docker-compose syntax..."
docker compose config > /dev/null
echo "  ✓ docker-compose.yml is valid"

echo ""
echo "✅ All validations passed!"
echo ""
echo "Next steps:"
echo "  1. Run 'make docker-up' to start PostgreSQL"
echo "  2. Run 'make migrate-up' to apply migrations"
echo "  3. Run 'make swagger' to generate API docs"
echo "  4. Run 'make run' to start the server"
echo "  5. Visit http://localhost:8080/swagger/index.html"
