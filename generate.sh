#!/bin/bash
# generate.sh - Generate Go API client from OpenAPI specification using ogen
#
# Prerequisites:
#   go install github.com/ogen-go/ogen/cmd/ogen@latest
#
# Usage:
#   ./generate.sh
#
# This script:
#   1. Downloads the latest ElevenLabs OpenAPI 3.1 spec (optional, with --fetch)
#   2. Converts OpenAPI 3.1 to 3.0.3 for ogen compatibility
#   3. Runs ogen to generate Go code
#   4. Runs go mod tidy to update dependencies
#   5. Verifies the build compiles

set -e

# Check if ogen is installed
if ! command -v ogen &> /dev/null; then
    echo "Error: ogen is not installed."
    echo "Install with: go install github.com/ogen-go/ogen/cmd/ogen@latest"
    exit 1
fi

# Optionally fetch the latest spec
if [ "$1" == "--fetch" ]; then
    echo "Fetching latest ElevenLabs OpenAPI spec..."
    curl -s https://api.elevenlabs.io/openapi.json -o openapi/openapi-v3.1.json
    echo "Saved to openapi/openapi-v3.1.json"
fi

# Check if OpenAPI 3.1 spec exists
if [ ! -f "openapi/openapi-v3.1.json" ]; then
    echo "Error: openapi/openapi-v3.1.json not found."
    echo "Download it with:"
    echo "  curl -o openapi/openapi-v3.1.json https://api.elevenlabs.io/openapi.json"
    echo "Or run: ./generate.sh --fetch"
    exit 1
fi

# Convert 3.1 to 3.0.3
echo "Converting OpenAPI 3.1 to 3.0.3..."
go run ./cmd/openapi-convert openapi/openapi-v3.1.json openapi/openapi-v3.0.json

# Generate API code
echo ""
echo "Generating API code with ogen..."
ogen --package api --target internal/api --clean openapi/openapi-v3.0.json

# Post-process: Fix ogen null handling bug (https://github.com/ogen-go/ogen/issues/1358)
echo ""
echo "Post-processing: Fixing Opt* null handling..."
go run github.com/agentplexus/ogen-tools/cmd/ogen-fixnull@latest internal/api/oas_json_gen.go

# Post-process: Fix error body preservation (body gets closed before caller can read it)
echo ""
echo "Post-processing: Fixing error body preservation..."
go run github.com/agentplexus/ogen-tools/cmd/ogen-fixerror@latest internal/api/oas_response_decoders_gen.go

echo ""
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "Verifying build..."
go build ./...

echo ""
echo "Done! API client regenerated successfully."
echo ""
echo "Next steps:"
echo "  1. Review changes in internal/api/"
echo "  2. Update SDK wrapper code if needed for new/changed endpoints"
echo "  3. Run tests: go test ./..."
echo "  4. Run linter: golangci-lint run"
