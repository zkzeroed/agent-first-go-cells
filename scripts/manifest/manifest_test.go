package manifest

import (
	"strings"
	"testing"
)

func TestParseRejectsUnknownFieldsAndMissingEntrypointFile(t *testing.T) {
	_, err := Parse("id: orders-create\npurpose: create an order\nentrypoints:\n  - symbol: Create\ndependencies: []\nvalidation:\n  - go test ./...\nunknown: value\n")
	if err == nil {
		t.Fatal("Parse() error = nil, want schema validation error")
	}
	if !strings.Contains(err.Error(), "field unknown") {
		t.Fatalf("Parse() error = %q, want unknown field error", err)
	}
}

func TestValidateRequiresExactKnownDependency(t *testing.T) {
	manifests := []Manifest{
		{ID: "orders-create", Dir: "internal/cells/orders-create"},
		{ID: "orders-read", Dir: "internal/cells/orders-read", Dependencies: []string{"orders"}},
	}
	err := Validate(manifests)
	if err == nil {
		t.Fatal("Validate() error = nil, want unknown dependency error")
	}
	if !strings.Contains(err.Error(), "exactly match an existing cell id") {
		t.Fatalf("Validate() error = %q, want exact dependency error", err)
	}
}

func TestParsePreservesNestedEntrypointAndValidation(t *testing.T) {
	m, err := Parse("id: orders/create\npurpose: create an order\nentrypoints:\n  - file: handler.go\n    symbol: Handle\ndependencies: []\nvalidation:\n  - go test ./internal/cells/orders/create/...\ninvariants:\n  - idempotent\n")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got, want := m.Entrypoints, []string{"handler.go"}; !equalStrings(got, want) {
		t.Fatalf("Entrypoints = %#v, want %#v", got, want)
	}
	if got, want := m.Validation, []string{"go test ./internal/cells/orders/create/..."}; !equalStrings(got, want) {
		t.Fatalf("Validation = %#v, want %#v", got, want)
	}
}

// FuzzParse verifies that arbitrary manifest input is rejected safely or
// produces a manifest without panicking.
func FuzzParse(f *testing.F) {
	for _, content := range []string{
		"",
		"id: orders-create\npurpose: create an order\nentrypoints:\n  - file: handler.go\ndependencies: []\nvalidation:\n  - go test ./...\n",
		"id: orders-create\nunknown: value\n",
		"---\nid: orders-create\n---\nid: orders-read\n",
	} {
		f.Add(content)
	}

	f.Fuzz(func(t *testing.T, content string) {
		_, _ = Parse(content)
	})
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
