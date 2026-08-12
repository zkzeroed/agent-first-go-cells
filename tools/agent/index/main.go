// Package main implements the cell index generator.
//
// Usage:
//
//	task index              # Regenerate gen/cells.json + context packs
//	task check-index        # Verify index is not stale (hash-based)
//	task index-json         # Output index as JSON without writing files
//
// This tool reads all cell.yaml manifests under internal/cells/, builds a
// JSON index, and generates per-cell context packs from cell.yaml + AGENTS.md.
// A SHA-256 hash of all manifest contents enables staleness detection without
// git dependencies.
//
// WHEN TO USE: After creating, modifying, or deleting any cell. Run `task index`
// to regenerate, or `task check-index` in CI to verify freshness.
//
// WHY IT HELPS AGENTS: The generated gen/cells.json lets an agent discover all
// cells in a single read. The context packs let an agent orient on a cell
// without reading 6-8 source files. The hash-based staleness check works in
// clean CI checkouts without git diff.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

type CellRecord struct {
	ID           string             `json:"id"`
	Path         string             `json:"path"`
	Package      string             `json:"package"`
	Kind         string             `json:"kind,omitempty"`
	Public       bool               `json:"public"`
	Purpose      string             `json:"purpose"`
	Entrypoints  []string           `json:"entrypoints"`
	Dependencies []string           `json:"dependencies"`
	Validation   []string           `json:"validation"`
	Conformance  *ConformanceRecord `json:"conformance,omitempty"`
}

// ConformanceRecord makes a public package's research provenance and known
// limitations available without reparsing its source manifest.
type ConformanceRecord struct {
	Basis     string           `json:"basis"`
	Status    string           `json:"status"`
	Evidence  string           `json:"evidence"`
	Citations []CitationRecord `json:"citations,omitempty"`
	Rationale string           `json:"rationale,omitempty"`
	Gaps      []string         `json:"gaps,omitempty"`
}

type CitationRecord struct {
	File    string        `json:"file"`
	Locator LocatorRecord `json:"locator"`
	Symbols []string      `json:"symbols"`
}

type LocatorRecord struct {
	Type    string `json:"type"`
	Pages   []uint `json:"pages,omitempty"`
	Heading string `json:"heading,omitempty"`
}

type CellIndex struct {
	SchemaVersion string       `json:"schemaVersion"`
	Hash          string       `json:"hash"`
	Cells         []CellRecord `json:"cells"`
}

const maxContextGuideBytes = 6_000

func main() {
	root, check, jsonOutput, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(2)
	}
	if check {
		if err := checkStale(root); err != nil {
			fmt.Fprintf(os.Stderr, "Index is stale: %v\n", err)
			fmt.Fprintln(os.Stderr, "Run 'task index' to regenerate.")
			os.Exit(1)
		}
		fmt.Println("Index is fresh.")
		return
	}

	if jsonOutput {
		if err := printJSON(root); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON index: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := generate(root); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating index: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs(values []string) (string, bool, bool, error) {
	root := "."
	var check, jsonOutput bool
	skipNext := false
	for i := range len(values) {
		if skipNext {
			skipNext = false
			continue
		}
		switch values[i] {
		case "--check":
			check = true
		case "--json":
			jsonOutput = true
		case "--root":
			if i+1 >= len(values) {
				return "", false, false, errors.New("--root requires a path")
			}
			root = values[i+1]
			skipNext = true
		default:
			return "", false, false, fmt.Errorf("unknown argument %q", values[i])
		}
	}
	return root, check, jsonOutput, nil
}

func generate(root string) error {
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		return err
	}
	if err := manifest.ValidateSourceAt(root, manifests); err != nil {
		return err
	}

	hash := computeHash(manifests)
	index := buildIndex(manifests)

	if err := writeIndex(root, index); err != nil {
		return err
	}

	if err := writeContextPacks(root, manifests); err != nil {
		return err
	}

	fmt.Printf("Generated index with %d cells (hash: sha256:%s)\n", len(manifests), hash[:12])
	return nil
}

