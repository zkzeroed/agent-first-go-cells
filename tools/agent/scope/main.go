// Command scope prints and verifies the explicit change boundary for one cell.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/zkzeroed/agent-first-go-cells/tools/agent/manifest"
	"github.com/zkzeroed/agent-first-go-cells/tools/agent/projectconfig"
)

type arguments struct {
	cellID string
	root   string
	with   string
	verify bool
}

type report struct {
	CellID       string
	Kind         string
	Public       bool
	Conformance  manifest.Conformance
	Purpose      string
	Owned        []string
	Entrypoints  []string
	Dependencies []string
	Dependents   []string
	Invariants   []string
	Validation   []string
	Allowed      []string
	Changed      []string
	OutOfScope   []string
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		fail(err)
	}
	report, err := buildReport(args)
	if err != nil {
		fail(err)
	}
	printReport(report, args.verify)
	if args.verify && len(report.OutOfScope) > 0 {
		os.Exit(1)
	}
}

func parseArgs(values []string) (arguments, error) {
	result := arguments{root: "."}
	for len(values) > 0 {
		value := values[0]
		values = values[1:]
		switch value {
		case "--verify":
			result.verify = true
		case "--root", "--with":
			if len(values) == 0 {
				return arguments{}, fmt.Errorf("%s requires a value", value)
			}
			if value == "--root" {
				result.root = values[0]
			} else {
				result.with = values[0]
			}
			values = values[1:]
		default:
			if strings.HasPrefix(value, "-") || result.cellID != "" {
				return arguments{}, fmt.Errorf("usage: scope [--verify] [--root PATH] [--with LIST] ID")
			}
			result.cellID = value
		}
	}
	if result.cellID == "" {
		return arguments{}, errors.New("cell ID is required")
	}
	return result, nil
}

func buildReport(args arguments) (report, error) {
	head, headErr := manifestAtHEAD(args.root, args.cellID)
	var extras []manifest.Manifest
	if headErr == nil {
		var err error
		extras, err = missingHeadManifests(args.root, *head)
		if err != nil {
			return report{}, err
		}
	}
	manifests, err := manifest.FindAllAtWith(args.root, extras)
	if err != nil {
		return report{}, err
	}
	target := findManifest(args.cellID, manifests)
	if target == nil {
		if headErr != nil {
			return report{}, fmt.Errorf("cell %q not found", args.cellID)
		}
		return report{}, fmt.Errorf("cell %q not found", args.cellID)
	}
	allowed, err := allowedPaths(*target, manifests, args.with)
	if err != nil {
		return report{}, err
	}
	result := report{
		CellID:       target.ID,
		Kind:         target.Kind,
		Public:       target.Public,
		Conformance:  target.Conformance,
		Purpose:      target.Purpose,
		Owned:        []string{target.Dir, filepath.ToSlash(filepath.Join("gen", "context", target.ID+".context.md")), "gen/cells.json"},
		Entrypoints:  entrypoints(*target),
		Dependencies: target.Dependencies,
		Dependents:   dependentIDs(target.ID, manifests),
		Invariants:   target.Invariants,
		Validation:   target.Validation,
		Allowed:      scopeLabels(allowed),
	}
	if !args.verify {
		return result, nil
	}
	changed, err := changedFiles(args.root)
	if err != nil {
		return report{}, err
	}
	result.Changed = changed
	result.OutOfScope = outOfScope(changed, allowed, *target, manifests)
	return result, nil
}

func missingHeadManifests(root string, target manifest.Manifest) ([]manifest.Manifest, error) {
	repository, err := gitRoot(root)
	if err != nil {
		return nil, err
	}
	project, err := projectPath(repository, root)
	if err != nil {
		return nil, err
	}
	path := filepath.ToSlash(filepath.Join(project, target.Dir))
	output, err := gitOutput(repository, "ls-tree", "-r", "--name-only", "HEAD", "--", path)
	if err != nil {
		return nil, err
	}
	files, err := filesUnderRoot(splitFiles(output), repository, root)
	if err != nil {
		return nil, err
	}
	var result []manifest.Manifest
	for _, file := range files {
		if file == filepath.ToSlash(filepath.Join(target.Dir, "cell.yaml")) {
			missing, err := missingFromWorktree(root, file)
			if err != nil {
				return nil, err
			}
			if missing {
				result = append(result, target)
			}
			continue
		}
		id, found := removedCellID(file)
		if !found || !isWithin(file, target.Dir) {
			continue
		}
		missing, err := missingFromWorktree(root, file)
		if err != nil {
			return nil, err
		}
		if !missing {
			continue
		}
		loaded, err := manifestAtHEAD(root, id)
		if err != nil {
			return nil, err
		}
		result = append(result, *loaded)
	}
	return result, nil
}

