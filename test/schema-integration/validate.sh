#!/bin/bash

# Schema Processing Integration Validation Script
# Validates that current outputs match expected outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Schema Processing Integration Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
        return 0
    else
        echo "❌ FAIL: $test_name - Output differs from expected"
        echo "Differences:"
        diff "$output_file" "$expected_file" || true
        return 1
    fi
}

# Check if test has been run
if [[ ! -d "$OUTPUT_DIR" ]]; then
    echo "❌ No output directory found. Please run test.sh first."
    exit 1
fi

total_tests=0
passed_tests=0

# Function to get test description
get_test_description() {
    case "$1" in
        "full-integration.yaml") echo "Full Integration Workflow (YAML)" ;;
        "full-integration.json") echo "Full Integration Workflow (JSON)" ;;
        "strict-validation.yaml") echo "Strict Schema Validation" ;;
        "input-only.yaml") echo "Input Schema Only" ;;
        "output-only.yaml") echo "Output Schema Only" ;;
        "with-variables.yaml") echo "Integration with Variables" ;;
        "multi-variables.yaml") echo "Multiple Variable Files" ;;
        "immutable-paths.yaml") echo "Immutable Paths Integration" ;;
        "full-integration.toml") echo "Integration to TOML" ;;
        "full-integration.env") echo "Integration to ENV" ;;
        "env-override.yaml") echo "Environment Variable Override" ;;
        *) echo "$1" ;;
    esac
}

# Validate each output file
for output_file in "$OUTPUT_DIR"/*; do
    if [[ -f "$output_file" ]]; then
        filename=$(basename "$output_file")
        expected_file="$EXPECTED_DIR/$filename"
        test_description=$(get_test_description "$filename")
        
        total_tests=$((total_tests + 1))
        
        if compare_files "$output_file" "$expected_file" "$test_description"; then
            passed_tests=$((passed_tests + 1))
        fi
    fi
done

# Count total files
output_files=$(find "$OUTPUT_DIR" -type f | wc -l | tr -d ' ')
expected_files=$(find "$EXPECTED_DIR" -type f | wc -l | tr -d ' ')

echo ""
echo "=== Validation Summary ==="
echo "Tests passed: $passed_tests/$total_tests"
echo "Output files generated: $output_files"
echo "Expected files: $expected_files"

if [[ $passed_tests -eq $total_tests ]] && [[ $output_files -eq $expected_files ]]; then
    echo "✅ All schema processing integration tests passed!"
    exit 0
else
    echo "❌ Some schema processing integration tests failed!"
    exit 1
fi
