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

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/projectconfig"
)

var cellIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*(/[a-z][a-z0-9]*(-[a-z0-9]+)*)?$`)

// Manifest is the validated source of truth for a cell.
type Manifest struct {
	ID            string
	Kind          string
	Public        bool
	Purpose       string
	Entrypoints   []string
	Symbols       []string
	Dependencies  []string
	Validation    []string
	Invariants    []string
	Conformance   Conformance
	RawContent    string
	AgentsContent string
	Dir           string
}

// Conformance records the research or engineering basis for a public library
// package, its implementation status, citations, and known limitations.
type Conformance struct {
	Basis     string
	Status    string
	Evidence  string
	Citations []Citation
	Rationale string
	Gaps      []string
}

// Citation identifies a local research source and the exported symbols it
// informs.
type Citation struct {
	File    string
	Locator CitationLocator
	Symbols []string
}

// CitationLocator identifies the relevant location in a local evidence file.
type CitationLocator struct {
	Type    string
	Pages   []uint
	Heading string
}

type document struct {
	ID          string `yaml:"id"`
	Kind        string `yaml:"kind"`
	Public      bool   `yaml:"public"`
	Purpose     string `yaml:"purpose"`
	Entrypoints []struct {
		File   string `yaml:"file"`
		Symbol string `yaml:"symbol"`
	} `yaml:"entrypoints"`
	Dependencies []string `yaml:"dependencies"`
	Validation   []string `yaml:"validation"`
	Invariants   []string `yaml:"invariants"`
	Conformance  struct {
		Basis     string `yaml:"basis"`
		Status    string `yaml:"status"`
		Evidence  string `yaml:"evidence"`
		Citations []struct {
			File    string `yaml:"file"`
			Locator struct {
				Type    string `yaml:"type"`
				Pages   []uint `yaml:"pages"`
				Heading string `yaml:"heading"`
			} `yaml:"locator"`
			Symbols []string `yaml:"symbols"`
		} `yaml:"citations"`
		Rationale string   `yaml:"rationale"`
		Gaps      []string `yaml:"gaps"`
	} `yaml:"conformance"`
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
	return FindAllAtWith(root, nil)
}

// FindAllAtWith reads manifests, adds extra manifests, then validates their
// combined dependency graph. Callers use extras only for removed cells that
// must remain visible while analyzing a Git diff.
func FindAllAtWith(root string, extras []Manifest) ([]Manifest, error) {
	config, err := projectconfig.Load(root)
	if err != nil {
		return nil, err
	}
	results, err := readAllAt(root, config)
	if err != nil {
		return nil, err
	}
	results = append(results, extras...)
	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	if err := ValidateWithConfig(results, config); err != nil {
		return nil, err
	}
	return results, nil
}

//nolint:cyclop // The walk callback validates each required manifest sidecar in one place.
func readAllAt(root string, config projectconfig.Config) ([]Manifest, error) {
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
	cellsDir := filepath.Join(root, projectconfig.CellsRoot)
	var results []Manifest
	if _, err := os.Stat(cellsDir); err != nil && !errors.Is(err, fs.ErrNotExist) {
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
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	for id, dir := range config.LibraryPackages {
		path := filepath.Join(root, dir, "cell.yaml")
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read library manifest %s: %w", path, err)
		}
		m, err := Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		if m.ID != id {
			return nil, fmt.Errorf("library package registry id %q does not match manifest %q", id, m.ID)
		}
		m.Dir = filepath.ToSlash(dir)
		agents, err := os.ReadFile(filepath.Join(root, dir, "AGENTS.md"))
		if err != nil {
			return nil, fmt.Errorf("read library guide %s: %w", dir, err)
		}
		m.AgentsContent = string(agents)
		results = append(results, m)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
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

	m := Manifest{ID: doc.ID, Kind: doc.Kind, Public: doc.Public, Purpose: doc.Purpose, Dependencies: doc.Dependencies, Validation: doc.Validation, Invariants: doc.Invariants, Conformance: Conformance{Basis: doc.Conformance.Basis, Status: doc.Conformance.Status, Evidence: doc.Conformance.Evidence, Rationale: doc.Conformance.Rationale, Gaps: doc.Conformance.Gaps}, RawContent: content}
	for _, citation := range doc.Conformance.Citations {
		m.Conformance.Citations = append(m.Conformance.Citations, Citation{File: citation.File, Locator: CitationLocator{Type: citation.Locator.Type, Pages: citation.Locator.Pages, Heading: citation.Locator.Heading}, Symbols: citation.Symbols})
	}
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
		if err := validateConformanceEvidence(root, m); err != nil {
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
	dirs := []string{filepath.Join(root, m.Dir)}
	if m.Kind != "library-package" {
		dirs = append(dirs, filepath.Join(root, m.Dir, "api"))
	}
	for _, dir := range dirs {
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
		path := filepath.ToSlash(m.Dir)
		if m.Kind != "library-package" {
			path = filepath.ToSlash(filepath.Join(m.Dir, "api"))
		}
		if importPath == modulePath+"/"+path {
			return m.ID, true
		}
		if m.Kind != "library-package" && importPath == modulePath+"/"+filepath.ToSlash(m.Dir) {
			return m.ID, true
		}
	}
	return "", false
}

// Validate validates relationships across a complete manifest collection.
func Validate(manifests []Manifest) error {
	return ValidateWithConfig(manifests, projectconfig.Config{})
}

// ValidateWithConfig validates private cells and registered public packages.
func ValidateWithConfig(manifests []Manifest, config projectconfig.Config) error {
	known := make(map[string]struct{}, len(manifests))
	for _, m := range manifests {
		if _, exists := known[m.ID]; exists {
			return fmt.Errorf("duplicate cell id %q", m.ID)
		}
		known[m.ID] = struct{}{}
		if m.Kind == "library-package" {
			path, ok := config.LibraryPackages[m.ID]
			if !ok || filepath.Clean(m.Dir) != filepath.Clean(path) || !m.Public {
				return fmt.Errorf("library cell %q must be public and located at its registered package", m.ID)
			}
		} else if m.Dir != "" && filepath.Clean(m.Dir) != filepath.Join(projectconfig.CellsRoot, filepath.FromSlash(m.ID)) {
			return fmt.Errorf("cell %q must be located at %s/%s, found %s", m.ID, projectconfig.CellsRoot, m.ID, m.Dir)
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
	if m.Kind != "" && m.Kind != "library-package" {
		return fmt.Errorf("unsupported kind %q", m.Kind)
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
	if m.Kind == "library-package" {
		if err := validateConformance(m.Conformance); err != nil {
			return err
		}
	}
	return nil
}

//nolint:cyclop // Each branch is a distinct schema rule with a specific diagnostic.
func validateConformance(value Conformance) error {
	if !slices.Contains([]string{"paper-defined-math", "fixed-profile-policy", "engineering-primitive"}, value.Basis) {
		return fmt.Errorf("conformance basis %q is unsupported", value.Basis)
	}
	if !slices.Contains([]string{"conformant", "simplified", "naive", "stub", "placeholder", "gap"}, value.Status) {
		return fmt.Errorf("conformance status %q is unsupported", value.Status)
	}
	if !slices.Contains([]string{"verified", "unverified"}, value.Evidence) {
		return fmt.Errorf("conformance evidence %q is unsupported", value.Evidence)
	}
	if value.Status == "conformant" && value.Evidence != "verified" {
		return errors.New("conformant status requires verified evidence")
	}
	if value.Basis == "paper-defined-math" && len(value.Citations) == 0 {
		return errors.New("paper-defined-math conformance requires citations")
	}
	if value.Basis != "paper-defined-math" && strings.TrimSpace(value.Rationale) == "" {
		return errors.New("non-paper conformance requires rationale")
	}
	if value.Status != "conformant" && len(value.Gaps) == 0 {
		return errors.New("non-conformant status requires gaps")
	}
	for _, citation := range value.Citations {
		if citation.File == "" || len(citation.Symbols) == 0 {
			return errors.New("citations require file and symbols")
		}
		switch citation.Locator.Type {
		case "pdf-pages":
			if len(citation.Locator.Pages) == 0 || citation.Locator.Heading != "" {
				return errors.New("pdf-pages citations require pages and no heading")
			}
			for index, page := range citation.Locator.Pages {
				if page == 0 || index > 0 && page <= citation.Locator.Pages[index-1] {
					return errors.New("citation pages must be positive and strictly increasing")
				}
			}
		case "markdown-heading":
			if citation.Locator.Heading == "" || len(citation.Locator.Pages) != 0 {
				return errors.New("markdown-heading citations require a heading and no pages")
			}
		default:
			return fmt.Errorf("citation locator type %q is unsupported", citation.Locator.Type)
		}
	}
	return nil
}

func validateConformanceEvidence(root string, m Manifest) error {
	if m.Conformance.Evidence != "verified" {
		return nil
	}
	exports, err := libraryExports(root, m)
	if err != nil {
		return err
	}
	for _, citation := range m.Conformance.Citations {
		path, err := evidencePath(root, citation.File)
		if err != nil {
			return fmt.Errorf("cell %q citation %q: %w", m.ID, citation.File, err)
		}
		switch citation.Locator.Type {
		case "pdf-pages":
			if strings.ToLower(filepath.Ext(path)) != ".pdf" {
				return fmt.Errorf("cell %q citation %q: pdf-pages locator requires a PDF file", m.ID, citation.File)
			}
		case "markdown-heading":
			if filepath.Ext(path) != ".md" {
				return fmt.Errorf("cell %q citation %q: markdown-heading locator requires a Markdown file", m.ID, citation.File)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cell %q citation %q: read: %w", m.ID, citation.File, err)
			}
			if !hasMarkdownHeading(string(content), citation.Locator.Heading) {
				return fmt.Errorf("cell %q citation %q: heading %q not found", m.ID, citation.File, citation.Locator.Heading)
			}
		}
		for _, symbol := range citation.Symbols {
			if _, found := exports[symbol]; !found {
				return fmt.Errorf("cell %q citation %q: symbol %q is not exported by the library package", m.ID, citation.File, symbol)
			}
		}
	}
	return nil
}

func libraryExports(root string, m Manifest) (map[string]struct{}, error) {
	if m.Kind != "library-package" {
		return nil, nil
	}
	entries, err := os.ReadDir(filepath.Join(root, m.Dir))
	if err != nil {
		return nil, fmt.Errorf("cell %q: read library package: %w", m.ID, err)
	}
	exports := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(root, m.Dir, entry.Name()), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("cell %q: parse library package: %w", m.ID, err)
		}
		for _, declaration := range file.Decls {
			for _, symbol := range declaredExportedSymbols(declaration) {
				exports[symbol] = struct{}{}
			}
		}
	}
	return exports, nil
}

func declaredExportedSymbols(declaration ast.Decl) []string {
	switch declaration := declaration.(type) {
	case *ast.FuncDecl:
		if declaration.Recv == nil && token.IsExported(declaration.Name.Name) {
			return []string{declaration.Name.Name}
		}
	case *ast.GenDecl:
		var result []string
		for _, spec := range declaration.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				if token.IsExported(spec.Name.Name) {
					result = append(result, spec.Name.Name)
				}
			case *ast.ValueSpec:
				for _, name := range spec.Names {
					if token.IsExported(name.Name) {
						result = append(result, name.Name)
					}
				}
			}
		}
		return result
	}
	return nil
}

func evidencePath(root, name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", errors.New("file must be a project-relative path")
	}
	path := filepath.Clean(filepath.Join(root, name))
	base := filepath.Clean(root) + string(filepath.Separator)
	if !strings.HasPrefix(path, base) {
		return "", errors.New("file must stay within the project")
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("evidence file does not exist: %w", err)
	}
	return path, nil
}

func hasMarkdownHeading(content, heading string) bool {
	for line := range strings.SplitSeq(content, "\n") {
		if strings.TrimSpace(strings.TrimLeft(line, "#")) == heading && strings.HasPrefix(line, "#") {
			return true
		}
	}
	return false
}

// SourceContent returns the ordered source inputs used to generate a context pack.
func (m Manifest) SourceContent() []byte {
	return bytes.Join([][]byte{[]byte(m.RawContent), []byte(m.AgentsContent)}, []byte("\x00"))
}
