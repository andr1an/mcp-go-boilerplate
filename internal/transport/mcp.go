package transport

import (
	"net/http"

	"github.com/andr1an/mcp-go-boilerplate/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPHandler struct {
	streamable *mcp.StreamableHTTPHandler
}

func NewMCPHandler(version string) *MCPHandler {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-go-boilerplate",
			Version: version,
		},
		nil,
	)

	tools.RegisterAll(server)

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	return &MCPHandler{
		streamable: handler,
	}
}

func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.streamable.ServeHTTP(w, r)
}
