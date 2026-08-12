// Package projectconfig reads project-wide architecture boundaries.
package projectconfig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	fileName  = "policy/architecture.yaml"
	CellsRoot = "internal/cells"
)

// Config declares public library-package locations. Private cells always live
// beneath CellsRoot.
type Config struct {
	LibraryPackages map[string]string `yaml:"libraryPackages"`
}

// Load returns explicit configuration or an empty library-package registry.
func Load(root string) (Config, error) {
	config := Config{}
	content, err := os.ReadFile(filepath.Join(root, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read %s: %w", fileName, err)
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(content)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", fileName, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return Config{}, fmt.Errorf("parse %s: must contain exactly one YAML document", fileName)
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

// RegisterLibrary records a public package path without replacing an existing mapping.
func RegisterLibrary(root, id, path string) error {
	config, err := Load(root)
	if err != nil {
		return err
	}
	if config.LibraryPackages == nil {
		config.LibraryPackages = make(map[string]string)
	}
	if _, exists := config.LibraryPackages[id]; exists {
		return fmt.Errorf("library package id %q is already registered", id)
	}
	for knownID, knownPath := range config.LibraryPackages {
		if knownPath == path {
			return fmt.Errorf("library package path %q is already registered to %q", path, knownID)
		}
	}
	config.LibraryPackages[id] = path
	if err := config.Validate(); err != nil {
		return err
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", fileName, err)
	}
	pathName := filepath.Join(root, fileName)
	if err := os.MkdirAll(filepath.Dir(pathName), 0o755); err != nil {
		return err
	}
	return os.WriteFile(pathName, content, 0o644)
}

// Validate rejects paths that could escape the project or expose private cells.
func (c Config) Validate() error {
	for id, path := range c.LibraryPackages {
		if id == "" {
			return errors.New("libraryPackages contains an empty id")
		}
		if err := validatePath("libraryPackages", path, true); err != nil {
			return err
		}
		if isWithin(path, CellsRoot) {
			return fmt.Errorf("library package %q must not be inside cellsRoot", path)
		}
		if hasInternalSegment(path) {
			return fmt.Errorf("library package %q must not be inside an internal directory", path)
		}
	}
	return nil
}

func validatePath(name, value string, allowRoot bool) error {
	if value == "." && allowRoot {
		return nil
	}
	if value == "" || filepath.IsAbs(value) || filepath.Clean(value) != value || value == "." || value == ".." || strings.HasPrefix(value, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%s contains invalid relative path %q", name, value)
	}
	return nil
}

func hasInternalSegment(path string) bool {
	for segment := range strings.SplitSeq(filepath.ToSlash(path), "/") {
		if segment == "internal" {
			return true
		}
	}
	return false
}

func isWithin(path, parent string) bool {
	path = strings.TrimSuffix(filepath.ToSlash(path), "/")
	parent = strings.TrimSuffix(filepath.ToSlash(parent), "/")
	return path == parent || strings.HasPrefix(path, parent+"/")
}

// IsPackageFile reports whether file belongs directly to an importable package.
func IsPackageFile(file, packageDir string) bool {
	file = filepath.ToSlash(file)
	packageDir = filepath.ToSlash(packageDir)
	if filepath.Ext(file) != ".go" {
		return false
	}
	if packageDir == "." {
		return filepath.Dir(file) == "."
	}
	return filepath.Dir(file) == packageDir
}
