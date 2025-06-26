# Changelog

## [Latest Commit] - 2025-06-27

### 🔧 **Release Script Enhancement**

#### **Modified Files**
- **Enhanced**: `release.sh` - Integrated automatic changelog extraction for GitHub releases

#### **New Features**
- **Added**: `extract_changelog_for_version()` function to automatically extract release notes from `CHANGELOG.md`
- **Enhanced**: Release creation process now automatically uses changelog content for GitHub releases
- **Improved**: Fallback mechanism when version-specific content isn't found in changelog
- **Added**: Automatic cleanup of temporary files during release process

#### **Configuration Changes**
- **Updated**: Default `notes_file` configuration to use `CHANGELOG.md`
- **Enhanced**: Release notes extraction supports various changelog formats (with/without brackets, version prefixes)

### 🏗️ **Internal Architecture Changes**

This section documents changes to the `internal` and `test` folders between the previous and current commit.

### 🏗️ **Internal Architecture Changes**

#### **New Packages and Modules**

##### **CLI Package (`internal/cli/`)**
- **Added**: `commands.go` - Command execution logic and coordination
- **Added**: `flags.go` - Flag definitions, parsing, and validation
- **Added**: `help.go` - Help text generation and display

##### **Configuration Package (`internal/config/`)**
- **Added**: `batch.go` - Batch processing configuration
- **Added**: `config.go` - Core configuration structures
- **Added**: `environment.go` - Environment handling

##### **Error Handling (`internal/errors/`)**
- **Added**: `errors.go` - Centralized error handling system

##### **Features Package (`internal/features/`)**
- **Added**: `generator/` - Code generation features
  - `concat.go` - Concatenation generator
  - `registry.go` - Generator registry
  - `types.go` - Generator type definitions
- **Added**: `input_schema/` - Input schema validation
  - `loader.go` - Schema loading functionality
  - `validator.go` - Schema validation logic
- **Added**: `transformer/` - Data transformation features
  - `add_key_prefix.go` - Key prefix transformation
  - `change_case.go` - Case transformation
  - `registry.go` - Transformer registry
  - `rename_key.go` - Key renaming transformation
  - `set_value.go` - Value setting transformation
  - `types.go` - Transformer type definitions
- **Added**: `validator/` - Validation engine
  - `engine.go` - Main validation engine
  - `number.go` - Number validation
  - `numeric_validator.go` - Numeric validation logic
  - `registry.go` - Validator registry
  - `string_validator.go` - String validation logic
  - `type_validator.go` - Type validation
  - `types.go` - Validator type definitions
- **Added**: `variables/` - Variable handling
  - `resolver.go` - Variable resolution logic (moved from `internal/schema/vars.go`)
  - `substitution.go` - Variable substitution
  - `types.go` - Variable type definitions

##### **Core Services**
- **Added**: `logger/logger.go` - Centralized logging system
- **Added**: `marshaller/` - Output formatting
  - `env.go` - Environment file marshalling
  - `json.go` - JSON marshalling
  - `marshaller.go` - Core marshalling interface
  - `registry.go` - Marshaller registry
  - `toml.go` - TOML marshalling
  - `yaml.go` - YAML marshalling
- **Added**: `merger/merger.go` - Configuration merging logic
- **Added**: `parser/` - Input parsing
  - `detector.go` - Format detection
  - `env.go` - Environment file parsing
  - `ini.go` - INI file parsing
  - `json.go` - JSON parsing
  - `parser.go` - Core parsing interface
  - `registry.go` - Parser registry
  - `toml.go` - TOML parsing
  - `yaml.go` - YAML parsing
- **Added**: `pipeline/` - Processing pipeline
  - `batch.go` - Batch processing pipeline
  - `coordinator.go` - Pipeline coordination
  - `optimized.go` - Optimized processing pipeline
  - `pipeline.go` - Core pipeline interface
  - `single.go` - Single file processing pipeline
- **Added**: `reader/` - File reading
  - `discovery.go` - File discovery logic (moved from `internal/loader/loader.go`)
  - `reader.go` - Core reading interface
  - `stream.go` - Stream reading functionality
