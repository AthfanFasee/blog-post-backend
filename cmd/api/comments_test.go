package main

import (
	"net/http"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestShowCommentsForPostHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/api/v1/posts/comments/1",
			wantCode: http.StatusOK,
			wantBody: "Mocked Comment",
		},
		{
			name:     "Empty ID",
			urlPath:  "/api/v1/posts/comments/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/api/v1/posts/comments/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/api/v1/posts/comments/1.1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/api/v1/posts/comments/one",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
