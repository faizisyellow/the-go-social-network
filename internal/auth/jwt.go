package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthenticator struct {
	secret, aud, iss string
}

func NewJwtAuthenticator(secret, aud, iss string) *JwtAuthenticator {

	return &JwtAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (j *JwtAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

func (j *JwtAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(j.aud),
		jwt.WithIssuer(j.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
