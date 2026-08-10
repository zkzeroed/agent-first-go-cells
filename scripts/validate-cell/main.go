// Package main implements the cell validation runner.
//
// Usage:
//
//	task validate-cell ID=user-authenticate
//
// Reads the validation commands from a cell's cell.yaml and executes them
// in sequence. Stops on first failure.
//
// WHEN TO USE: After modifying a cell, before integrating into wiring.go.
// Also useful in CI for targeted validation.
//
// WHY IT HELPS AGENTS: An agent doesn't need to remember which commands to
// run for each cell — they're declared in cell.yaml and executed by this
// tool. This ensures consistent validation across all cells and all agents.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/scripts/manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: task validate-cell ID=<cell-id>")
		os.Exit(1)
	}

	cellID := os.Args[1]
	manifestPath := filepath.Join("internal/cells", cellID, "cell.yaml")

	content, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot read %s: %v\n", manifestPath, err)
		os.Exit(1)
	}

	m, err := manifest.Parse(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid manifest %s: %v\n", manifestPath, err)
		os.Exit(1)
	}
	commands := m.Validation

	if len(commands) == 0 {
		fmt.Fprintf(os.Stderr, "No validation commands found in %s\n", manifestPath)
		os.Exit(1)
	}

	fmt.Printf("=== Validating cell: %s ===\n", cellID)
	failed := 0

	for i, cmd := range commands {
		fmt.Printf("\n[%d/%d] Running: %s\n", i+1, len(commands), cmd)

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		c := exec.Command(parts[0], parts[1:]...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("  ✗ FAILED: %s\n", cmd)
			failed++
		} else {
			fmt.Printf("  ✓ PASSED\n")
		}
	}

	if failed > 0 {
		fmt.Printf("\n%d/%d validation commands failed.\n", failed, len(commands))
		os.Exit(1)
	}

	fmt.Printf("\nAll %d validation commands passed.\n", len(commands))
}
