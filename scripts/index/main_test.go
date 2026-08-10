package main

import (
	"strings"
	"testing"

	"github.com/sploitzberg/go-agent-first-arch/scripts/manifest"
)

func TestContextPackIncludesAgentsSourceAndHashChanges(t *testing.T) {
	base := manifest.Manifest{ID: "orders-create", Purpose: "create", RawContent: "id: orders-create\n", AgentsContent: "# Guide\n\nUse the service.\n"}
	changed := base
	changed.AgentsContent = "# Guide\n\nUse the handler.\n"

	if computeHash([]manifest.Manifest{base}) == computeHash([]manifest.Manifest{changed}) {
		t.Fatal("computeHash() did not include AGENTS.md content")
	}
	pack := buildContextPack(base)
	if !strings.Contains(pack, "## Cell Guide") || !strings.Contains(pack, "Use the service.") {
		t.Fatalf("context pack did not include AGENTS.md content: %q", pack)
	}
}

func TestBoundedGuidePreservesSmallGuidesAndMarksLargeOnes(t *testing.T) {
	if got := boundedGuide("small"); got != "small" {
		t.Fatalf("boundedGuide(small) = %q", got)
	}
	large := strings.Repeat("x", maxContextGuideBytes+1)
	if got := boundedGuide(large); !strings.Contains(got, "Guide truncated") {
		t.Fatal("boundedGuide(large) did not include truncation marker")
	}
}
