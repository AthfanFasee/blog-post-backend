package main

import (
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert.go"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/api/v1/healthcheck")

	assert.Equal(t, code, 200)
	
	if body != "" {
		assert.StringContains(t, body, "available")
	}
	
}