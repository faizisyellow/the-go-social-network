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
	app.logger.Warnw("not found error:", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusNotFound, err.Error())

}

func (app *application) unAuthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorize", "path", r.URL, "method", r.Method, "error", err.Error())

	WriteJSONError(w, http.StatusUnauthorized, "unauthorize")
}

func (app *application) unAuthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorize Basic error", "path", r.URL, "method", r.Method, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONError(w, http.StatusUnauthorized, "unauthorize")
}

func (app *application) forbiddenErrorResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden access", "path", r.URL, "method", r.Method)

	WriteJSONError(w, http.StatusForbidden, "forbidden")
}
