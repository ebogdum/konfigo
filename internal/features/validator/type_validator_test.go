package validator

import (
	"strings"
	"testing"
)

func TestTypeValidator_NilValue(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(nil, "test.path", Rule{Type: "string"})
	if err == nil {
		t.Fatal("expected error for nil value, got nil")
	}
	if !strings.Contains(err.Error(), "got null") {
		t.Errorf("expected error message to contain 'got null', got: %s", err.Error())
	}
}

func TestTypeValidator_NoTypeRule(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate("anything", "test.path", Rule{Type: ""})
	if err != nil {
		t.Errorf("expected no error for empty type rule, got: %v", err)
	}
}

func TestTypeValidator_StringMatch(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate("hello", "test.path", Rule{Type: "string"})
	if err != nil {
		t.Errorf("expected no error for string value with string type, got: %v", err)
	}
}

func TestTypeValidator_StringMismatch(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(42, "test.path", Rule{Type: "string"})
	if err == nil {
		t.Fatal("expected error for int value with string type, got nil")
	}
}

func TestTypeValidator_NumberFromInt(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(42, "test.path", Rule{Type: "number"})
	if err != nil {
		t.Errorf("expected no error for int value with number type, got: %v", err)
	}
}

func TestTypeValidator_NumberFromFloat(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(3.14, "test.path", Rule{Type: "number"})
	if err != nil {
		t.Errorf("expected no error for float value with number type, got: %v", err)
	}
}

func TestTypeValidator_NumberMismatch(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate("not_a_number", "test.path", Rule{Type: "number"})
	if err == nil {
		t.Fatal("expected error for string value with number type, got nil")
	}
}

func TestTypeValidator_BoolMatch(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(true, "test.path", Rule{Type: "bool"})
	if err != nil {
		t.Errorf("expected no error for bool value with bool type, got: %v", err)
	}
}

func TestTypeValidator_NilNumberType(t *testing.T) {
	tv := &TypeValidator{}
	err := tv.Validate(nil, "db.port", Rule{Type: "number"})
	if err == nil {
		t.Fatal("expected error for nil value with number type, got nil")
	}
	if !strings.Contains(err.Error(), "got null") {
		t.Errorf("expected error to contain 'got null', got: %s", err.Error())
	}
}

// TestTypeValidator_IntegerAcrossParserWidths guards against a regression where
// the validator compared reflect.Kind names, so an "integer" rule passed for
// YAML input (which yields int) and failed for the identical JSON or TOML value
// (which yields int64) with "expected type integer, got int64".
func TestTypeValidator_IntegerAcrossParserWidths(t *testing.T) {
	tv := &TypeValidator{}
	values := []interface{}{
		int(8080),   // yaml.v3
		int64(8080), // encoding/json (normalized) and BurntSushi/toml
		int32(8080),
	}
	for _, ruleType := range []string{"int", "integer"} {
		for _, v := range values {
			if err := tv.Validate(v, "port", Rule{Type: ruleType}); err != nil {
				t.Errorf("Validate(%T(%v), type=%q) = %v, want nil", v, v, ruleType, err)
			}
		}
	}
}

// TestTypeValidator_IntegerRejectsNonIntegers ensures widening the accepted kinds
// did not make the integer rule accept everything.
func TestTypeValidator_IntegerRejectsNonIntegers(t *testing.T) {
	tv := &TypeValidator{}
	for _, v := range []interface{}{"8080", 3.14, true, []interface{}{1}, map[string]interface{}{"a": 1}} {
		if err := tv.Validate(v, "port", Rule{Type: "integer"}); err == nil {
			t.Errorf("Validate(%T(%v), type=integer) = nil, want an error", v, v)
		}
	}
}

// TestTypeValidator_ObjectAndArrayAliases covers the map/slice aliases across kinds.
func TestTypeValidator_ObjectAndArrayAliases(t *testing.T) {
	tv := &TypeValidator{}
	if err := tv.Validate(map[string]interface{}{"a": 1}, "p", Rule{Type: "object"}); err != nil {
		t.Errorf("object alias rejected a map: %v", err)
	}
	if err := tv.Validate([]interface{}{1, 2}, "p", Rule{Type: "array"}); err != nil {
		t.Errorf("array alias rejected a slice: %v", err)
	}
	if err := tv.Validate([]interface{}{1}, "p", Rule{Type: "object"}); err == nil {
		t.Error("object alias accepted a slice, want an error")
	}
}
