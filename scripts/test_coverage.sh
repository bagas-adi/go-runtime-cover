#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Default values
BASE_URL="http://localhost:8080"
COVERAGE_DIR="./coverage-data"
TEST_INTERVAL=2  # seconds between tests

# Function to make HTTP requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -z "$data" ]; then
        curl -s -X "$method" "$BASE_URL$endpoint"
    else
        curl -s -X "$method" -H "Content-Type: application/json" -d "$data" "$BASE_URL$endpoint"
    fi
}

# Function to test an endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    print_status "Testing $description..."
    response=$(make_request "$method" "$endpoint" "$data")
    
    if [ $? -eq 0 ]; then
        print_status "✓ $description successful"
        echo "Response: $response"
    else
        print_error "✗ $description failed"
    fi
    
    # Force a coverage dump after each test
    make_request "GET" "/coverage" > /dev/null
    sleep $TEST_INTERVAL
}

# Function to run a test suite
run_test_suite() {
    local suite_name=$1
    shift
    
    print_status "\nRunning test suite: $suite_name"
    echo "----------------------------------------"
    
    for test in "$@"; do
        $test
    done
}

# Test cases
test_home_endpoint() {
    test_endpoint "GET" "/" "" "Home endpoint"
}

test_get_endpoint() {
    test_endpoint "GET" "/test" "" "GET /test endpoint"
}

test_post_endpoint() {
    test_endpoint "POST" "/test" '{"data": "test"}' "POST /test endpoint"
}

test_invalid_method() {
    test_endpoint "PUT" "/test" "" "Invalid method (PUT) on /test endpoint"
}

test_delete_endpoint() {
    test_endpoint "DELETE" "/test" "" "DELETE /test endpoint"
}

test_coverage_endpoint() {
    test_endpoint "GET" "/coverage" "" "Coverage dump endpoint"
}

# Main test execution
main() {
    print_status "Starting coverage test suite"
    print_status "Base URL: $BASE_URL"
    print_status "Coverage directory: $COVERAGE_DIR"
    
    # Check if the server is running
    if ! curl -s "$BASE_URL" > /dev/null; then
        print_error "Server is not running at $BASE_URL"
        exit 1
    fi
    
    # Run test suites
    run_test_suite "Basic Endpoints" \
        test_home_endpoint \
        test_get_endpoint \
        test_post_endpoint
    
    run_test_suite "Error Cases" \
        test_invalid_method \
        test_delete_endpoint
    
    run_test_suite "Coverage Operations" \
        test_coverage_endpoint
    
    # Final coverage dump
    print_status "\nPerforming final coverage dump..."
    make_request "GET" "/coverage" > /dev/null
    
    print_status "\nTest suite completed"
    print_status "Coverage data has been generated in $COVERAGE_DIR"
    
    # List coverage files
    print_status "\nGenerated coverage files:"
    ls -l "$COVERAGE_DIR"
}

# Run the main function
main 