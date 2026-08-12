# Research Library Project

This runnable example uses an exportable `modulus` package whose manifest
records a research basis, implementation status, local source location, and
the exported symbol governed by that source. The package composes the private
`modulus-reduce` cell.

```bash
task context ID=modulus ROOT=examples/research-library-project
task deps ID=modulus ROOT=examples/research-library-project
task scope ID=modulus ROOT=examples/research-library-project
task validate-cell ID=modulus ROOT=examples/research-library-project
task ready ROOT=examples/research-library-project
```

The example deliberately keeps the mathematics small: `modulus.Reduce` returns
the canonical Euclidean remainder modulo 97. Its provenance is in
`modulus/cell.yaml` and `docs/research/modular-reduction.md`.
