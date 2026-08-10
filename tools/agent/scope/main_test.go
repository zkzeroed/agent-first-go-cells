package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func TestAllowedPathsRejectsInvalidScopeExpansion(t *testing.T) {
	manifests := fixtures{{id: "orders"}, {id: "orders-create"}}.writeAll(t)
	if _, err := allowedPaths(manifests[0], manifests, "orders,orders"); err == nil {
		t.Fatal("allowedPaths() error = nil, want duplicate error")
	}
	if _, err := allowedPaths(manifests[0], manifests, "missing"); err == nil {
		t.Fatal("allowedPaths() error = nil, want unknown error")
	}
	if _, err := allowedPaths(manifests[0], manifests, "@unknown"); err == nil {
		t.Fatal("allowedPaths() error = nil, want unknown scope error")
	}
}

func TestVerifyScopeAllowsExactCellAndGeneratedFiles(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}, {id: "orders-create"}})
	write(t, root, "internal/cells/orders/service.go", "package orders\n")
	write(t, root, "gen/cells.json", "{}\n")
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{}; !slices.Equal(result.OutOfScope, want) {
		t.Errorf("OutOfScope = %v, want %v", result.OutOfScope, want)
	}
}

func TestEntrypointsIncludeDeclaredSymbol(t *testing.T) {
	result := entrypoints(manifest.Manifest{Entrypoints: []string{"api/api.go"}, Symbols: []string{"Service"}})
	if want := []string{"api/api.go: Service"}; !slices.Equal(result, want) {
		t.Errorf("entrypoints() = %v, want %v", result, want)
	}
}

func TestVerifyScopeExcludesNestedCellsAndRequiresExplicitSharedScope(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}, {id: "orders/create"}})
	write(t, root, "internal/cells/orders/create/service.go", "package create\n")
	write(t, root, "internal/contracts/clock.go", "package contracts\n")
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"internal/cells/orders/create/service.go", "internal/contracts/clock.go"}
	if !slices.Equal(result.OutOfScope, want) {
		t.Errorf("OutOfScope = %v, want %v", result.OutOfScope, want)
	}
	result, err = buildReport(arguments{cellID: "orders", root: root, with: "orders/create,@contracts", verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.OutOfScope) != 0 {
		t.Errorf("OutOfScope = %v, want none", result.OutOfScope)
	}
}

func TestVerifyScopeTracksUntrackedFilesOfAnyExtension(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}})
	write(t, root, "notes.md", "out of scope\n")
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"notes.md"}; !slices.Equal(result.OutOfScope, want) {
		t.Errorf("OutOfScope = %v, want %v", result.OutOfScope, want)
	}
}

func TestVerifyScopeRejectsRootFileNamedLikeTheTargetCell(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}})
	write(t, root, "orders", "not a cell file\n")
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"orders"}; !slices.Equal(result.OutOfScope, want) {
		t.Errorf("OutOfScope = %v, want %v", result.OutOfScope, want)
	}
}

func TestVerifyScopeAllowsRemovingItsTargetCell(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}, {id: "orders-render", dependencies: []string{"orders"}}})
	for _, name := range []string{"internal/cells/orders/cell.yaml", "internal/cells/orders/AGENTS.md"} {
		if err := os.Remove(filepath.Join(root, name)); err != nil {
			t.Fatal(err)
		}
	}
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.OutOfScope) != 0 {
		t.Errorf("OutOfScope = %v, want none", result.OutOfScope)
	}
}

func TestVerifyScopeRequiresChildScopeWhenRemovingNestedCells(t *testing.T) {
	root := scopeFixture(t, []manifestFixture{{id: "orders"}, {id: "orders/create"}})
	if err := os.RemoveAll(filepath.Join(root, "internal", "cells", "orders")); err != nil {
		t.Fatal(err)
	}
	result, err := buildReport(arguments{cellID: "orders", root: root, verify: true})
	if err != nil {
		t.Fatal(err)
	}
	childManifest := "internal/cells/orders/create/cell.yaml"
	if !slices.Contains(result.OutOfScope, childManifest) {
		t.Errorf("OutOfScope = %v, want %s", result.OutOfScope, childManifest)
	}
	result, err = buildReport(arguments{cellID: "orders", root: root, with: "orders/create", verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.OutOfScope) != 0 {
		t.Errorf("OutOfScope = %v, want none", result.OutOfScope)
	}
}

type manifestFixture struct {
	id           string
	dependencies []string
}

type fixtures []manifestFixture

func (fixtures fixtures) writeAll(t *testing.T) []manifest.Manifest {
	t.Helper()
	root := scopeFixture(t, fixtures)
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	return manifests
}

func scopeFixture(t *testing.T, fixtures []manifestFixture) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init")
	git(t, root, "config", "user.email", "test@example.com")
	git(t, root, "config", "user.name", "Scope Test")
	for _, fixture := range fixtures {
		writeCell(t, root, fixture.id, fixture.dependencies)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "initial")
	return root
}

func git(t *testing.T, directory string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func writeCell(t *testing.T, root, id string, dependencies []string) {
	t.Helper()
	dependenciesYAML := "dependencies: []\n"
	if len(dependencies) > 0 {
		dependenciesYAML = "dependencies:\n"
		for _, dependency := range dependencies {
			dependenciesYAML += "  - " + dependency + "\n"
		}
	}
	write(t, root, filepath.Join("internal", "cells", id, "cell.yaml"), "id: "+id+"\npurpose: Test cell\nentrypoints:\n  - file: api/api.go\n    symbol: Service\n"+dependenciesYAML+"validation:\n  - go test ./...\n")
	write(t, root, filepath.Join("internal", "cells", id, "AGENTS.md"), "# Guide\n")
}

func write(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