- **Added**: `util/` - Utility functions
  - `type_inference.go` - Type inference utilities
  - `util.go` - General utility functions
- **Added**: `writer/` - Output writing
  - `directory.go` - Directory output writing
  - `target.go` - Target specification
  - `writer.go` - Core writing interface

#### **Modified Schema Package (`internal/schema/`)**
- **Removed**: `generator.go` - Functionality moved to `internal/features/generator/`
- **Removed**: `transformer.go` - Functionality moved to `internal/features/transformer/`
- **Removed**: `validator.go` - Functionality moved to `internal/features/validator/`
- **Modified**: `vars.go` → moved to `internal/features/variables/resolver.go`
- **Added**: `processor.go` - Schema processing logic
- **Modified**: `schema.go` - Updated schema handling

### 🧪 **Test Infrastructure Changes**

#### **Test Documentation**
- **Modified**: `test/README.md` - Updated test documentation

#### **Batch Processing Tests (`test/batch/`)**
- **Added**: `README.md` - Batch testing documentation
- **Added**: `config/` - Test configuration files
  - `deployment-schema.yaml` - Deployment schema configuration
  - `service-schema.yaml` - Service schema configuration
- **Added**: `expected/` - Expected output files
  - `complex/` - Complex scenario outputs
  - `deployments/` - Deployment scenario outputs
  - `envs/` - Environment configuration outputs
  - `services/` - Service configuration outputs
- **Added**: `input/` - Test input files
  - `base-config.json` - Base configuration
  - `empty-base.json` - Empty base configuration
- **Added**: `output/` - Actual test outputs
- **Added**: `test.sh` - Batch test execution script
- **Added**: `validate.sh` - Batch test validation script
- **Added**: `variables/` - Variable test files
  - Various environment and batch configuration files

#### **Environment Integration Tests (`test/env-integration/`)**
- **Added**: `README.md` - Environment integration test documentation
- **Added**: `config/` - Schema configuration files
- **Added**: `expected/` - Expected outputs for environment tests
- **Added**: `input/` - Test input files for environment scenarios
- **Added**: `output/` - Actual test outputs
- **Added**: `test.sh` - Environment integration test script
- **Added**: `validate.sh` - Environment integration validation script

#### **Format Conversion Tests (`test/format-conversion/`)**
- **Added**: `expected/` - Expected outputs for format conversion
- **Added**: `input/` - Input files in various formats (env, ini, json, toml, yaml)
- **Added**: `output/` - Actual conversion outputs
- **Added**: `test.sh` - Format conversion test script
- **Added**: `validate.sh` - Format conversion validation script

#### **Generator Tests (`test/generators/`)**
- **Added**: `README.md` - Generator test documentation
- **Added**: `config/` - Generator configuration files
  - Various schema files for testing generator functionality
  - Error scenario configurations
- **Added**: Variable test files

#### **Test Utilities**
- **Modified**: `test/common_functions.sh` - Updated common test functions

### 📋 **Summary of Changes**

#### **Added**
- **58 new files** in `internal/` package
- **200+ new test files** across multiple test suites
- Complete modularization of the codebase with feature-based architecture
- Comprehensive test coverage for batch processing, environment integration, format conversion, and generators

#### **Removed**
- **3 legacy files** from `internal/schema/` package:
  - `generator.go`
  - `transformer.go` 
  - `validator.go`

#### **Moved/Refactored**
- `internal/schema/vars.go` → `internal/features/variables/resolver.go`
- `internal/loader/loader.go` → `internal/reader/discovery.go`

### 🎯 **Key Improvements**

1. **Modular Architecture**: Complete restructuring into feature-based modules
2. **Enhanced Testing**: Comprehensive test suites for all major features
3. **Better Separation of Concerns**: Clear separation between CLI, processing, and feature logic
4. **Improved Maintainability**: Registry-based pattern for extensible functionality
5. **Robust Error Handling**: Centralized error management system
