package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type EchoTool struct{}

type EchoInput struct {
	Message string `json:"message"`
	Upper   bool   `json:"upper,omitempty"`
}

type EchoResult struct {
	Echo string `json:"echo"`
}

func init() {
	MustRegisterTool(NewEchoTool)
}

func NewEchoTool() Tool {
	return &EchoTool{}
}

func (t *EchoTool) Name() string {
	return "echo_message"
}

func (t *EchoTool) Description() string {
	return "Returns the provided message"
}

func (t *EchoTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "Message to echo back",
			},
			"upper": map[string]any{
				"type":        "boolean",
				"description": "Whether to uppercase the message",
			},
		},
		"required":             []string{"message"},
		"additionalProperties": false,
	}
}

func (t *EchoTool) Invoke(ctx context.Context, input []byte) (any, error) {
	_ = ctx

	var req EchoInput
	if err := json.Unmarshal(input, &req); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON body: %v", ErrInvalidArgument, err)
	}

	if req.Message == "" {
		return nil, fmt.Errorf("%w: message is required", ErrInvalidArgument)
	}

	out := req.Message
	if req.Upper {
		out = strings.ToUpper(out)
	}

	return EchoResult{
		Echo: out,
	}, nil
}
