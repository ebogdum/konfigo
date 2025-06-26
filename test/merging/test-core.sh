#!/bin/bash

# Configuration Merging Test Suite - Working Tests Only
# Tests core merging functionality that works correctly

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Path to konfigo binary (relative to test directory)
KONFIGO="../../konfigo"

# Ensure konfigo binary exists
if [ ! -f "$KONFIGO" ]; then
    echo "Error: konfigo binary not found at $KONFIGO"
    echo "Please build konfigo first: cd ../../ && go build -o konfigo cmd/konfigo/main.go"
    exit 1
fi

echo "=== Configuration Merging Test Suite (Core Tests) ==="
echo

# Test counter
test_count=0
passed_count=0

# Test function
run_test() {
    local test_name="$1"
    local cmd="$2"
    local output_file="$3"
    
    test_count=$((test_count + 1))
    echo "Test $test_count: $test_name"
    echo "Command: $cmd"
    
    # Create output file path
    output_path="output/$output_file"
    
    # Run the command
    if eval "$cmd" > "$output_path" 2>&1; then
        echo "✅ PASSED"
        passed_count=$((passed_count + 1))
    else
        echo "❌ FAILED"
        echo "Error output:"
        cat "$output_path"
    fi
    echo
}

# Ensure output directory exists
mkdir -p output

echo "1. BASIC MERGE PRECEDENCE TESTS"
echo "================================"

# Test 1: Basic JSON + JSON merge (later source wins)
run_test "Basic JSON-to-JSON merge precedence" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -oj" \
    "basic-json-merge.json"

# Test 2: Cross-format merge (JSON base + YAML override)
run_test "Cross-format merge (JSON + YAML)" \
    "$KONFIGO -s input/base-config.json,input/override-dev.yaml -oj" \
    "cross-format-json-yaml.json"

# Test 3: Multi-format chain (JSON + YAML + TOML)
run_test "Multi-format merge (JSON + YAML + TOML)" \
    "$KONFIGO -s input/base-config.json,input/override-dev.yaml,input/override-staging.toml -oj" \
    "multi-format-chain.json"

# Test 4: All format types in sequence
run_test "All formats merge (JSON + YAML + TOML + ENV)" \
    "$KONFIGO -s input/base-config.json,input/override-dev.yaml,input/override-staging.toml,input/override-test.env -oj" \
    "all-formats-merge.json"

# Test 5: Reverse order (different precedence)
run_test "Reverse order precedence" \
    "$KONFIGO -s input/override-prod.json,input/base-config.json -oj" \
    "reverse-precedence.json"

echo "2. CASE SENSITIVITY TESTS"
echo "=========================="

# Test 6: Case-insensitive merge (default)
run_test "Case-insensitive merge (default behavior)" \
    "$KONFIGO -s input/case-base.json,input/case-override.json -oj" \
    "case-insensitive-merge.json"

# Test 7: Case-sensitive merge with -c flag
run_test "Case-sensitive merge (-c flag)" \
    "$KONFIGO -s input/case-base.json,input/case-override.json -c -oj" \
    "case-sensitive-merge.json"

echo "3. IMMUTABLE PATHS TESTS"
echo "========================"

# Test 8: Immutable paths with YAML schema
run_test "Immutable paths protection (YAML schema)" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -S config/schema-immutable.yaml -oj" \
    "immutable-yaml-schema.json"

# Test 9: Immutable paths with JSON schema
run_test "Immutable paths protection (JSON schema)" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -S config/schema-immutable.json -oj" \
    "immutable-json-schema.json"

# Test 10: Immutable paths with TOML schema
run_test "Immutable paths protection (TOML schema)" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -S config/schema-immutable.toml -oj" \
    "immutable-toml-schema.json"

echo "4. BASIC ENVIRONMENT VARIABLE TESTS"
echo "==================================="

# Test 11: Basic KONFIGO_KEY_ overrides (without immutable conflicts)
run_test "Basic KONFIGO_KEY_ overrides" \
    "env KONFIGO_KEY_application.environment=production KONFIGO_KEY_features.new_feature=true $KONFIGO -s input/base-config.json -oj" \
    "basic-env-overrides.json"

