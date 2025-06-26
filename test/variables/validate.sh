#!/bin/bash

# Variables Testing Validation Script
# Compares current output with expected results

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Variables Test Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
        echo
        return 1
    fi
}

# Run the tests first to generate fresh output
echo "Running variable tests to generate fresh output..."
./test.sh > /dev/null 2>&1 || true

echo "Comparing outputs with expected results..."
echo

# Compare all test output files
if [ -d "expected" ] && [ -d "output" ]; then
    for expected_file in expected/*.yaml expected/*.json expected/*.toml expected/*.env; do
        if [ -f "$expected_file" ]; then
            filename=$(basename "$expected_file")
            output_file="output/$filename"
            test_name=$(echo "$filename" | sed 's/\.[^.]*$//')
            
            compare_files "$expected_file" "$output_file" "$test_name"
        fi
    done
else
    echo -e "${RED}Missing expected or output directories${NC}"
    exit 1
fi

echo
echo -e "${YELLOW}=== Validation Results ===${NC}"
echo "Total files: $TOTAL_FILES"
echo "Passed: $PASSED_FILES"
echo "Failed: $((TOTAL_FILES - PASSED_FILES))"

if [ $PASSED_FILES -eq $TOTAL_FILES ]; then
    echo -e "${GREEN}All validations passed!${NC}"
    exit 0
else
    echo -e "${RED}Some validations failed.${NC}"
    exit 1
fi
