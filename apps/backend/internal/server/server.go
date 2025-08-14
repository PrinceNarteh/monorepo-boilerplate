package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PrinceNarteh/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	logger     *zerolog.Logger
}

// New creates a new HTTP server instance
func New(cfg *config.Config, handler http.Handler, logger *zerolog.Logger) *Server {
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return &Server{
		httpServer: srv,
		logger:     logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info().Msgf("Starting HTTP server on port %s", s.httpServer.Addr)
	
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	
	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down HTTP server...")
	
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	
	s.logger.Info().Msg("HTTP server stopped")
	return nil
}
