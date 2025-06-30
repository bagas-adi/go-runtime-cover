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

# Stop coverage dumper
print_status "Stopping coverage dumper..."
pkill -f "coverage-dump" || true

# Stop example integration
print_status "Stopping example integration..."
pkill -f "integration" || true

# Wait for processes to stop
sleep 2

# Verify processes are stopped
if pgrep -f "coverage-dump" > /dev/null; then
    print_error "Failed to stop coverage dumper"
else
    print_status "Coverage dumper stopped"
fi

if pgrep -f "integration" > /dev/null; then
    print_error "Failed to stop example integration"
else
    print_status "Example integration stopped"
fi

print_status "All services stopped" 