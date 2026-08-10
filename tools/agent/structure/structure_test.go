// Package structure contains AST-based architecture invariant tests.
//
// These tests enforce the structural rules of the Agent-First Go architecture.
// They run via `task structure-test` and as part of `task doctor`.
//
// WHEN TO USE: In CI and before committing. These are the guardrail tests
// that prevent architecture drift.
//
// WHY IT HELPS AGENTS: An agent can't accidentally create a 500-line file,
// add an init() function, or skip required cell files — these tests catch
// all of those violations deterministically.
package structure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// projectRoot finds the repository root by searching for go.mod.
func projectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// shouldSkip returns true for files that should not be checked for size/function limits.
func shouldSkip(path string) bool {
	return strings.HasPrefix(path, "vendor/") ||
		strings.HasPrefix(path, "tools/agent/") ||
		strings.HasSuffix(path, "_test.go")
}

// TestCellImportsUseAPIPackages verifies that cell-to-cell imports cross an
// explicit contract boundary instead of reaching implementation packages.
func TestCellImportsUseAPIPackages(t *testing.T) {
	root := projectRoot()
	if !dirExists(filepath.Join(root, "internal/cells")) {
		t.Skip("no internal/cells directory")
	}

	violations := 0

	filepath.WalkDir(filepath.Join(root, "internal/cells"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			if strings.Contains(importPath, "/internal/cells/") && !strings.HasSuffix(importPath, "/api") {
				t.Errorf("%s imports %s; cell dependencies must use a cell api package", path, importPath)
				violations++
			}
		}
		return nil
	})

	if violations > 0 {
		t.Errorf("Found %d cell API import violation(s)", violations)
	}
}

// TestNoInitFunctions verifies that no init() functions exist anywhere.
func TestNoInitFunctions(t *testing.T) {
	root := projectRoot()
	violations := 0

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		if strings.HasPrefix(relPath, "vendor/") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fn.Name.Name == "init" && fn.Recv == nil {
				t.Errorf("%s contains init() function — no init() allowed", relPath)
				violations++
			}
		}
		return nil
	})

	if violations > 0 {
		t.Errorf("Found %d init() function(s)", violations)
	}
}

// TestFilesUnder300LOC verifies that no Go file exceeds 300 lines.
func TestFilesUnder300LOC(t *testing.T) {
	root := projectRoot()
	violations := 0

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		if shouldSkip(relPath) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Count(string(content), "\n") + 1
		if lines > 300 {
			t.Errorf("%s is %d lines (max 300)", relPath, lines)
			violations++
		}
		return nil
	})

	if violations > 0 {
		t.Errorf("Found %d file(s) exceeding 300 LOC", violations)
	}
}

// TestNoFunctionOver40Lines verifies that no function exceeds 40 lines (AST-based).
func TestNoFunctionOver40Lines(t *testing.T) {
	root := projectRoot()
	violations := 0

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		if shouldSkip(relPath) {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			start := fset.Position(fn.Body.Pos()).Line
			end := fset.Position(fn.Body.End()).Line
			bodyLines := end - start - 1 // exclude braces

			if bodyLines > 40 {
				overridden := false
				if fn.Doc != nil {
					for _, comment := range fn.Doc.List {
						overridden = overridden || strings.Contains(comment.Text, "AGENT_OVERRIDE")
					}
				}
				if overridden {
					continue
				}
				t.Errorf("%s: function %s is %d lines (max 40)", relPath, fn.Name.Name, bodyLines)
				violations++
			}
		}
		return nil
	})

	if violations > 0 {
		t.Errorf("Found %d function(s) exceeding 40 lines", violations)
	}
}

// TestEveryCellHasAPIPackage verifies that every manifest-backed cell exposes
// a concrete, importable public contract package.
func TestEveryCellHasAPIPackage(t *testing.T) {
	root := projectRoot()
	if !dirExists(filepath.Join(root, "internal/cells")) {
		t.Skip("no internal/cells directory")
	}

	missing := 0
	filepath.WalkDir(filepath.Join(root, "internal/cells"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "cell.yaml" {
			return nil
		}
		apiPath := filepath.Join(filepath.Dir(path), "api", "api.go")
		if _, err := os.Stat(apiPath); err != nil {
			relPath, _ := filepath.Rel(root, apiPath)
			t.Errorf("%s is missing", relPath)
			missing++
		}
		return nil
	})
	if missing > 0 {
		t.Errorf("Found %d cell(s) missing an api package", missing)
	}
}

