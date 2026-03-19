package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

type ToolRegistrar func(s *mcp.Server)

var registrars []ToolRegistrar

func MustRegister(r ToolRegistrar) {
	if r == nil {
		panic("tools: nil registrar")
	}
	registrars = append(registrars, r)
}

func RegisterAll(s *mcp.Server) {
	for _, r := range registrars {
		r(s)
	}
}
