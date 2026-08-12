package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRejectsUnknownFieldsAndMissingEntrypointFile(t *testing.T) {
	_, err := Parse("id: orders-create\npurpose: create an order\nentrypoints:\n  - symbol: Create\ndependencies: []\nvalidation:\n  - go test ./...\nunknown: value\n")
	if err == nil {
		t.Fatal("Parse() error = nil, want schema validation error")
	}
	if !strings.Contains(err.Error(), "field unknown") {
		t.Fatalf("Parse() error = %q, want unknown field error", err)
	}
}

func TestValidateRequiresExactKnownDependency(t *testing.T) {
	manifests := []Manifest{
		{ID: "orders-create", Dir: "internal/cells/orders-create"},
		{ID: "orders-read", Dir: "internal/cells/orders-read", Dependencies: []string{"orders"}},
	}
	err := Validate(manifests)
	if err == nil {
		t.Fatal("Validate() error = nil, want unknown dependency error")
	}
	if !strings.Contains(err.Error(), "exactly match an existing cell id") {
		t.Fatalf("Validate() error = %q, want exact dependency error", err)
	}
}

func TestParsePreservesNestedEntrypointAndValidation(t *testing.T) {
	m, err := Parse("id: orders/create\npurpose: create an order\nentrypoints:\n  - file: handler.go\n    symbol: Handle\ndependencies: []\nvalidation:\n  - go test ./internal/cells/orders/create/...\ninvariants:\n  - idempotent\n")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got, want := m.Entrypoints, []string{"handler.go"}; !equalStrings(got, want) {
		t.Fatalf("Entrypoints = %#v, want %#v", got, want)
	}
	if got, want := m.Symbols, []string{"Handle"}; !equalStrings(got, want) {
		t.Fatalf("Symbols = %#v, want %#v", got, want)
	}
	if got, want := m.Validation, []string{"go test ./internal/cells/orders/create/..."}; !equalStrings(got, want) {
		t.Fatalf("Validation = %#v, want %#v", got, want)
	}
}

func TestParseValidatesLibraryConformance(t *testing.T) {
	valid := "id: field\nkind: library-package\npublic: true\npurpose: public field\nentrypoints:\n  - file: field.go\n    symbol: Service\ndependencies: []\nvalidation:\n  - go test ./field/...\nconformance:\n  basis: paper-defined-math\n  status: simplified\n  evidence: unverified\n  citations:\n    - file: docs/paper.pdf\n      locator:\n        type: pdf-pages\n        pages: [4, 5]\n      symbols: [Service]\n  gaps:\n    - missing optimized implementation\n"
	manifest, err := Parse(valid)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if manifest.Conformance.Basis != "paper-defined-math" || len(manifest.Conformance.Citations) != 1 {
		t.Fatalf("Conformance = %#v, want parsed paper record", manifest.Conformance)
	}

	missingCitation := strings.Replace(valid, "      symbols: [Service]\n", "", 1)
	if _, err := Parse(missingCitation); err == nil || !strings.Contains(err.Error(), "citations require") {
		t.Fatalf("Parse() error = %v, want citation validation error", err)
	}
	missingRationale := strings.Replace(valid, "basis: paper-defined-math", "basis: engineering-primitive", 1)
	if _, err := Parse(missingRationale); err == nil || !strings.Contains(err.Error(), "rationale") {
		t.Fatalf("Parse() error = %v, want rationale validation error", err)
	}
	withoutConformance := strings.Split(valid, "conformance:")[0]
	if _, err := Parse(withoutConformance); err == nil || !strings.Contains(err.Error(), "conformance basis") {
		t.Fatalf("Parse() error = %v, want required conformance error", err)
	}
}

func TestValidateSourceAtRejectsInvalidEntrypointsAndDependencies(t *testing.T) {
	root := t.TempDir()
	writeManifestFixture(t, root, "orders-create", nil, "package orderscreate\n\ntype Creator struct{}\n")
	writeManifestFixture(t, root, "orders-render", []string{"orders-create"}, "package ordersrender\n\ntype Creator struct{}\n")

	manifests, err := FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateSourceAt(root, manifests)
	if err == nil || !strings.Contains(err.Error(), "without importing") {
		t.Fatalf("ValidateSourceAt() error = %v, want undeclared import error", err)
	}

	writeManifestFixture(t, root, "orders-render", nil, "package ordersrender\n\nimport _ \"example.com/test/internal/cells/orders-create/api\"\n\ntype Creator struct{}\n")
	manifests, err = FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateSourceAt(root, manifests)
	if err == nil || !strings.Contains(err.Error(), "without declaring") {
		t.Fatalf("ValidateSourceAt() error = %v, want missing dependency error", err)
	}

	writeManifestFixture(t, root, "orders-render", []string{"orders-create"}, "package ordersrender\n\nimport _ \"example.com/test/internal/cells/orders-create/api\"\n\ntype Creator struct{}\n")
	manifests, err = FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSourceAt(root, manifests); err != nil {
		t.Fatalf("ValidateSourceAt() error = %v, want nil", err)
	}
}