// TestEveryCellHasCellYaml verifies that every cell directory has a cell.yaml.
func TestEveryCellHasCellYaml(t *testing.T) {
	root := projectRoot()
	if !dirExists(filepath.Join(root, "internal/cells")) {
		t.Skip("no internal/cells directory")
	}

	missing := 0

	filepath.WalkDir(filepath.Join(root, "internal/cells"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		if path == filepath.Join(root, "internal/cells") {
			return nil
		}

		// Check if this is a leaf cell dir (has .go files) or a domain dir (has subdirs)
		hasGoFiles := false
		hasSubdirs := false
		filepath.WalkDir(path, func(subpath string, subd fs.DirEntry, suberr error) error {
			if suberr != nil {
				return nil
			}
			if subpath == path {
				return nil
			}
			if subd.IsDir() {
				hasSubdirs = true
			}
			if strings.HasSuffix(subd.Name(), ".go") {
				hasGoFiles = true
			}
			return nil
		})

		// Only check dirs that have .go files (actual cell dirs)
		if hasGoFiles {
			yamlPath := filepath.Join(path, "cell.yaml")
			if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
				relPath, _ := filepath.Rel(root, path)
				t.Errorf("%s is missing cell.yaml", relPath)
				missing++
			}
		}

		if hasSubdirs {
			return filepath.SkipDir
		}
		return nil
	})

	if missing > 0 {
		t.Errorf("Found %d cell(s) missing cell.yaml", missing)
	}
}

// TestEveryCellHasDocGo verifies that every cell directory has a doc.go.
func TestEveryCellHasDocGo(t *testing.T) {
	root := projectRoot()
	if !dirExists(filepath.Join(root, "internal/cells")) {
		t.Skip("no internal/cells directory")
	}

	missing := 0

	filepath.WalkDir(filepath.Join(root, "internal/cells"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		if path == filepath.Join(root, "internal/cells") {
			return nil
		}

		hasGoFiles := false
		hasSubdirs := false
		filepath.WalkDir(path, func(subpath string, subd fs.DirEntry, suberr error) error {
			if suberr != nil {
				return nil
			}
			if subpath == path {
				return nil
			}
			if subd.IsDir() {
				hasSubdirs = true
			}
			if strings.HasSuffix(subd.Name(), ".go") {
				hasGoFiles = true
			}
			return nil
		})

		if hasGoFiles {
			docPath := filepath.Join(path, "doc.go")
			if _, err := os.Stat(docPath); os.IsNotExist(err) {
				relPath, _ := filepath.Rel(root, path)
				t.Errorf("%s is missing doc.go", relPath)
				missing++
			}
		}

		if hasSubdirs {
			return filepath.SkipDir
		}
		return nil
	})

	if missing > 0 {
		t.Errorf("Found %d cell(s) missing doc.go", missing)
	}
}

// TestEveryCellHasAgentsMd verifies that every cell directory has an AGENTS.md.
func TestEveryCellHasAgentsMd(t *testing.T) {
	root := projectRoot()
	if !dirExists(filepath.Join(root, "internal/cells")) {
		t.Skip("no internal/cells directory")
	}

	missing := 0

	filepath.WalkDir(filepath.Join(root, "internal/cells"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		if path == filepath.Join(root, "internal/cells") {
			return nil
		}

		hasGoFiles := false
		hasSubdirs := false
		filepath.WalkDir(path, func(subpath string, subd fs.DirEntry, suberr error) error {
			if suberr != nil {
				return nil
			}
			if subpath == path {
				return nil
			}
			if subd.IsDir() {
				hasSubdirs = true
			}
			if strings.HasSuffix(subd.Name(), ".go") {
				hasGoFiles = true
			}
			return nil
		})

		if hasGoFiles {
			agentsPath := filepath.Join(path, "AGENTS.md")
			if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
				relPath, _ := filepath.Rel(root, path)
				t.Errorf("%s is missing AGENTS.md", relPath)
				missing++
			}
		}

		if hasSubdirs {
			return filepath.SkipDir
		}
		return nil
	})

	if missing > 0 {
		t.Errorf("Found %d cell(s) missing AGENTS.md", missing)
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
