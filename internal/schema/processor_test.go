package schema

import (
	"testing"

	"konfigo/internal/features/generator"
)

// TestProcess_ImmutableMapPathProtectedFromGenerator guards against a regression
// where the immutable snapshot stored the live value instead of a deep copy.
// Maps are reference types, so the snapshot aliased the config: a generator
// mutating a child key changed both sides and the equality check compared the
// value against itself, silently skipping the restore. Protection then only
// worked for scalars, while the documented guidance is to mark map-valued paths
// such as "database.credentials" immutable.
func TestProcess_ImmutableMapPathProtectedFromGenerator(t *testing.T) {
	config := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "ORIGINAL-HOST",
			"port": 5432,
		},
		"scalar": "ORIGINAL-SCALAR",
	}

	sch := &Schema{
		Immutable: []string{"database", "scalar"},
		Generators: []generator.Definition{
			{
				Type:       "concat",
				TargetPath: "database.host",
				Format:     "HIJACKED-{X}",
				Sources:    map[string]string{"X": "database.port"},
			},
			{
				Type:       "concat",
				TargetPath: "scalar",
				Format:     "HIJACKED-{X}",
				Sources:    map[string]string{"X": "database.port"},
			},
		},
	}

	out, err := Process(config, sch, nil, nil)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	db, ok := out["database"].(map[string]interface{})
	if !ok {
		t.Fatalf("database is not a map: %#v", out["database"])
	}
	if got := db["host"]; got != "ORIGINAL-HOST" {
		t.Errorf("immutable map-valued path was modified by a generator: database.host = %v, want %q",
			got, "ORIGINAL-HOST")
	}
	if got := out["scalar"]; got != "ORIGINAL-SCALAR" {
		t.Errorf("immutable scalar path was modified by a generator: scalar = %v, want %q",
			got, "ORIGINAL-SCALAR")
	}
}

// TestProcess_NonImmutableGeneratorStillApplies ensures the deep-copy snapshot
// did not turn every generator into a no-op.
func TestProcess_NonImmutableGeneratorStillApplies(t *testing.T) {
	config := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "ORIGINAL-HOST",
			"port": 5432,
		},
	}

	sch := &Schema{
		Generators: []generator.Definition{
			{
				Type:       "concat",
				TargetPath: "database.host",
				Format:     "REWRITTEN-{X}",
				Sources:    map[string]string{"X": "database.port"},
			},
		},
	}

	out, err := Process(config, sch, nil, nil)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	db := out["database"].(map[string]interface{})
	if got := db["host"]; got != "REWRITTEN-5432" {
		t.Errorf("generator did not apply to a mutable path: database.host = %v, want %q",
			got, "REWRITTEN-5432")
	}
}
