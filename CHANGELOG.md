# Changelog

## [2.0.2] - 2026-04-02

### 🔧 **Enhancements**
- Type validator now recognizes common aliases: `boolean`, `integer`, `array`, `object`, `float`, `double`
- Immutable path protection now extends to all child paths (marking `database` immutable also protects `database.host`, `database.port`, etc.)
- Generators and transformers can no longer overwrite immutable paths; values are snapshotted and restored if modified
- Warning logged when `${ITEM_FILE_BASENAME}` is used in `filenamePattern` with `items` mode (resolves to empty string)
- Debug log when case-insensitive merge changes key casing
- Large array union merge (>1,000 elements per side) automatically skips deduplication for performance
- Regex validation cache bounded to 500 entries with safe eviction
- Duplicate variable names in schema `vars` section are now detected and rejected
- CLI uses `flag.NewFlagSet` instead of global state for better testability

### 🐛 **Bug Fixes**
- Fixed path traversal bypass in batch filename patterns via schema variable values
- Fixed symlink bypass in `itemFile` path containment checks (now resolves symlinks before validation)
- Fixed `splitParams` out-of-bounds panic on short format strings
- Fixed silent fallback to predictable random seed when crypto/rand fails (now returns error)
- Fixed nil panic in variable substitution when config is nil
- Fixed concat generator false-positive unresolved placeholder detection on `${VAR}` patterns
- Fixed sequential ID generator polluting user config with `_internal` namespace
- Fixed `resolveItemFilePath` rejecting absolute paths when no `-V` flag is provided
- Fixed regex cache eviction data race under concurrent access

### 🏗️ **Internal Changes**
- Added 50 MiB file size limit on configuration file reads to prevent OOM
- Removed dead code: `StringsPool`, `ReadFileBuffered`, `ReadFileStream`, `ReadFiles`, `FileExists`, `WriteMultipleFiles`, `ValidateOutputConfiguration`
- Removed unnecessary `Coordinator` abstraction layer
- Renamed logger `Init` parameter from `verbose` to `debug` for clarity

### 🧪 **Tests**
- Fixed schema-integration test schema ref paths to use correct relative paths from schema directory
- Updated `immutable-paths.yaml` expected output to reflect correct immutable path protection behavior

### 📚 **Documentation**
- Updated type validator docs with common type aliases
- Updated immutable fields docs with child path protection and generator/transformer behavior
- Added relative path resolution note for inputSchema/outputSchema
- Added duplicate variable name detection warning
- Added case-insensitive merge key casing note
- Added array merge performance note for large arrays
- Added ITEM_FILE_BASENAME warning for items mode
- Updated file size limits and performance notes in reference docs

---
