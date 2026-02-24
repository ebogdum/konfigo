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
