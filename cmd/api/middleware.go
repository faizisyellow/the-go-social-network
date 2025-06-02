package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. get auth header
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// 2. parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unAuthorizedBasicErrorResponse(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unAuthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. get auth header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		// 2. parsing header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// 3. claims the token
		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, int(userID))
		if err != nil {
			app.unAuthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) CheckPostOwnership(Requiredrole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		post := getPostFromContext(r)
		user := getUserFromContext(r)

		// if it is the user's post
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// role presedence check
		allowed, err := app.checkRolePresedence(r.Context(), user, Requiredrole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenErrorResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePresedence(ctx context.Context, user *store.User, rolename string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, rolename)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