func printJSON(root string) error {
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		return err
	}
	if err := manifest.ValidateSourceAt(root, manifests); err != nil {
		return err
	}

	data, err := json.MarshalIndent(buildIndex(manifests), "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func buildIndex(manifests []manifest.Manifest) CellIndex {
	index := CellIndex{
		SchemaVersion: "agent-first/v5",
		Hash:          "sha256:" + computeHash(manifests),
		Cells:         []CellRecord{},
	}
	for _, m := range manifests {
		index.Cells = append(index.Cells, CellRecord{
			ID:           m.ID,
			Path:         m.Dir,
			Package:      strings.ReplaceAll(m.ID, "-", ""),
			Kind:         m.Kind,
			Public:       m.Public,
			Purpose:      m.Purpose,
			Entrypoints:  m.Entrypoints,
			Dependencies: m.Dependencies,
			Validation:   m.Validation,
			Conformance:  conformanceRecord(m.Conformance),
		})
	}
	return index
}

func conformanceRecord(value manifest.Conformance) *ConformanceRecord {
	if value.Basis == "" {
		return nil
	}
	record := &ConformanceRecord{Basis: value.Basis, Status: value.Status, Evidence: value.Evidence, Rationale: value.Rationale, Gaps: slices.Clone(value.Gaps)}
	for _, citation := range value.Citations {
		record.Citations = append(record.Citations, CitationRecord{File: citation.File, Locator: LocatorRecord{Type: citation.Locator.Type, Pages: slices.Clone(citation.Locator.Pages), Heading: citation.Locator.Heading}, Symbols: slices.Clone(citation.Symbols)})
	}
	return record
}

func checkStale(root string) error {
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		return err
	}
	if err := manifest.ValidateSourceAt(root, manifests); err != nil {
		return err
	}

	currentHash := "sha256:" + computeHash(manifests)

	existing, err := os.ReadFile(filepath.Join(root, "gen", "cells.json"))
	if err != nil {
		return fmt.Errorf("cannot read existing index: %w", err)
	}

	var idx CellIndex
	if err := json.Unmarshal(existing, &idx); err != nil {
		return fmt.Errorf("cannot parse existing index: %w", err)
	}

	if idx.Hash != currentHash {
		return fmt.Errorf("hash mismatch: existing=%s current=%s", idx.Hash, currentHash)
	}
	return checkContextPacks(root, manifests)
}

