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

func TestToStructuredContentWithMap(t *testing.T) {
	input := map[string]any{"key": "value"}
	result := toStructuredContent(input)

	if result["key"] != "value" {
		t.Fatalf("expected key=value, got %#v", result)
	}
}

func TestToStructuredContentWithNonMap(t *testing.T) {
	result := toStructuredContent("plain string")

	if result["result"] != "plain string" {
		t.Fatalf("expected result='plain string', got %#v", result)
	}
}

func TestToStructuredContentWithSlice(t *testing.T) {
	input := []string{"a", "b", "c"}
	result := toStructuredContent(input)

	items, ok := result["result"].([]string)
	if !ok {
		t.Fatalf("expected result to be []string, got %T", result["result"])
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
}

func TestMCPHandlerCallToolWithError(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(&testTool{
		name:        "failing_tool",
		description: "A tool that fails",
		schema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		err: tools.ErrToolNotFound,
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
			Name: "failing_tool",
		},
	})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}

	if !callRes.IsError {
		t.Fatalf("expected error result, got success")
	}
}

func TestMCPHandlerMultipleTools(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(&testTool{
		name:        "tool_a",
		description: "First tool",
		schema:      map[string]any{"type": "object"},
		result:      map[string]any{"name": "a"},
	})
	reg.Register(&testTool{
		name:        "tool_b",
		description: "Second tool",
		schema:      map[string]any{"type": "object"},
		result:      map[string]any{"name": "b"},
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

	if len(res.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(res.Tools))
	}

	// Verify both tools can be called
	for _, toolName := range []string{"tool_a", "tool_b"} {
		callRes, err := c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: toolName,
			},
		})
		if err != nil {
			t.Fatalf("call tool %s failed: %v", toolName, err)
		}
		if callRes.IsError {
			t.Fatalf("expected success for %s", toolName)
		}
	}
}

func TestMCPHandlerCallToolWithNilArguments(t *testing.T) {
	reg := tools.NewRegistry()
	reg.Register(&testTool{
		name:        "no_args_tool",
		description: "Tool with no arguments",
		schema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		result: "success",
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

	// Call without passing Arguments at all
	callRes, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "no_args_tool",
		},
	})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}

	if callRes.IsError {
		t.Fatalf("expected successful result")
	}

	// Verify the non-map result is wrapped in "result" key
	out, ok := callRes.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("expected structured content object, got %#v", callRes.StructuredContent)
	}

	if out["result"] != "success" {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestToMCPToolWithValidSchema(t *testing.T) {
	info := tools.ToolInfo{
		Name:        "valid_tool",
		Description: "A valid tool",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
			"required": []string{"name"},
		},
	}

	tool := toMCPTool(info)

	if tool.Name != "valid_tool" {
		t.Fatalf("expected name 'valid_tool', got %q", tool.Name)
	}

	if tool.Description != "A valid tool" {
		t.Fatalf("expected description 'A valid tool', got %q", tool.Description)
	}

	var schema map[string]any
	if err := json.Unmarshal(tool.RawInputSchema, &schema); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("expected properties to be object, got %T", schema["properties"])
	}

	if _, exists := props["name"]; !exists {
		t.Fatalf("expected 'name' property in schema")
	}
}