func scopeLabels(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = strings.TrimPrefix(value, "@cell:")
	}
	return result
}

func manifestAtHEAD(root, id string) (*manifest.Manifest, error) {
	config, err := projectconfig.Load(root)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(projectconfig.CellsRoot, filepath.FromSlash(id))
	if libraryDir, found := config.LibraryPackages[id]; found {
		dir = libraryDir
	}
	repository, err := gitRoot(root)
	if err != nil {
		return nil, err
	}
	project, err := projectPath(repository, root)
	if err != nil {
		return nil, err
	}
	path := filepath.ToSlash(filepath.Join(project, dir, "cell.yaml"))
	content, err := gitOutput(repository, "show", "HEAD:"+path)
	if err != nil {
		return nil, err
	}
	result, err := manifest.Parse(content)
	if err != nil {
		return nil, err
	}
	result.Dir = filepath.ToSlash(dir)
	return &result, nil
}

func entrypoints(cell manifest.Manifest) []string {
	result := make([]string, len(cell.Entrypoints))
	for i, file := range cell.Entrypoints {
		result[i] = file + ": " + cell.Symbols[i]
	}
	return result
}

func findManifest(id string, manifests []manifest.Manifest) *manifest.Manifest {
	for i := range manifests {
		if manifests[i].ID == id {
			return &manifests[i]
		}
	}
	return nil
}

func allowedPaths(target manifest.Manifest, manifests []manifest.Manifest, value string) ([]string, error) {
	allowed := []string{"@cell:" + target.ID, "gen/cells.json", filepath.ToSlash(filepath.Join("gen", "context", target.ID+".context.md"))}
	if value == "" {
		return allowed, nil
	}
	seen := map[string]struct{}{target.ID: {}}
	for item := range strings.SplitSeq(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, errors.New("WITH cannot contain an empty value")
		}
		if _, exists := seen[item]; exists {
			return nil, fmt.Errorf("WITH contains duplicate %q", item)
		}
		seen[item] = struct{}{}
		switch item {
		case "@contracts", "@platform", "@wiring":
			allowed = append(allowed, item)
		default:
			if findManifest(item, manifests) == nil {
				return nil, fmt.Errorf("WITH contains unknown cell or scope %q", item)
			}
			allowed = append(allowed, "@cell:"+item)
		}
	}
	return allowed, nil
}

func dependentIDs(id string, manifests []manifest.Manifest) []string {
	var result []string
	for _, m := range manifests {
		if slices.Contains(m.Dependencies, id) {
			result = append(result, m.ID)
		}
	}
	return result
}

