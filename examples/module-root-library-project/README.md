# Module-Root Library Project

This runnable example exposes its module root, `example.com/canonical`, as the
public `canonical` package. It was scaffolded with:

```bash
task new-package ID=canonical PATH=. PACKAGE=canonical ROOT=examples/module-root-library-project
```

The package composes the lean private `text-fold` cell under `internal/cells/`.

```bash
task cells ROOT=examples/module-root-library-project
task context ID=canonical ROOT=examples/module-root-library-project
task deps ID=canonical ROOT=examples/module-root-library-project
task conformance ID=canonical ROOT=examples/module-root-library-project
task ready ROOT=examples/module-root-library-project
```
