// Package manifest provides parsing and validation for cell manifests.
package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var cellIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*(/[a-z][a-z0-9]*(-[a-z0-9]+)*)?$`)

// Manifest is the validated source of truth for a cell.
type Manifest struct {
	ID            string
	Purpose       string
	Entrypoints   []string
	Symbols       []string
	Dependencies  []string
	Validation    []string
	Invariants    []string
	RawContent    string
	AgentsContent string
	Dir           string
}

type document struct {
	ID          string `yaml:"id"`
	Purpose     string `yaml:"purpose"`
	Entrypoints []struct {
		File   string `yaml:"file"`
		Symbol string `yaml:"symbol"`
	} `yaml:"entrypoints"`
	Dependencies []string `yaml:"dependencies"`
	Validation   []string `yaml:"validation"`
	Invariants   []string `yaml:"invariants"`
}

// FindAll reads every manifest and its adjacent AGENTS.md, then validates the
// collection-level rules that make dependency references unambiguous.
func FindAll() ([]Manifest, error) {
	return FindAllAt(".")
}

// FindAllAt reads and validates every manifest beneath root. Manifest paths in
// the returned collection remain relative to root so generated indexes are
// portable within the project they describe.
func FindAllAt(root string) ([]Manifest, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve project root: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat project root %s: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project root %s is not a directory", root)
	}
	cellsDir := filepath.Join(root, "internal", "cells")
	var results []Manifest
	if _, err := os.Stat(cellsDir); errors.Is(err, fs.ErrNotExist) {
		return results, nil
	} else if err != nil {
		return nil, fmt.Errorf("stat %s: %w", cellsDir, err)
	}
	err = filepath.WalkDir(cellsDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || d.Name() != "cell.yaml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		m, err := Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		cellDir := filepath.Dir(path)
		m.Dir, err = filepath.Rel(root, cellDir)
		if err != nil {
			return fmt.Errorf("make manifest path relative to project root: %w", err)
		}
		agentsPath := filepath.Join(cellDir, "AGENTS.md")
		agents, err := os.ReadFile(agentsPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", agentsPath, err)
		}
		m.AgentsContent = string(agents)
		results = append(results, m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	if err := Validate(results); err != nil {
		return nil, err
	}
	return results, nil
}

// Parse parses one manifest and validates rules that do not require other cells.
func Parse(content string) (Manifest, error) {
	decoder := yaml.NewDecoder(strings.NewReader(content))
	decoder.KnownFields(true)
	var doc document
	if err := decoder.Decode(&doc); err != nil {
		return Manifest{}, fmt.Errorf("decode YAML: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return Manifest{}, errors.New("manifest must contain exactly one YAML document")
		}
		return Manifest{}, fmt.Errorf("decode YAML: %w", err)
	}

	m := Manifest{ID: doc.ID, Purpose: doc.Purpose, Dependencies: doc.Dependencies, Validation: doc.Validation, Invariants: doc.Invariants, RawContent: content}
	for _, entrypoint := range doc.Entrypoints {
		if entrypoint.File == "" {
			return Manifest{}, errors.New("entrypoints must each specify file")
		}
		if entrypoint.Symbol == "" {
			return Manifest{}, errors.New("entrypoints must each specify symbol")
		}
		m.Entrypoints = append(m.Entrypoints, entrypoint.File)
		m.Symbols = append(m.Symbols, entrypoint.Symbol)
	}
	if err := validateOne(m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// ValidateSourceAt verifies that manifests match the cell source they describe.
// It checks declared entrypoints and direct imports of other cell API packages.
func ValidateSourceAt(root string, manifests []Manifest) error {
	if len(manifests) == 0 {
		return nil
	}
	modulePath, err := modulePathAt(root)
	if err != nil {
		return err
	}
	for _, m := range manifests {
		if err := validateEntrypoints(root, m); err != nil {
			return err
		}
		if err := validateImports(root, modulePath, m, manifests); err != nil {
			return err
		}
	}
	return nil
}

func modulePathAt(root string) (string, error) {
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("read module path: %w", err)
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", errors.New("module directive not found in go.mod")
}

func validateEntrypoints(root string, m Manifest) error {
	if len(m.Entrypoints) != len(m.Symbols) {
		return fmt.Errorf("cell %q has mismatched entrypoint files and symbols", m.ID)
	}
	for i, file := range m.Entrypoints {
		path, err := cellFilePath(root, m.Dir, file)
		if err != nil {
			return fmt.Errorf("cell %q entrypoint %q: %w", m.ID, file, err)
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return fmt.Errorf("cell %q entrypoint %q: parse Go source: %w", m.ID, file, err)
		}
		if !declaresSymbol(parsed, m.Symbols[i]) {
			return fmt.Errorf("cell %q entrypoint %q does not declare symbol %q", m.ID, file, m.Symbols[i])
		}
	}
	return nil
}

func cellFilePath(root, cellDir, name string) (string, error) {
	if filepath.IsAbs(name) || filepath.Ext(name) != ".go" {
		return "", errors.New("file must be a relative Go source path")
	}
	path := filepath.Clean(filepath.Join(root, cellDir, name))
	base := filepath.Clean(filepath.Join(root, cellDir)) + string(filepath.Separator)
	if !strings.HasPrefix(path, base) {
		return "", errors.New("file must stay within the cell directory")
	}
	return path, nil
}

func declaresSymbol(file *ast.File, symbol string) bool {
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil && decl.Name.Name == symbol {
				return true
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if namedSpec(spec, symbol) {
					return true
				}
			}
		}
	}
	return false
}

func namedSpec(spec ast.Spec, symbol string) bool {
	switch spec := spec.(type) {
	case *ast.TypeSpec:
		return spec.Name.Name == symbol
	case *ast.ValueSpec:
		return slices.ContainsFunc(spec.Names, func(name *ast.Ident) bool { return name.Name == symbol })
	default:
		return false
	}
}

func validateImports(root, modulePath string, m Manifest, manifests []Manifest) error {
	imports, err := cellImports(root, m)
	if err != nil {
		return err
	}
	used := make(map[string]struct{}, len(m.Dependencies))
	for importPath := range imports {
		dependency, found := cellDependencyID(modulePath, importPath, manifests)
		if !found || dependency == m.ID {
			continue
		}
		used[dependency] = struct{}{}
		if slices.Contains(m.Dependencies, dependency) {
			continue
		}
		return fmt.Errorf("cell %q imports dependency %q without declaring it", m.ID, dependency)
	}
	for _, dependency := range m.Dependencies {
		if _, found := used[dependency]; !found {
			return fmt.Errorf("cell %q declares dependency %q without importing its api package", m.ID, dependency)
		}
	}
	return nil
}

func cellImports(root string, m Manifest) (map[string]struct{}, error) {
	result := map[string]struct{}{}
	for _, dir := range []string{filepath.Join(root, m.Dir), filepath.Join(root, m.Dir, "api")} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("read cell source %s: %w", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
				continue
			}
			file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, entry.Name()), nil, parser.ImportsOnly)
			if err != nil {
				return nil, fmt.Errorf("parse cell source %s: %w", entry.Name(), err)
			}
			for _, imp := range file.Imports {
				result[strings.Trim(imp.Path.Value, "\"")] = struct{}{}
			}
		}
	}
	return result, nil
}

func cellDependencyID(modulePath, importPath string, manifests []Manifest) (string, bool) {
	for _, m := range manifests {
		if importPath == modulePath+"/"+filepath.ToSlash(filepath.Join(m.Dir, "api")) {
			return m.ID, true
		}
	}
	return "", false
}

// Validate validates relationships across a complete manifest collection.
func Validate(manifests []Manifest) error {
	known := make(map[string]struct{}, len(manifests))
	for _, m := range manifests {
		if _, exists := known[m.ID]; exists {
			return fmt.Errorf("duplicate cell id %q", m.ID)
		}
		known[m.ID] = struct{}{}
		if m.Dir != "" && filepath.Clean(m.Dir) != filepath.Join("internal", "cells", filepath.FromSlash(m.ID)) {
			return fmt.Errorf("cell %q must be located at internal/cells/%s, found %s", m.ID, m.ID, m.Dir)
		}
	}
	for _, m := range manifests {
		seen := map[string]struct{}{}
		for _, dependency := range m.Dependencies {
			if dependency == m.ID {
				return fmt.Errorf("cell %q cannot depend on itself", m.ID)
			}
			if _, exists := known[dependency]; !exists {
				return fmt.Errorf("cell %q dependency %q must exactly match an existing cell id", m.ID, dependency)
			}
			if _, duplicate := seen[dependency]; duplicate {
				return fmt.Errorf("cell %q lists dependency %q more than once", m.ID, dependency)
			}
			seen[dependency] = struct{}{}
		}
	}
	return nil
}

func validateOne(m Manifest) error {
	if !cellIDPattern.MatchString(m.ID) {
		return fmt.Errorf("id %q must be kebab-case, optionally with one domain/action path", m.ID)
	}
	if strings.TrimSpace(m.Purpose) == "" {
		return errors.New("purpose is required")
	}
	if len(m.Entrypoints) == 0 {
		return errors.New("at least one entrypoint is required")
	}
	if len(m.Validation) == 0 {
		return errors.New("at least one validation command is required")
	}
	if slices.Contains(m.Dependencies, "") {
		return errors.New("dependencies cannot contain an empty value")
	}
	return nil
}

// SourceContent returns the ordered source inputs used to generate a context pack.
func (m Manifest) SourceContent() []byte {
	return bytes.Join([][]byte{[]byte(m.RawContent), []byte(m.AgentsContent)}, []byte("\x00"))
}
