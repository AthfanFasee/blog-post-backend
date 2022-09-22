package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Converting our err helpers as handlers and using them instead of default err handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/posts", app.requireAuthenticatedUser(app.showPostsHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/posts/:id", app.requireAuthenticatedUser(app.showSinglePostHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/posts", app.requireActivatedUser(app.createPostHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/posts/:id", app.requireActivatedUser(app.updatePostHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/posts/:id", app.requireActivatedUser(app.deletePostHandler))

	router.HandlerFunc(http.MethodPost, "/api/v1/auth/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/auth/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/auth/login", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
