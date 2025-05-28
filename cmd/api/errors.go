package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s %s", err, r.Method, r.URL.Path)

	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s %s", err, r.Method, r.URL.Path)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s %s", err, r.Method, r.URL.Path)

	WriteJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s %s", err, r.Method, r.URL.Path)

	WriteJSONError(w, http.StatusNotFound, "not found")
}
