package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AthfanFasee/blog-post-backend/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// Encode data into JSON
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// MarshalIndent will return a []byte containing the encoded JSON with any prefix and indent added
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to the JSON to make it easier to view in terminal
	js = append(js, '\n')

	// Go will not loop over if the map is nil
	// The reason for not using Set here is that Set takes string as value, but in here our value is a []string
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Decode JSON values
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit the size of body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// This will return error if any field of JSON cannot be mapped to dst, instead of ignoring
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)

	// Using a triage to handle errors properly
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		// fmt.Errorf, errors.New both will return an error type which is an interface with Error() method attached
		switch {
		// This err comes when we pass a nil pointer. Unexpected errors from Server better be handled by panic()
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.As(err, &unmarshalTypeError):
			// Field means the key or field in JSON object and It could be empty as well
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body cannot be larger than %d bytes", maxBytes)

		default:
			return err
		}
	}

	// Second call to decode will return io.EOF error if there's only one JSON value in req body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body cannot contain more than one JSON value")
	}

	return nil
}

// Read ID param from request url
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Retrieve a slice containing req parameter names and values
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Return a string value from query string map or a default value
func (app *application) readString(queryString url.Values, key string, defaultValue string) string {
	stringValue := queryString.Get(key)

	// Get() will return empty string if value is'nt found. Then we return default value
	if stringValue == "" {
		return defaultValue
	}

	return stringValue
}

// Return an int value from query string map or a default value
func (app *application) readInt(queryString url.Values, key string, defaultValue int, v *validator.Validator) int {
	stringValue := queryString.Get(key)

	if stringValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return intValue
}

// takes an artbitary function as parameter
func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		// Recover panic
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
