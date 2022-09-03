package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Retrieve a slice containing req parameter names and values
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// Marshal will return a []byte containing the encoded JSON
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to the JSON to make it easier to view in terminal
	js = append(js, '\n')

	// Go will not loop over if the map is nil
	for key, value := range headers {
		fmt.Println(key, value, "jfjfj")
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(w.Header())
	w.WriteHeader(status)
	w.Write(js)

	return nil
}