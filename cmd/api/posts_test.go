package main

import (
	"net/http"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestShowPostsHandler(t *testing.T) {
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
			name:     "Without Params",
			urlPath:  "/api/v1/posts",
			wantCode: http.StatusOK,
			wantBody: "Mocked Post Title",
		},
		{
			name:     "Valid Title Param",
			urlPath:  "/api/v1/posts?title=title",
			wantCode: http.StatusOK,
			wantBody: "Title",
		},
		{
			name:     "Valid ID Param",
			urlPath:  "/api/v1/posts?id=1",
			wantCode: http.StatusOK,
			wantBody: "Mocked Post Title",
		},
		{
			name:     "Valid ID, Title Param",
			urlPath:  "/api/v1/posts?title=title&id=1",
			wantCode: http.StatusOK,
			wantBody: "Title",
		},
		{
			name:     "Valid ID, Invalid Title Param",
			urlPath:  "/api/v1/posts?title=invalid&id=1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty Title Param",
			urlPath:  "/api/v1/posts?title=",
			wantCode: http.StatusOK,
		},
		{
			name:     "Empty ID Param",
			urlPath:  "/api/v1/posts?id=",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent title Param",
			urlPath:  "/api/v1/posts?title=invalid",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Number Title Param",
			urlPath:  "/api/v1/posts?title=2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID Param",
			urlPath:  "/api/v1/posts?id=1.1",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "String ID Param",
			urlPath:  "/api/v1/posts?id=one",
			wantCode: http.StatusUnprocessableEntity,
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

func TestShowSinglePostHandler(t *testing.T) {
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
			urlPath:  "/api/v1/post/1",
			wantCode: http.StatusOK,
			wantBody: "Mocked Post Title",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/api/v1/post/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/api/v1/post/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/api/v1/post/1.1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/api/v1/post/one",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/api/v1/post/",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Valid ID returns user name",
			urlPath:  "/api/v1/post/1",
			wantCode: http.StatusOK,
			wantBody: "Mocked User",
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
