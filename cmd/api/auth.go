package main

import (
	"net/http"

	"faizisyellow.github.com/thegosocialnetwork/internal/store"
)

type registerUserPayload struct {
	Username string `json:"user" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// RegisterUserHandler godoc
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"User Credentials"
//	@Success		201		{object}	store.User			"User Registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Route			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payloadRegister registerUserPayload

	if err := readJSON(w, r, &payloadRegister); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payloadRegister); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	user := store.User{
		Username: payloadRegister.Username,
		Email:    payloadRegister.Email,
	}

	if err := user.Password.Set(payloadRegister.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// hash password

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
