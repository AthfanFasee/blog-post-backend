package main

import (
	"net/http"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
	"github.com/AthfanFasee/blog-post-backend/internal/data"
)

func TestContextUser(t *testing.T) {
	app := newTestApplication(t)

	req := &http.Request{}
	user := &data.User{
		ID: 1,
	}

	reqWithUser := app.contextSetUser(req, user)
	userFromReq := app.contextGetUser(reqWithUser)

	assert.Equal(t, userFromReq.ID, user.ID)
}
