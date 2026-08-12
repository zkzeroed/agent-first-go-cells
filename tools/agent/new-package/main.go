// Command new-package scaffolds an exportable, manifest-backed Go package.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/projectconfig"
)

var (
	idPattern      = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	packagePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
)

func main() {
	root := flag.String("root", ".", "project root")
	id := flag.String("id", "", "library package id")
	path := flag.String("path", "", "package path relative to root")
	packageName := flag.String("package", "", "Go package name")
	flag.Parse()
	if err := scaffold(*root, *id, *path, *packageName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//nolint:cyclop // Scaffolding is intentionally an ordered set of safety checks and writes.
func scaffold(root, id, path, packageName string) error {
	if !idPattern.MatchString(id) {
		return errors.New("id must be kebab-case")
	}
	if path == "" {
		path = id
	}
	if packageName == "" && path != "." {
		packageName = strings.ReplaceAll(filepath.Base(path), "-", "")
	}
	if !packagePattern.MatchString(packageName) {
		return errors.New("package must be a lowercase Go identifier; specify --package for PATH=")
	}
	validationPath := "./" + path + "/..."
	if path == "." {
		validationPath = "./..."
	}
	config, err := projectconfig.Load(root)
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return err
	}
	if path == "" || filepath.IsAbs(path) || filepath.Clean(path) != path || path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator)) {
		return fmt.Errorf("invalid package path %q", path)
	}
	dir := filepath.Join(root, path)
	for _, name := range []string{"cell.yaml", "AGENTS.md", "doc.go", "api.go", "api_test.go"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return fmt.Errorf("refusing to overwrite %s", filepath.Join(path, name))
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	files := map[string]string{
		"cell.yaml":   fmt.Sprintf("id: %s\nkind: library-package\npublic: true\npurpose: \"TODO: describe this public library package\"\nentrypoints:\n  - file: api.go\n    symbol: Service\ndependencies: []\nvalidation:\n  - go test %s\nconformance:\n  basis: engineering-primitive\n  status: placeholder\n  evidence: unverified\n  rationale: \"TODO: state the engineering basis and research boundary.\"\n  gaps:\n    - \"TODO: replace this placeholder with the remaining conformance work.\"\n", id, validationPath),
		"doc.go":      fmt.Sprintf("// Package %s implements the %s library package.\npackage %s\n", packageName, id, packageName),
		"api.go":      fmt.Sprintf("package %s\n\n// Service is the public contract for this package.\ntype Service = any\n", packageName),
		"api_test.go": fmt.Sprintf("package %s\n\nimport \"testing\"\n\nfunc TestPackage(t *testing.T) {}\n", packageName),
		"AGENTS.md":   fmt.Sprintf("# Cell: %s\n\n## Purpose\n\nTODO: describe this public package.\n\n## Start Here\n\nRead `cell.yaml`, `api.go`, and tests.\n\n## Invariants\n\n- Keep exported API stable and do not expose private implementation types.\n\n## Common Tasks\n\nDeclare every public-package or private-cell dependency in `cell.yaml`.\n\n## Reliability & Concurrency\n\nDocument concurrency and error behavior for exported operations.\n\n## Validation\n\n- `go test %s`\n", id, validationPath),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	if err := projectconfig.RegisterLibrary(root, id, path); err != nil {
		return err
	}
	fmt.Printf("Created library package %s at %s\n", id, path)
	return nil
}
