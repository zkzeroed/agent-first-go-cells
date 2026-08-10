package main

import "testing"

func TestValidateImportEnforcesInternalAllowList(t *testing.T) {
	rule := Rule{From: "internal/cells/**", Allow: []string{"internal/cells/**/api"}}
	if got := validateImport(rule, "internal/cells/orders/api"); got != "" {
		t.Fatalf("validateImport(api) = %q", got)
	}
	if got := validateImport(rule, "internal/cells/orders"); got == "" {
		t.Fatal("validateImport(implementation) = nil, want violation")
	}
	if got := validateImport(rule, "net/http"); got != "" {
		t.Fatalf("validateImport(external) = %q", got)
	}
}

func TestNormalizeImport(t *testing.T) {
	if got, want := normalizeImport("example.com/project/internal/app", "example.com/project"), "internal/app"; got != want {
		t.Fatalf("normalizeImport() = %q, want %q", got, want)
	}
}
