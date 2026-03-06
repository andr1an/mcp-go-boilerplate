package tools

import "testing"

func TestAllToolsHaveBasicSchemaShape(t *testing.T) {
	reg := NewRegistry()
	reg.Register(NewEchoTool())

	for _, info := range reg.List() {
		if info.Name == "" {
			t.Fatal("tool name must not be empty")
		}
		if info.Description == "" {
			t.Fatalf("tool %q has empty description", info.Name)
		}
		if info.InputSchema == nil {
			t.Fatalf("tool %q has nil input schema", info.Name)
		}

		typ, ok := info.InputSchema["type"].(string)
		if !ok {
			t.Fatalf("tool %q schema.type must be string", info.Name)
		}
		if typ != "object" {
			t.Fatalf("tool %q schema.type must be object, got %q", info.Name, typ)
		}

		props, ok := info.InputSchema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool %q schema.properties must be map[string]any", info.Name)
		}
		if props == nil {
			t.Fatalf("tool %q schema.properties must not be nil", info.Name)
		}

		if _, ok := info.InputSchema["additionalProperties"]; !ok {
			t.Fatalf("tool %q schema.additionalProperties must be set", info.Name)
		}
	}
}
