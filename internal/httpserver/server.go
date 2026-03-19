package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/andr1an/mcp-go-boilerplate/internal/auth"
	"github.com/andr1an/mcp-go-boilerplate/internal/config"
	"github.com/andr1an/mcp-go-boilerplate/internal/middleware"
	"github.com/andr1an/mcp-go-boilerplate/internal/transport"
)

type Server struct {
	*http.Server
}

func New(cfg config.Config, logger *slog.Logger, version string) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/mcp", chain(
		transport.NewMCPHandler(version),
		middleware.RequestID,
		middleware.Logging(logger),
		buildAuthMiddleware(cfg),
	))

	srv := &http.Server{
		Addr:           cfg.Address(),
		Handler:        mux,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	return &Server{Server: srv}, nil
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func buildAuthMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	switch cfg.AuthMode {
	case config.AuthDisabled:
		return middleware.AuthDisabled
	case config.AuthJWT:
		validator, err := auth.NewJWTValidator(cfg.JWTPublicKey)
		if err != nil {
			panic(err)
		}
		return middleware.AuthJWT(validator)
	default:
		return middleware.AuthDisabled
	}
}

func chain(final http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	h := final
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
