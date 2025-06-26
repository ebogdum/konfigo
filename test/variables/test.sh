#!/bin/bash

# Variables Testing Script
# Tests all aspects of Konfigo variable substitution

set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Variables Testing Suite"

# Function to run test with output capture
# Ensure output directory exists
mkdir -p output

# Test 1: Basic variable substitution with schema only
run_test "Basic schema variables (JSON->YAML)" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -oy -of output/test1-schema-only.yaml"

# Test 2: Variables file overrides schema variables
run_test "Variables file overrides schema (JSON->YAML)" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oy -of output/test2-vars-file-override.yaml"

# Test 3: Environment variables override everything (KONFIGO_VAR_)
export KONFIGO_VAR_API_HOST="env-override.example.com"
export KONFIGO_VAR_NESTED_VAR="from-environment"
run_test "Environment variables override (highest precedence)" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oy -of output/test3-env-override.yaml"

# Test 4: fromEnv variable resolution
export DB_PASS="secret123"
export SERVICE_PORT="9090"
run_test "fromEnv variable resolution" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -oy -of output/test4-fromenv.yaml"

# Test 5: Different input formats with variables
run_test "YAML input with variables" \
    "$KONFIGO -s input/base-config.yaml -S config/schema-basic.yaml -oy -of output/test5-yaml-input.yaml"

run_test "TOML input with variables" \
    "$KONFIGO -s input/base-config.toml -S config/schema-basic.yaml -oy -of output/test6-toml-input.yaml"

# Test 6: Different output formats
run_test "Variables with JSON output" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oj -of output/test7-json-output.json"

run_test "Variables with TOML output" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -ot -of output/test8-toml-output.toml"

run_test "Variables with ENV output" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oe -of output/test9-env-output.env"

# Test 7: Schema with different formats
run_test "JSON schema with variables" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.json -V config/variables-basic.json -oy -of output/test10-json-schema.yaml"

# Test 8: Variable substitution without schema (basic mode)
unset KONFIGO_VAR_API_HOST KONFIGO_VAR_NESTED_VAR
run_test "Variables without schema (basic substitution)" \
    "$KONFIGO -s input/base-config.json -V config/variables-basic.yaml -oy -of output/test11-no-schema.yaml"

# Test 9: Missing variable handling (should fail gracefully)
run_test "Missing required variable (should fail)" \
    "$KONFIGO -s input/base-config.json -S config/schema-error-test.yaml -oy -of output/test12-missing-var.yaml" \
    "false"

# Test 10: Environment variable integration (KONFIGO_KEY_)
export KONFIGO_KEY_app_version="2.0.0"
export KONFIGO_KEY_features_newFeature="true"
run_test "KONFIGO_KEY_ environment integration with variables" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oy -of output/test13-konfigo-key.yaml"

# Test 11: Complex nested variable substitution
export KONFIGO_VAR_API_HOST="nested.api.com"
export KONFIGO_VAR_API_PORT="443"
run_test "Complex nested variable substitution" \
    "$KONFIGO -s input/base-config.json -S config/schema-basic.yaml -oy -of output/test14-nested-complex.yaml"

# Clean up environment variables
unset KONFIGO_VAR_API_HOST KONFIGO_VAR_NESTED_VAR KONFIGO_VAR_API_PORT
unset DB_PASS SERVICE_PORT
unset KONFIGO_KEY_app_version KONFIGO_KEY_features_newFeature

# Print test summary
print_test_summary
exit $?
