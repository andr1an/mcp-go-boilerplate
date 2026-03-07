# CLAUDE.md

## Repository Context
- This is a Go-based MCP tool server boilerplate.
- Main executable: `main.go`.
- Internal modules live under `internal/` (`tools`, `transport`, `httpserver`, `middleware`, etc.).

## Working Expectations
- Make focused, minimal edits.
- Match existing code style and naming conventions.
- Keep MCP protocol behavior stable unless explicitly asked to change it.
- Add or update tests when modifying behavior.

## Common Commands
- Run app: `go run .`
- Run full test suite: `go test ./...`
- Build: `make build`
- Cross-build: `make build-all`

## Environment
- Example env file: `env.example`
- JWT auth is optional and configured via environment variables.

## Documentation
- If APIs, routes, tool schemas, or usage flow change, update `README.md`.

## Safety
- Do not remove or rewrite unrelated user changes.
- Avoid destructive commands unless explicitly requested.
