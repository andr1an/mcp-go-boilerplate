package tools

import (
	"regexp"
	"testing"
)

var toolNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func TestToolNamingConvention(t *testing.T) {
	reg := NewRegistry()
	reg.Register(NewEchoTool())

	for name := range reg.tools {
		if !toolNamePattern.MatchString(name) {
			t.Fatalf(
				"tool name %q violates naming convention (lowercase_with_underscores)",
				name,
			)
		}
	}
}

func TestToolNamesAreUniqueInList(t *testing.T) {
	reg := NewRegistry()
	reg.Register(NewEchoTool())

	seen := make(map[string]struct{})
	for _, info := range reg.List() {
		if _, ok := seen[info.Name]; ok {
			t.Fatalf("duplicate tool name in list: %q", info.Name)
		}
		seen[info.Name] = struct{}{}
	}
}
