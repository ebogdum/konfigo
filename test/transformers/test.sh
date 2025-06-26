#!/bin/bash

# Transformers & Data Transformation Test Suite
# Tests all transformer types: renameKey, changeCase, addKeyPrefix, setValue

set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Transformers & Data Transformation Test Suite"

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

# === RENAME KEY TESTS ===
echo "--- Rename Key Tests ---"

run_test "rename-basic-json" "base-config.json" "schema-rename-basic.json" "" "false"
run_test "rename-basic-yaml" "base-config.yaml" "schema-rename-basic.yaml" "" "false"
run_test "rename-basic-toml" "base-config.toml" "schema-rename-basic.yaml" "" "false"

# === CHANGE CASE TESTS ===
echo "--- Change Case Tests ---"

run_test "changecase-basic-json" "base-config.json" "schema-changecase-basic.yaml" "" "false"
run_test "changecase-basic-yaml" "base-config.yaml" "schema-changecase-basic.yaml" "" "false"
run_test "changecase-basic-toml" "base-config.toml" "schema-changecase-basic.yaml" "" "false"

run_test "changecase-all-types-json" "base-config.json" "schema-changecase-all-types.yaml" "" "false"
run_test "changecase-all-types-yaml" "base-config.yaml" "schema-changecase-all-types.yaml" "" "false"
run_test "changecase-all-types-toml" "base-config.toml" "schema-changecase-all-types.yaml" "" "false"

# === ADD KEY PREFIX TESTS ===
echo "--- Add Key Prefix Tests ---"

run_test "addprefix-basic-json" "base-config.json" "schema-addprefix-basic.yaml" "" "false"
run_test "addprefix-basic-yaml" "base-config.yaml" "schema-addprefix-basic.yaml" "" "false"
run_test "addprefix-basic-toml" "base-config.toml" "schema-addprefix-basic.yaml" "" "false"

run_test "addprefix-vars-json" "base-config.json" "schema-addprefix-vars.yaml" "" "false"
run_test "addprefix-vars-yaml" "base-config.yaml" "schema-addprefix-vars.yaml" "" "false"
run_test "addprefix-vars-toml" "base-config.toml" "schema-addprefix-vars.yaml" "" "false"

# === SET VALUE TESTS ===
echo "--- Set Value Tests ---"

run_test "setvalue-basic-json" "base-config.json" "schema-setvalue-basic.yaml" "" "false"
run_test "setvalue-basic-yaml" "base-config.yaml" "schema-setvalue-basic.yaml" "" "false"
run_test "setvalue-basic-toml" "base-config.toml" "schema-setvalue-basic.yaml" "" "false"

run_test "setvalue-vars-json" "base-config.json" "schema-setvalue-vars.yaml" "" "false"
run_test "setvalue-vars-yaml" "base-config.yaml" "schema-setvalue-vars.yaml" "" "false"
run_test "setvalue-vars-toml" "base-config.toml" "schema-setvalue-vars.yaml" "" "false"

run_test "setvalue-complex-json" "base-config.json" "schema-setvalue-complex.yaml" "" "false"
run_test "setvalue-complex-yaml" "base-config.yaml" "schema-setvalue-complex.yaml" "" "false"
run_test "setvalue-complex-toml" "base-config.toml" "schema-setvalue-complex.yaml" "" "false"

# === COMBINED TRANSFORMATIONS TESTS ===
echo "--- Combined Transformations Tests ---"

run_test "combined-json" "base-config.json" "schema-combined.yaml" "" "false"
run_test "combined-yaml" "base-config.yaml" "schema-combined.yaml" "" "false"
run_test "combined-toml" "base-config.toml" "schema-combined.yaml" "" "false"

# === ERROR HANDLING TESTS ===
echo "--- Error Handling Tests ---"

run_test "error-rename-missing" "base-config.json" "schema-error-rename-missing.yaml" "" "true"
run_test "error-changecase-nonstring" "base-config.json" "schema-error-changecase-nonstring.yaml" "" "true"
run_test "error-changecase-invalid" "base-config.json" "schema-error-changecase-invalid.yaml" "" "true"
run_test "error-addprefix-nonmap" "base-config.json" "schema-error-addprefix-nonmap.yaml" "" "true"

# === ENVIRONMENT VARIABLE TESTS ===
echo "--- Environment Variable Tests ---"

echo "TEST $((total_tests + 1)): env-var-substitution"
total_tests=$((total_tests + 1))
echo "  Testing KONFIGO_VAR_ environment variable substitution in transformers"

# Set environment variables
export KONFIGO_VAR_ENV_PREFIX="test"

env_test_passed=true
cmd="$KONFIGO -s $INPUT_DIR/base-config.yaml -S $CONFIG_DIR/schema-combined.yaml -oy -of $OUTPUT_DIR/env-substitution-test.yaml"

if $cmd >/dev/null 2>&1; then
    echo "  ✅ PASS: Environment variable substitution test completed"
    passed_tests=$((passed_tests + 1))
else
    echo "  ❌ FAIL: Environment variable substitution test failed"
    failed_tests=$((failed_tests + 1))
fi

# Clean up environment variables
unset KONFIGO_VAR_ENV_PREFIX
echo

# === EDGE CASE TESTS ===
echo "--- Edge Case Tests ---"

echo "TEST $((total_tests + 1)): stdin-input-transform"
total_tests=$((total_tests + 1))
echo "  Testing transformer with stdin input"

if echo '{"user":{"name":"StdinUser","id":999},"settings":{"timeout":60}}' | \
   $KONFIGO -s - -sj -S $CONFIG_DIR/schema-rename-basic.yaml -oy -of $OUTPUT_DIR/stdin-transform-test.yaml >/dev/null 2>&1; then
    echo "  ✅ PASS: Stdin transformer test"
    passed_tests=$((passed_tests + 1))
else
    echo "  ❌ FAIL: Stdin transformer test"
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
