package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EchoInput struct {
	Message string `json:"message" jsonschema:"Message to echo back"`
	Upper   bool   `json:"upper,omitempty" jsonschema:"Whether to uppercase the message"`
}

type EchoOutput struct {
	Echo string `json:"echo"`
}

func init() {
	MustRegister(func(s *mcp.Server) {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "echo_message",
			Description: "Returns the provided message",
		}, Echo)
	})
}

func Echo(ctx context.Context, req *mcp.CallToolRequest, input EchoInput) (*mcp.CallToolResult, EchoOutput, error) {
	out := input.Message
	if input.Upper {
		out = strings.ToUpper(out)
	}
	return nil, EchoOutput{Echo: out}, nil
}
