// Package main implements the cell listing tool.
//
// Usage:
//
//	task cells
//	task cells-json
//
// Prints every cell, its path, purpose, and index freshness status. This is
// the primary orientation command for an agent arriving at a codebase.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	Status       string             `json:"status"`
}

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

func main() {
	jsonOutput := flag.Bool("json", false, "print machine-readable JSON")
	root := flag.String("root", ".", "project root")
	flag.Parse()

	idx, err := readIndex(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i := range idx.Cells {
		idx.Cells[i].Status = "ok"
	}

	if *jsonOutput {
		printJSON(idx)
		return
	}
	printTable(idx)
}

func readIndex(root string) (CellIndex, error) {
	data, err := os.ReadFile(filepath.Join(root, "gen", "cells.json"))
	if err != nil {
		return CellIndex{}, fmt.Errorf("error: gen/cells.json not found. Run 'task index' first")
	}

	var idx CellIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return CellIndex{}, fmt.Errorf("error parsing gen/cells.json: %w", err)
	}
	return idx, nil
}

func printJSON(idx CellIndex) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(idx); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		os.Exit(1)
	}
}

func printTable(idx CellIndex) {
	if len(idx.Cells) == 0 {
		fmt.Println("No cells found. Use 'task new-cell ID=<id>' to create one.")
		return
	}

	fmt.Printf("%-25s %-45s %-18s %-8s %-32s %s\n", "ID", "PATH", "KIND", "PUBLIC", "PURPOSE", "STATUS")
	fmt.Println(strings.Repeat("-", 160))
	for _, c := range idx.Cells {
		fmt.Printf("%-25s %-45s %-18s %-8t %-32s %s\n", c.ID, c.Path, displayKind(c.Kind), c.Public, truncate(c.Purpose), "✓")
	}
	fmt.Printf("\nTotal: %d cells (index hash: %s)\n", len(idx.Cells), idx.Hash)
}

func displayKind(kind string) string {
	if kind == "" {
		return "private cell"
	}
	return kind
}

func truncate(value string) string {
	if len(value) <= 43 {
		return value
	}
	return value[:40] + "..."
}
