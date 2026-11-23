// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package planner_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jaja360/evoyage/internal/planner"
)

// StubClient must keep satisfying Client.
var _ planner.Client = planner.StubClient{}

func TestStubClientPlanRoute(t *testing.T) {
	t.Parallel()

	response, err := planner.StubClient{}.PlanRoute(context.Background(), sampleRequest())

	if !errors.Is(err, planner.ErrNotImplemented) {
		t.Fatalf("err = %v, want %v", err, planner.ErrNotImplemented)
	}
	if !reflect.DeepEqual(response, planner.RouteResponse{}) {
		t.Fatalf("response = %+v, want the zero value", response)
	}
}
