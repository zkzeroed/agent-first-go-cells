# 13. Concrete Examples

Use the runnable [`examples/reference-project`](../../examples/reference-project/) or the live scaffolders instead of copying examples from documentation:

```bash
task new-cell ID=user-authenticate
task scaffold-domain ID=users
task new-cell ID=users/user-invite
```

Each generated cell contains the current schema, local guide, manifest, and `api/api.go` contract boundary. Implement it inside-out, then run:

```bash
task index
task test-cell ID=<id>
task ready
```

The scaffold is deliberately the example: it stays executable and is covered by validation, unlike a static code listing.
