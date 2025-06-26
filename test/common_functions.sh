#!/bin/bash

# Common functions and utilities for test and validation scripts
# Source this file in other test scripts for consistent behavior

# Colors for consistent output formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Global counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Common setup function
setup_test_environment() {
    local script_name="$1"
    local konfigo_path="${2:-../../konfigo}"
    
    # Set script directory and navigate to it
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[1]}")" && pwd)"
    cd "$SCRIPT_DIR"
    
    # Check if konfigo binary exists
    if [[ ! -f "$konfigo_path" ]]; then
        echo -e "${RED}❌ ERROR: konfigo binary not found at $konfigo_path${NC}"
        echo -e "${YELLOW}Please build konfigo first: cd ../../ && go build -o konfigo cmd/konfigo/main.go${NC}"
        exit 1
    fi
    
    # Set global konfigo path
    KONFIGO="$konfigo_path"
    
    # Create output directory
    OUTPUT_DIR="output"
    mkdir -p "$OUTPUT_DIR"
    
    echo -e "${BLUE}${BOLD}=== $script_name ===${NC}"
    echo
}

# Function to run a single test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_success="${3:-true}"  # Default to expecting success
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${CYAN}Test $TOTAL_TESTS: $test_name${NC}"
    echo "Command: $command"
    
    if eval "$command" >/dev/null 2>&1; then
        if [[ "$expected_success" == "true" ]]; then
            echo -e "${GREEN}✓ PASSED${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
            return 0
        else
            echo -e "${RED}✗ FAILED (expected failure but succeeded)${NC}"
            FAILED_TESTS=$((FAILED_TESTS + 1))
            return 1
        fi
    else
        if [[ "$expected_success" == "false" ]]; then
            echo -e "${GREEN}✓ PASSED (expected failure)${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
            return 0
        else
            echo -e "${RED}✗ FAILED (unexpected failure)${NC}"
            echo "Error output:"
            eval "$command" 2>&1 | head -5
            FAILED_TESTS=$((FAILED_TESTS + 1))
            return 1
        fi
    fi
}

# Function to print test summary
print_test_summary() {
    echo
    echo -e "${BOLD}=== Test Summary ===${NC}"
    echo -e "Total tests: ${BLUE}$TOTAL_TESTS${NC}"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [[ $FAILED_TESTS -eq 0 ]]; then
        echo -e "${GREEN}${BOLD}🎉 All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}${BOLD}❌ Some tests failed.${NC}"
        return 1
    fi
}

# Common validation setup function
setup_validation_environment() {
    local script_name="$1"
    
    # Set script directory and navigate to it
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[1]}")" && pwd)"
    cd "$SCRIPT_DIR"
    
    OUTPUT_DIR="output"
    EXPECTED_DIR="expected"
    
    echo -e "${BLUE}${BOLD}=== $script_name ===${NC}"
    echo
    
    # Check if directories exist
    if [[ ! -d "$OUTPUT_DIR" ]]; then
        echo -e "${RED}❌ ERROR: Output directory not found. Run test.sh first to generate outputs.${NC}"
        exit 1
    fi
    
    if [[ ! -d "$EXPECTED_DIR" ]]; then
        echo -e "${RED}❌ ERROR: Expected directory not found. Run test.sh first to generate expected outputs.${NC}"
        exit 1
    fi
}

# Global validation counters
TOTAL_VALIDATIONS=0
PASSED_VALIDATIONS=0
FAILED_VALIDATIONS=0

