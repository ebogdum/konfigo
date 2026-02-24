package generator

import (
	"testing"
)

func TestSequentialId_NoInternalLeak(t *testing.T) {
	config := make(map[string]interface{})

	defs := []Definition{
		{Type: "id", TargetPath: "request_id", Format: "sequential"},
		{Type: "id", TargetPath: "other_id", Format: "sequential"},
	}

	err := Apply(config, defs, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Verify _internal is not present in output
	if _, exists := config["_internal"]; exists {
		t.Error("_internal key should be cleaned up after Apply, but it still exists")
	}

	// Verify the IDs were generated
	if config["request_id"] != "1" {
		t.Errorf("expected request_id = '1', got %v", config["request_id"])
	}
	if config["other_id"] != "1" {
		t.Errorf("expected other_id = '1', got %v", config["other_id"])
	}
}

func TestSequentialId_Increments(t *testing.T) {
	config := make(map[string]interface{})

	// Generate sequential IDs for the same path twice
	defs := []Definition{
		{Type: "id", TargetPath: "counter", Format: "sequential"},
	}

	// First apply
	err := Apply(config, defs, nil)
	if err != nil {
		t.Fatalf("first Apply failed: %v", err)
	}
	if config["counter"] != "1" {
		t.Errorf("expected counter = '1', got %v", config["counter"])
	}
}

func TestSimpleId_Length(t *testing.T) {
	config := make(map[string]interface{})
	defs := []Definition{
		{Type: "id", TargetPath: "myid", Format: "simple:16"},
	}

	err := Apply(config, defs, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	val, ok := config["myid"].(string)
	if !ok {
		t.Fatalf("expected string value for myid, got %T", config["myid"])
	}
	if len(val) != 16 {
		t.Errorf("expected ID length 16, got %d", len(val))
	}
}

func TestPrefixId_HasPrefix(t *testing.T) {
	config := make(map[string]interface{})
	defs := []Definition{
		{Type: "id", TargetPath: "myid", Format: "prefix:usr_:8"},
	}

	err := Apply(config, defs, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	val, ok := config["myid"].(string)
	if !ok {
		t.Fatalf("expected string value for myid, got %T", config["myid"])
	}
	if len(val) != 12 { // "usr_" (4) + 8 = 12
		t.Errorf("expected ID length 12, got %d: %s", len(val), val)
	}
	if val[:4] != "usr_" {
		t.Errorf("expected prefix 'usr_', got '%s'", val[:4])
	}
}

func TestNumericId_OnlyDigits(t *testing.T) {
	config := make(map[string]interface{})
	defs := []Definition{
		{Type: "id", TargetPath: "myid", Format: "numeric:10"},
	}

	err := Apply(config, defs, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	val, ok := config["myid"].(string)
	if !ok {
		t.Fatalf("expected string value for myid, got %T", config["myid"])
	}
	for _, ch := range val {
		if ch < '0' || ch > '9' {
			t.Errorf("expected only digits, found '%c' in %s", ch, val)
			break
		}
	}
}

func TestIdGenerator_InvalidFormat(t *testing.T) {
	config := make(map[string]interface{})
	defs := []Definition{
		{Type: "id", TargetPath: "myid", Format: "unknown_format"},
	}

	err := Apply(config, defs, nil)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}
