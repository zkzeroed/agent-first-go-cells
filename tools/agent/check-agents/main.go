// Package main implements the AGENTS.md validation tool.
//
// Usage:
//
//	task check-agents
//
// Checks that every cell-local AGENTS.md file has the required sections:
//   - Purpose
//   - Start Here (or Start)
//   - Invariants
//   - Common Tasks (or Tasks)
//   - Reliability & Concurrency
//   - Validation
//
// WHEN TO USE: After creating or editing an AGENTS.md. Also runs as part
// of `task doctor`.
//
// WHY IT HELPS AGENTS: Prevents drift in agent documentation. An agent
// creating a new cell via `task new-cell` gets a template with all sections,
// but later edits might remove sections. This check enforces completeness
// so the next agent arriving at the cell has all the context it needs.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

var requiredSections = []string{
	"## Purpose",
	"## Start",
	"## Invariants",
	"## Common Tasks",
	"## Reliability",
	"## Validation",
}

func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	exitCode := 0

	fmt.Println("=== AGENTS.md Validation ===")

	manifests, err := manifest.FindAllAt(*root)
	if err != nil {
		fmt.Printf("  ✗ reading manifests: %v\n", err)
		exitCode = 1
	}
	count := 0
	for _, cell := range manifests {
		path := filepath.Join(*root, cell.Dir, "AGENTS.md")
		if err := checkFile(path, filepath.ToSlash(filepath.Join(cell.Dir, "AGENTS.md"))); err != nil {
			fmt.Printf("  ✗ %s\n", err)
			exitCode = 1
		} else {
			fmt.Printf("  ✓ %s — all sections present\n", filepath.ToSlash(filepath.Join(cell.Dir, "AGENTS.md")))
			count++
		}
	}

	if exitCode == 0 {
		fmt.Printf("\nAll cell AGENTS.md files valid (%d files checked).\n", count)
	} else {
		fmt.Printf("\nSome cell AGENTS.md files are missing required sections.\n")
	}
	os.Exit(exitCode)
}

func checkFile(path, label string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s — file not found", label)
	}

	var missing []string
	lower := strings.ToLower(string(content))

	for _, section := range requiredSections {
		if !strings.Contains(lower, strings.ToLower(section)) {
			missing = append(missing, section)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s — missing: %s", label, strings.Join(missing, ", "))
	}

	return nil
}
