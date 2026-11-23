// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

// Package server implements the EVoyage web service: HTTP routing, middleware,
// and handlers for the HTMX frontend.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server timeouts and limits.
const (
	// readHeaderTimeout bounds how long a client may take to send headers.
	readHeaderTimeout = 5 * time.Second
	// idleTimeout bounds how long an idle keep-alive connection is kept.
	idleTimeout = 60 * time.Second
	// shutdownTimeout bounds the graceful shutdown of in-flight requests.
	shutdownTimeout = 10 * time.Second
)

// Server is the EVoyage HTTP server.
type Server struct {
	logger  *slog.Logger
	mux     *http.ServeMux
	now     func() time.Time
	started time.Time
}

// New returns a Server that logs through logger. A nil logger discards logs.
// Uptime is measured from the moment New is called.
func New(logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	srv := &Server{
		logger: logger,
		mux:    http.NewServeMux(),
		now:    time.Now,
	}
	srv.started = srv.now()
	srv.routes()

	return srv
}

// Handler returns the server's handler, with request logging applied.
func (s *Server) Handler() http.Handler {
	return s.withRequestLogging(s.mux)
}

// routes registers every route of the web service.
func (s *Server) routes() {
	s.mux.HandleFunc("GET /{$}", s.handleIndex)
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
}

// Run listens on addr and serves HTTP until ctx is canceled, then shuts the
// server down gracefully.
func (s *Server) Run(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %q: %w", addr, err)
	}

	return s.Serve(ctx, listener)
}

// Serve serves HTTP on listener until ctx is canceled, then stops accepting
// connections and waits for in-flight requests to complete.
func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	httpServer := &http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- httpServer.Serve(listener)
	}()

	s.logger.Info("http server listening", "addr", listener.Addr().String())

	select {
	case err := <-serveErr:
		return wrapServeError(err)
	case <-ctx.Done():
	}

	s.logger.Info("http server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return wrapServeError(<-serveErr)
}

// wrapServeError normalizes the errors returned by http.Server.Serve.
func wrapServeError(err error) error {
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return fmt.Errorf("serve http: %w", err)
}
