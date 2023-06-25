package main

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	app := newTestApplication(t)

	t.Run("valid id parameter", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/test/1", nil)
		params := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)
		r = r.WithContext(ctx)

		id, err := app.readIDParam(r)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, id, int64(1))
	})

	t.Run("invalid id parameter", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/test/abc", nil)
		params := httprouter.Params{httprouter.Param{Key: "id", Value: "abc"}}
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)
		r = r.WithContext(ctx)

		_, err := app.readIDParam(r)
		assert.StringContains(t, err.Error(), "invalid id parameter")
	})
}

func TestReadString(t *testing.T) {
	app := newTestApplication(t)

	t.Run("key exists in query string", func(t *testing.T) {
		query := url.Values{}
		query.Add("test", "value")
		value := app.readString(query, "test", "default")

		assert.Equal(t, value, "value")
	})

	t.Run("key does not exist in query string", func(t *testing.T) {
		query := url.Values{}
		value := app.readString(query, "test", "default")

		assert.Equal(t, value, "default")
	})
}

func TestReadInt(t *testing.T) {
	app := newTestApplication(t)

	t.Run("valid integer in query string", func(t *testing.T) {
		query := url.Values{}
		query.Add("test", "5")
		v := validator.New()
		value := app.readInt(query, "test", 1, v)

		assert.Equal(t, value, 5)
		// Expecting no validation errors.
		assert.Equal(t, len(v.Errors), 0)
	})

	t.Run("invalid integer in query string", func(t *testing.T) {
		query := url.Values{}
		query.Add("test", "abc")
		v := validator.New()
		value := app.readInt(query, "test", 1, v)

		// Should return default value.
		assert.Equal(t, value, 1)
		// Expecting one validation error.
		assert.Equal(t, len(v.Errors), 1)
		assert.StringContains(t, v.Errors["test"], "must be an integer value")
	})
}

func TestBackground(t *testing.T) {
	app := newTestApplication(t)

	t.Run("function without panic", func(t *testing.T) {
		var testVal int
		fn := func() { testVal = 5 }
		app.background(fn)

		// Wait until all goroutines finish, otherwise the test might finish before goroutine is done.
		app.wg.Wait()

		assert.Equal(t, testVal, 5)
	})

	t.Run("function with panic", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("unexpected panic: %v", err)
			}
			app.wg.Wait()
		}()

		fn := func() { panic("test panic") }
		app.background(fn)
	})
}
