package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Converting our err helpers as handlers and using them instead of default err handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/posts", app.createPostHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/posts", app.showPostsHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/posts/:id", app.showSinglePostHandler)
	router.HandlerFunc(http.MethodPatch, "/api/v1/posts/:id", app.updatePostHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/posts/:id", app.deletePostHandler)

	return app.recoverPanic(app.rateLimit(router))
}