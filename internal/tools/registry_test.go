package tools

import (
	"context"
	"errors"
	"testing"
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

func TestRegistryRegisterAndGet(t *testing.T) {
	r := NewRegistry()
	tool := &testTool{
		name:        "test_echo",
		description: "test tool",
		schema:      map[string]any{"type": "object"},
	}

	r.Register(tool)

	got, ok := r.Get("test_echo")
	if !ok {
		t.Fatal("expected tool to be registered")
	}
	if got != tool {
		t.Fatal("expected same tool instance")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()
	r.Register(&testTool{
		name:        "test_echo",
		description: "test tool",
		schema:      map[string]any{"type": "object"},
	})

	list := r.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(list))
	}

	if list[0].Name != "test_echo" {
		t.Fatalf("unexpected tool name: %q", list[0].Name)
	}
	if list[0].Description != "test tool" {
		t.Fatalf("unexpected description: %q", list[0].Description)
	}
}

func TestRegistryInvoke(t *testing.T) {
	r := NewRegistry()
	r.Register(&testTool{
		name:   "test_echo",
		result: map[string]any{"ok": true},
	})

	got, err := r.Invoke(context.Background(), "test_echo", []byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("unexpected result type: %T", got)
	}
	if m["ok"] != true {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestRegistryInvokeToolNotFound(t *testing.T) {
	r := NewRegistry()

	_, err := r.Invoke(context.Background(), "missing_tool", []byte(`{}`))
	if !errors.Is(err, ErrToolNotFound) {
		t.Fatalf("expected ErrToolNotFound, got %v", err)
	}
}

func TestRegistryRegisterOverridesSameName(t *testing.T) {
	r := NewRegistry()

	first := &testTool{name: "test_echo", description: "first"}
	second := &testTool{name: "test_echo", description: "second"}

	r.Register(first)
	r.Register(second)

	got, ok := r.Get("test_echo")
	if !ok {
		t.Fatal("expected tool to exist")
	}
	if got != second {
		t.Fatal("expected later registration to override earlier one")
	}
}
