package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AthfanFasee/blog-post-backend/internal/data"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

func (app *application) showPostsHandler (w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Genres []string
		data.Filters
	}

	v := validator.New()

	// Get the url.Values map containing the query string data
	queryString := r.URL.Query()

	input.Filters.Sort = app.readString(queryString, "sort", "-id")
	input.Title = app.readString(queryString, "title", "")
	input.Genres = app.readCSV(queryString, "genres", []string{})
	input.Filters.Page = app.readInt(queryString, "page", 1, v)
	input.Filters.ID = app.readInt(queryString, "id", 0, v)
	input.Filters.Limit = app.readInt(queryString, "limit", 6, v)

	// Add the supported sort values for this endpoint to the sort safelist.
	// NEED TO FIGURE OUT MULTIPLE SORTING. ONE THIMNG IS U CAN USE CSV HELPER FOR SORTINGS TOO
	input.Filters.SortSafeList = []string{"id", "title", "readtime", "likescount", "-id", "-title", "-readtime", "-likescount"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	posts, metadata, err := app.models.Posts.GetAll(input.Title, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"posts": posts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w,r, err)
	}
}

func (app *application) showSinglePostHandler (w http.ResponseWriter, r *http.Request) {
	
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
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

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		PostText string `json:"postText"`
		ReadTime data.ReadTime `json:"readTime"`
		Img string `json:"img"`
		CreatedBy int
	}
	
	// Decoding JSON values in to input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
	app.badRequestResponse(w, r, err)
	return
	}

	post := &data.Post{
		Title: input.Title,
		PostText: input.PostText,
		Img: input.Img,
		ReadTime: input.ReadTime,
		CreatedBy: 1,
	}
	// Initialize a new Validator
	v := validator.New()

	if data.ValidatePost(v, post); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Location header mentioning from which URL client can find the newly-created resource at
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/v1/posts/%d", post.ID))
	
	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler (w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	
	// Check if a post with provided id exists
	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r , err)
		}
		return
	}

	// Define the input struct in a way, all the field got zero value 'nil'
	var input struct {
		Title *string `json:"title"`
		PostText *string `json:"postText"`
		Img *string `json:"img"`
		ReadTime *data.ReadTime `json:"readTime"`
	}

	// Decoding JSON values in to input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	// Copy values from req body to appropriate fields of post record
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.PostText != nil {
		post.PostText = *input.PostText
	}
	if input.Img != nil {
		post.Img = *input.Img
	}
	if input.ReadTime != nil {
		post.ReadTime = *input.ReadTime
	}

	v := validator.New()

	// Title and PostText must be provided by the client (other fields are optional)
	if input.Title == nil {
		v.AddError("title", "must be provided")
	}
	if input.PostText == nil {
		v.AddError("postText", "must be provided")
	}

	if data.ValidatePost(v, post); !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Posts.Update(post)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "post deleted succesfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}