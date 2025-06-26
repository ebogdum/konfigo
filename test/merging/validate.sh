#!/bin/bash

# Configuration Merging Validation Script
# Validates that test outputs match expected results

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Configuration Merging Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
