#!/bin/bash

# Format Conversion Test Script
# Tests conversion between all supported formats

set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "Format Conversion Test Suite"

INPUT_DIR="input"

# Test all format conversions
declare -a input_formats=("json" "yaml" "toml" "env" "ini")
declare -a output_formats=("json" "yaml" "toml" "env")

for input_format in "${input_formats[@]}"; do
    for output_format in "${output_formats[@]}"; do
        input_file="$INPUT_DIR/input.$input_format"
        output_file="$OUTPUT_DIR/${input_format}-to-${output_format}.${output_format}"
        
        # Skip if input file doesn't exist
        if [[ ! -f "$input_file" ]]; then
            continue
        fi
        
        # Determine input and output flags
        input_flag=""
        case "$input_format" in
            "env") input_flag="-se" ;;
        esac
        
        output_flag=""
        case "$output_format" in
            "json") output_flag="-oj" ;;
            "yaml") output_flag="-oy" ;;
            "toml") output_flag="-ot" ;;
            "env") output_flag="-oe" ;;
        esac
        
        # Run the test
        run_test "${input_format} to ${output_format} conversion" \
            "$KONFIGO -s $input_file $input_flag $output_flag -of $output_file"
    done
done

# Print test summary
print_test_summary
exit $?
