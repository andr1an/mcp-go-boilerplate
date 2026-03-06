package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andr1an/mcp-go-boilerplate/internal/tools"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type testTool struct {
	name        string
	description string
	schema      map[string]any
	result      any
	err         error
}

func (t *testTool) Name() string {
	return t.name
}

func (t *testTool) Description() string {
	return t.description
}

func (t *testTool) InputSchema() map[string]any {
	return t.schema
}

func (t *testTool) Invoke(ctx context.Context, input []byte) (any, error) {
	_ = ctx
	_ = input
	return t.result, t.err
}

func TestMCPHandlerListTools(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(&testTool{
		name:        "echo_message",
		description: "Returns a message",
		schema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	})

	h := NewMCPHandler(reg)
	s := httptest.NewServer(h)
	defer s.Close()

	c, err := client.NewStreamableHttpClient(s.URL)
	if err != nil {
		t.Fatalf("create client failed: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Start(ctx); err != nil {
		t.Fatalf("start client failed: %v", err)
	}

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "transport-test-client",
				Version: "1.0.0",
			},
		},
	})
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	res, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}

	if len(res.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(res.Tools))
	}

	if res.Tools[0].Name != "echo_message" {
		t.Fatalf("unexpected tool name: %q", res.Tools[0].Name)
	}
}

func TestMCPHandlerCallTool(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(&testTool{
		name:        "echo_message",
		description: "Returns a message",
		schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{"type": "string"},
			},
		},
		result: map[string]any{
			"echo": "hello",
		},
	})

	h := NewMCPHandler(reg)
	s := httptest.NewServer(h)
	defer s.Close()

	c, err := client.NewStreamableHttpClient(s.URL)
	if err != nil {
		t.Fatalf("create client failed: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Start(ctx); err != nil {
		t.Fatalf("start client failed: %v", err)
	}

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "transport-test-client",
				Version: "1.0.0",
			},
		},
	})
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	callRes, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "echo_message",
			Arguments: map[string]any{
				"message": "hello",
			},
		},
	})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}

	if callRes.IsError {
		t.Fatalf("expected successful tool result")
	}

	out, ok := callRes.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("expected structured content object, got %#v", callRes.StructuredContent)
	}

	if out["echo"] != "hello" {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestMCPHandlerMethodNotAllowed(t *testing.T) {
	reg := tools.NewRegistry()
	h := NewMCPHandler(reg)

	req := httptest.NewRequest(http.MethodPut, "/mcp", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestToMCPToolUsesFallbackSchemaOnMarshalError(t *testing.T) {
	tool := toMCPTool(tools.ToolInfo{
		Name:        "bad_schema_tool",
		Description: "bad schema",
		InputSchema: map[string]any{
			"type": "object",
			"bad":  make(chan int),
		},
	})

	var schema map[string]any
	if err := json.Unmarshal(tool.RawInputSchema, &schema); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	if schema["type"] != "object" {
		t.Fatalf("expected fallback schema type object, got %#v", schema["type"])
	}
}

func TestToJSONStringFallbackOnMarshalError(t *testing.T) {
	if got := toJSONString(make(chan int)); got != "{}" {
		t.Fatalf("expected fallback {}, got %q", got)
	}
}