# Function to compare two files with detailed diff output
compare_files() {
    local expected_file="$1"
    local output_file="$2"
    local test_name="$3"
    
    TOTAL_VALIDATIONS=$((TOTAL_VALIDATIONS + 1))
    
    if [[ ! -f "$expected_file" ]]; then
        echo -e "${RED}✗ $test_name: Expected file missing: $expected_file${NC}"
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
    
    if [[ ! -f "$output_file" ]]; then
        echo -e "${RED}✗ $test_name: Output file missing: $output_file${NC}"
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
    
    if diff -q "$expected_file" "$output_file" >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $test_name: Files match${NC}"
        PASSED_VALIDATIONS=$((PASSED_VALIDATIONS + 1))
        return 0
    else
        echo -e "${RED}✗ $test_name: Files differ${NC}"
        echo -e "${YELLOW}  Expected: $expected_file${NC}"
        echo -e "${YELLOW}  Actual:   $output_file${NC}"
        echo -e "${CYAN}  Differences:${NC}"
        
        # Show contextual diff with line numbers
        diff -u --label="Expected" --label="Actual" "$expected_file" "$output_file" | head -20
        
        local total_diff_lines=$(diff "$expected_file" "$output_file" | wc -l)
        if [[ $total_diff_lines -gt 20 ]]; then
            echo -e "${YELLOW}  (showing first 20 lines of diff - total: $total_diff_lines lines)${NC}"
        fi
        echo
        
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
}

# Function to compare directories recursively
compare_directories() {
    local expected_dir="$1"
    local output_dir="$2"
    local test_name="$3"
    
    TOTAL_VALIDATIONS=$((TOTAL_VALIDATIONS + 1))
    
    if [[ ! -d "$expected_dir" ]]; then
        echo -e "${RED}✗ $test_name: Expected directory missing: $expected_dir${NC}"
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
    
    if [[ ! -d "$output_dir" ]]; then
        echo -e "${RED}✗ $test_name: Output directory missing: $output_dir${NC}"
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
    
    # Compare directories recursively
    if diff -r -q "$expected_dir" "$output_dir" >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $test_name: All files match${NC}"
        PASSED_VALIDATIONS=$((PASSED_VALIDATIONS + 1))
        return 0
    else
        echo -e "${RED}✗ $test_name: Directories differ${NC}"
        echo -e "${CYAN}  Differences:${NC}"
        diff -r --brief "$expected_dir" "$output_dir" | head -10
        
        local total_diffs=$(diff -r --brief "$expected_dir" "$output_dir" | wc -l)
        if [[ $total_diffs -gt 10 ]]; then
            echo -e "${YELLOW}  (showing first 10 differences - total: $total_diffs differences)${NC}"
        fi
        echo
        
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
}

# Function to validate all files in output directory against expected
validate_all_outputs() {
    local file_patterns=("*.json" "*.yaml" "*.yml" "*.toml" "*.env")
    
    echo -e "${CYAN}Comparing output files with expected results...${NC}"
    echo
    
    # Check for missing expected files
    for pattern in "${file_patterns[@]}"; do
        for output_file in $OUTPUT_DIR/$pattern; do
            [[ -f "$output_file" ]] || continue
            
            local filename=$(basename "$output_file")
            local expected_file="$EXPECTED_DIR/$filename"
            
            compare_files "$expected_file" "$output_file" "$filename"
        done
    done
    
    # Check for missing output files
    for pattern in "${file_patterns[@]}"; do
        for expected_file in $EXPECTED_DIR/$pattern; do
            [[ -f "$expected_file" ]] || continue
            
            local filename=$(basename "$expected_file")
            local output_file="$OUTPUT_DIR/$filename"
            
            if [[ ! -f "$output_file" ]]; then
                echo -e "${RED}⚠️  $filename: MISSING from output directory${NC}"
                TOTAL_VALIDATIONS=$((TOTAL_VALIDATIONS + 1))
                FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
            fi
        done
    done
}

# Function to print validation summary
print_validation_summary() {
    echo
    echo -e "${BOLD}=== Validation Summary ===${NC}"
    echo -e "Total validations: ${BLUE}$TOTAL_VALIDATIONS${NC}"
    echo -e "Passed: ${GREEN}$PASSED_VALIDATIONS${NC}"
    echo -e "Failed: ${RED}$FAILED_VALIDATIONS${NC}"
    
    if [[ $FAILED_VALIDATIONS -eq 0 ]]; then
        echo -e "${GREEN}${BOLD}🎉 All validations passed!${NC}"
        echo -e "${GREEN}Current outputs match expected results exactly.${NC}"
        return 0
    else
        echo -e "${RED}${BOLD}❌ Some validations failed.${NC}"
        echo -e "${YELLOW}Review the differences above and update expected results if changes are intentional.${NC}"
        return 1
    fi
}
