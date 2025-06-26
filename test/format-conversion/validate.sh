#!/bin/bash

# Format Conversion Validation Script
# Validates that actual outputs match expected outputs

set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "Format Conversion Validation"

# Validate all output files against expected
validate_all_outputs

# Print validation summary and exit with appropriate code
print_validation_summary
exit $?
        return 0
    else
        echo -e "${RED}  ✗ Output differs from expected${NC}"
        echo -e "${YELLOW}    Differences:${NC}"
        diff "$output_file" "$expected_file" | head -10
        FAILED_VALIDATIONS=$((FAILED_VALIDATIONS + 1))
        return 1
    fi
}

# Validate all format conversions
echo ""
echo "Validating JSON conversions..."
validate_output "json-to-json.json"
validate_output "json-to-yaml.yaml"
validate_output "json-to-toml.toml"
validate_output "json-to-env.env"

echo ""
echo "Validating YAML conversions..."
validate_output "yaml-to-json.json"
validate_output "yaml-to-yaml.yaml"
validate_output "yaml-to-toml.toml"
validate_output "yaml-to-env.env"

echo ""
echo "Validating TOML conversions..."
validate_output "toml-to-json.json"
validate_output "toml-to-yaml.yaml"
validate_output "toml-to-toml.toml"
validate_output "toml-to-env.env"

echo ""
echo "Validating ENV conversions..."
validate_output "env-to-json.json"
validate_output "env-to-yaml.yaml"
validate_output "env-to-toml.toml"
validate_output "env-to-env.env"

echo ""
echo "Validating INI conversions..."
validate_output "ini-to-json.json"
validate_output "ini-to-yaml.yaml"
validate_output "ini-to-toml.toml"
validate_output "ini-to-env.env"

echo ""
echo "=============================================="
echo -e "${BLUE}Format Conversion Validation Results:${NC}"
echo "  Total validations: $TOTAL_VALIDATIONS"
echo -e "  ${GREEN}Passed: $PASSED_VALIDATIONS${NC}"
echo -e "  ${RED}Failed: $FAILED_VALIDATIONS${NC}"

if [ $FAILED_VALIDATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All format conversion validations passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some format conversion validations failed.${NC}"
    exit 1
fi
