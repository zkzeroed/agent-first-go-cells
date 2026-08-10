// Package manifest provides parsing and validation for cell manifests.
package manifest

import (
	"bytes"
	"errors"
	"fmt"
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
		m.Entrypoints = append(m.Entrypoints, entrypoint.File)
	}
	if err := validateOne(m); err != nil {
		return Manifest{}, err
	}
	return m, nil
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
