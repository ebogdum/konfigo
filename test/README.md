# Konfigo Test Suite

This directory contains comprehensive test suites for the Konfigo configuration management tool.

## Structure

The test suite is organized into the following directories:

- `batch/` - Tests for batch processing and `konfigo_forEach` functionality
- `env-integration/` - Tests for environment variable integration
- `format-conversion/` - Tests for converting between different configuration formats
- `generators/` - Tests for configuration generation features
- `merging/` - Tests for configuration merging functionality
- `recursive-discovery/` - Tests for recursive configuration discovery
- `schema-integration/` - Tests for schema processing and validation
- `transformers/` - Tests for data transformation features
- `validators/` - Tests for configuration validation
- `variables/` - Tests for variable substitution

## Usage

### Running All Tests

To run all test suites:

```bash
./test.sh
```

This will:
- Check that the konfigo binary exists
- Run all individual test suites in sequence
- Provide a summary of passed/failed tests

### Validating All Outputs

To validate that all test outputs match expected results:

```bash
./validate.sh
```

This will:
- Compare all generated outputs with expected results
- Show detailed differences where files don't match
- Provide a summary of validation results

### Running Individual Test Suites

To run a specific test suite:

```bash
cd <test-directory>
./test.sh
```

To validate a specific test suite:

```bash
cd <test-directory>
./validate.sh
```

## Test Script Standards

All test and validation scripts follow consistent patterns:

### Test Scripts (`test.sh`)
- Use the common functions from `../common_functions.sh`
- Follow the same output formatting and error handling
- Use the `run_test()` function for individual tests
- Print comprehensive summaries

### Validation Scripts (`validate.sh`)
- Use the common functions from `../common_functions.sh`
- Show detailed differences when files don't match
- Use contextual diffs with line numbers
- Provide clear pass/fail indicators

## Common Functions

The `common_functions.sh` file provides:

- **Consistent colors and formatting** for all output
- **Standardized test execution** with `run_test()`
- **Detailed file comparison** with `compare_files()`
- **Directory comparison** with `compare_directories()`
- **Comprehensive validation** with `validate_all_outputs()`
- **Summary reporting** for tests and validations

## Adding New Tests

When adding new test suites:

1. Create a new directory under `test/`
2. Add both `test.sh` and `validate.sh` scripts
3. Source `../common_functions.sh` in both scripts
4. Use the standardized functions for consistency
5. Add the new directory to the main test and validation scripts

## Example Test Structure

```bash
#!/bin/bash
set -e

# Source common functions
source "../common_functions.sh"

# Setup test environment
setup_test_environment "My Test Suite"

# Run tests
run_test "Test description" \
    "$KONFIGO -s input.json -of output.json"

# Print summary
print_test_summary
exit $?
```

## Example Validation Structure

```bash
#!/bin/bash
set -e

# Source common functions
source "../common_functions.sh"

# Setup validation environment
setup_validation_environment "My Test Validation"

# Validate outputs
validate_all_outputs

# Print summary
print_validation_summary
exit $?
```

## Features

### Detailed Difference Reporting
When validations fail, the scripts show:
- Contextual diffs with line numbers
- Clear indication of which files differ
- First 20 lines of differences (configurable)
- Suggestions for reviewing changes

### Comprehensive Error Handling
- Checks for missing binaries and directories
- Clear error messages with suggested actions
- Proper exit codes for CI/CD integration
- Graceful handling of edge cases

### Consistent Output Format
- Color-coded results (✓ green for pass, ✗ red for fail)
- Standardized section headers and summaries
- Progress indicators for long-running operations
- Uniform emoji and formatting across all scripts
