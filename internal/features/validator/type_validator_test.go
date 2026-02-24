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
