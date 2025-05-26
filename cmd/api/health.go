package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "200",
		"version": version,
		"env":     app.config.env,
	}

	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}

}
