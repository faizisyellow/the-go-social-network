package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := NewTestApplication(t)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated request", func(t *testing.T) {

		// the set up
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// the action
		rr := ExecuteRequest(req, mux)

		CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticate request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/448", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", "Bearer "+testToken)

		rr := ExecuteRequest(req, mux)

		CheckResponseCode(t, http.StatusOK, rr.Code)
	})
}
