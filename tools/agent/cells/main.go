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

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/cellindex"
)

type outputCell struct {
	cellindex.Cell
	Status string `json:"status"`
}

type outputIndex struct {
	SchemaVersion string       `json:"schemaVersion"`
	Hash          string       `json:"hash"`
	Cells         []outputCell `json:"cells"`
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

	if *jsonOutput {
		printJSON(idx)
		return
	}
	printTable(idx)
}

func readIndex(root string) (cellindex.Index, error) {
	data, err := os.ReadFile(filepath.Join(root, "gen", "cells.json"))
	if err != nil {
		return cellindex.Index{}, fmt.Errorf("error: gen/cells.json not found. Run 'task index' first")
	}

	idx, err := cellindex.Decode(data)
	if err != nil {
		return cellindex.Index{}, fmt.Errorf("error parsing gen/cells.json: %w", err)
	}
	return idx, nil
}

func printJSON(idx cellindex.Index) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(withStatus(idx)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		os.Exit(1)
	}
}

func withStatus(idx cellindex.Index) outputIndex {
	cells := make([]outputCell, len(idx.Cells))
	for i, cell := range idx.Cells {
		cells[i] = outputCell{Cell: cell, Status: "ok"}
	}
	return outputIndex{SchemaVersion: idx.SchemaVersion, Hash: idx.Hash, Cells: cells}
}

func printTable(idx cellindex.Index) {
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
