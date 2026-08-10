package agent_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOperationalDocsAvoidRetiredArchitectureTerms(t *testing.T) {
	retired := []string{"agent-first/v1", ".devin/", ".rules", "interface files", "interface file only"}
	root := projectRoot(t)
	paths := []string{filepath.Join(root, "README.md"), filepath.Join(root, "AGENTS.md")}
	err := filepath.WalkDir(filepath.Join(root, "docs/architecture"), func(path string, d os.DirEntry, err error) error {
		if path == filepath.Join(root, "docs/architecture/18-architecture-feedback-log.md") {
			return err
		}
		if err != nil || d.IsDir() {
			return err
		}
		if strings.HasSuffix(path, ".md") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range paths {
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		for _, term := range retired {
			if strings.Contains(strings.ToLower(string(content)), term) {
				t.Errorf("%s contains retired term %q", filename, term)
			}
		}
	}
}

func TestOperationalDocsDescribeCurrentBoundaries(t *testing.T) {
	root := projectRoot(t)
	checks := map[string][]string{
		"README.md": {"api/api.go", "API contract"},
		"docs/architecture/05-manifest-schema.md":  {"exact IDs of other cells"},
		"docs/architecture/14-taskfile-targets.md": {"index-json", "cells-json", "deps-json", "impact-json", "context` operates"},
		"examples/reference-project/README.md":     {"greeting-compose", "token-issue"},
	}
	for path, required := range checks {
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		for _, phrase := range required {
			if !strings.Contains(string(content), phrase) {
				t.Errorf("%s must document %q", path, phrase)
			}
		}
	}
}

func TestOperationalDocsDescribeCurrentLinting(t *testing.T) {
	root := projectRoot(t)
	content, err := os.ReadFile(filepath.Join(root, "docs/architecture/02-core-principles.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, required := range []string{"AST-based structural tests", "cyclop", "maximum 16"} {
		if !strings.Contains(text, required) {
			t.Errorf("core principles must document %q", required)
		}
	}
	for _, retired := range []string{"funlen", "gocognit", "200 LOC warning", "25 LOC warning"} {
		if strings.Contains(text, retired) {
			t.Errorf("core principles contains retired linting claim %q", retired)
		}
	}
}

func projectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