func changedFiles(root string) ([]string, error) {
	repository, err := gitRoot(root)
	if err != nil {
		return nil, err
	}
	tracked, err := gitOutput(repository, "diff", "--name-only", "HEAD")
	if err != nil {
		return nil, err
	}
	untracked, err := gitOutput(repository, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	files, err := filesUnderRoot(append(splitFiles(tracked), splitFiles(untracked)...), repository, root)
	if err != nil {
		return nil, err
	}
	return unique(files), nil
}

func gitRoot(root string) (string, error) {
	return gitOutput(root, "rev-parse", "--show-toplevel")
}

func gitOutput(root string, args ...string) (string, error) {
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func splitFiles(output string) []string {
	if output == "" {
		return nil
	}
	return slices.Collect(strings.SplitSeq(output, "\n"))
}

func removedCellID(file string) (string, bool) {
	file = filepath.ToSlash(file)
	prefix := "internal/cells/"
	if !strings.HasPrefix(file, prefix) || !strings.HasSuffix(file, "/cell.yaml") {
		return "", false
	}
	return strings.TrimSuffix(strings.TrimPrefix(file, prefix), "/cell.yaml"), true
}

func missingFromWorktree(root, file string) (bool, error) {
	_, err := os.Stat(filepath.Join(root, file))
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	return false, err
}

func filesUnderRoot(files []string, repository, root string) ([]string, error) {
	relative, err := projectPath(repository, root)
	if err != nil {
		return nil, err
	}
	prefix := filepath.ToSlash(filepath.Clean(relative))
	if prefix == "." {
		return files, nil
	}
	prefix += "/"
	var result []string
	for _, file := range files {
		if path, found := strings.CutPrefix(filepath.ToSlash(file), prefix); found {
			result = append(result, path)
		}
	}
	return result, nil
}

func projectPath(repository, root string) (string, error) {
	project, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(repository, project)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("project root must be inside the Git repository")
	}
	return relative, nil
}

func unique(values []string) []string {
	var result []string
	for _, value := range values {
		if !slices.Contains(result, value) {
			result = append(result, value)
		}
	}
	return result
}

func outOfScope(files, allowed []string, target manifest.Manifest, manifests []manifest.Manifest) []string {
	var result []string
	for _, file := range files {
		if !isAllowed(file, allowed, target, manifests) {
			result = append(result, file)
		}
	}
	return result
}

func isAllowed(file string, allowed []string, target manifest.Manifest, manifests []manifest.Manifest) bool {
	for _, item := range allowed {
		if item == "@contracts" && strings.HasPrefix(file, "internal/contracts/") {
			return true
		}
		if item == "@platform" && strings.HasPrefix(file, "internal/platform/") {
			return true
		}
		if item == "@platform" && strings.HasPrefix(file, "internal/") && !strings.HasPrefix(file, "internal/cells/") && !strings.HasPrefix(file, "internal/contracts/") {
			return true
		}
		if item == "@wiring" && isWiringFile(file) {
			return true
		}
		if id, found := strings.CutPrefix(item, "@cell:"); found && belongsToCell(file, id, manifests) {
			return true
		}
		if file == item {
			return true
		}
	}
	return false
}

func belongsToCell(file, id string, manifests []manifest.Manifest) bool {
	cell := findManifest(id, manifests)
	if cell == nil || !isWithin(file, cell.Dir) {
		return false
	}
	for _, other := range manifests {
		if other.ID != cell.ID && isWithin(other.Dir, cell.Dir) && isWithin(file, other.Dir) {
			return false
		}
	}
	return true
}

func isWithin(path, directory string) bool {
	path = filepath.ToSlash(path)
	directory = strings.TrimSuffix(filepath.ToSlash(directory), "/")
	if directory == "." {
		return filepath.Dir(path) == "."
	}
	return path == directory || strings.HasPrefix(path, directory+"/")
}

func isWiringFile(path string) bool {
	return filepath.Dir(path) == "internal/app" && strings.HasPrefix(filepath.Base(path), "wiring") && strings.HasSuffix(path, ".go")
}

func printReport(report report, verify bool) {
	title := "Scope"
	if verify {
		title = "Scope Verification"
	}
	fmt.Printf("=== %s: %s ===\n", title, report.CellID)
	printValues("Kind", []string{displayKind(report.Kind)})
	printValues("Public", []string{fmt.Sprintf("%t", report.Public)})
	printConformance(report.Conformance)
	printValues("Purpose", []string{report.Purpose})
	printValues("Owned paths", report.Owned)
	printValues("Declared scope", report.Allowed)
	printValues("Public entrypoints", report.Entrypoints)
	printValues("Dependencies", report.Dependencies)
	printValues("Direct dependents", report.Dependents)
	printValues("Invariants", report.Invariants)
	printValues("Validation", report.Validation)
	if verify {
		printValues("Changed files", report.Changed)
		printValues("Out-of-scope files", report.OutOfScope)
	}
}

func printConformance(value manifest.Conformance) {
	if value.Basis == "" {
		return
	}
	printValues("Conformance basis", []string{value.Basis})
	printValues("Conformance status", []string{value.Status})
	printValues("Conformance evidence", []string{value.Evidence})
	if value.Rationale != "" {
		printValues("Conformance rationale", []string{value.Rationale})
	}
	if len(value.Gaps) > 0 {
		printValues("Conformance gaps", value.Gaps)
	}
	for _, citation := range value.Citations {
		printValues("Conformance citation", []string{fmt.Sprintf("%s, %s; symbols: %s", citation.File, displayLocator(citation.Locator), strings.Join(citation.Symbols, ", "))})
	}
}

func displayLocator(locator manifest.CitationLocator) string {
	if locator.Type == "markdown-heading" {
		return "heading " + fmt.Sprintf("%q", locator.Heading)
	}
	return fmt.Sprintf("PDF pages %v", locator.Pages)
}

func displayKind(kind string) string {
	if kind == "" {
		return "private cell"
	}
	return kind
}

func printValues(title string, values []string) {
	fmt.Printf("\n%s:\n", title)
	if len(values) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, value := range values {
		fmt.Printf("  %s\n", value)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
