package http

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/evzpav/user-auth/pkg/log"
)

// Server ...
type Server struct {
	server *http.Server
	log    log.Logger
}

// New ...
func New(handler http.Handler, host, port string, log log.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:         host + ":" + port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 55 * time.Second,
			Handler:      handler,
		},
		log: log,
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	go func() {
		s.log.Info().Sendf("user-auth is running on %s!", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error().Err(err).Sendf("Error on ListenAndServe: %q", err)
		}
	}()
}

// Shutdown ...
func (s *Server) Shutdown() {
	s.log.Info().Send("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		s.log.Error().Err(err).Sendf("Could not shutdown in 60s: %q", err)
		return
	}

	s.log.Info().Sendf("Server gracefully stopped")
}
