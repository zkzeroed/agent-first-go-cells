// Command conformance validates and prints a public package's research record.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	if flag.NArg() != 1 {
		fail("usage: conformance [--root PATH] ID")
	}
	manifests, err := manifest.FindAllAt(*root)
	if err != nil {
		fail(err.Error())
	}
	if err := manifest.ValidateSourceAt(*root, manifests); err != nil {
		fail(err.Error())
	}
	for _, cell := range manifests {
		if cell.ID != flag.Arg(0) {
			continue
		}
		if cell.Kind != "library-package" {
			fail(fmt.Sprintf("cell %q is not a public library package", cell.ID))
		}
		printConformance(cell)
		return
	}
	fail(fmt.Sprintf("cell %q not found", flag.Arg(0)))
}

func printConformance(cell manifest.Manifest) {
	value := cell.Conformance
	fmt.Printf("=== Conformance: %s ===\n", cell.ID)
	fmt.Printf("Basis: %s\nStatus: %s\nEvidence: %s\n", value.Basis, value.Status, value.Evidence)
	if value.Rationale != "" {
		fmt.Printf("Rationale: %s\n", value.Rationale)
	}
	for _, gap := range value.Gaps {
		fmt.Printf("Gap: %s\n", gap)
	}
	for _, citation := range value.Citations {
		fmt.Printf("Citation: %s — %s; symbols: %s\n", citation.File, locator(citation.Locator), strings.Join(citation.Symbols, ", "))
	}
}

func locator(value manifest.CitationLocator) string {
	if value.Type == "markdown-heading" {
		return "heading " + fmt.Sprintf("%q", value.Heading)
	}
	return fmt.Sprintf("PDF pages %v", value.Pages)
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, "Error:", message)
	os.Exit(1)
}
