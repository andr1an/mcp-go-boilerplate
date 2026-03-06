package tools

import "fmt"

type ToolFactory func() Tool

var builtinFactories []ToolFactory

func MustRegisterTool(factory ToolFactory) {
	if factory == nil {
		panic("tools: nil factory")
	}

	tool := factory()
	if tool == nil {
		panic("tools: factory returned nil tool")
	}
	if tool.Name() == "" {
		panic("tools: tool name must not be empty")
	}

	for _, existing := range builtinFactories {
		if existing().Name() == tool.Name() {
			panic(fmt.Sprintf("tools: duplicate tool registration: %s", tool.Name()))
		}
	}

	builtinFactories = append(builtinFactories, factory)
}

func RegisterBuiltins(registry *Registry) {
	for _, factory := range builtinFactories {
		registry.Register(factory())
	}
}