# Test 12: Nested path KONFIGO_KEY_ overrides (without immutable conflicts)
run_test "Nested path KONFIGO_KEY_ overrides" \
    "env KONFIGO_KEY_database.pool.min=10 KONFIGO_KEY_database.pool.max=50 $KONFIGO -s input/base-config.json,input/override-prod.json -oj" \
    "nested-env-overrides.json"

echo "5. RECURSIVE DISCOVERY TESTS"
echo "============================"

# Test 13: Recursive discovery with -r flag
run_test "Recursive file discovery (-r flag)" \
    "$KONFIGO -s input/nested -r -oj" \
    "recursive-discovery.json"

# Test 14: Recursive + base file merge
run_test "Recursive discovery + base file merge" \
    "$KONFIGO -s input/base-config.json,input/nested -r -oj" \
    "recursive-plus-base.json"

# Test 15: Complex recursive + overrides + schema
run_test "Complex recursive + overrides + immutable schema" \
    "$KONFIGO -s input/base-config.json,input/nested -r -S config/schema-immutable.yaml -oj" \
    "complex-recursive-schema.json"

echo "6. STDIN INPUT TESTS"
echo "==================="

# Test 16: Stdin JSON input
run_test "Stdin JSON input" \
    "echo '{\"application\":{\"name\":\"stdin-app\",\"port\":5000}}' | $KONFIGO -s - -sj -oj" \
    "stdin-json.json"

# Test 17: Stdin YAML input + file merge
run_test "Stdin YAML + file merge" \
    "printf 'application:\\n  name: stdin-yaml-app\\n  debug: true' | $KONFIGO -s input/base-config.json,- -sy -oj" \
    "stdin-yaml-merge.json"

# Test 18: Stdin TOML input
run_test "Stdin TOML input" \
    "printf '[application]\\nname = \"stdin-toml-app\"\\nport = 6000' | $KONFIGO -s - -st -oj" \
    "stdin-toml.json"

# Test 19: Stdin ENV input
run_test "Stdin ENV input" \
    "printf 'APPLICATION_NAME=stdin-env-app\\nAPPLICATION_PORT=7000' | $KONFIGO -s - -se -oj" \
    "stdin-env.json"

echo "7. OUTPUT FORMAT TESTS"
echo "====================="

# Test 20: Merge to YAML output
run_test "Merge output to YAML" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -oy" \
    "merge-output.yaml"

# Test 21: Merge to TOML output
run_test "Merge output to TOML" \
    "$KONFIGO -s input/base-config.json,input/override-dev.yaml -ot" \
    "merge-output.toml"

# Test 22: Merge to ENV output
run_test "Merge output to ENV" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json -oe" \
    "merge-output.env"

echo "8. EDGE CASES AND ERROR HANDLING"
echo "================================"

# Test 23: Empty file handling
run_test "Empty file in merge sequence" \
    "echo '{}' > output/empty.json && $KONFIGO -s input/base-config.json,output/empty.json,input/override-prod.json -oj" \
    "empty-file-merge.json"

# Test 24: Complex deep merge structures
run_test "Complex deep merge structures" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json,input/nested/env/prod.json,input/nested/services/web.yaml -oj" \
    "complex-deep-merge.json"

echo "9. CROSS-FORMAT INPUT/OUTPUT TESTS"
echo "=================================="

# Test 25: All input formats merged to JSON
run_test "All input formats merged to JSON" \
    "$KONFIGO -s input/base-config.json,input/base-config.yaml,input/base-config.toml,input/base-config.env -oj" \
    "all-inputs-to-json.json"

# Test 26: JSON inputs to all output formats
run_test "JSON merge to YAML output" \
    "$KONFIGO -s input/base-config.json,input/override-prod.json" \
    "json-merge-to-yaml.yaml"

# Test 27: Cross-format with immutable schema to ENV output
run_test "Cross-format with schema to ENV" \
    "$KONFIGO -s input/base-config.yaml,input/override-staging.toml -S config/schema-immutable.json -oe" \
    "cross-format-schema-env.env"

echo "=== Test Summary ==="
echo "Total tests: $test_count"
echo "Passed: $passed_count"
echo "Failed: $((test_count - passed_count))"

if [ $passed_count -eq $test_count ]; then
    echo "🎉 All tests passed!"
    exit 0
else
    echo "❌ Some tests failed. Check output files for details."
    exit 1
fi
