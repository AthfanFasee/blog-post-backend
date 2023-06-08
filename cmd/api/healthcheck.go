package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	serverData := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"environment": app.config.env,
			"version":     version,
			"build_time":  buildTime,
		},
	}

	err := app.writeJSON(w, http.StatusOK, serverData, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
