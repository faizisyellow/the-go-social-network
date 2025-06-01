package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"faizisyellow.github.com/thegosocialnetwork/internal/helpers"
	"faizisyellow.github.com/thegosocialnetwork/internal/mailer"
	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type registerUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type userWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUserHandler godoc
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"User Credentials"
//	@Success		201		{object}	userWithToken		"User Registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
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

	user := &store.User{
		Username: payloadRegister.Username,
		Email:    payloadRegister.Email,
	}

	if err := user.Password.Set(payloadRegister.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Token invitation
	plainToken := uuid.New().String()

	token, err := helpers.HashToken(plainToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Users.CreateAndInvite(ctx, user, token, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// TODO: the user is_active is not suppose to be true yet
	userWToken := userWithToken{
		User:  user,
		Token: plainToken,
	}

	isDevEnv := app.config.env == "Development"

	// the links is from the frontend router (http://localhost:5173/confirm/{plaintoken})
	activationUrl := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	// struct literal
	vars := struct {
		Username      string
		ActivationUrl string
	}{
		Username:      user.Username,
		ActivationUrl: activationUrl,
	}

	// send email
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isDevEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err.Error())

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			log.Printf("error deleting user while rollback, error: %v", err.Error())
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWToken); err != nil {
		app.internalServerError(w, r, err)

		return
	}

}

type createTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// CreateTokenHandler godoc
//
//	@Summary		Create a token for a user
//	@Description	Create a token authentication for user
//	@Tags			authentication
//	@Accept			json
//	@Param			payload	body		createTokenPayload	true	"User Credentials"
//	@Success		201		{string}	string				"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	// 1. get payload
	var payloadCreateToken createTokenPayload

	if err := readJSON(w, r, &payloadCreateToken); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payloadCreateToken); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	// 2. fetch the user
	user, err := app.store.Users.GetByEmail(r.Context(), payloadCreateToken.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unAuthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)

		return
	}
}
