# Changelog

## [2.0.3] - 2026-07-31

### 🔒 **Security**
- Resolved all 6 open Dependabot advisories in the documentation toolchain (`npm audit` now reports 0 vulnerabilities)
- Vite pinned to `6.4.3` via npm `overrides`, fixing `server.fs.deny` bypass on Windows alternate paths (high), path traversal in optimized deps `.map` handling (medium), and NTLMv2 hash disclosure through `launch-editor` UNC path handling on Windows (medium)
- Rollup pinned to `4.62.3` via npm `overrides`, fixing arbitrary file write via path traversal (high)
- PostCSS pinned to `8.5.25` via npm `overrides`, fixing XSS via unescaped `</style>` in stringify output and arbitrary `.map` file disclosure via attacker-controlled `sourceMappingURL` (high)
- `esbuild` bumped to `0.28.1`, fixing arbitrary file read when running the development server on Windows; the transitive copy pulled in by Vite now resolves to `0.25.12`, outside all affected ranges
- No vulnerability affected the shipped `konfigo` binary — every advisory above is confined to the VitePress build and dev-server toolchain
- `govulncheck` reports no vulnerabilities in the Go module or its dependencies

### 🔧 **Enhancements**
- `gopkg.in/ini.v1` updated from `1.67.1` to `1.67.3`

### 🏗️ **Internal Changes**
- Added an `overrides` block to `package.json` so transitive documentation dependencies stay on patched versions regardless of what VitePress `1.6.4` requests

### 📚 **Documentation**
- Rebuilt the published documentation site with the updated toolchain

---

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
