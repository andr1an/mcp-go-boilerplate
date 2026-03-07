# AGENTS.md

## Project
- Go MCP server boilerplate.
- Entry point: `main.go`.
- Core packages are under `internal/`.

## Setup
- Use Go `1.25+` (per `go.mod`).
- Copy env file if needed: `cp env.example env`.

## Development Commands
- Run server: `go run .`
- Run tests: `go test ./...`
- Build binary: `make build`
- Build all targets: `make build-all`
- Print build metadata: `make print-version`

## Coding Rules
- Keep changes minimal and scoped to the task.
- Follow existing package/layout patterns in `internal/*`.
- Prefer table-driven tests for new behavior.
- Do not add dependencies unless necessary.
- Preserve backward-compatible tool names and MCP method handling.

## Tooling/Validation
- Run `go test ./...` after code changes.
- If behavior changes, update `README.md` examples/docs.

## Git Hygiene
- Do not revert unrelated local changes.
- Avoid destructive git commands unless explicitly requested.
- Keep commits focused and small.
