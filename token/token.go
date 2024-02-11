package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/storvik/goshrt"
)

var _ goshrt.Authorizer = &Auth{}

type Auth struct {
	token string // The master token
}

type claims struct {
	// Possible to add additional fields here
	jwt.RegisteredClaims
}

// NewAuth returns new jwt tokne authenticator.
func NewAuth(t string) *Auth {
	return &Auth{
		token: t,
	}
}

// Create creates a valid token to be used when authorizing clients.
// The string `id` is client name, and can by just about anything.
func (a *Auth) Create(id string) (string, error) {
	c := claims{
		jwt.RegisteredClaims{
			// Omitting Expire
			// ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
			ID:       "goshrt-jwt",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject:  id,
			Issuer:   "goshrt-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString([]byte(a.token))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate validates a token string `t` using the secret
// included in auth struct.
func (a *Auth) Validate(t string) (bool, error) {
	c := &claims{}

	token, err := jwt.ParseWithClaims(t, c, func(_ *jwt.Token) (interface{}, error) {
		return []byte(a.token), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return false, nil
		}

		return false, errors.New("could not parse jwt token")
	}

	if !token.Valid {
		return false, nil
	}

	return true, nil
}
