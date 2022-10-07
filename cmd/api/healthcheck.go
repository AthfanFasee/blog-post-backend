package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// %q will automatically wrap the variables in "". %s cannot do this
	serverData := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"environment": app.config.env,
			"version":     version,
			"build_time":  buildTime,
		},
	}

	// Marshal will return a []byte containing the encoded JSON
	err := app.writeJSON(w, http.StatusOK, serverData, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
