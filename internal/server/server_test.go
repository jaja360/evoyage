// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jaja360/evoyage/internal/server"
)

// discardLogger returns a logger that drops every record.
func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	server.New(discardLogger()).Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	res := rec.Result()
	t.Cleanup(func() { _ = res.Body.Close() })

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got := res.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var body struct {
		Status string `json:"status"`
		Uptime string `json:"uptime"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	if body.Status != "ok" {
		t.Fatalf("status field = %q, want %q", body.Status, "ok")
	}
	if _, err := time.ParseDuration(body.Uptime); err != nil {
		t.Fatalf("uptime %q is not a duration: %v", body.Uptime, err)
	}
}

func TestHealthzRejectsOtherMethods(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	server.New(discardLogger()).Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if allow := rec.Header().Get("Allow"); !strings.Contains(allow, http.MethodGet) {
		t.Fatalf("Allow = %q, want it to contain %q", allow, http.MethodGet)
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	server.New(discardLogger()).Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if want := "EVoyage — EV route planner (scaffold)"; !strings.Contains(string(body), want) {
		t.Fatalf("body does not contain %q:\n%s", want, body)
	}
}

func TestUnknownPathReturnsNotFound(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	server.New(discardLogger()).Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestRequestLogging(t *testing.T) {
	t.Parallel()

	var logs bytes.Buffer
	srv := server.New(slog.New(slog.NewJSONHandler(&logs, nil)))

	for _, path := range []string{"/healthz", "/missing"} {
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	}

	lines := strings.Split(strings.TrimSpace(logs.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("logged %d lines, want 2:\n%s", len(lines), logs.String())
	}

	tests := []struct {
		line       string
		wantPath   string
		wantStatus float64
	}{
		{line: lines[0], wantPath: "/healthz", wantStatus: http.StatusOK},
		{line: lines[1], wantPath: "/missing", wantStatus: http.StatusNotFound},
	}

	for _, test := range tests {
		var entry map[string]any
		if err := json.Unmarshal([]byte(test.line), &entry); err != nil {
			t.Fatalf("decode log line %q: %v", test.line, err)
		}

		if entry["msg"] != "http request" {
			t.Errorf("msg = %v, want http request", entry["msg"])
		}
		if entry["method"] != http.MethodGet {
			t.Errorf("method = %v, want %s", entry["method"], http.MethodGet)
		}
		if entry["path"] != test.wantPath {
			t.Errorf("path = %v, want %s", entry["path"], test.wantPath)
		}
		if entry["status"] != test.wantStatus {
			t.Errorf("status = %v, want %v", entry["status"], test.wantStatus)
		}
		if _, ok := entry["duration"]; !ok {
			t.Errorf("log line %q has no duration", test.line)
		}
	}
}

func TestServeGracefulShutdown(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- server.New(discardLogger()).Serve(ctx, listener)
	}()

	baseURL := "http://" + listener.Addr().String()
	waitForHealthy(t, baseURL)

	cancel()

	select {
	case serveErr := <-done:
		if serveErr != nil {
			t.Fatalf("Serve returned %v, want nil", serveErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return after context cancellation")
	}

	res, err := http.Get(baseURL + "/healthz")
	if err == nil {
		_ = res.Body.Close()
		t.Fatal("server still answers requests after shutdown")
	}
}

func TestRunReturnsListenError(t *testing.T) {
	t.Parallel()

	err := server.New(nil).Run(context.Background(), "not-a-valid-address")
	if err == nil {
		t.Fatal("Run returned nil, want a listen error")
	}
	if !strings.Contains(err.Error(), "listen") {
		t.Fatalf("error = %q, want it to mention listen", err)
	}
}

func TestRunStopsOnCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := server.New(nil).Run(ctx, "127.0.0.1:0"); err != nil {
		t.Fatalf("Run returned %v, want nil", err)
	}
}

// waitForHealthy polls baseURL until the server answers /healthz successfully.
func waitForHealthy(t *testing.T, baseURL string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		res, err := http.Get(baseURL + "/healthz")
		if err == nil {
			_ = res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("server at %s never became healthy", baseURL)
}
