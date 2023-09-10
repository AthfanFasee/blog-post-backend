package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AthfanFasee/blog-post-backend/internal/data"
	"github.com/AthfanFasee/blog-post-backend/internal/dto"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

func (app *application) showPostsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		data.Filters
	}

	v := validator.New()

	// Get the url.Values map containing the query string data.
	queryString := r.URL.Query()

	input.Filters.Sort = app.readString(queryString, "sort", "-id")
	input.Title = app.readString(queryString, "title", "")
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.ID = app.readInt(queryString, "id", 0, v)
	input.Filters.Limit = app.readInt(queryString, "limit", 6, v)

	// Add the supported sort values for this endpoint to the sort safelist.
	input.Filters.SortSafeList = []string{"id", "title", "readtime", "likescount", "-id", "-title", "-readtime", "-likescount"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	posts, metadata, err := app.models.Posts.GetAll(input.Title, input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"posts": posts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showSinglePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	post, userName, err := app.models.Posts.GetWithUserName(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	PostResponseBody := dto.PostResponseBody{
		ID:        post.ID,
		Title:     post.Title,
		PostText:  post.PostText,
		Img:       post.Img,
		ReadTime:  post.ReadTime,
		CreatedAt: post.CreatedAt,
		CreatedBy: post.CreatedBy,
		UserName:  *userName,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": PostResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePostRequestBody
	// Decoding JSON values in to input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	post := &data.Post{
		Title:     strings.TrimSpace(input.Title),
		PostText:  strings.TrimSpace(input.PostText),
		Img:       input.Img,
		ReadTime:  input.ReadTime,
		CreatedBy: user.ID,
	}

	v := validator.New()

	if data.ValidatePost(v, post); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	PostResponseBody := dto.PostResponseBody{
		ID:        post.ID,
		Title:     post.Title,
		PostText:  post.PostText,
		Img:       post.Img,
		ReadTime:  post.ReadTime,
		CreatedAt: post.CreatedAt,
		CreatedBy: user.ID,
		UserName:  user.Name,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": PostResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// Check if a post with provided id exists.
	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input dto.UpdatePostRequestBody

	// Decoding JSON values in to input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy values from req body to appropriate fields of post record only if they are not nil.
	if input.Title != nil {
		post.Title = strings.TrimSpace(*input.Title)
	}
	if input.PostText != nil {
		post.PostText = strings.TrimSpace(*input.PostText)
	}
	if input.Img != nil {
		post.Img = *input.Img
	}
	if input.ReadTime != nil {
		post.ReadTime = *input.ReadTime
	}

	v := validator.New()

	// Title and PostText must be provided by the client (other fields are optional when updating).
	if nil == input.Title {
		v.AddError("title", "must be provided")
	}
	if nil == input.PostText {
		v.AddError("postText", "must be provided")
	}

	if data.ValidatePost(v, post); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Update(post)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	PostResponseBody := dto.PostResponseBody{
		ID:        post.ID,
		Title:     post.Title,
		PostText:  post.PostText,
		Img:       post.Img,
		ReadTime:  post.ReadTime,
		LikedBy:   post.LikedBy,
		CreatedAt: post.CreatedAt,
		CreatedBy: post.CreatedBy,
		UserName:  user.Name,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": PostResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Posts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "post deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) likePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	// Check if a post with provided id exists. Return username of the user who created the post as well.
	post, userName, err := app.models.Posts.GetWithUserName(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Posts.AddLike(post, user.ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	PostResponseBody := dto.PostResponseBody{
		ID:        post.ID,
		Title:     post.Title,
		PostText:  post.PostText,
		Img:       post.Img,
		ReadTime:  post.ReadTime,
		LikedBy:   post.LikedBy,
		CreatedAt: post.CreatedAt,
		CreatedBy: post.CreatedBy,
		UserName:  *userName,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": PostResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) dislikePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	post, userName, err := app.models.Posts.GetWithUserName(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Posts.RemoveLike(post, user.ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	PostResponseBody := dto.PostResponseBody{
		ID:        post.ID,
		Title:     post.Title,
		PostText:  post.PostText,
		Img:       post.Img,
		ReadTime:  post.ReadTime,
		LikedBy:   post.LikedBy,
		CreatedAt: post.CreatedAt,
		CreatedBy: post.CreatedBy,
		UserName:  *userName,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": PostResponseBody}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
