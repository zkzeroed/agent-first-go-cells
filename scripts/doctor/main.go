// Package main implements the architecture health check tool.
//
// Usage:
//
//	task doctor
//
// Runs all architecture-level checks and prints a summary. Does NOT run
// tests or build — this is the fast structural check.
//
// WHEN TO USE: Before committing. Catches architecture violations without
// the overhead of `task all` (which also builds and runs tests).
//
// WHY IT HELPS AGENTS: An agent can run `task doctor` in seconds to verify
// its work doesn't violate architecture invariants. If it passes, the agent
// can be confident the architecture is intact before producing a handoff.
//
// CHECKS PERFORMED:
//  1. Manifests valid (all cell.yaml have required fields)
//  2. Index fresh (gen/cells.json hash matches current manifests)
//  3. Structure tests pass (file size, function size, no init, cell completeness)
//  4. AGENTS.md sections complete (all required sections present)
//  5. Import policy clean (no forbidden imports)
package main

import (
	"fmt"
	"os"
	"os/exec"
)

type check struct {
	name string
	fn   func() (bool, string)
}

func main() {
	checks := []check{
		{"Manifests valid", checkManifests},
		{"Index fresh", checkIndex},
		{"Structure tests pass", checkStructureTests},
		{"AGENTS.md sections complete", checkAgents},
		{"Import policy clean", checkPolicy},
	}

	passed := 0
	failed := 0

	fmt.Println("=== architecture Health Check ===")

	for _, c := range checks {
		ok, detail := c.fn()
		if ok {
			fmt.Printf("  ✓ %s\n", c.name)
			passed++
		} else {
			fmt.Printf("  ✗ %s\n", c.name)
			if detail != "" {
				fmt.Printf("    %s\n", detail)
			}
			failed++
		}
	}

	fmt.Printf("\nResult: %d/%d checks passed", passed, passed+failed)
	if failed > 0 {
		fmt.Printf(". Fix %d issue(s) before committing.\n", failed)
		os.Exit(1)
	}
	fmt.Println(". All checks passed.")
}

func checkManifests() (bool, string) {
	cmd := exec.Command("task", "check-manifests")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false, "Run 'task check-manifests' for details"
	}
	return true, ""
}

func checkIndex() (bool, string) {
	cmd := exec.Command("go", "run", "./scripts/index/", "--check")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false, "Index is stale. Run 'task index' to regenerate."
	}
	return true, ""
}

func checkStructureTests() (bool, string) {
	cmd := exec.Command("go", "test", "-run", "Test(Cell|No|Files|Every)", "./scripts/structure/...")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false, "Run 'task structure-test' for details"
	}
	return true, ""
}

func checkAgents() (bool, string) {
	cmd := exec.Command("go", "run", "./scripts/check-agents/")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false, "Run 'task check-agents' for details"
	}
	return true, ""
}

func checkPolicy() (bool, string) {
	cmd := exec.Command("go", "run", "./scripts/imports/")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false, "Run 'task policy' for details"
	}
	return true, ""
}
