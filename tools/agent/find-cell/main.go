// Package main implements the structured cell search tool.
//
// Usage:
//
//	task find-cell QUERY=payment
//
// Searches cell IDs, purposes (from gen/cells.json), and AGENTS.md content
// for a text query. Returns matching cells with context.
//
// WHEN TO USE: When an agent needs to find which cell handles a concept.
// Faster and more precise than `grep -r "payment" internal/` because it only
// searches cell metadata and documentation, not implementation files.
//
// WHY IT HELPS AGENTS: Eliminates search ambiguity. An agent searching for
// "payment" gets high-signal results (cell IDs, purposes, AGENTS.md matches)
// instead of low-signal grep hits across all .go files.
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
	ID      string `json:"id"`
	Path    string `json:"path"`
	Purpose string `json:"purpose"`
}

type CellIndex struct {
	Cells []CellRecord `json:"cells"`
}

func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: task find-cell QUERY=<text>")
		os.Exit(1)
	}

	queryText := flag.Arg(0)
	query := strings.ToLower(queryText)

	data, err := os.ReadFile(filepath.Join(*root, "gen", "cells.json"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: gen/cells.json not found. Run 'task index' first.")
		os.Exit(1)
	}

	var idx CellIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing index: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=== Search: \"%s\" ===\n", queryText)
	found := false

	for _, c := range idx.Cells {
		matched := false
		matchSource := ""

		if strings.Contains(strings.ToLower(c.ID), query) {
			matched = true
			matchSource = "ID"
		}
		if strings.Contains(strings.ToLower(c.Purpose), query) {
			matched = true
			if matchSource == "" {
				matchSource = "purpose"
			}
		}

		agentsMatch := searchAgentsMD(filepath.Join(*root, c.Path), query)
		if agentsMatch != "" {
			matched = true
			matchSource = "AGENTS.md: " + agentsMatch
		}

		if matched {
			purpose := c.Purpose
			if len(purpose) > 50 {
				purpose = purpose[:47] + "..."
			}
			fmt.Printf("  %-20s %-45s \"%s\"", c.ID, c.Path, purpose)
			if matchSource != "" && matchSource != "ID" && matchSource != "purpose" {
				fmt.Printf("  (%s)", matchSource)
			}
			fmt.Println()
			found = true
		}
	}

	if !found {
		fmt.Println("  No matches found.")
	}
}

func searchAgentsMD(cellPath, query string) string {
	agentsPath := filepath.Join(cellPath, "AGENTS.md")
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.Contains(strings.ToLower(line), query) {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 60 {
				trimmed = trimmed[:57] + "..."
			}
			return "\"" + trimmed + "\""
		}
	}
	return ""
}
