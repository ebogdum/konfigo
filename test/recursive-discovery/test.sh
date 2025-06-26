#!/bin/bash

# Recursive Discovery Test Suite
# Tests konfigo's -r (recursive) flag for discovering configuration files in subdirectories

set -e  # Exit on any error

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Recursive Discovery Test Suite"

INPUT_DIR="input"
CONFIG_DIR="config"
mkdir -p "$OUTPUT_DIR"

echo "=== Recursive Discovery Test Suite ==="
echo "Testing konfigo's recursive file discovery functionality..."

# Test 1: Basic recursive discovery
echo -e "\n--- Test 1: Basic recursive discovery ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/basic-recursive.yaml"

# Test 2: Recursive discovery with schema
echo -e "\n--- Test 2: Recursive discovery with schema ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -S "$CONFIG_DIR/discovery-schema.yaml" \
    -of "$OUTPUT_DIR/recursive-with-schema.yaml"

# Test 3: Recursive discovery different output formats
echo -e "\n--- Test 3: Recursive discovery - JSON output ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-discovery.json"

# Test 4: Recursive discovery - TOML output
echo -e "\n--- Test 4: Recursive discovery - TOML output ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-discovery.toml"

# Test 5: Recursive discovery - ENV output
echo -e "\n--- Test 5: Recursive discovery - ENV output ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-discovery.env"

# Test 6: Recursive vs non-recursive comparison
echo -e "\n--- Test 6: Non-recursive (single file) comparison ---"
$KONFIGO \
    -s "$INPUT_DIR/configs/base.json" \
    -of "$OUTPUT_DIR/non-recursive-single.yaml"

# Test 7: Non-recursive with specific files
echo -e "\n--- Test 7: Non-recursive with specific files ---"
$KONFIGO \
    -s "$INPUT_DIR/configs/base.json,$INPUT_DIR/configs/app/app.yaml,$INPUT_DIR/configs/database/database.toml" \
    -of "$OUTPUT_DIR/non-recursive-specific.yaml"

# Test 8: Recursive discovery with environment variables
echo -e "\n--- Test 8: Recursive discovery with environment variables ---"
env "KONFIGO_KEY_app.environment=production" \
    "KONFIGO_KEY_global.debug=false" \
    "KONFIGO_KEY_metadata.override_test=env-override" \
    $KONFIGO \
    -r \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-with-env.yaml"

# Test 9: Recursive discovery from subdirectory
echo -e "\n--- Test 9: Recursive discovery from subdirectory ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs/services" \
    -of "$OUTPUT_DIR/recursive-services-only.yaml"

# Test 10: Multiple directory recursive discovery
echo -e "\n--- Test 10: Multiple directory recursive discovery ---"
$KONFIGO \
    -r \
    -s "$INPUT_DIR/configs/app,$INPUT_DIR/configs/database" \
    -of "$OUTPUT_DIR/recursive-multi-dirs.yaml"

# Test 11: Case sensitivity with recursive discovery
echo -e "\n--- Test 11: Case sensitive recursive discovery ---"
$KONFIGO \
    -r \
    -c \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-case-sensitive.yaml"

# Test 12: Debug mode with recursive discovery
echo -e "\n--- Test 12: Debug mode recursive discovery ---"
$KONFIGO \
    -r \
    -d \
    -s "$INPUT_DIR/configs" \
    -of "$OUTPUT_DIR/recursive-debug.yaml" \
    2> "$OUTPUT_DIR/recursive-debug.log"

echo -e "\n=== Recursive Discovery Tests Completed ==="
echo "Output files generated in: $OUTPUT_DIR"
echo "Run 'validate.sh' to compare with expected outputs."
