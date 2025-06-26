#!/bin/bash

# Batch Processing Validation Script
# Validates that current outputs match expected outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Batch Processing Validation"

# Validate main directory comparison
compare_directories "$EXPECTED_DIR" "$OUTPUT_DIR" "Batch Processing Outputs"

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
