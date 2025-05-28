package main

import (
	"context"
	"net/http"
	"strconv"

	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type FollowerUser struct {
	UserID int `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followUser := getUserFromContext(r)

	// TODO: change the user's ID payload from auth middleware
	payload := FollowerUser{}
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.Follow(r.Context(), followUser.ID, payload.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, "follow user successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowUser := getUserFromContext(r)

	// TODO: change the user's ID payload from auth middleware
	payload := FollowerUser{}
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.UnFollow(r.Context(), unFollowUser.ID, payload.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, "unfollow user successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getUserFromContext(r *http.Request) *store.User {

	return r.Context().Value(userCtx).(*store.User)
}
