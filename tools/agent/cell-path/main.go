// Command cell-path prints the configured directory of one manifest-backed cell.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
)

func main() {
	root := flag.String("root", ".", "project root")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: cell-path [--root PATH] ID")
		os.Exit(2)
	}
	manifests, err := manifest.FindAllAt(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, cell := range manifests {
		if cell.ID == flag.Arg(0) {
			fmt.Println(cell.Dir)
			return
		}
	}
	fmt.Fprintf(os.Stderr, "cell %q not found\n", flag.Arg(0))
	os.Exit(1)
}
