#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create necessary directories
print_status "Creating directories..."
timestamp=$(date +%Y%m%d_%H%M%S)
coverage_dir=./coverage-data/$timestamp
mkdir -p $coverage_dir

# Build coverage dumper
# print_status "Building coverage dumper..."
# go build -cover -covermode=atomic -o bin/coverage-dump ./cmd/coverage-dump
# if [ $? -ne 0 ]; then
#     print_error "Failed to build coverage dumper"
#     exit 1
# fi

# Build example integration
print_status "Building example integration..."
go build -cover -covermode=atomic -o bin/integration ./examples/integration
if [ $? -ne 0 ]; then
    print_error "Failed to build example integration"
    exit 1
fi

# Build coverage analyzer
# print_status "Building coverage analyzer..."
# go build -o bin/coverage-analyzer ./cmd/coverage-analyzer
# if [ $? -ne 0 ]; then
#     print_error "Failed to build coverage analyzer"
#     exit 1
# fi

# Kill any existing processes
print_status "Cleaning up existing processes..."
# pkill -f "coverage-dump" || true
pkill -f "integration" || true

# Start coverage dumper in background
# print_status "Starting coverage dumper..."
# ./bin/coverage-dump -output ./coverage-data -interval 1m &
# COVERAGE_DUMPER_PID=$!

# Wait a bit for the dumper to start
sleep 2

# Start example integration in background
print_status "Starting example integration..."
# ./bin/integration -port 8080 -coverage-dir $coverage_dir -dump-interval 1m &
./bin/integration -coverage-dir $coverage_dir &
INTEGRATION_PID=$!

# Wait a bit for the integration to start
sleep 2

# Run coverage generation script
# print_status "Generating coverage data..."
# ./scripts/generate_coverage.sh

# Print status
print_status "Services are running:"
# echo "Coverage dumper (PID: $COVERAGE_DUMPER_PID)"
echo "Example integration (PID: $INTEGRATION_PID)"
echo ""
echo "Test endpoints:"
echo "  http://localhost:8080/test (GET/POST)"
echo "  http://localhost:8080/coverage (GET)"
echo ""
echo "Coverage data directory: $coverage_dir"
echo ""
echo "To stop services, run: ./scripts/stop_services.sh" 