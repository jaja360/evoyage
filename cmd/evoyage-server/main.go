// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

// Command evoyage-server runs the EVoyage web service.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jaja360/evoyage/internal/server"
)

// defaultAddr is used when EVOYAGE_PORT is unset.
const defaultAddr = ":8080"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("evoyage-server exited with error", "error", err)
		os.Exit(1)
	}
}

// run starts the HTTP server and blocks until ctx is canceled.
func run(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	addr := listenAddr()
	logger.Info("starting evoyage-server", "addr", addr)

	if err := server.New(logger).Run(ctx, addr); err != nil {
		return err
	}

	logger.Info("evoyage-server stopped")

	return nil
}

// listenAddr returns the address to listen on, derived from EVOYAGE_PORT.
// Accepted forms are "8080", ":8080", and "host:8080".
func listenAddr() string {
	port := strings.TrimSpace(os.Getenv("EVOYAGE_PORT"))
	if port == "" {
		return defaultAddr
	}
	if strings.Contains(port, ":") {
		return port
	}

	return ":" + port
}
