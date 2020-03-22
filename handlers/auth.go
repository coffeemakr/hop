package handlers

import (
	"context"
	"errors"
	http_error "github.com/coffeemakr/go-http-error"
	"log"
	"net/http"
	"strings"
)

const (
	bearerTokenPrefix            = "Bearer "
	httpHeaderAuthorization = "Authorization"
	ContextUserName   ContextKey = "userID"
)

type ContextKey string

type Authenticator struct {
	Verifier *JwtTokenVerifier
}

func (a Authenticator) MiddleWare(next http.Handler) http.Handler {
	if next == nil {
		panic("next handler is nil")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(httpHeaderAuthorization)
		if authHeader == "" {
			http_error.ErrUnauthorized.
				CauseString("missing 'Authorization' header").
				Write(w, r)
			return
		}
		if !strings.HasPrefix(authHeader, bearerTokenPrefix) {
			http_error.ErrUnauthorized.
				Causef("missing '%s' prefix in auth header: %s", bearerTokenPrefix, authHeader).
				Write(w, r)
			return
		}
		// remove prefix
		authHeader = authHeader[len(bearerTokenPrefix):]
		decodedToken, err := a.Verifier.VerifyToken(authHeader)
		if err != nil {
			http_error.ErrBadRequest.Cause(err).Write(w, r)
			return
		}
		log.Printf("Got user ID: '%s'\n", decodedToken.UserName)
		ctx := context.WithValue(r.Context(), ContextUserName, decodedToken.UserName)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

var ErrNoUserName = errors.New("no username in request")

func GetUserNameFromRequest(r *http.Request) (string, error) {
	name, ok := r.Context().Value(ContextUserName).(string)
	if !ok || name == "" {
		return "", ErrNoUserName
	}
	return name, nil
}
