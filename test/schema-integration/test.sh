#!/bin/bash

# Schema Processing Integration Test Suite
# Tests complex workflows with input/output schemas and integrated processing

set -e  # Exit on any error

KONFIGO="../../konfigo"
INPUT_DIR="input"
CONFIG_DIR="config"
SCHEMAS_DIR="schemas"
VARIABLES_DIR="variables"
OUTPUT_DIR="output"

# Ensure konfigo binary exists
if [[ ! -f "$KONFIGO" ]]; then
    echo "Error: konfigo binary not found at $KONFIGO"
    echo "Please build konfigo first: go build -o ../../konfigo ../cmd/konfigo"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "=== Schema Processing Integration Test Suite ==="
echo "Testing complex schema workflows with input/output schemas..."

# Test 1: Full Integration - Input + Output Schemas + All Features
echo -e "\n--- Test 1: Full Integration Workflow ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml,$INPUT_DIR/additional-config.toml" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -of "$OUTPUT_DIR/full-integration.yaml"

# Test 2: Full Integration to JSON
echo -e "\n--- Test 2: Full Integration to JSON ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml,$INPUT_DIR/additional-config.toml" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -of "$OUTPUT_DIR/full-integration.json"

# Test 3: Strict Input and Output Schema Validation
echo -e "\n--- Test 3: Strict Schema Validation ---"
$KONFIGO \
    -s "$INPUT_DIR/simple-config.json" \
    -S "$CONFIG_DIR/strict-schema.yaml" \
    -of "$OUTPUT_DIR/strict-validation.yaml"

# Test 4: Input Schema Only
echo -e "\n--- Test 4: Input Schema Only ---"
$KONFIGO \
    -s "$INPUT_DIR/simple-config.json" \
    -S "$CONFIG_DIR/input-only-schema.yaml" \
    -of "$OUTPUT_DIR/input-only.yaml"

# Test 5: Output Schema Only
echo -e "\n--- Test 5: Output Schema Only ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml" \
    -S "$CONFIG_DIR/output-only-schema.yaml" \
    -of "$OUTPUT_DIR/output-only.yaml"

# Test 6: Integration with Variables
echo -e "\n--- Test 6: Integration with Variables ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -V "$VARIABLES_DIR/complex-vars.yaml" \
    -of "$OUTPUT_DIR/with-variables.yaml"

# Test 7: Multiple Variable Files
echo -e "\n--- Test 7: Multiple Variable Files Integration ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -V "$VARIABLES_DIR/complex-vars.yaml" \
    -V "$VARIABLES_DIR/override-vars.yaml" \
    -of "$OUTPUT_DIR/multi-variables.yaml"

# Test 8: Immutable Paths with Schemas
echo -e "\n--- Test 8: Immutable Paths Integration ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml" \
    -S "$CONFIG_DIR/immutable-schema.yaml" \
    -of "$OUTPUT_DIR/immutable-paths.yaml"

# Test 9: Different Output Formats
echo -e "\n--- Test 9a: Integration to TOML ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -of "$OUTPUT_DIR/full-integration.toml"

echo -e "\n--- Test 9b: Integration to ENV ---"
$KONFIGO \
    -s "$INPUT_DIR/base-config.json,$INPUT_DIR/env-override.yaml" \
    -S "$CONFIG_DIR/full-integration-schema.yaml" \
    -of "$OUTPUT_DIR/full-integration.env"

# Test 10: Error Case - Strict Input Schema Violation (Expected to fail)
echo -e "\n--- Test 10: Error Case - Strict Input Schema Violation (Expected to fail) ---"
# Create a config that violates strict input schema
cat > "$INPUT_DIR/invalid-config.json" << 'EOF'
{
  "service": {
    "name": "test-service",
    "port": 8080,
    "protocol": "http",
    "invalid_field": "should-not-be-here"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp",
    "ssl": true
  },
  "unexpected_section": {
    "should_fail": true
  }
}
EOF

if $KONFIGO \
    -s "$INPUT_DIR/invalid-config.json" \
    -S "$CONFIG_DIR/strict-schema.yaml" \
    -of "$OUTPUT_DIR/should-fail.yaml" 2>/dev/null; then
    echo "ERROR: Expected failure but command succeeded"
else
    echo "EXPECTED: Command failed as expected for strict input schema violation"
fi

# Test 11: Error Case - Strict Output Schema Violation (Expected to fail)
echo -e "\n--- Test 11: Error Case - Strict Output Schema Violation (Expected to fail) ---"
# Create a schema that generates fields not in strict output schema
cat > "$CONFIG_DIR/invalid-output-schema.yaml" << 'EOF'
apiVersion: "konfigo/v1alpha1"

outputSchema:
  path: "schemas/output-schema-strict.json"
  strict: true

transform:
  - type: "setValue"
    path: "service.name"
    value: "test-service"
  - type: "setValue"
    path: "service.port"
    value: 8080
  - type: "setValue"
    path: "database.host"
    value: "localhost"
  - type: "setValue"
    path: "database.name"  
    value: "myapp"
  - type: "setValue"
    path: "features.cache"
    value: true
  - type: "setValue"
    path: "extra_field"
    value: "not-in-output-schema"
EOF

if $KONFIGO \
    -s "$INPUT_DIR/base-config.json" \
    -S "$CONFIG_DIR/invalid-output-schema.yaml" \
    -of "$OUTPUT_DIR/should-fail-output.yaml" 2>/dev/null; then
    echo "ERROR: Expected failure but command succeeded"
else
    echo "EXPECTED: Command failed as expected for strict output schema violation"
fi

# Test 12: Environment Variable Override with Schemas
echo -e "\n--- Test 12: Environment Variable Override with Schemas ---"
env "KONFIGO_KEY_service.name=env-override-service" \
    "KONFIGO_KEY_database.port=6543.0" \
    "KONFIGO_VAR_ENVIRONMENT=development" \
    $KONFIGO \
    -s "$INPUT_DIR/base-config.json" \
    -S "$CONFIG_DIR/env-override-schema.yaml" \
    -of "$OUTPUT_DIR/env-override.yaml"

echo -e "\n=== Schema Processing Integration Tests Completed ==="
echo "Output files generated in: $OUTPUT_DIR"
echo "Run 'validate.sh' to compare with expected outputs."
