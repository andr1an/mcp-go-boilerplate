package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHandlerListTools(t *testing.T) {
	h := NewMCPHandler("test")
	s := httptest.NewServer(h)
	defer s.Close()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: s.URL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer session.Close()

	res, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}

	if len(res.Tools) < 1 {
		t.Fatalf("expected at least 1 tool, got %d", len(res.Tools))
	}

	found := false
	for _, tool := range res.Tools {
		if tool.Name == "echo_message" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find echo_message tool")
	}
}

func TestMCPHandlerCallTool(t *testing.T) {
	h := NewMCPHandler("test")
	s := httptest.NewServer(h)
	defer s.Close()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: s.URL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer session.Close()

	callRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "echo_message",
		Arguments: map[string]any{
			"message": "hello",
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
		t.Fatalf("expected structured content object, got %T", callRes.StructuredContent)
	}

	if out["echo"] != "hello" {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestMCPHandlerCallToolUpper(t *testing.T) {
	h := NewMCPHandler("test")
	s := httptest.NewServer(h)
	defer s.Close()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: s.URL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer session.Close()

	callRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "echo_message",
		Arguments: map[string]any{
			"message": "hello",
			"upper":   true,
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
		t.Fatalf("expected structured content object, got %T", callRes.StructuredContent)
	}

	if out["echo"] != "HELLO" {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestMCPHandlerMethodNotAllowed(t *testing.T) {
	h := NewMCPHandler("test")

	req := httptest.NewRequest(http.MethodPut, "/mcp", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	// SDK returns 400 Bad Request for unsupported methods
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
