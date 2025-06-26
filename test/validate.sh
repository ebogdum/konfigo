#!/bin/bash

# Main Validation Script
# Runs all validation scripts across all test directories

set -e

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

echo -e "${BLUE}${BOLD}=== Konfigo Main Validation Suite ===${NC}"
echo -e "${CYAN}Validating all test outputs against expected results...${NC}"
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
TOTAL_VALIDATIONS=0
PASSED_VALIDATIONS=0
FAILED_VALIDATIONS=0

# Function to run a validation suite
run_validation_suite() {
    local test_dir="$1"
    local test_name="$(echo $test_dir | tr '-' ' ' | sed 's/\b\w/\u&/g')"
    
    TOTAL_VALIDATIONS=$((TOTAL_VALIDATIONS + 1))
    
    echo -e "${YELLOW}${BOLD}=== Validating $test_name Outputs ===${NC}"
    
    if [[ -d "$test_dir" ]] && [[ -f "$test_dir/validate.sh" ]]; then
        cd "$test_dir"
        
        # Check if output directory exists (test was run)
        if [[ ! -d "output" ]]; then
            echo -e "${RED}✗ $test_name validation: No output directory found${NC}"
            echo -e "${YELLOW}  Run 'cd $test_dir && ./test.sh' first to generate outputs${NC}"
            FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
            cd ..
            return 1
        fi
        
        # Check if expected directory exists
        if [[ ! -d "expected" ]]; then
            echo -e "${RED}✗ $test_name validation: No expected directory found${NC}"
            echo -e "${YELLOW}  Run 'cd $test_dir && ./test.sh' first to generate baseline${NC}"
            FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
            cd ..
            return 1
        fi
        
        if ./validate.sh; then
            echo -e "${GREEN}✓ $test_name validation PASSED${NC}"
            PASSED_VALIDATIONS=$((PASSED_VALIDATIONS + 1))
            cd ..
            return 0
        else
            echo -e "${RED}✗ $test_name validation FAILED${NC}"
            FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
            cd ..
            return 1
        fi
    else
        echo -e "${RED}✗ $test_name validation: validate.sh not found in $test_dir${NC}"
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
}

# Run all validation suites
for test_dir in "${TEST_DIRS[@]}"; do
    run_validation_suite "$test_dir"
    echo
done

# Print overall summary
echo -e "${BOLD}=== Overall Validation Summary ===${NC}"
echo -e "Total validation suites: ${BLUE}$TOTAL_VALIDATIONS${NC}"
echo -e "Passed: ${GREEN}$PASSED_VALIDATIONS${NC}"
echo -e "Failed: ${RED}$FAILED_VALIDATIONS${NC}"

if [[ $FAILED_VALIDATIONS -eq 0 ]]; then
    echo -e "${GREEN}${BOLD}🎉 All validations passed!${NC}"
    echo -e "${CYAN}All test outputs match expected results exactly.${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}❌ Some validations failed.${NC}"
    echo -e "${YELLOW}Check the failed validations above for details.${NC}"
    echo -e "${CYAN}Use 'diff' commands to see specific differences between expected and actual outputs.${NC}"
    exit 1
fi
