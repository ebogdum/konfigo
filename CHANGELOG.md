# Changelog

## v1.0.4

### 🚀 **New Features**
- **Added**: `addKeySuffix` transformer - Adds suffixes to all keys within a map object
  - **Usage**: Transforms map keys by appending specified suffix to each key name
  - **Fields**: `type: "addKeySuffix"`, `path: "path.to.map"`, `suffix: "_suffix"`
  - **Example**: Transform `{host: "localhost", port: 5432}` → `{host_prod: "localhost", port_prod: 5432}`
- **Added**: `deleteKey` transformer - Removes specified keys from configuration
  - **Usage**: Deletes configuration keys at specified paths (useful for removing sensitive data)
  - **Fields**: `type: "deleteKey"`, `path: "path.to.key"`
  - **Example**: Remove secrets, temporary values, or deprecated configuration keys
- **Added**: `trim` transformer - Trims whitespace or custom patterns from string values
  - **Usage**: Cleans up string values by removing unwanted characters from start/end
  - **Fields**: `type: "trim"`, `path: "path.to.string"`, `pattern: "characters"` (optional)
  - **Default**: Trims whitespace if no pattern specified
  - **Example**: `"  value  "` → `"value"` or `"---token---"` → `"token"` (with pattern: "-")
- **Added**: `replaceKey` transformer - Replaces value with content from another path
  - **Usage**: Takes value from target path, places it at destination path, then deletes target
  - **Fields**: `type: "replaceKey"`, `path: "destination.path"`, `target: "source.path"`
  - **Example**: Move temporary/staged values to their final configuration locations

### 🔧 **Enhancements**
- **Enhanced**: `changeCase` transformer now supports additional case formats
  - **Added**: `kebab` case support for kebab-case conversions
  - **Added**: `pascal` case support for PascalCase conversions  
  - **Supported formats**: upper, lower, snake, camel, kebab, pascal
- **Enhanced**: Transformer Definition structure with new fields
  - **Added**: `suffix` field for addKeySuffix transformer
  - **Added**: `pattern` field for trim transformer  
  - **Added**: `target` field for replaceKey transformer
  - **Improved**: Variable substitution now supports all new transformer fields

### 🐛 **Bug Fixes**
- _No bug fixes yet_

### 🏗️ **Internal Changes**
- _No internal changes yet_

### 🧪 **Tests**
- **Enhanced**: Transformer test suite to include new transformer types
  - **Added**: Test schemas for addKeySuffix, deleteKey, trim, and replaceKey transformers
  - **Added**: Test cases covering all new transformer functionality and edge cases
  - **Added**: Error handling tests for new transformers (missing paths, type mismatches)
  - **Added**: Variable substitution tests for new transformer fields
  - **Updated**: Test documentation to reflect new transformer coverage

### 📚 **Documentation**
- **Enhanced**: Transformation documentation with comprehensive coverage of new transformers
  - **Added**: Complete documentation for `addKeySuffix` transformer with examples
  - **Added**: Complete documentation for `deleteKey` transformer with use cases
  - **Added**: Complete documentation for `trim` transformer with pattern examples
  - **Added**: Complete documentation for `replaceKey` transformer with workflow examples
  - **Updated**: Transformation overview to include all eight available transformer types
  - **Enhanced**: Combined transformation examples showing new transformers in action

---
