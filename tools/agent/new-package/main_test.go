package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldCreatesConformancePlaceholder(t *testing.T) {
	root := t.TempDir()
	if err := scaffold(root, "field", "field", "field"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, "field", "cell.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"conformance:",
		"basis: engineering-primitive",
		"status: placeholder",
		"rationale:",
		"gaps:",
	} {
		if !strings.Contains(string(content), expected) {
			t.Errorf("scaffolded manifest does not contain %q:\n%s", expected, content)
		}
	}
}

func TestScaffoldUsesCanonicalValidationPathAtModuleRoot(t *testing.T) {
	root := t.TempDir()
	if err := scaffold(root, "canonical", ".", "canonical"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, "cell.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "- go test ./...\n") {
		t.Fatalf("module-root validation = %q, want go test ./...", content)
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agents), "`go test ./...`") {
		t.Fatalf("module-root guide = %q, want go test ./...", agents)
	}
}
