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
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

type impactArgs struct {
	base string
	json bool
	root string
}

type impactReport struct {
	Base                  string   `json:"base"`
	Changed               []string `json:"changedFiles"`
	OwningCells           []string `json:"owningCells"`
	Affected              []string `json:"affectedCells"`
	Validation            []string `json:"validation"`
	SharedSurfaces        []string `json:"sharedSurfaces"`
	FullProjectValidation bool     `json:"fullProjectValidation"`
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

	removed, err := removedManifestsAtHEAD(root, changedFiles)
	if err != nil {
		return impactReport{}, fmt.Errorf("reading removed manifests: %w", err)
	}
	manifests, err := manifest.FindAllAtWith(root, removed)
	if err != nil {
		return impactReport{}, fmt.Errorf("reading manifests: %w", err)
	}

	sharedSurfaces := findSharedSurfaces(changedFiles)
	owningCells := mapCells(changedFiles, manifests)
	fullProjectValidation := len(sharedSurfaces) > 0
	affectedCells := findAffected(owningCells, manifests)
	if fullProjectValidation {
		affectedCells = allCellIDs(manifests)
	}
	return impactReport{
		Base:                  base,
		Changed:               emptyIfNil(changedFiles),
		OwningCells:           emptyIfNil(owningCells),
		Affected:              emptyIfNil(affectedCells),
		Validation:            emptyIfNil(validationCommands(owningCells, affectedCells, manifests, fullProjectValidation)),
		SharedSurfaces:        emptyIfNil(sharedSurfaces),
		FullProjectValidation: fullProjectValidation,
	}, nil
}

func removedManifestsAtHEAD(root string, files []string) ([]manifest.Manifest, error) {
	var manifests []manifest.Manifest
	for _, file := range files {
		id, found := removedCellID(file)
		if !found || manifestByID(id, manifests) != nil {
			continue
		}
		missing, err := missingFromWorktree(root, file)
		if err != nil {
			return nil, err
		}
		if !missing {
			continue
		}
		removed, err := manifestAtHEAD(root, id)
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, *removed)
	}
	return manifests, nil
}

func missingFromWorktree(root, file string) (bool, error) {
	_, err := os.Stat(filepath.Join(root, file))
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	return false, err
}

func removedCellID(file string) (string, bool) {
	file = filepath.ToSlash(file)
	prefix := "internal/cells/"
	if !strings.HasPrefix(file, prefix) || !strings.HasSuffix(file, "/cell.yaml") {
		return "", false
	}
	return strings.TrimSuffix(strings.TrimPrefix(file, prefix), "/cell.yaml"), true
}

func manifestAtHEAD(root, id string) (*manifest.Manifest, error) {
	repository, err := gitRoot(root)
	if err != nil {
		return nil, err
	}
	project, err := projectPath(repository, root)
	if err != nil {
		return nil, err
	}
	path := filepath.ToSlash(filepath.Join(project, "internal", "cells", id, "cell.yaml"))
	content, err := gitOutput(repository, "show", "HEAD:"+path)
	if err != nil {
		return nil, err
	}
	result, err := manifest.Parse(content)
	if err != nil {
		return nil, err
	}
	result.Dir = filepath.Join("internal", "cells", filepath.FromSlash(id))
	return &result, nil
}

func manifestByID(id string, manifests []manifest.Manifest) *manifest.Manifest {
	for i := range manifests {
		if manifests[i].ID == id {
			return &manifests[i]
		}
	}
	return nil
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
	return gitOutput(root, "rev-parse", "--show-toplevel")
}

func gitOutput(root string, args ...string) (string, error) {
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func filesUnderRoot(files []string, repositoryRoot, root string) ([]string, error) {
	rootRelative, err := projectPath(repositoryRoot, root)
	if err != nil {
		return nil, err
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

func projectPath(repositoryRoot, root string) (string, error) {
	projectRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	rootRelative, err := filepath.Rel(repositoryRoot, projectRoot)
	if err != nil || isOutsideRoot(rootRelative) {
		return "", fmt.Errorf("project root %s is outside Git repository %s", projectRoot, repositoryRoot)
	}
	return rootRelative, nil
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
	return path != "" && (strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") || isCellAgentsFile(path) || sharedSurface(path) != "")
}

func isCellAgentsFile(path string) bool {
	path = filepath.ToSlash(path)
	return filepath.Base(path) == "AGENTS.md" && strings.HasPrefix(path, "internal/cells/")
}

func mapCells(files []string, manifests []manifest.Manifest) []string {
	var result []string
	for _, file := range files {
		owner := ""
		ownerLength := 0
		for _, m := range manifests {
			if isWithinCell(file, m.Dir) && len(m.Dir) > ownerLength {
				owner = m.ID
				ownerLength = len(m.Dir)
			}
		}
		if owner != "" && !slices.Contains(result, owner) {
			result = append(result, owner)
		}
		if surface := sharedSurface(file); surface != "" && !slices.Contains(result, surface+" (shared)") {
			result = append(result, surface+" (shared)")
		}
	}
	return result
}

func findSharedSurfaces(files []string) []string {
	var result []string
	for _, file := range files {
		if surface := sharedSurface(file); surface != "" && !slices.Contains(result, surface) {
			result = append(result, surface)
		}
	}
	return result
}

func sharedSurface(file string) string {
	file = filepath.ToSlash(file)
	switch {
	case strings.HasPrefix(file, "internal/contracts/"):
		return "contracts"
	case strings.HasPrefix(file, "internal/platform/"):
		return "platform"
	case filepath.Dir(file) == "internal/app" && strings.HasPrefix(filepath.Base(file), "wiring") && strings.HasSuffix(file, ".go"):
		return "wiring"
	default:
		return ""
	}
}

func isWithinCell(file, cellDir string) bool {
	file = filepath.ToSlash(file)
	cellDir = strings.TrimSuffix(filepath.ToSlash(cellDir), "/")
	return file == cellDir || strings.HasPrefix(file, cellDir+"/")
}

func findAffected(owningCells []string, manifests []manifest.Manifest) []string {
	affected := make(map[string]bool)
	frontier := unique(slices.Clone(owningCells))
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, m := range manifests {
			if !slices.Contains(m.Dependencies, current) || slices.Contains(owningCells, m.ID) || affected[m.ID] {
				continue
			}
			affected[m.ID] = true
			frontier = append(frontier, m.ID)
		}
	}
	result := slices.Collect(maps.Keys(affected))
	slices.Sort(result)
	return result
}

func allCellIDs(manifests []manifest.Manifest) []string {
	result := make([]string, 0, len(manifests))
	for _, m := range manifests {
		result = append(result, m.ID)
	}
	slices.Sort(result)
	return result
}

func validationCommands(owningCells []string, affected []string, manifests []manifest.Manifest, fullProject bool) []string {
	if fullProject {
		return []string{"go test ./..."}
	}
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
	printSection("Shared surfaces", report.SharedSurfaces)
	printSection("Affected cells (downstream of changed cells)", report.Affected)
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
