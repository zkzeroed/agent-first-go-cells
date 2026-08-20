package main

import (
	"encoding/json"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/cellindex"
)

func TestWithStatusPreservesCellsJSONShape(t *testing.T) {
	index := cellindex.Index{
		SchemaVersion: cellindex.SchemaVersion,
		Hash:          "sha256:test",
		Cells:         []cellindex.Cell{{ID: "orders", Dependencies: []string{}}},
	}
	data, err := json.Marshal(withStatus(index))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		SchemaVersion string `json:"schemaVersion"`
		Hash          string `json:"hash"`
		Cells         []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"cells"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.SchemaVersion != cellindex.SchemaVersion || decoded.Hash != index.Hash {
		t.Fatalf("index metadata = %#v", decoded)
	}
	if len(decoded.Cells) != 1 || decoded.Cells[0].ID != "orders" || decoded.Cells[0].Status != "ok" {
		t.Fatalf("cells = %#v", decoded.Cells)
	}
}
