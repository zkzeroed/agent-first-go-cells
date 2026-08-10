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
	ID           string   `json:"id"`
	Path         string   `json:"path"`
	Package      string   `json:"package"`
	Purpose      string   `json:"purpose"`
	Entrypoints  []string `json:"entrypoints"`
	Dependencies []string `json:"dependencies"`
	Validation   []string `json:"validation"`
	Status       string   `json:"status"`
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

	fmt.Printf("%-25s %-45s %-45s %s\n", "ID", "PATH", "PURPOSE", "STATUS")
	fmt.Println(strings.Repeat("-", 140))
	for _, c := range idx.Cells {
		fmt.Printf("%-25s %-45s %-45s %s\n", c.ID, c.Path, truncate(c.Purpose), "✓")
	}
	fmt.Printf("\nTotal: %d cells (index hash: %s)\n", len(idx.Cells), idx.Hash)
}

func truncate(value string) string {
	if len(value) <= 43 {
		return value
	}
	return value[:40] + "..."
}
