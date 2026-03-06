package tools

import (
	"context"
	"errors"
	"testing"
)

func TestEchoToolMetadata(t *testing.T) {
	tool := NewEchoTool()

	if tool.Name() != "echo_message" {
		t.Fatalf("unexpected name: %q", tool.Name())
	}

	if tool.Description() == "" {
		t.Fatal("expected non-empty description")
	}

	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Fatalf("unexpected schema type: %#v", schema["type"])
	}
}

func TestEchoToolInvoke(t *testing.T) {
	tool := NewEchoTool()

	got, err := tool.Invoke(context.Background(), []byte(`{"message":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, ok := got.(EchoResult)
	if !ok {
		t.Fatalf("unexpected result type: %T", got)
	}

	if res.Echo != "hello" {
		t.Fatalf("unexpected echo: %q", res.Echo)
	}
}

func TestEchoToolInvokeUpper(t *testing.T) {
	tool := NewEchoTool()

	got, err := tool.Invoke(context.Background(), []byte(`{"message":"hello","upper":true}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, ok := got.(EchoResult)
	if !ok {
		t.Fatalf("unexpected result type: %T", got)
	}

	if res.Echo != "HELLO" {
		t.Fatalf("unexpected echo: %q", res.Echo)
	}
}

func TestEchoToolInvokeInvalidJSON(t *testing.T) {
	tool := NewEchoTool()

	_, err := tool.Invoke(context.Background(), []byte(`{"message":`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestEchoToolInvokeMissingMessage(t *testing.T) {
	tool := NewEchoTool()

	_, err := tool.Invoke(context.Background(), []byte(`{}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}