func computeHash(manifests []manifest.Manifest) string {
	h := sha256.New()
	sorted := slices.Clone(manifests)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	for _, m := range sorted {
		h.Write([]byte(m.ID))
		h.Write([]byte{0})
		h.Write(m.SourceContent())
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeIndex(root string, index CellIndex) error {
	if err := os.MkdirAll(filepath.Join(root, "gen"), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(root, "gen", "cells.json"), data, 0o644)
}

func writeContextPacks(root string, manifests []manifest.Manifest) error {
	if err := os.MkdirAll(filepath.Join(root, "gen", "context"), 0o755); err != nil {
		return err
	}

	expected := make(map[string]string, len(manifests))
	for _, m := range manifests {
		pack := buildContextPack(m)
		path := filepath.Join(root, contextPackPath(m.ID))
		expected[filepath.Clean(path)] = pack
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(pack), 0o644); err != nil {
			return err
		}
	}
	return removeOrphanContextPacks(filepath.Join(root, "gen", "context"), expected)
}

func buildContextPack(m manifest.Manifest) string {
	var sb strings.Builder

	h := sha256.New()
	h.Write(m.SourceContent())
	sourceHash := hex.EncodeToString(h.Sum(nil))

	sb.WriteString(fmt.Sprintf("<!-- source-hash: sha256:%s -->\n", sourceHash))
	sb.WriteString("<!-- generated-from: cell.yaml, AGENTS.md -->\n\n")
	sb.WriteString(fmt.Sprintf("# %s\n\n", m.ID))
	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", m.Purpose))
	sb.WriteString(fmt.Sprintf("**Kind:** %s\n\n", displayKind(m.Kind)))
	sb.WriteString(fmt.Sprintf("**Public:** %t\n\n", m.Public))
	writeConformance(&sb, m.Conformance)

	if len(m.Entrypoints) > 0 {
		sb.WriteString("**Entrypoints:**\n\n")
		for _, e := range m.Entrypoints {
			sb.WriteString(fmt.Sprintf("- `%s`\n", e))
		}
		sb.WriteString("\n")
	}

	if len(m.Dependencies) > 0 {
		sb.WriteString("**Dependencies:**\n\n")
		for _, d := range m.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s\n", d))
		}
		sb.WriteString("\n")
	}

	if len(m.Invariants) > 0 {
		sb.WriteString("**Invariants:**\n\n")
		for _, inv := range m.Invariants {
			sb.WriteString(fmt.Sprintf("- %s\n", inv))
		}
		sb.WriteString("\n")
	}

	if len(m.Validation) > 0 {
		sb.WriteString("**Validation:**\n\n")
		for _, v := range m.Validation {
			sb.WriteString(fmt.Sprintf("- `%s`\n", v))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Cell Guide\n\n")
	guide := boundedGuide(m.AgentsContent)
	sb.WriteString(guide)
	if !strings.HasSuffix(guide, "\n") {
		sb.WriteString("\n")
	}

	return sb.String()
}

func writeConformance(sb *strings.Builder, value manifest.Conformance) {
	if value.Basis == "" {
		return
	}
	sb.WriteString("**Conformance:**\n\n")
	sb.WriteString(fmt.Sprintf("- Basis: %s\n- Status: %s\n- Evidence: %s\n", value.Basis, value.Status, value.Evidence))
	if value.Rationale != "" {
		sb.WriteString(fmt.Sprintf("- Rationale: %s\n", value.Rationale))
	}
	for _, gap := range value.Gaps {
		sb.WriteString(fmt.Sprintf("- Gap: %s\n", gap))
	}
	for _, citation := range value.Citations {
		sb.WriteString(fmt.Sprintf("- Citation: `%s`, %s; symbols: %s\n", citation.File, displayLocator(citation.Locator), strings.Join(citation.Symbols, ", ")))
	}
	sb.WriteString("\n")
}

func displayLocator(locator manifest.CitationLocator) string {
	if locator.Type == "markdown-heading" {
		return "heading " + fmt.Sprintf("%q", locator.Heading)
	}
	parts := make([]string, len(locator.Pages))
	for index, page := range locator.Pages {
		parts[index] = fmt.Sprintf("%d", page)
	}
	return "PDF pages " + strings.Join(parts, ", ")
}

func displayKind(kind string) string {
	if kind == "" {
		return "private cell"
	}
	return kind
}

func boundedGuide(guide string) string {
	if len(guide) <= maxContextGuideBytes {
		return guide
	}
	return guide[:maxContextGuideBytes] + "\n\n[Guide truncated; read the cell AGENTS.md for the full instructions.]\n"
}

func contextPackPath(id string) string {
	return filepath.Join("gen", "context", filepath.FromSlash(id)+".context.md")
}

func checkContextPacks(root string, manifests []manifest.Manifest) error {
	expected := make(map[string]string, len(manifests))
	for _, m := range manifests {
		expected[filepath.Clean(filepath.Join(root, contextPackPath(m.ID)))] = buildContextPack(m)
	}

	for path, want := range expected {
		got, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read context pack %s: %w", path, err)
		}
		if string(got) != want {
			return fmt.Errorf("context pack %s does not match its manifest and AGENTS.md sources", path)
		}
	}
	return findOrphanContextPack(filepath.Join(root, "gen", "context"), expected)
}

func removeOrphanContextPacks(contextDir string, expected map[string]string) error {
	var orphans []string
	if err := filepath.WalkDir(contextDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".context.md") {
			if _, exists := expected[filepath.Clean(path)]; !exists {
				orphans = append(orphans, path)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	for _, path := range orphans {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove orphan context pack %s: %w", path, err)
		}
	}
	return nil
}

func findOrphanContextPack(contextDir string, expected map[string]string) error {
	err := filepath.WalkDir(contextDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".context.md") {
			if _, exists := expected[filepath.Clean(path)]; !exists {
				return fmt.Errorf("orphan context pack %s", path)
			}
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) && len(expected) == 0 {
		return nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return errors.New("context pack directory is missing")
	}
	return err
}
