package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/scripts/manifest"
)

func TestImpactFilesIncludesRelevantTrackedAndUntrackedFilesWithoutDuplicates(t *testing.T) {
	tracked := "internal/cells/example/service.go\ninternal/cells/example/cell.yaml\nREADME.md\n"
	untracked := "internal/cells/example/service.go\ninternal/cells/example/handler.go\ninternal/cells/example/AGENTS.md\nAGENTS.md\ninternal/cells/example/notes.txt\ninternal/cells/other/cell.yml\n"

	got := impactFiles(tracked, untracked)
	want := []string{
		"internal/cells/example/service.go",
		"internal/cells/example/cell.yaml",
		"internal/cells/example/handler.go",
		"internal/cells/example/AGENTS.md",
		"internal/cells/other/cell.yml",
	}
	if !slices.Equal(got, want) {
		t.Errorf("impactFiles() = %v, want %v", got, want)
	}
}

func TestMapCellsUsesDirectoryBoundaries(t *testing.T) {
	manifests := []manifest.Manifest{
		{ID: "orders", Dir: "internal/cells/orders"},
		{ID: "orders-create", Dir: "internal/cells/orders-create"},
	}

	got := mapCells([]string{
		"internal/cells/orders/AGENTS.md",
		"internal/cells/orders-create/service.go",
		"internal/cells/orders-archive/service.go",
	}, manifests)
	want := []string{"orders", "orders-create"}
	if !slices.Equal(got, want) {
		t.Errorf("mapCells() = %v, want %v", got, want)
	}
}

func TestFindAffectedIncludesTransitiveDependentsInSortedOrder(t *testing.T) {
	manifests := []manifest.Manifest{
		{ID: "render", Dependencies: []string{"compose"}},
		{ID: "publish", Dependencies: []string{"render", "validate"}},
		{ID: "compose", Dependencies: []string{"source"}},
		{ID: "validate", Dependencies: []string{"source"}},
		{ID: "source"},
	}

	got := findAffected([]string{"source"}, manifests)
	want := []string{"compose", "publish", "render", "validate"}
	if !slices.Equal(got, want) {
		t.Errorf("findAffected() = %v, want %v", got, want)
	}
}

func TestGetChangedFilesIncludesNonIgnoredUntrackedFiles(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.com")
	git(t, dir, "config", "user.name", "Impact Test")
	writeFile(t, dir, "tracked.go", "package example\n")
	writeFile(t, dir, ".gitignore", "ignored.go\n")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "initial")

	writeFile(t, dir, "tracked.go", "package example\n// changed\n")
	writeFile(t, dir, "untracked.yaml", "name: example\n")
	writeFile(t, dir, "ignored.go", "package ignored\n")
	t.Chdir(dir)

	got, err := getChangedFiles("HEAD")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"tracked.go", "untracked.yaml"}
	if !slices.Equal(got, want) {
		t.Errorf("getChangedFiles() = %v, want %v", got, want)
	}
}

func TestBuildReportMapsNestedProjectChangesToCells(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init")
	git(t, repository, "config", "user.email", "test@example.com")
	git(t, repository, "config", "user.name", "Impact Test")
	project := filepath.Join(repository, "examples", "reference-project")
	writeCell(t, project, "greeting-compose", nil)
	writeCell(t, project, "greeting-render", []string{"greeting-compose"})
	git(t, repository, "add", ".")
	git(t, repository, "commit", "-m", "initial")

	writeFile(t, project, "internal/cells/greeting-compose/service.go", "package greetingcompose\n")
	report, err := buildReport("HEAD", project)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"internal/cells/greeting-compose/service.go"}; !slices.Equal(report.Changed, want) {
		t.Errorf("Changed = %v, want %v", report.Changed, want)
	}
	if want := []string{"greeting-compose"}; !slices.Equal(report.OwningCells, want) {
		t.Errorf("OwningCells = %v, want %v", report.OwningCells, want)
	}
	if want := []string{"greeting-render"}; !slices.Equal(report.Affected, want) {
		t.Errorf("Affected = %v, want %v", report.Affected, want)
	}
}

func TestBuildReportMapsCellAgentsChangesToTheOwningCell(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init")
	git(t, repository, "config", "user.email", "test@example.com")
	git(t, repository, "config", "user.name", "Impact Test")
	writeCell(t, repository, "greeting-compose", nil)
	writeCell(t, repository, "greeting-render", []string{"greeting-compose"})
	git(t, repository, "add", ".")
	git(t, repository, "commit", "-m", "initial")

	writeFile(t, repository, "internal/cells/greeting-compose/AGENTS.md", "# Updated guidance\n")
	report, err := buildReport("HEAD", repository)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"greeting-compose"}; !slices.Equal(report.OwningCells, want) {
		t.Errorf("OwningCells = %v, want %v", report.OwningCells, want)
	}
	if want := []string{"greeting-render"}; !slices.Equal(report.Affected, want) {
		t.Errorf("Affected = %v, want %v", report.Affected, want)
	}
}

func TestFilesUnderRootFiltersRepositoryPaths(t *testing.T) {
	got, err := filesUnderRoot(
		[]string{"scripts/impact/main.go", "examples/reference-project/internal/cells/example/service.go"},
		"/repo",
		"/repo/examples/reference-project",
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"internal/cells/example/service.go"}
	if !slices.Equal(got, want) {
		t.Errorf("filesUnderRoot() = %v, want %v", got, want)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func writeFile(t *testing.T, dir, name, contents string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeCell(t *testing.T, project, id string, dependencies []string) {
	t.Helper()
	dependencyLines := "dependencies: []\n"
	if len(dependencies) > 0 {
		dependencyLines = "dependencies:\n"
		for _, dependency := range dependencies {
			dependencyLines += "  - " + dependency + "\n"
		}
	}
	dir := filepath.Join(project, "internal", "cells", id)
	manifest := "id: " + id + "\npurpose: Test cell\nentrypoints:\n  - file: service.go\n    symbol: Service\n" + dependencyLines + "validation:\n  - go test ./...\ninvariants: []\n"
	writeFile(t, dir, "cell.yaml", manifest)
	writeFile(t, dir, "AGENTS.md", "# Test cell\n")
}
