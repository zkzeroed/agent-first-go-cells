package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func TestImpactFilesIncludesRelevantTrackedAndUntrackedFilesWithoutDuplicates(t *testing.T) {
	tracked := "internal/cells/example/service.go\ninternal/cells/example/cell.yaml\nREADME.md\n"
	untracked := "internal/cells/example/service.go\ninternal/cells/example/handler.go\ninternal/cells/example/AGENTS.md\nAGENTS.md\ninternal/cells/example/notes.txt\ninternal/cells/other/cell.yml\ninternal/contracts/clock.md\n"

	got := impactFiles(tracked, untracked)
	want := []string{
		"internal/cells/example/service.go",
		"internal/cells/example/cell.yaml",
		"internal/cells/example/handler.go",
		"internal/cells/example/AGENTS.md",
		"internal/cells/other/cell.yml",
		"internal/contracts/clock.md",
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

func TestMapCellsUsesOnlyTheMostSpecificCellOwner(t *testing.T) {
	manifests := []manifest.Manifest{
		{ID: "orders", Dir: "internal/cells/orders"},
		{ID: "orders-create", Dir: "internal/cells/orders/create"},
	}

	got := mapCells([]string{"internal/cells/orders/create/service.go"}, manifests)
	want := []string{"orders-create"}
	if !slices.Equal(got, want) {
		t.Errorf("mapCells() = %v, want %v", got, want)
	}
}

func TestFindSharedSurfacesClassifiesContractsPlatformAndWiring(t *testing.T) {
	got := findSharedSurfaces([]string{
		"internal/contracts/clock.go",
		"internal/platform/http/server.go",
		"internal/app/wiring.go",
		"internal/app/wiring_test.go",
		"internal/app/wirings.go",
		"internal/app/wiring/helpers.go",
	})
	want := []string{"contracts", "platform", "wiring"}
	if !slices.Equal(got, want) {
		t.Errorf("findSharedSurfaces() = %v, want %v", got, want)
	}
}

func TestImpactReportJSONIncludesSharedSurfaceFields(t *testing.T) {
	report := impactReport{
		SharedSurfaces:        []string{"platform"},
		FullProjectValidation: true,
	}

	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if got := decoded["sharedSurfaces"]; got == nil {
		t.Error("JSON does not include sharedSurfaces")
	}
	if got, ok := decoded["fullProjectValidation"].(bool); !ok || !got {
		t.Errorf("fullProjectValidation = %v, want true", got)
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
	project := filepath.Join(repository, "nested", "project")
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

func TestBuildReportExpandsSharedSurfaceChangesToAllCells(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init")
	git(t, repository, "config", "user.email", "test@example.com")
	git(t, repository, "config", "user.name", "Impact Test")
	writeCell(t, repository, "greeting-compose", nil)
	writeCell(t, repository, "greeting-render", []string{"greeting-compose"})
	writeCell(t, repository, "token-issue", nil)
	writeFile(t, repository, "internal/contracts/clock.go", "package contracts\n")
	git(t, repository, "add", ".")
	git(t, repository, "commit", "-m", "initial")

	writeFile(t, repository, "internal/contracts/clock.go", "package contracts\n// changed\n")
	report, err := buildReport("HEAD", repository)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"contracts"}; !slices.Equal(report.SharedSurfaces, want) {
		t.Errorf("SharedSurfaces = %v, want %v", report.SharedSurfaces, want)
	}
	if !report.FullProjectValidation {
		t.Error("FullProjectValidation = false, want true")
	}
	if want := []string{"greeting-compose", "greeting-render", "token-issue"}; !slices.Equal(report.Affected, want) {
		t.Errorf("Affected = %v, want %v", report.Affected, want)
	}
	if want := []string{"go test ./..."}; !slices.Equal(report.Validation, want) {
		t.Errorf("Validation = %v, want %v", report.Validation, want)
	}
}

func TestBuildReportPreservesDeletedCellImpactFromHEAD(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init")
	git(t, repository, "config", "user.email", "test@example.com")
	git(t, repository, "config", "user.name", "Impact Test")
	writeCell(t, repository, "greeting-compose", nil)
	writeCell(t, repository, "greeting-render", []string{"greeting-compose"})
	git(t, repository, "add", ".")
	git(t, repository, "commit", "-m", "initial")

	if err := os.RemoveAll(filepath.Join(repository, "internal", "cells", "greeting-compose")); err != nil {
		t.Fatal(err)
	}
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

func TestBuildReportAcceptsModifiedAndNewManifests(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init")
	git(t, repository, "config", "user.email", "test@example.com")
	git(t, repository, "config", "user.name", "Impact Test")
	writeCell(t, repository, "greeting-compose", nil)
	git(t, repository, "add", ".")
	git(t, repository, "commit", "-m", "initial")

	writeFile(t, repository, "internal/cells/greeting-compose/cell.yaml", "id: greeting-compose\npurpose: Updated\nentrypoints:\n  - file: service.go\n    symbol: Service\ndependencies: []\nvalidation:\n  - go test ./...\n")
	writeCell(t, repository, "token-issue", nil)
	if _, err := buildReport("HEAD", repository); err != nil {
		t.Fatalf("buildReport() error = %v, want nil", err)
	}
}

func TestFilesUnderRootFiltersRepositoryPaths(t *testing.T) {
	got, err := filesUnderRoot(
		[]string{"tools/agent/impact/main.go", "nested/project/internal/cells/example/service.go"},
		"/repo",
		"/repo/nested/project",
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
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
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
