package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestEcho(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	res, out, err := Echo(ctx, req, EchoInput{Message: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
	if out.Echo != "hello" {
		t.Fatalf("unexpected echo: %q", out.Echo)
	}
}

func TestEchoUpper(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	_, out, err := Echo(ctx, req, EchoInput{Message: "hello", Upper: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Echo != "HELLO" {
		t.Fatalf("unexpected echo: %q", out.Echo)
	}
}

func TestEchoEmptyMessage(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	_, out, err := Echo(ctx, req, EchoInput{Message: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Echo != "" {
		t.Fatalf("unexpected echo: %q", out.Echo)
	}
}
