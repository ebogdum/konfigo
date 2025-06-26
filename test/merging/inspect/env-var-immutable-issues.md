# Configuration Merging Issues

## Environment Variable Override Issues

### Issue 1: KONFIGO_KEY_ vs Immutable Paths

**Expected Behavior (per documentation)**: 
- KONFIGO_KEY_ environment variables should override immutable paths
- Documentation states: "KONFIGO_KEY_... environment variables *can* override values at immutable paths"

**Actual Behavior**:
- KONFIGO_KEY_ environment variables are not overriding immutable paths
- Immutable paths from schema are preventing KONFIGO_KEY_ overrides

**Test Command**:
```bash
env KONFIGO_KEY_application.name=env-override-app KONFIGO_KEY_database.port=9999 \
./konfigo -s test/merging/input/base-config.json,test/merging/input/override-prod.json \
-S test/merging/config/schema-immutable.yaml -oj
```

**Expected Output**: 
- application.name should be "env-override-app" 
- database.port should be 9999

**Actual Output**:
- application.name remains "my-app" (from base, protected by immutable)
- database.port remains 5432 (from base, protected by immutable)

**Status**: This appears to be a bug or documentation discrepancy that needs investigation.

### Issue 2: Environment Variable Processing Order

**Investigation Needed**: 
- When are KONFIGO_KEY_ variables applied in the merge pipeline?
- Are they applied before or after immutable path checking?
- Should immutable paths only apply to file-based sources?

### Recommendation

Move these tests to `inspect/` until the environment variable vs immutable path behavior is clarified and potentially fixed in the codebase.
