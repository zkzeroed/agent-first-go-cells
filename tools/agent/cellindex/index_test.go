package cellindex

import (
	"strings"
	"testing"
)

func TestDecodeAcceptsCurrentIndex(t *testing.T) {
	index, err := Decode([]byte(`{"schemaVersion":"agent-first/v1","hash":"sha256:test","cells":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if index.SchemaVersion != SchemaVersion || index.Hash != "sha256:test" || len(index.Cells) != 0 {
		t.Fatalf("Decode() = %#v", index)
	}
}

func TestDecodeRejectsAmbiguousOrUnknownJSON(t *testing.T) {
	invalidUTF8 := append([]byte(`{"schemaVersion":"agent-first/v1","hash":"`), 0xff)
	invalidUTF8 = append(invalidUTF8, []byte(`","cells":[]}`)...)

	tests := map[string][]byte{
		"duplicate member":  []byte(`{"schemaVersion":"agent-first/v1","hash":"first","hash":"second","cells":[]}`),
		"unknown member":    []byte(`{"schemaVersion":"agent-first/v1","hash":"sha256:test","cells":[],"extra":true}`),
		"unknown cell data": []byte(`{"schemaVersion":"agent-first/v1","hash":"sha256:test","cells":[{"extra":true}]}`),
		"invalid UTF-8":     invalidUTF8,
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := Decode(data); err == nil {
				t.Fatal("Decode() error = nil")
			}
		})
	}
}

func TestDecodeRejectsIncompleteOrIncompatibleIndex(t *testing.T) {
	tests := map[string]struct {
		data []byte
		want string
	}{
		"schema": {[]byte(`{"schemaVersion":"agent-first/v0","hash":"sha256:test","cells":[]}`), "schema version mismatch"},
		"hash":   {[]byte(`{"schemaVersion":"agent-first/v1","cells":[]}`), "hash is required"},
		"cells":  {[]byte(`{"schemaVersion":"agent-first/v1","hash":"sha256:test"}`), "cells are required"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Decode(test.data)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Decode() error = %v, want %q", err, test.want)
			}
		})
	}
}
