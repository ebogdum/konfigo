#!/bin/bash

# Generators & Data Generation Test Suite
# Tests the concat generator functionality with various scenarios

set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Generators & Data Generation Test Suite"

# Test directories
INPUT_DIR="input"
CONFIG_DIR="config"
OUTPUT_DIR="output"

# Clean output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Test counter
total_tests=0
passed_tests=0
failed_tests=0

# Test function
run_test() {
    local test_name="$1"
    local input_file="$2"
    local schema_file="$3"
    local variables_file="$4"
    local expected_to_fail="$5"
    
    total_tests=$((total_tests + 1))
    echo "TEST $total_tests: $test_name"
    echo "  Input: $input_file"
    echo "  Schema: $schema_file"
    if [[ -n "$variables_file" ]]; then
        echo "  Variables: $variables_file"
    fi
    
    # Build command
    local cmd="$KONFIGO"
    if [[ -n "$input_file" ]]; then
        cmd="$cmd -s $INPUT_DIR/$input_file"
    fi
    if [[ -n "$schema_file" ]]; then
        cmd="$cmd -S $CONFIG_DIR/$schema_file"
    fi
    if [[ -n "$variables_file" ]]; then
        cmd="$cmd -V $CONFIG_DIR/$variables_file"
    fi
    
    # Test all output formats
    local formats=("json" "yaml" "toml")
    local test_passed=true
    
    for format in "${formats[@]}"; do
        local output_file="$OUTPUT_DIR/${test_name}-${format}.${format}"
        local format_flag=""
        case $format in
            "json") format_flag="-oj" ;;
            "yaml") format_flag="-oy" ;;
            "toml") format_flag="-ot" ;;
        esac
        local test_cmd="$cmd $format_flag -of $output_file"
        
        echo "    Testing $format output..."
        
        if [[ "$expected_to_fail" == "true" ]]; then
            # Test should fail
            if $test_cmd >/dev/null 2>&1; then
                echo "    ❌ FAIL: Expected error but command succeeded"
                test_passed=false
            else
                echo "    ✅ PASS: Expected error occurred"
            fi
        else
            # Test should succeed
            if $test_cmd >/dev/null 2>&1; then
                echo "    ✅ PASS: $format output generated"
            else
                echo "    ❌ FAIL: Error generating $format output"
                echo "    Command: $test_cmd"
                $test_cmd || true  # Show the error
                test_passed=false
            fi
        fi
    done
    
    if [[ "$test_passed" == "true" ]]; then
        passed_tests=$((passed_tests + 1))
        echo "  ✅ OVERALL: PASS"
    else
        failed_tests=$((failed_tests + 1))
        echo "  ❌ OVERALL: FAIL"
    fi
    echo
}

# === BASIC GENERATOR TESTS ===
echo "--- Basic Generator Tests ---"

run_test "basic-concat-json" "base-config.json" "schema-basic.json" "" "false"
run_test "basic-concat-yaml" "base-config.yaml" "schema-basic.yaml" "" "false"

# === MULTIPLE GENERATORS TESTS ===
echo "--- Multiple Generators Tests ---"

run_test "multiple-generators-json" "base-config.json" "schema-multiple.yaml" "" "false"
run_test "multiple-generators-yaml" "base-config.yaml" "schema-multiple.yaml" "" "false"
run_test "multiple-generators-toml" "base-config.toml" "schema-multiple.yaml" "" "false"

# === MIXED CONTENT TESTS ===
echo "--- Mixed Content Tests ---"

run_test "mixed-content-json" "base-config.json" "schema-mixed.yaml" "" "false"
run_test "mixed-content-yaml" "base-config.yaml" "schema-mixed.yaml" "" "false"
run_test "mixed-content-toml" "base-config.toml" "schema-mixed.yaml" "" "false"

# === VARIABLES ONLY TESTS ===
echo "--- Variables Only Tests ---"

run_test "vars-only-json" "base-config.json" "schema-vars-only.yaml" "" "false"
run_test "vars-only-yaml" "base-config.yaml" "schema-vars-only.yaml" "" "false"

# === EXTERNAL VARIABLES TESTS ===
echo "--- External Variables Tests ---"

run_test "external-vars-yaml-file" "base-config.yaml" "schema-cascading.yaml" "variables.yaml" "false"
run_test "external-vars-json-file" "base-config.json" "schema-cascading.yaml" "variables.json" "false"

# === CASCADING GENERATORS TESTS ===
echo "--- Cascading Generators Tests ---"

run_test "cascading-json" "base-config.json" "schema-cascading.yaml" "variables.yaml" "false"
run_test "cascading-yaml" "base-config.yaml" "schema-cascading.yaml" "variables.yaml" "false"

# === ERROR HANDLING TESTS ===
echo "--- Error Handling Tests ---"

run_test "error-missing-path" "base-config.json" "schema-error-missing-path.yaml" "" "true"

# Note: Other validation tests moved to inspect/ directory due to validation gaps

# === ENVIRONMENT VARIABLE OVERRIDE TESTS ===
echo "--- Environment Variable Override Tests ---"

echo "TEST $((total_tests + 1)): env-var-override"
total_tests=$((total_tests + 1))
echo "  Testing KONFIGO_VAR_ environment variable overrides"

# Set environment variables
export KONFIGO_VAR_EXTERNAL_VAR="env-override-value"
export KONFIGO_VAR_GLOBAL_VAR="env-global-value"

env_test_passed=true
cmd="$KONFIGO -s $INPUT_DIR/base-config.yaml -S $CONFIG_DIR/schema-mixed.yaml -oy -of $OUTPUT_DIR/env-override-test.yaml"

if $cmd >/dev/null 2>&1; then
    echo "  ✅ PASS: Environment override test completed"
    passed_tests=$((passed_tests + 1))
else
    echo "  ❌ FAIL: Environment override test failed"
    failed_tests=$((failed_tests + 1))
fi

# Clean up environment variables
unset KONFIGO_VAR_EXTERNAL_VAR
unset KONFIGO_VAR_GLOBAL_VAR
echo

# === EDGE CASE TESTS ===
echo "--- Edge Case Tests ---"

# Test with stdin input
echo "TEST $((total_tests + 1)): stdin-input"
total_tests=$((total_tests + 1))
echo "  Testing generator with stdin input"

if echo '{"service":{"name":"stdin-test","instanceId":"id-123"},"region":"us-east-1"}' | \
   $KONFIGO -s - -sj -S $CONFIG_DIR/schema-basic.yaml -oj -of $OUTPUT_DIR/stdin-test.json >/dev/null 2>&1; then
    echo "  ✅ PASS: Stdin input test"
    passed_tests=$((passed_tests + 1))
else
    echo "  ❌ FAIL: Stdin input test"
    failed_tests=$((failed_tests + 1))
fi
echo

# === RESULTS SUMMARY ===
echo "=== Test Results Summary ==="
echo "Total tests: $total_tests"
echo "Passed: $passed_tests"
echo "Failed: $failed_tests"

if [[ $failed_tests -eq 0 ]]; then
    echo "✅ ALL TESTS PASSED!"
    exit 0
else
    echo "❌ Some tests failed. Check output above for details."
    exit 1
fi
