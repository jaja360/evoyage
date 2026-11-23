// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package server

import (
	"encoding/json"
	"io"
	"net/http"
)

// indexHTML is the placeholder page served at / until the HTMX UI lands.
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>EVoyage</title>
</head>
<body>
<main>
<h1>EVoyage</h1>
<p>EVoyage — EV route planner (scaffold)</p>
</main>
</body>
</html>
`

// healthResponse is the body of GET /healthz.
type healthResponse struct {
	// Status is "ok" while the process serves requests.
	Status string `json:"status"`
	// Uptime is the time elapsed since the server was created.
	Uptime string `json:"uptime"`
}

// handleHealthz reports that the process is alive.
func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	body := healthResponse{
		Status: "ok",
		Uptime: s.now().Sub(s.started).String(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.logger.Warn("encode health response", "error", err)
	}
}

// handleIndex serves the placeholder landing page.
func (s *Server) handleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, indexHTML); err != nil {
		s.logger.Warn("write index page", "error", err)
	}
}

// withRequestLogging logs one line per request, including its status and
// duration.
func (s *Server) withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := s.now()
		recorder := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		s.logger.InfoContext(r.Context(), "http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode(),
			"bytes", recorder.bytes,
			"duration", s.now().Sub(start).String(),
		)
	})
}

// responseRecorder captures the status code and body size written by a handler.
type responseRecorder struct {
	http.ResponseWriter
	bytes int
	// status is zero until the handler writes a header or a body.
	status int
}

// WriteHeader records the first status code written and forwards it.
func (r *responseRecorder) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

// Write records the number of bytes written and forwards them.
func (r *responseRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	written, err := r.ResponseWriter.Write(body)
	r.bytes += written

	return written, err
}

// statusCode returns the recorded status, defaulting to 200 like net/http.
func (r *responseRecorder) statusCode() int {
	if r.status == 0 {
		return http.StatusOK
	}
	return r.status
}
