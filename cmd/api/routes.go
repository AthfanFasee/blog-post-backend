package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// Converting our err helpers as handlers and using them instead of default err handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/posts", app.createPostHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/posts", app.showPostsHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/posts/:id", app.showSinglePostHandler)

	return router
}