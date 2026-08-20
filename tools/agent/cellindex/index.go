// Package cellindex defines the persisted cell-index contract shared by agent tools.
package cellindex

import (
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
)

// SchemaVersion identifies the current generated index format.
const SchemaVersion = "agent-first/v1"

// Index is the persisted representation of gen/cells.json.
type Index struct {
	SchemaVersion string `json:"schemaVersion"`
	Hash          string `json:"hash"`
	Cells         []Cell `json:"cells"`
}

// Cell describes one manifest-backed capability in the generated index.
type Cell struct {
	ID           string       `json:"id"`
	Path         string       `json:"path"`
	Package      string       `json:"package"`
	Kind         string       `json:"kind,omitempty"`
	Public       bool         `json:"public"`
	Purpose      string       `json:"purpose"`
	Entrypoints  []string     `json:"entrypoints"`
	Dependencies []string     `json:"dependencies"`
	Validation   []string     `json:"validation"`
	Conformance  *Conformance `json:"conformance,omitempty"`
}

// Conformance records the provenance and known limitations of a public package.
type Conformance struct {
	Basis     string     `json:"basis"`
	Status    string     `json:"status"`
	Evidence  string     `json:"evidence"`
	Citations []Citation `json:"citations,omitempty"`
	Rationale string     `json:"rationale,omitempty"`
	Gaps      []string   `json:"gaps,omitempty"`
}

// Citation identifies local evidence governing exported symbols.
type Citation struct {
	File    string   `json:"file"`
	Locator Locator  `json:"locator"`
	Symbols []string `json:"symbols"`
}

// Locator identifies a precise location within a citation.
type Locator struct {
	Type    string `json:"type"`
	Pages   []uint `json:"pages,omitempty"`
	Heading string `json:"heading,omitempty"`
}

// Decode strictly decodes and validates a persisted cell index.
func Decode(data []byte) (Index, error) {
	var index Index
	if err := jsonv2.Unmarshal(data, &index, jsonv2.RejectUnknownMembers(true)); err != nil {
		return Index{}, fmt.Errorf("decode cell index: %w", err)
	}
	if index.SchemaVersion != SchemaVersion {
		return Index{}, fmt.Errorf("schema version mismatch: existing=%s expected=%s", index.SchemaVersion, SchemaVersion)
	}
	if index.Hash == "" {
		return Index{}, errors.New("cell index hash is required")
	}
	if index.Cells == nil {
		return Index{}, errors.New("cell index cells are required")
	}
	return index, nil
}
