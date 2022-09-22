package main

import (
	"context"
	"net/http"

	"github.com/AthfanFasee/blog-post-backend/internal/data"
)

// Prevent naming collisions in request context by defining a custom type
type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("user value is missing in request context")
	}

	return user
}
