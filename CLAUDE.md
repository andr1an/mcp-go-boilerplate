# CLAUDE.md

## Repository Context
- This is a Go-based MCP tool server boilerplate.
- Uses the official MCP Go SDK: `github.com/modelcontextprotocol/go-sdk`
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

## Adding Tools
Tools are in `internal/tools/`. Each tool file should:
1. Define input/output structs with `json` and `jsonschema` tags
2. Register via `init()` using `MustRegister`
3. Implement handler with signature: `func(ctx, *mcp.CallToolRequest, Input) (*mcp.CallToolResult, Output, error)`

Example:
```go
type MyInput struct {
    Field string `json:"field" jsonschema:"Field description"`
}

type MyOutput struct {
    Result string `json:"result"`
}

func init() {
    MustRegister(func(s *mcp.Server) {
        mcp.AddTool(s, &mcp.Tool{
            Name:        "my_tool",
            Description: "What the tool does",
        }, MyHandler)
    })
}

func MyHandler(ctx context.Context, req *mcp.CallToolRequest, input MyInput) (*mcp.CallToolResult, MyOutput, error) {
    return nil, MyOutput{Result: input.Field}, nil
}
```

Key points:
- Fields without `omitempty` in json tag are required
- `jsonschema` tag provides field description
- Return `nil` for `*mcp.CallToolResult` unless customizing response
- Return error to signal tool failure (SDK wraps it automatically)

## Environment
- Example env file: `env.example`
- JWT auth is optional and configured via environment variables.

## Documentation
- If APIs, routes, tool schemas, or usage flow change, update `README.md`.

## Safety
- Do not remove or rewrite unrelated user changes.
- Avoid destructive commands unless explicitly requested.
