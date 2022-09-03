package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// %q will automatically wrap the variables in "". %s cannot do this
	data := map[string]string{
		"status": "available",
		"environment": app.config.env,
		"version": version,
	}

	// Marshal will return a []byte containing the encoded JSON
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Printf(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}