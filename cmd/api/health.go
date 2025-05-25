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

	if r.Header.Get("User-Agent") == "lizzy" {
		err := WriteJSONError(w, http.StatusBadRequest, "user ga harus lizzy coy")
		if err != nil {
			// TODO: refactor
			println(err)
		}
		return
	}

	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		// TODO: refactor
		println(err)
	}

}
