#!/bin/bash

# Transformers Validation Script
# Validates that current transformer outputs match expected outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Konfigo Transformers Validation Suite"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?

for expected_file in "$EXPECTED_DIR"/*.{json,yaml,toml}; do
    # Skip if file doesn't exist (handles glob expansion issues)
    [[ -f "$expected_file" ]] || continue
    
    filename=$(basename "$expected_file")
    output_file="$OUTPUT_DIR/$filename"
    
    total_comparisons=$((total_comparisons + 1))
    
    if [[ ! -f "$output_file" ]]; then
        echo "❌ MISSING: $filename (expected file exists but output file missing)"
        failed_comparisons=$((failed_comparisons + 1))
        continue
    fi
    
    # Compare files
    if cmp -s "$expected_file" "$output_file"; then
        echo "✅ MATCH: $filename"
        successful_comparisons=$((successful_comparisons + 1))
    else
        echo "❌ DIFF: $filename"
        echo "  Expected: $expected_file"
        echo "  Actual:   $output_file"
        echo "  Use 'diff $expected_file $output_file' to see differences"
        failed_comparisons=$((failed_comparisons + 1))
    fi
done

echo
echo "=== Validation Results ==="
echo "Total comparisons: $total_comparisons"
echo "Successful matches: $successful_comparisons"
echo "Failed comparisons: $failed_comparisons"

if [[ $failed_comparisons -eq 0 ]]; then
    echo "✅ ALL VALIDATIONS PASSED!"
    echo "Current transformer outputs match expected results exactly."
    exit 0
else
    echo "❌ Some validations failed."
    echo "Review the differences above and update expected results if the changes are intentional."
    exit 1
fi
