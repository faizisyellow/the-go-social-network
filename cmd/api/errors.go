package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusNotFound, err.Error())
}
