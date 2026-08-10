// Package main implements the dependency query tool.
//
// Usage:
//
//	task deps ID=user-authenticate
//	task deps-json ID=user-authenticate
//
// Shows what a cell depends on and what cells depend on it. This answers the
// static blast-radius question without needing a git diff.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/scripts/manifest"
)

type dependencyReport struct {
	ID           string      `json:"id"`
	Path         string      `json:"path"`
	Dependencies []string    `json:"dependencies"`
	Dependents   []dependent `json:"dependents"`
}

type dependent struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type args struct {
	cellID string
	json   bool
	root   string
}

func main() {
	args := parseArgs(os.Args[1:])
	if args.cellID == "" {
		fmt.Fprintln(os.Stderr, "Usage: task deps ID=<cell-id>")
		os.Exit(1)
	}

	report, err := buildReport(args.root, args.cellID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if args.json {
		printJSON(report)
		return
	}
	printText(report)
}

func parseArgs(values []string) args {
	result := args{root: "."}
	skipNext := false
	for i := range len(values) {
		if skipNext {
			skipNext = false
			continue
		}
		switch values[i] {
		case "--json":
			result.json = true
		case "--root":
			if i+1 >= len(values) {
				return args{}
			}
			result.root = values[i+1]
			skipNext = true
		default:
			if result.cellID == "" {
				result.cellID = values[i]
			}
		}
	}
	return result
}

func buildReport(root, targetID string) (dependencyReport, error) {
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		return dependencyReport{}, fmt.Errorf("error: %w", err)
	}

	target := findTarget(targetID, manifests)
	if target == nil {
		return dependencyReport{}, missingCellError(targetID, manifests)
	}

	return dependencyReport{
		ID:           target.ID,
		Path:         filepath.ToSlash(target.Dir),
		Dependencies: target.Dependencies,
		Dependents:   findDependents(targetID, manifests),
	}, nil
}

func findTarget(targetID string, manifests []manifest.Manifest) *manifest.Manifest {
	for i := range manifests {
		if manifests[i].ID == targetID {
			return &manifests[i]
		}
	}
	return nil
}

func missingCellError(targetID string, manifests []manifest.Manifest) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Cell %q not found.\nAvailable cells:\n", targetID)
	for _, m := range manifests {
		fmt.Fprintf(&builder, "  %s\n", m.ID)
	}
	return fmt.Errorf("%s", strings.TrimRight(builder.String(), "\n"))
}

func findDependents(targetID string, manifests []manifest.Manifest) []dependent {
	var dependents []dependent
	for _, m := range manifests {
		if manifestDependsOn(m, targetID) {
			dependents = append(dependents, dependent{ID: m.ID, Path: m.Dir})
		}
	}
	return dependents
}

func manifestDependsOn(m manifest.Manifest, targetID string) bool {
	return slices.Contains(m.Dependencies, targetID)
}

func printJSON(report dependencyReport) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		os.Exit(1)
	}
}

func printText(report dependencyReport) {
	fmt.Printf("=== Dependencies of %s ===\n", report.ID)
	printStringList(report.Dependencies, "  →")

	fmt.Printf("\n=== Dependents of %s ===\n", report.ID)
	if len(report.Dependents) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, dep := range report.Dependents {
		fmt.Printf("  ← %s (%s)\n", dep.ID, dep.Path)
	}
}

func printStringList(values []string, prefix string) {
	if len(values) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, value := range values {
		fmt.Printf("%s %s\n", prefix, value)
	}
}
