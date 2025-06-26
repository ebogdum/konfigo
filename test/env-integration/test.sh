#!/bin/bash

# Environment Variable Integration Test Suite  
# Tests KONFIGO_KEY_ direct configuration overrides and KONFIGO_VAR_ variable integration

set -e  # Exit on any error

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Environment Variable Integration Test Suite"

INPUT_DIR="input"
CONFIG_DIR="config"

# Test 1: Basic KONFIGO_KEY_ overrides
run_test "Basic KONFIGO_KEY_ overrides" \
    "env 'KONFIGO_KEY_app.name=env-override-app' 'KONFIGO_KEY_app.port=9090' 'KONFIGO_KEY_database.ssl=true' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/basic-overrides.yaml"

# Test 2: Nested path KONFIGO_KEY_ overrides
run_test "Nested path KONFIGO_KEY_ overrides" \
    "env 'KONFIGO_KEY_database.connection.timeout=60' 'KONFIGO_KEY_database.connection.pool_size=25' 'KONFIGO_KEY_nested.deep.very.deep.value=env-modified' 'KONFIGO_KEY_logging.outputs.0=syslog' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/nested-overrides.yaml"

# Test 3: KONFIGO_KEY_ with multiple input files
run_test "KONFIGO_KEY_ with multiple input files" \
    "env 'KONFIGO_KEY_app.environment=development' 'KONFIGO_KEY_security.tls_version=1.2' 'KONFIGO_KEY_metrics.enabled=false' $KONFIGO -s $INPUT_DIR/base-config.json,$INPUT_DIR/override-config.yaml,$INPUT_DIR/additional-config.toml -of $OUTPUT_DIR/multi-file-overrides.yaml"

# Test 4: KONFIGO_VAR_ variable integration
run_test "KONFIGO_VAR_ variable integration" \
    "env 'KONFIGO_VAR_KONFIGO_SERVICE_NAME=test-service' 'KONFIGO_VAR_KONFIGO_ENVIRONMENT=staging' $KONFIGO -s $INPUT_DIR/base-config.json -S $CONFIG_DIR/var-integration-schema.yaml -of $OUTPUT_DIR/var-integration.yaml"

# Test 5: Combined KONFIGO_KEY_ and KONFIGO_VAR_
run_test "Combined KONFIGO_KEY_ and KONFIGO_VAR_" \
    "env 'KONFIGO_KEY_app.version=2.0.0' 'KONFIGO_KEY_logging.level=debug' 'KONFIGO_VAR_KONFIGO_SERVICE_NAME=production-service' 'KONFIGO_VAR_KONFIGO_ENVIRONMENT=production' $KONFIGO -s $INPUT_DIR/base-config.json -S $CONFIG_DIR/var-integration-schema.yaml -of $OUTPUT_DIR/combined-env-vars.yaml"

# Test 6: KONFIGO_KEY_ with schema validation
run_test "KONFIGO_KEY_ with schema validation" \
    "env 'KONFIGO_KEY_app.port=3000' 'KONFIGO_KEY_database.connection.pool_size=50' 'KONFIGO_KEY_features.auth=false' $KONFIGO -s $INPUT_DIR/base-config.json -S $CONFIG_DIR/env-friendly-schema.yaml -of $OUTPUT_DIR/env-with-schema.yaml"

# Test 7: KONFIGO_KEY_ vs immutable paths
run_test "KONFIGO_KEY_ vs immutable paths" \
    "env 'KONFIGO_KEY_app.name=should-override-immutable' 'KONFIGO_KEY_database.host=should-override-immutable' 'KONFIGO_KEY_security.enabled=should-override-immutable' 'KONFIGO_KEY_app.version=4.0.0' $KONFIGO -s $INPUT_DIR/base-config.json,$INPUT_DIR/override-config.yaml -S $CONFIG_DIR/immutable-schema.yaml -of $OUTPUT_DIR/immutable-override.yaml"

# Test 8: New key creation with KONFIGO_KEY_
run_test "New key creation with KONFIGO_KEY_" \
    "env 'KONFIGO_KEY_runtime.environment=staging' 'KONFIGO_KEY_runtime.deployment_id=deploy-123' 'KONFIGO_KEY_new_section.enabled=true' 'KONFIGO_KEY_new_section.config.timeout=45' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/new-keys.yaml"

# Test 9: Array index overrides with KONFIGO_KEY_
run_test "Array index overrides with KONFIGO_KEY_" \
    "env 'KONFIGO_KEY_logging.outputs.0=stderr' 'KONFIGO_KEY_logging.outputs.1=journald' 'KONFIGO_KEY_logging.outputs.2=custom' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/array-overrides.yaml"

# Test 10: Type conversion with KONFIGO_KEY_
run_test "Type conversion with KONFIGO_KEY_" \
    "env 'KONFIGO_KEY_app.version=8443' 'KONFIGO_KEY_database.ssl=false' 'KONFIGO_KEY_features.cache=true' 'KONFIGO_KEY_database.connection.timeout=120.5' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/type-conversion.yaml"

# Test 11a: Env vars to JSON
run_test "Env vars to JSON" \
    "env 'KONFIGO_KEY_app.format_test=json-output' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/env-override.json"

# Test 11b: Env vars to TOML
run_test "Env vars to TOML" \
    "env 'KONFIGO_KEY_app.format_test=toml-output' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/env-override.toml"

# Test 11c: Env vars to ENV
run_test "Env vars to ENV" \
    "env 'KONFIGO_KEY_app.format_test=env-output' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/env-override.env"

# Test 12: Precedence testing - env vars vs file values
echo -e "\n--- Test 12: Precedence testing ---"
run_test "Precedence testing" \
    "env 'KONFIGO_KEY_app.environment=env-wins' 'KONFIGO_KEY_logging.level=env-debug' 'KONFIGO_KEY_security.tls_version=env-tls-1.3' $KONFIGO -s $INPUT_DIR/base-config.json,$INPUT_DIR/override-config.yaml -of $OUTPUT_DIR/precedence-test.yaml"

# Test 13: Edge case - Complex key names
run_test "Edge case - Complex key names" \
    "env 'KONFIGO_KEY_complex-key-name.with_underscore.and-dash=complex-value' 'KONFIGO_KEY_123numeric.start=numeric-key' $KONFIGO -s $INPUT_DIR/base-config.json -of $OUTPUT_DIR/complex-keys.yaml"

# Print test summary
print_test_summary
exit $?
