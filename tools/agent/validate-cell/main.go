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
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: task validate-cell ID=<cell-id>")
		os.Exit(1)
	}

	cellID := flag.Arg(0)
	manifests, err := manifest.FindAllAt(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: read manifests: %v\n", err)
		os.Exit(1)
	}
	var m *manifest.Manifest
	for i := range manifests {
		if manifests[i].ID == cellID {
			m = &manifests[i]
			break
		}
	}
	if m == nil {
		fmt.Fprintf(os.Stderr, "Error: cell %q not found\n", cellID)
		os.Exit(1)
	}
	commands := m.Validation

	if len(commands) == 0 {
		fmt.Fprintf(os.Stderr, "No validation commands found for %s\n", cellID)
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
		c.Dir = *root
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
