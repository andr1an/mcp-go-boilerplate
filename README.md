# 🚀 MCP Go Boilerplate

A minimal, production-ready **Model Context Protocol (MCP) tool server** written in Go.

This project provides a clean foundation for building MCP-compatible tool servers that can be used by AI agents (Claude, local agents, etc.) or any client capable of making HTTP requests.

It includes:

- structured tool system
- automatic tool discovery
- JSON schema for tool inputs
- HTTP transport
- optional JWT authentication
- structured logging
- test suite
- production-safe server configuration
- reproducible cross-platform builds via `Makefile`

The goal is to provide a **simple and extensible starting point** for MCP servers.

## Project Structure

```
main.go
Makefile
internal/
  auth/
    jwt.go
  config/
    config.go
  httpserver/
    server.go
  middleware/
    auth.go
    logging.go
    requestid.go
  tools/
    tool.go
    registry.go
    builtins.go
    errors.go
    echo.go
  transport/
    mcp.go
```

## Quick Start

### 1. Clone project

```bash
git clone https://github.com/andr1an/mcp-go-boilerplate
cd mcp-go-boilerplate
```

### 2. Run server

```bash
go run .
```

Server will start on:

```
http://localhost:8080
```

### 3. Print build metadata

```bash
go run . version
```

Example output:

```text
version=dev commit=none date=unknown
```

For release binaries, metadata is injected at build time via `-ldflags` in the `Makefile`.

The `version` value is also reported by the MCP server in protocol responses (e.g., during `initialize`).

## Endpoints

### Health Check

```
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### MCP Transport (`/mcp`)

`/mcp` uses MCP Streamable HTTP transport (`jsonrpc` requests via `POST`).

Example initialize request:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "latest supported MCP version",
    "clientInfo": {
      "name": "example-client",
      "version": "1.0.0"
    },
    "capabilities": {}
  }
}
```

List tools:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}
```

Call tool:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "echo_message",
    "arguments": {
      "message": "hello",
      "upper": true
    }
  }
}
```

## Adding a Tool

Tools are implemented inside:

```
internal/tools/
```

Each tool must implement the Tool interface.

```go
type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]any
	Invoke(ctx context.Context, input []byte) (any, error)
}
```

## Example Tool

Example implementation:

```go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

type ExampleTool struct{}

type ExampleInput struct {
	Message string `json:"message"`
}

type ExampleResult struct {
	Response string `json:"response"`
}

func init() {
	MustRegisterTool(NewExampleTool)
}

func NewExampleTool() Tool {
	return &ExampleTool{}
}

func (t *ExampleTool) Name() string {
	return "example_tool"
}

func (t *ExampleTool) Description() string {
	return "Example tool"
}

func (t *ExampleTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"message"},
		"additionalProperties": false,
	}
}

func (t *ExampleTool) Invoke(ctx context.Context, input []byte) (any, error) {
	var req ExampleInput

	if err := json.Unmarshal(input, &req); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON", ErrInvalidArgument)
	}

	return ExampleResult{
		Response: req.Message,
	}, nil
}
```

## Automatic Tool Registration

Tools register themselves automatically via:

```go
func init() {
	MustRegisterTool(NewExampleTool)
}
```

This means:

- no changes to server.go
- dropping a file into `internal/tools/` automatically exposes the tool


## Tool Naming Convention

Tool names must follow:

```
lowercase_with_underscores
```

Regex enforced by tests:

```
^[a-z][a-z0-9_]*$
```

## Authentication (Optional)

Server supports **JWT authentication**.

Enable with environment variables:

```bash
AUTH_MODE=jwt
JWT_PUBLIC_KEY=public.pem
```

Start server:

```bash
AUTH_MODE=jwt JWT_PUBLIC_KEY=public.pem go run .
```

Clients must include:

```
Authorization: Bearer <token>
```

## Configuration

Server configuration is controlled via environment variables.

|**Variable**|**Default**|**Description**|
|---|---|---|
|`LISTEN_ADDR`|127.0.0.1:8080|Listen address (e.g., `127.0.0.1:8080`, `0.0.0.0:8080`, `[::]:8080`)|
|`AUTH_MODE`|disabled|disabled or jwt|
|`JWT_PUBLIC_KEY`|empty|RSA public key for JWT|
|`LOG_LEVEL`|info|log level|
|`READ_TIMEOUT`|15s|request timeout|
|`WRITE_TIMEOUT`|60s|response timeout|
|`IDLE_TIMEOUT`|60s|keepalive timeout|

## Build

Build local binary:

```bash
make build
```

Cross-compile for common OS/ARCH targets:

```bash
make build-all
```

## Testing

Run the full test suite:

```bash
make test
```

Tests cover:

- tool registry
- tool naming conventions
- JSON schema validation
- transport handler
- HTTP server

## License

MIT License
