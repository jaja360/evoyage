// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package planner

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned by clients that cannot plan routes yet.
var ErrNotImplemented = errors.New("planner: plan route is not implemented")

// Client plans routes. The web layer depends on this interface only, so the
// concrete sidecar transport can be swapped in without touching callers.
type Client interface {
	// PlanRoute returns the route for req or an error if none can be planned.
	PlanRoute(ctx context.Context, req RouteRequest) (RouteResponse, error)
}

// StubClient is a Client that always fails with ErrNotImplemented. It stands in
// until the planner sidecar exists.
type StubClient struct{}

// PlanRoute implements Client. It always returns ErrNotImplemented.
func (StubClient) PlanRoute(_ context.Context, _ RouteRequest) (RouteResponse, error) {
	return RouteResponse{}, ErrNotImplemented
}
