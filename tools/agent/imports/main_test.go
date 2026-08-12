package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/projectconfig"
)

func TestValidateImportEnforcesInternalAllowList(t *testing.T) {
	rule := Rule{From: "internal/cells/**", Allow: []string{"internal/cells/**/api"}}
	if got := validateImport(rule, "internal/cells/orders/api"); got != "" {
		t.Fatalf("validateImport(api) = %q", got)
	}
	if got := validateImport(rule, "internal/cells/orders"); got == "" {
		t.Fatal("validateImport(implementation) = nil, want violation")
	}
	if got := validateImport(rule, "net/http"); got != "" {
		t.Fatalf("validateImport(external) = %q", got)
	}
}

func TestPolicyWithoutCustomRulesStillRejectsPrivateImplementationImport(t *testing.T) {
	root := t.TempDir()
	writeImportFixture(t, root, "go.mod", "module example.com/test\n\ngo 1.26.5\n")
	writeImportFixture(t, root, "cmd/tool/main.go", "package main\n\nimport _ \"example.com/test/internal/cells/helper\"\n")
	writeImportFixture(t, root, "internal/cells/helper/helper.go", "package helper\n")

	modulePath, err := readModulePath(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if got := validateCellBoundary("cmd/tool/main.go", normalizeImport("example.com/test/internal/cells/helper", modulePath), projectconfig.Config{}); got == "" {
		t.Fatal("built-in boundary = nil, want private implementation violation")
	}
}

func TestValidateCellBoundaryAllowsTestConstruction(t *testing.T) {
	got := validateCellBoundary("internal/cells/render/handler_test.go", "internal/cells/compose", projectconfig.Config{})
	if got != "" {
		t.Fatalf("validateCellBoundary(test source) = %q, want allowed", got)
	}
}

func TestValidateCellBoundaryAllowsWiringConstruction(t *testing.T) {
	got := validateCellBoundary("internal/app/wiring_http.go", "internal/cells/compose", projectconfig.Config{})
	if got != "" {
		t.Fatalf("validateCellBoundary(wiring source) = %q, want allowed", got)
	}
	got = validateCellBoundary("internal/app/service.go", "internal/cells/compose", projectconfig.Config{})
	if got == "" {
		t.Fatal("validateCellBoundary(non-wiring app source) = nil, want violation")
	}
}

func writeImportFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeImport(t *testing.T) {
	if got, want := normalizeImport("example.com/project/internal/app", "example.com/project"), "internal/app"; got != want {
		t.Fatalf("normalizeImport() = %q, want %q", got, want)
	}
}

func TestValidateCellBoundaryTreatsEveryRegisteredLibraryPathTheSame(t *testing.T) {
	config := projectconfig.Config{LibraryPackages: map[string]string{
		"root":   ".",
		"widget": "pkg/widget",
		"field":  "field",
	}}
	for _, source := range []string{"api.go", "pkg/widget/api.go", "field/api.go"} {
		if got := validateCellBoundary(source, "internal/cells/helper", config); got != "" {
			t.Errorf("validateCellBoundary(%q, private cell) = %q, want allowed for manifest validation", source, got)
		}
		if got := validateCellBoundary(source, "internal/app", config); got == "" {
			t.Errorf("validateCellBoundary(%q, internal/app) = nil, want violation", source)
		}
		if got := validateCellBoundary(source, "internal/platform/logging", config); got == "" {
			t.Errorf("validateCellBoundary(%q, internal/platform) = nil, want violation", source)
		}
		if got := validateCellBoundary(source, "internal/contracts/clock", config); got == "" {
			t.Errorf("validateCellBoundary(%q, internal/contracts) = nil, want violation", source)
		}
	}
	if got := validateCellBoundary("cmd/tool/main.go", "internal/cells/helper", config); got == "" {
		t.Fatal("non-library package import of private cell implementation = nil, want violation")
	}
}
