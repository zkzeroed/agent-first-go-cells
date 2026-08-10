// Package main implements the impact analysis tool.
//
// Usage:
//
//	task impact
//	task impact-json
//
// Maps changed files to owning cells and affected downstream cells, then prints
// the validation commands that should be run.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sploitzberg/go-agent-first-arch/scripts/manifest"
)

type impactArgs struct {
	base string
	json bool
	root string
}

type impactReport struct {
	Base        string   `json:"base"`
	Changed     []string `json:"changedFiles"`
	OwningCells []string `json:"owningCells"`
	Affected    []string `json:"affectedCells"`
	Validation  []string `json:"validation"`
}

func main() {
	args := parseArgs(os.Args[1:])
	report, err := buildReport(args.base, args.root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if args.json {
		printJSON(report)
		return
	}
	printText(report)
}

func parseArgs(values []string) impactArgs {
	result := impactArgs{base: "HEAD", root: "."}
	skipNext := false
	for i := range len(values) {
		if skipNext {
			skipNext = false
			continue
		}
		switch values[i] {
		case "--json":
			result.json = true
		case "--base":
			if i+1 < len(values) {
				result.base = values[i+1]
				skipNext = true
			}
		case "--root":
			if i+1 < len(values) {
				result.root = values[i+1]
				skipNext = true
			}
		}
	}
	return result
}

func buildReport(base, root string) (impactReport, error) {
	repositoryRoot, err := gitRoot(root)
	if err != nil {
		return impactReport{}, fmt.Errorf("finding Git repository: %w", err)
	}

	changedFiles, err := getChangedFilesAt(base, repositoryRoot)
	if err != nil {
		return impactReport{}, fmt.Errorf("getting changed files: %w", err)
	}

	changedFiles, err = filesUnderRoot(changedFiles, repositoryRoot, root)
	if err != nil {
		return impactReport{}, fmt.Errorf("filtering changed files: %w", err)
	}

	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		return impactReport{}, fmt.Errorf("reading manifests: %w", err)
	}

	owningCells := mapCells(changedFiles, manifests)
	affectedCells := findAffected(owningCells, manifests)
	return impactReport{
		Base:        base,
		Changed:     emptyIfNil(changedFiles),
		OwningCells: emptyIfNil(owningCells),
		Affected:    emptyIfNil(affectedCells),
		Validation:  emptyIfNil(validationCommands(owningCells, affectedCells, manifests)),
	}, nil
}

func getChangedFiles(base string) ([]string, error) {
	return getChangedFilesAt(base, ".")
}

func getChangedFilesAt(base, directory string) ([]string, error) {
	cmd := exec.Command("git", "-C", directory, "diff", "--name-only", base)
	tracked, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("git", "-C", directory, "ls-files", "--others", "--exclude-standard")
	untracked, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return impactFiles(string(tracked), string(untracked)), nil
}

func gitRoot(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func filesUnderRoot(files []string, repositoryRoot, root string) ([]string, error) {
	projectRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	rootRelative, err := filepath.Rel(repositoryRoot, projectRoot)
	if err != nil {
		return nil, err
	}
	if isOutsideRoot(rootRelative) {
		return nil, fmt.Errorf("project root %s is outside Git repository %s", projectRoot, repositoryRoot)
	}

	rootPrefix := filepath.ToSlash(filepath.Clean(rootRelative))
	if rootPrefix == "." {
		return files, nil
	}
	rootPrefix += "/"
	var result []string
	for _, file := range files {
		if relative, found := strings.CutPrefix(filepath.ToSlash(file), rootPrefix); found {
			result = append(result, relative)
		}
	}
	return result, nil
}

func isOutsideRoot(path string) bool {
	return path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator))
}

func impactFiles(outputs ...string) []string {
	var files []string
	for _, output := range outputs {
		for line := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
			line = strings.TrimSpace(line)
			if isTrackedForImpact(line) && !slices.Contains(files, line) {
				files = append(files, line)
			}
		}
	}
	return files
}

func isTrackedForImpact(path string) bool {
	return path != "" && (strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"))
}

func mapCells(files []string, manifests []manifest.Manifest) []string {
	var result []string
	for _, file := range files {
		for _, m := range manifests {
			if strings.HasPrefix(file, m.Dir) && !slices.Contains(result, m.ID) {
				result = append(result, m.ID)
			}
		}
		if strings.HasPrefix(file, "internal/contracts/") && !slices.Contains(result, "contracts (shared)") {
			result = append(result, "contracts (shared)")
		}
	}
	return result
}

func findAffected(owningCells []string, manifests []manifest.Manifest) []string {
	var result []string
	for _, m := range manifests {
		if dependsOnAny(m, owningCells) && !slices.Contains(owningCells, m.ID) {
			result = append(result, m.ID)
		}
	}
	return result
}

func dependsOnAny(m manifest.Manifest, owningCells []string) bool {
	for _, owner := range owningCells {
		if slices.Contains(m.Dependencies, owner) {
			return true
		}
	}
	return false
}

func validationCommands(owningCells []string, affected []string, manifests []manifest.Manifest) []string {
	cellIDs := unique(append(slices.Clone(owningCells), affected...))
	var commands []string
	for _, cellID := range cellIDs {
		for _, m := range manifests {
			if m.ID == cellID {
				commands = append(commands, m.Validation...)
			}
		}
	}
	return commands
}

func unique(values []string) []string {
	var result []string
	for _, value := range values {
		if !slices.Contains(result, value) {
			result = append(result, value)
		}
	}
	return result
}

func emptyIfNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func printJSON(report impactReport) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		os.Exit(1)
	}
}

func printText(report impactReport) {
	if len(report.Changed) == 0 {
		fmt.Println("No changed files detected.")
		return
	}
	fmt.Println("=== Impact Analysis ===")
	printSection("Changed files", report.Changed)
	printSection("Owning cells", report.OwningCells)
	printSection("Affected cells (depend on changed cells)", report.Affected)
	printSection("Validation commands to run", report.Validation)
}

func printSection(title string, values []string) {
	fmt.Printf("\n%s:\n", title)
	if len(values) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, value := range values {
		fmt.Printf("  %s\n", value)
	}
}
