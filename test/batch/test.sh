#!/bin/bash

# Batch Processing Test Suite
# Tests konfigo_forEach functionality across all supported formats

set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Batch Processing Test Suite"

# Define directories
INPUT_DIR="input"
CONFIG_DIR="config"
VARIABLES_DIR="variables"

# Test 1: Basic Services Batch (items)
run_test "Basic Services Batch (items)" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/services-batch.yaml"

# Test 2: Basic Services Batch to JSON
run_test "Basic Services Batch to JSON" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/services-batch.json"

# Test 3: Deployments Batch (items)
run_test "Deployments Batch (items)" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/deployments-batch.yaml -S $CONFIG_DIR/deployment-schema.yaml -of $OUTPUT_DIR/deployments-batch.yaml"

# Test 4: Complex Multi-Level Batch
run_test "Complex Multi-Level Batch" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/complex-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/complex-batch.yaml"

# Test 5: Environment Files Batch (itemFiles)
run_test "Environment Files Batch (itemFiles)" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/envs-itemfiles-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/envs-itemfiles-batch.yaml"

# Test 6: Empty Base Config with Batch
run_test "Empty Base Config with Batch" \
    "$KONFIGO -s $INPUT_DIR/empty-base.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/empty-base-services.yaml"

# Test 7: Batch with Schema Validation
run_test "Batch with Schema Validation" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/services-with-schema.yaml"

# Test 8: Batch with Different Output Formats
run_test "Batch to TOML" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/services-batch.toml"

run_test "Batch to ENV" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -S $CONFIG_DIR/service-schema.yaml -oe -of $OUTPUT_DIR/services-batch.env"

# Test 9: Multiple variable files (advanced use case)
run_test "Multiple variable files" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/services-batch.yaml -V $VARIABLES_DIR/deployments-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/multi-batch.yaml"

# Test 10: Error Case - Missing itemFiles (should fail)
# Create a batch config with invalid itemFiles reference
echo 'konfigo_forEach:
  item: service
  itemFiles: 
    - "nonexistent/file1.yaml"
    - "nonexistent/file2.yaml"
  output: "services-{service.name}.yaml"' > "$VARIABLES_DIR/error-batch.yaml"

run_test "Error handling test (missing itemFiles)" \
    "$KONFIGO -s $INPUT_DIR/base-config.json -V $VARIABLES_DIR/error-batch.yaml -S $CONFIG_DIR/service-schema.yaml -of $OUTPUT_DIR/error-batch.yaml" \
    "false"

# Print test summary
print_test_summary
exit $?
