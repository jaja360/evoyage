// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Jaël Champagne Gareau

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseRecorder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		write      func(*responseRecorder)
		wantStatus int
		wantBytes  int
	}{
		{
			name:       "implicit 200 on body write",
			write:      func(r *responseRecorder) { _, _ = r.Write([]byte("hello")) },
			wantStatus: http.StatusOK,
			wantBytes:  5,
		},
		{
			name: "explicit status is kept",
			write: func(r *responseRecorder) {
				r.WriteHeader(http.StatusTeapot)
				_, _ = r.Write([]byte("tea"))
			},
			wantStatus: http.StatusTeapot,
			wantBytes:  3,
		},
		{
			name: "first status wins",
			write: func(r *responseRecorder) {
				r.WriteHeader(http.StatusCreated)
				r.WriteHeader(http.StatusInternalServerError)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "status defaults to 200 without writes",
			write:      func(*responseRecorder) {},
			wantStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			recorder := &responseRecorder{ResponseWriter: httptest.NewRecorder()}
			test.write(recorder)

			if got := recorder.statusCode(); got != test.wantStatus {
				t.Errorf("statusCode() = %d, want %d", got, test.wantStatus)
			}
			if recorder.bytes != test.wantBytes {
				t.Errorf("bytes = %d, want %d", recorder.bytes, test.wantBytes)
			}
		})
	}
}
