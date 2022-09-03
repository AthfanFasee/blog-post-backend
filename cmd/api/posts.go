package main

import (
	"fmt"
	"net/http"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Post")
}

func (app *application) showPostsHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Show all movies")
}

func (app *application) showSinglePostHandler (w http.ResponseWriter, r *http.Request) {
	
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of post %d\n", id)
}