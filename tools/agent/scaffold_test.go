package agent_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func TestScaffoldersRejectExistingCellPaths(t *testing.T) {
	for _, test := range []struct {
		name   string
		script string
		id     string
	}{
		{name: "flat cell", script: "new-cell.sh", id: "sample-create"},
		{name: "domain cell", script: "new-domain.sh", id: "users"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			runScaffolder(t, root, test.script, test.id)
			assertScaffoldManifestValid(t, root)
			path := filepath.Join(root, "internal", "cells", test.id, "cell.yaml")
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := scaffolder(t, root, test.script, test.id).Run(); err == nil {
				t.Fatal("second scaffold succeeded, want an existing-path error")
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatal("second scaffold changed existing cell content")
			}
		})
	}
}

func assertScaffoldManifestValid(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/scaffold-test\n\ngo 1.27.0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifests, err := manifest.FindAllAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := manifest.ValidateSourceAt(root, manifests); err != nil {
		t.Fatal(err)
	}
}

func runScaffolder(t *testing.T, root, script, id string) {
	t.Helper()
	if output, err := scaffolder(t, root, script, id).CombinedOutput(); err != nil {
		t.Fatalf("scaffold %s: %v\n%s", id, err, output)
	}
}

func scaffolder(t *testing.T, root, script, id string) *exec.Cmd {
	t.Helper()
	command := exec.Command("bash", filepath.Join(projectRoot(t), "tools", "agent", script), id)
	command.Dir = root
	return command
}
