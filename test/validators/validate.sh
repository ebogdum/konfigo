#!/bin/bash

# Validation Script for Validators Test Suite
# Compares expected outputs with actual outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Validators Test Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
            if diff -q "$expected_file" "$output_file" >/dev/null 2>&1; then
                echo "✓ $filename: MATCH"
                match_count=$((match_count + 1))
            else
                echo "✗ $filename: DIFFERENT"
                echo "  Expected: $(wc -c < "$expected_file") bytes"
                echo "  Actual:   $(wc -c < "$output_file") bytes"
            fi
        else
            echo "✗ $filename: MISSING in output/"
        fi
    fi
done

echo
echo "=== Validation Summary ==="
echo "Total files validated: $validation_count"
echo "Matching files: $match_count"
echo "Different/missing files: $((validation_count - match_count))"

if [ $match_count -eq $validation_count ]; then
    echo "✓ All outputs match expected results!"
    exit 0
else
    echo "✗ Some outputs don't match expected results."
    exit 1
fi
