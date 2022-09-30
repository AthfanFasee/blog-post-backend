package main

import (
	"net/http"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert.go"
)

func TestShowSinglePostHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name string
		urlPath string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID",
			urlPath: "/api/v1/post/1",
			wantCode: http.StatusOK,
			wantBody: "Mocked Post Title",
		},
		{
			name: "Non-existent ID",
			urlPath: "/api/v1/post/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/api/v1/post/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/api/v1/post/1.1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/api/v1/post/one",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/api/v1/post/",
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