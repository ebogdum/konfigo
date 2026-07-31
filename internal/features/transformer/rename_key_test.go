package transformer

import "testing"

// TestRenameKey_RejectsOverlappingPaths guards against silent data loss from
// renames where one path is nested inside the other. Renaming a -> a.b used to
// empty the config entirely, and a.b -> a used to drop the parent's other keys,
// both without any error.
func TestRenameKey_RejectsOverlappingPaths(t *testing.T) {
	tr := &RenameKeyTransformer{}

	cases := []struct {
		name string
		from string
		to   string
	}{
		{"into own descendant", "a", "a.b"},
		{"onto own ancestor", "a.b", "a"},
		{"deep descendant", "app", "app.db.host"},
		{"deep ancestor", "app.db.host", "app"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tr.ValidateDefinition(Definition{From: tc.from, To: tc.to})
			if err == nil {
				t.Errorf("ValidateDefinition(from=%q, to=%q) = nil, want an error", tc.from, tc.to)
			}
		})
	}
}

// TestRenameKey_AllowsNonOverlappingPaths ensures the new guard did not reject
// ordinary renames, including ones that merely share a prefix string.
func TestRenameKey_AllowsNonOverlappingPaths(t *testing.T) {
	tr := &RenameKeyTransformer{}

	cases := []struct {
		name string
		from string
		to   string
	}{
		{"sibling leaves", "a.b", "a.c"},
		{"unrelated roots", "old", "new"},
		{"shared name prefix but distinct paths", "app", "application"},
		{"nested to unrelated", "app.db.host", "legacy.host"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tr.ValidateDefinition(Definition{From: tc.from, To: tc.to}); err != nil {
				t.Errorf("ValidateDefinition(from=%q, to=%q) = %v, want nil", tc.from, tc.to, err)
			}
		})
	}
}

// TestReplaceKey_RejectsOverlappingPaths covers the same overlap hazard in
// replaceKey, which writes to 'path' and then deletes 'target'.
func TestReplaceKey_RejectsOverlappingPaths(t *testing.T) {
	tr := &ReplaceKeyTransformer{}

	for _, tc := range []struct{ name, path, target string }{
		{"path inside target", "a.b", "a"},
		{"target inside path", "a", "a.b"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tr.ValidateDefinition(Definition{Path: tc.path, Target: tc.target}); err == nil {
				t.Errorf("ValidateDefinition(path=%q, target=%q) = nil, want an error", tc.path, tc.target)
			}
		})
	}
}

// TestReplaceKey_AllowsNonOverlappingPaths ensures ordinary replacements still validate.
func TestReplaceKey_AllowsNonOverlappingPaths(t *testing.T) {
	tr := &ReplaceKeyTransformer{}

	for _, tc := range []struct{ name, path, target string }{
		{"siblings", "a.b", "a.c"},
		{"unrelated", "dest", "src"},
		{"shared name prefix", "app", "application"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tr.ValidateDefinition(Definition{Path: tc.path, Target: tc.target}); err != nil {
				t.Errorf("ValidateDefinition(path=%q, target=%q) = %v, want nil", tc.path, tc.target, err)
			}
		})
	}
}
