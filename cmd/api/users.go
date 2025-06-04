package main

import (
	"net/http"
	"strconv"

	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@summary		Fetch a user profile
//	@Description	Fetch a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"userID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Securiy		ApikeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@summary		Follows a user
//	@Description	Follows a user profile by ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int		true	"userID"
//	@Success		201	{object}	string	"follow user successfully"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Securiy		ApikeyAuth
//	@Router			/users/{id}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	followedID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Followers.Follow(r.Context(), followedID, followerUser.ID)
	if err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "follow user successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnFollowUser godoc
//
//	@summary		Unfollows a user
//	@Description	Unfollows a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"userID"
//	@Success		201	{object}	string	"unfollow user successfully"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Securiy		ApikeyAuth
//	@Router			/users/{id}/unfollow [put]
func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	unfollowUserID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Followers.UnFollow(r.Context(), unfollowUserID, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "unfollow user successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getUserFromContext(r *http.Request) *store.User {

	return r.Context().Value(userCtx).(*store.User)
}

// ActivateUser godoc
//
//	@Summary		Activate/Register a user
//	@Description	Activate/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, "User activated"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
