# Contributing

## Development setup

- Go 1.22+
- Clone repository and run commands from repo root

## Local checks

```bash
go test ./...
go vet ./...
go build ./cmd/nc-guard
```

## Coding conventions

- Keep functions small and testable
- Prefer deterministic output (stable ordering)
- Return robust errors instead of panics for expected failures
- Avoid introducing network dependencies in v0.1
