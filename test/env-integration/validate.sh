#!/bin/bash

# Environment Variable Integration Validation Script
# Compares test outputs with expected results

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Environment Variable Integration Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
    if [[ -f "$output_file" ]]; then
        filename=$(basename "$output_file")
        expected_file="$EXPECTED_DIR/$filename"
        
        total_files=$((total_files + 1))
        
        if [[ ! -f "$expected_file" ]]; then
            echo "❌ Missing expected file: $filename"
            differing_files=$((differing_files + 1))
            continue
        fi
        
        if diff -q "$output_file" "$expected_file" >/dev/null; then
            echo "✅ $filename"
            matching_files=$((matching_files + 1))
        else
            echo "❌ $filename"
            echo "   Differences found:"
            diff "$output_file" "$expected_file" | head -10
            echo "   ..."
            differing_files=$((differing_files + 1))
        fi
    fi
done

echo -e "\n=== Validation Summary ==="
echo "Total files: $total_files"
echo "Matching: $matching_files"
echo "Differing: $differing_files"

if [[ $differing_files -eq 0 ]]; then
    echo -e "\n🎉 All environment variable integration tests passed!"
    exit 0
else
    echo -e "\n❌ $differing_files test(s) failed validation"
    exit 1
fi
