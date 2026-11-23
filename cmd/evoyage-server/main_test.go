// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package main

import (
	"context"
	"strings"
	"testing"
)

func TestListenAddr(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want string
	}{
		{name: "default when unset", env: "", want: ":8080"},
		{name: "bare port", env: "9090", want: ":9090"},
		{name: "port with colon", env: ":9091", want: ":9091"},
		{name: "host and port", env: "127.0.0.1:9092", want: "127.0.0.1:9092"},
		{name: "surrounding whitespace is trimmed", env: " 9093 ", want: ":9093"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("EVOYAGE_PORT", test.env)

			if got := listenAddr(); got != test.want {
				t.Fatalf("listenAddr() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRunStopsOnCancelledContext(t *testing.T) {
	t.Setenv("EVOYAGE_PORT", "127.0.0.1:0")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := run(ctx); err != nil {
		t.Fatalf("run() = %v, want nil", err)
	}
}

func TestRunReturnsListenError(t *testing.T) {
	t.Setenv("EVOYAGE_PORT", "not-a-valid-address")

	if err := run(context.Background()); err == nil {
		t.Fatal("run() = nil, want a listen error")
	} else if !strings.Contains(err.Error(), "listen") {
		t.Fatalf("run() = %q, want it to mention listen", err)
	}
}
