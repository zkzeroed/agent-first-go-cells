package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func TestCheckContextPacksAllowsMissingDirectoryWithoutCells(t *testing.T) {
	if err := checkContextPacks(t.TempDir(), nil); err != nil {
		t.Fatalf("checkContextPacks() error = %v, want nil", err)
	}
}

func TestCheckStaleAllowsMissingContextDirectoryWithoutCells(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "gen"), 0o755); err != nil {
		t.Fatal(err)
	}
	index, err := json.Marshal(buildIndex(nil))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "gen", "cells.json"), index, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkStale(root); err != nil {
		t.Fatalf("checkStale() error = %v, want nil", err)
	}
}

func TestCheckStaleRejectsSchemaVersionMismatch(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "gen"), 0o755); err != nil {
		t.Fatal(err)
	}
	index := buildIndex(nil)
	index.SchemaVersion = "agent-first/v0"
	data, err := json.Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "gen", "cells.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkStale(root); err == nil || !strings.Contains(err.Error(), "schema version mismatch") {
		t.Fatalf("checkStale() error = %v, want schema-version mismatch", err)
	}
}

func TestContextPackIncludesAgentsSourceAndHashChanges(t *testing.T) {
	base := manifest.Manifest{ID: "orders-create", Purpose: "create", RawContent: "id: orders-create\n", AgentsContent: "# Guide\n\nUse the service.\n"}
	changed := base
	changed.AgentsContent = "# Guide\n\nUse the handler.\n"

	if computeHash([]manifest.Manifest{base}) == computeHash([]manifest.Manifest{changed}) {
		t.Fatal("computeHash() did not include AGENTS.md content")
	}
	pack := buildContextPack(base)
	if !strings.Contains(pack, "## Cell Guide") || !strings.Contains(pack, "Use the service.") {
		t.Fatalf("context pack did not include AGENTS.md content: %q", pack)
	}
}

func TestIndexAndContextPackExposeLibraryMetadata(t *testing.T) {
	cell := manifest.Manifest{ID: "greeting", Kind: "library-package", Public: true, Purpose: "greet", Conformance: manifest.Conformance{Basis: "paper-defined-math", Status: "conformant", Evidence: "verified", Citations: []manifest.Citation{{File: "docs/paper.pdf", Locator: manifest.CitationLocator{Type: "pdf-pages", Pages: []uint{4}}, Symbols: []string{"New"}}}}, RawContent: "id: greeting\n", AgentsContent: "# Guide\n"}
	index := buildIndex([]manifest.Manifest{cell})
	if got := index.Cells[0]; got.Kind != "library-package" || !got.Public {
		t.Fatalf("library metadata = kind %q, public %t; want library-package, true", got.Kind, got.Public)
	}
	if got := index.Cells[0].Conformance; got == nil || got.Basis != "paper-defined-math" || len(got.Citations) != 1 {
		t.Fatalf("conformance = %#v, want generated paper record", got)
	}
	pack := buildContextPack(cell)
	if !strings.Contains(pack, "**Kind:** library-package") || !strings.Contains(pack, "**Public:** true") || !strings.Contains(pack, "docs/paper.pdf") {
		t.Fatalf("context pack did not include library metadata: %q", pack)
	}
}

func TestBoundedGuidePreservesSmallGuidesAndMarksLargeOnes(t *testing.T) {
	if got := boundedGuide("small"); got != "small" {
		t.Fatalf("boundedGuide(small) = %q", got)
	}
	large := strings.Repeat("x", maxContextGuideBytes+1)
	if got := boundedGuide(large); !strings.Contains(got, "Guide truncated") {
		t.Fatal("boundedGuide(large) did not include truncation marker")
	}
}
