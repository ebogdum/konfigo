#!/bin/bash

# Generators Validation Script
# Validates that current generator outputs match expected outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Konfigo Generators Validation Suite"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
