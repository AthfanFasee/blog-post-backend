package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AthfanFasee/blog-post-backend/internal/data"
	"github.com/AthfanFasee/blog-post-backend/internal/dto"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

func (app *application) showCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comments, err := app.models.Comments.GetAllForPost(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comments": comments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.CommentRequestBody

	// Decoding JSON values in to input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	comment := &data.Comment{
		Text:      strings.TrimSpace(input.Text),
		PostID:    input.PostID,
		CreatedBy: user.ID,
	}

	v := validator.New()

	if data.ValidateComment(v, comment); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	err = app.models.Comments.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	CommentResponseBody := dto.CommentResponseBody{
		ID:        comment.ID,
		Text:      comment.Text,
		CreatedBy: comment.CreatedBy,
		PostID:    comment.PostID,
		UserName:  user.Name,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": CommentResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Comments.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "comment deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
