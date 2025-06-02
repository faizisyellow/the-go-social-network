package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"faizisyellow.github.com/thegosocialnetwork/internal/auth"
	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"go.uber.org/zap"
)

func NewTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()

	mockStore := store.NewMockStore()

	testAuth := auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		authenticator: &testAuth,
	}
}

func ExecuteRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("expected the response to be %d but we got %d", expected, actual)
	}
}
