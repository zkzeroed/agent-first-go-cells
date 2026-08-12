// Package main implements the import policy validator.
//
// Usage:
//
//	task policy
//
// Parses policy/imports.yaml and walks the import graph of all .go files
// using go/ast. Fails on any import that violates the declared rules.
//
// WHEN TO USE: In CI and before committing. Also runs as part of `task doctor`.
//
// WHY IT HELPS AGENTS: Prevents an agent from accidentally importing
// internal/app/wiring.go from a cell, or importing internal/ from pkg/.
// The policy is machine-readable (YAML) and enforced programmatically (AST),
// so violations are caught deterministically — not by code review.
package main

import (
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/projectconfig"
)

type Rule struct {
	From  string
	Allow []string
	Deny  []string
}

//nolint:cyclop // This is the command's linear validation workflow; helpers would only obscure its stages.
func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	rules, err := parsePolicy(filepath.Join(*root, "policy", "imports.yaml"))
	if errors.Is(err, os.ErrNotExist) {
		rules = nil
		err = nil
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing policy: %v\n", err)
		os.Exit(1)
	}

	modulePath, err := readModulePath(filepath.Join(*root, "go.mod"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading module path: %v\n", err)
		os.Exit(1)
	}
	config, err := projectconfig.Load(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading architecture config: %v\n", err)
		os.Exit(1)
	}
	violations := 0

	if err := filepath.WalkDir(*root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		relPath, err := filepath.Rel(*root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if strings.HasPrefix(relPath, "vendor/") || strings.HasPrefix(relPath, "tools/agent/") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			normalized := normalizeImport(importPath, modulePath)
			if violation := validateCellBoundary(relPath, normalized, config); violation != "" {
				fmt.Printf("VIOLATION: %s imports %s (%s)\n", relPath, importPath, violation)
				violations++
			}
			for _, rule := range rules {
				if matchGlob(rule.From, relPath) {
					if violation := validateImport(rule, normalized); violation != "" {
						fmt.Printf("VIOLATION: %s imports %s (%s)\n", relPath, importPath, violation)
						violations++
					}
				}
			}
		}
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error walking source files: %v\n", err)
		os.Exit(1)
	}

	if violations > 0 {
		fmt.Printf("\n%d import policy violation(s) found.\n", violations)
		os.Exit(1)
	}

	if len(rules) == 0 {
		fmt.Println("✓ Built-in cell and library import boundaries clean (no custom policy).")
		return
	}
	fmt.Println("✓ Import policy clean.")
}

func validateCellBoundary(source, importPath string, config projectconfig.Config) string {
	// Tests may construct implementations to exercise a unit in isolation.
	// Production source remains subject to the explicit package boundaries below.
	if strings.HasSuffix(source, "_test.go") {
		return ""
	}
	cellsPrefix := projectconfig.CellsRoot + "/"
	if strings.HasPrefix(source, cellsPrefix) {
		if !strings.HasPrefix(importPath, cellsPrefix) {
			return ""
		}
		if strings.HasSuffix(importPath, "/api") {
			return ""
		}
		return "cell dependencies must use an api package"
	}
	if isWiringFile(source) && strings.HasPrefix(importPath, cellsPrefix) {
		return ""
	}
	for _, packageDir := range config.LibraryPackages {
		if projectconfig.IsPackageFile(source, packageDir) {
			if !strings.HasPrefix(importPath, "internal/") {
				return ""
			}
			if strings.HasPrefix(importPath, cellsPrefix) {
				return ""
			}
			return "library packages may import only declared private cell implementations"
		}
	}
	if !strings.HasPrefix(importPath, cellsPrefix) {
		return ""
	}
	return "only configured library packages may import private cell implementations"
}

func isWiringFile(path string) bool {
	return filepath.Dir(path) == "internal/app" &&
		strings.HasPrefix(filepath.Base(path), "wiring") &&
		strings.HasSuffix(path, ".go")
}

func parsePolicy(path string) ([]Rule, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	var rules []Rule
	current := Rule{}
	inAllow := false
	inDeny := false

	for line := range strings.SplitSeq(string(content), "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "- from:") {
			if current.From != "" {
				rules = append(rules, current)
			}
			current = Rule{
				From: strings.TrimSpace(strings.TrimPrefix(trimmed, "- from:")),
			}
			inAllow = false
			inDeny = false
		} else if strings.HasPrefix(trimmed, "allow:") {
			inAllow = true
			inDeny = false
		} else if strings.HasPrefix(trimmed, "deny:") {
			inDeny = true
			inAllow = false
		} else if val, ok := strings.CutPrefix(trimmed, "- "); ok {
			val = strings.TrimSpace(val)
			val = strings.TrimSpace(strings.Split(val, "#")[0])
			if val == "" {
				continue
			}
			if inAllow {
				current.Allow = append(current.Allow, val)
			} else if inDeny {
				current.Deny = append(current.Deny, val)
			}
		} else if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			inAllow = false
			inDeny = false
		}
	}

	if current.From != "" {
		rules = append(rules, current)
	}

	return rules, nil
}

func matchGlob(pattern, path string) bool {
	if before, after, ok := strings.Cut(pattern, "**"); ok {
		return strings.HasPrefix(path, before) && strings.HasSuffix(path, after)
	}
	matched, err := pathpkg.Match(pattern, path)
	return err == nil && matched
}

func validateImport(rule Rule, importPath string) string {
	if matches(rule.Deny, importPath) {
		return fmt.Sprintf("denied by rule from %s", rule.From)
	}
	if strings.HasPrefix(importPath, "internal/") && len(rule.Allow) > 0 && !matches(rule.Allow, importPath) {
		return fmt.Sprintf("not allowed by rule from %s", rule.From)
	}
	return ""
}

func matches(patterns []string, value string) bool {
	for _, pattern := range patterns {
		if matchGlob(pattern, value) {
			return true
		}
	}
	return false
}

func normalizeImport(importPath, modulePath string) string {
	if rest, ok := strings.CutPrefix(importPath, modulePath+"/"); ok {
		return rest
	}
	return importPath
}

func readModulePath(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("module directive not found in %s", filename)
}