func TestFindAllAtDiscoversRegisteredLibraryPackageAndPrivateDependency(t *testing.T) {
	root := t.TempDir()
	writeManifestFixture(t, root, "helper", nil, "package helper\n\ntype Creator struct{}\n")
	if err := os.MkdirAll(filepath.Join(root, "field"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "policy"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "policy", "architecture.yaml"), []byte("libraryPackages:\n  field: field\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "field", "cell.yaml"), []byte("id: field\nkind: library-package\npublic: true\npurpose: public field\nentrypoints:\n  - file: field.go\n    symbol: Service\ndependencies:\n  - helper\nvalidation:\n  - go test ./field/...\nconformance:\n  basis: engineering-primitive\n  status: conformant\n  evidence: verified\n  rationale: test package\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "field", "AGENTS.md"), []byte("# field\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "field", "field.go"), []byte("package field\n\nimport _ \"example.com/test/internal/cells/helper\"\n\ntype Service struct{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifests, err := FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 2 || manifests[0].ID != "field" || manifests[0].Dir != "field" {
		t.Fatalf("FindAllAt() = %#v, want registered library package", manifests)
	}
	if err := ValidateSourceAt(root, manifests); err != nil {
		t.Fatalf("ValidateSourceAt() error = %v", err)
	}
}

func TestValidateSourceAtVerifiesConformanceEvidence(t *testing.T) {
	root := t.TempDir()
	for _, directory := range []string{"field", "docs", "policy"} {
		if err := os.MkdirAll(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.com/test\n\ngo 1.26.5\n")
	write("policy/architecture.yaml", "libraryPackages:\n  field: field\n")
	write("field/AGENTS.md", "# Guide\n")
	write("field/field.go", "package field\n\nfunc Reduce(value int) int { return value }\n")
	write("docs/reduction.md", "# Euclidean remainder\n")
	write("field/cell.yaml", "id: field\nkind: library-package\npublic: true\npurpose: public field\nentrypoints:\n  - file: field.go\n    symbol: Reduce\ndependencies: []\nvalidation:\n  - go test ./field/...\nconformance:\n  basis: paper-defined-math\n  status: conformant\n  evidence: verified\n  citations:\n    - file: docs/reduction.md\n      locator:\n        type: markdown-heading\n        heading: Euclidean remainder\n      symbols: [Reduce]\n")

	manifests, err := FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSourceAt(root, manifests); err != nil {
		t.Fatalf("ValidateSourceAt() error = %v, want nil", err)
	}
	write("docs/reduction.md", "# Different heading\n")
	if err := ValidateSourceAt(root, manifests); err == nil || !strings.Contains(err.Error(), "heading") {
		t.Fatalf("ValidateSourceAt() error = %v, want missing heading error", err)
	}
	manifests[0].Conformance.Citations[0].Symbols = []string{"Missing"}
	write("docs/reduction.md", "# Euclidean remainder\n")
	if err := ValidateSourceAt(root, manifests); err == nil || !strings.Contains(err.Error(), "not exported") {
		t.Fatalf("ValidateSourceAt() error = %v, want missing symbol error", err)
	}
}

func TestValidateSourceAtRejectsMissingAndEscapingEntrypoints(t *testing.T) {
	root := t.TempDir()
	writeManifestFixture(t, root, "orders-create", nil, "package orderscreate\n\ntype Creator struct{}\n")
	manifests, err := FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	manifests[0].Entrypoints[0] = "missing.go"
	if err := ValidateSourceAt(root, manifests); err == nil || !strings.Contains(err.Error(), "parse Go source") {
		t.Fatalf("ValidateSourceAt() error = %v, want missing file error", err)
	}
	manifests[0].Entrypoints[0] = "../outside.go"
	if err := ValidateSourceAt(root, manifests); err == nil || !strings.Contains(err.Error(), "stay within") {
		t.Fatalf("ValidateSourceAt() error = %v, want path containment error", err)
	}
	manifests[0].Entrypoints[0] = "api/api.go"
	manifests[0].Symbols[0] = "Missing"
	if err := ValidateSourceAt(root, manifests); err == nil || !strings.Contains(err.Error(), "does not declare") {
		t.Fatalf("ValidateSourceAt() error = %v, want missing symbol error", err)
	}
}

func TestValidateSourceAtRejectsMismatchedEntrypointMetadata(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := Manifest{ID: "orders-create", Entrypoints: []string{"api/api.go"}}
	err := ValidateSourceAt(root, []Manifest{manifest})
	if err == nil || !strings.Contains(err.Error(), "mismatched entrypoint") {
		t.Fatalf("ValidateSourceAt() error = %v, want mismatched metadata error", err)
	}
}

func writeManifestFixture(t *testing.T, root, id string, dependencies []string, source string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/test\n\ngo 1.26.5\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "internal", "cells", id)
	if err := os.MkdirAll(filepath.Join(dir, "api"), 0o755); err != nil {
		t.Fatal(err)
	}
	dependencyLines := "dependencies: []\n"
	if len(dependencies) > 0 {
		dependencyLines = "dependencies:\n"
		for _, dependency := range dependencies {
			dependencyLines += "  - " + dependency + "\n"
		}
	}
	manifest := "id: " + id + "\npurpose: Test cell\nentrypoints:\n  - file: api/api.go\n    symbol: Creator\n" + dependencyLines + "validation:\n  - go test ./...\n"
	if err := os.WriteFile(filepath.Join(dir, "cell.yaml"), []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# Guide\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "api", "api.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
}

// FuzzParse verifies that arbitrary manifest input is rejected safely or
// produces a manifest without panicking.
func FuzzParse(f *testing.F) {
	for _, content := range []string{
		"",
		"id: orders-create\npurpose: create an order\nentrypoints:\n  - file: handler.go\ndependencies: []\nvalidation:\n  - go test ./...\n",
		"id: orders-create\nunknown: value\n",
		"---\nid: orders-create\n---\nid: orders-read\n",
	} {
		f.Add(content)
	}

	f.Fuzz(func(t *testing.T, content string) {
		_, _ = Parse(content)
	})
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
