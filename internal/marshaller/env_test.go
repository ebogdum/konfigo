package marshaller

import (
	"strings"
	"testing"
)

// TestENVMarshaller_NoLineInjection guards against a regression where non-string
// values were formatted with a raw %v and written unescaped. A newline inside a
// slice element terminated the line and injected an attacker-controlled
// KEY=VALUE pair into the generated .env file.
func TestENVMarshaller_NoLineInjection(t *testing.T) {
	em := &ENVMarshaller{}
	data := map[string]interface{}{
		"injected": []interface{}{"a\nEVIL=pwned", "b"},
	}

	out, err := em.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.HasPrefix(line, "EVIL=") {
			t.Errorf("env line injection: produced an unexpected variable line %q\nfull output:\n%s", line, out)
		}
	}

	// The value must survive as a single escaped line rather than a raw newline.
	if strings.Contains(string(out), "[a\nEVIL") {
		t.Errorf("raw newline was emitted inside a value:\n%s", out)
	}
}

// TestENVMarshaller_QuotesValuesWithSpaces ensures slice values, which render
// with spaces, are quoted so dotenv parsers read the whole value.
func TestENVMarshaller_QuotesValuesWithSpaces(t *testing.T) {
	em := &ENVMarshaller{}
	data := map[string]interface{}{
		"outputs": []interface{}{"stdout", "file"},
	}

	out, err := em.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if !strings.Contains(string(out), `OUTPUTS="[stdout file]"`) {
		t.Errorf("space-containing value was not quoted:\n%s", out)
	}
}

// TestENVMarshaller_LeavesSimpleScalarsBare ensures the escaping fix did not
// start quoting values that never needed it.
func TestENVMarshaller_LeavesSimpleScalarsBare(t *testing.T) {
	em := &ENVMarshaller{}
	data := map[string]interface{}{
		"port":    8080,
		"enabled": true,
		"name":    "plain",
	}

	out, err := em.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	for _, want := range []string{"PORT=8080", "ENABLED=true", "NAME=plain"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("expected bare %q in output:\n%s", want, out)
		}
	}
}

// TestENVMarshaller_KeyCollisionDetected ensures a nested map and a scalar that
// flatten to the same key are reported rather than silently overwriting.
func TestENVMarshaller_KeyCollisionDetected(t *testing.T) {
	em := &ENVMarshaller{}
	data := map[string]interface{}{
		"db_host": "scalar",
		"db":      map[string]interface{}{"host": "nested"},
	}

	if _, err := em.Marshal(data); err == nil {
		t.Error("expected a key collision error for db_host vs db.host, got nil")
	}
}
