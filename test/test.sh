#!/bin/bash

# Main Test Script
# Runs all test suites across all test directories

# Note: Not using 'set -e' to allow continuing after individual test failures

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Navigate to the test directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}${BOLD}=== Konfigo Main Test Suite ===${NC}"
echo -e "${CYAN}Running all test suites...${NC}"
echo

# List of all test directories
TEST_DIRS=(
    "batch"
    "env-integration"
    "format-conversion"
    "generators"
    "merging"
    "recursive-discovery"
    "schema-integration"
    "transformers"
    "validators"
    "variables"
)

# Counters
TOTAL_SUITES=0
PASSED_SUITES=0
FAILED_SUITES=0

# Function to run a test suite
run_test_suite() {
    local test_dir="$1"
    local test_name="$(echo $test_dir | tr '-' ' ' | sed 's/\b\w/\u&/g')"
    
    TOTAL_SUITES=$((TOTAL_SUITES + 1))
    
    echo -e "${YELLOW}${BOLD}=== Running $test_name Tests ===${NC}"
    
    if [[ -d "$test_dir" ]] && [[ -f "$test_dir/test.sh" ]]; then
        cd "$test_dir"
        
        if ./test.sh 2>&1; then
            echo -e "${GREEN}✓ $test_name test suite PASSED${NC}"
            PASSED_SUITES=$((PASSED_SUITES + 1))
            cd ..
            return 0
        else
            local exit_code=$?
            echo -e "${RED}✗ $test_name test suite FAILED (exit code: $exit_code)${NC}"
            echo -e "${YELLOW}  Continuing with remaining test suites...${NC}"
            FAILED_SUITES=$((FAILED_SUITES + 1))
            cd ..
            return 1
        fi
    else
        echo -e "${RED}✗ $test_name test suite: test.sh not found in $test_dir${NC}"
        FAILED_SUITES=$((FAILED_SUITES + 1))
        return 1
    fi
}

# Check if konfigo binary exists
KONFIGO_BINARY="../konfigo"
if [[ ! -f "$KONFIGO_BINARY" ]]; then
    echo -e "${RED}❌ ERROR: konfigo binary not found at $KONFIGO_BINARY${NC}"
    echo -e "${YELLOW}Please build konfigo first: cd .. && go build -o konfigo cmd/konfigo/main.go${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Found konfigo binary at $KONFIGO_BINARY${NC}"
echo

# Run all test suites
for test_dir in "${TEST_DIRS[@]}"; do
    run_test_suite "$test_dir"
    echo
done

# Print overall summary
echo -e "${BOLD}=== Overall Test Summary ===${NC}"
echo -e "Total test suites: ${BLUE}$TOTAL_SUITES${NC}"
echo -e "Passed: ${GREEN}$PASSED_SUITES${NC}"
echo -e "Failed: ${RED}$FAILED_SUITES${NC}"

if [[ $FAILED_SUITES -eq 0 ]]; then
    echo -e "${GREEN}${BOLD}🎉 All test suites passed!${NC}"
    echo -e "${CYAN}Run './validate.sh' to validate outputs against expected results.${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}❌ Some test suites failed.${NC}"
    echo -e "${YELLOW}Check the failed test suites above for details.${NC}"
    exit 1
fi
