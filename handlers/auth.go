package handlers

import (
	"context"
	"encoding/json"
	"errors"
	http_error "github.com/coffeemakr/go-http-error"
	"log"
	"net/http"
	"strings"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func getKeySetFromURL(jwksURL string) (*jose.JSONWebKeySet, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}

	jwks := new(jose.JSONWebKeySet)
	err = json.NewDecoder(resp.Body).Decode(jwks)
	return jwks, err
}

func check(toCheck bool, message string) {
	if !toCheck {
		log.Fatal(message)
	}
}

type jwtVerifier struct {
	JSONWebKeySet *jose.JSONWebKeySet
	Expected      jwt.Expected
}

var ErrInvalidToken = errors.New("invalid token")

func (v jwtVerifier) validateToken(raw string) (*jwt.Claims, error) {
	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return nil, ErrInvalidToken
	}

	cl := jwt.Claims{}
	if err := tok.Claims(v.JSONWebKeySet, &cl); err != nil {
		return nil, ErrInvalidToken
	}

	err = cl.Validate(v.Expected)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return &cl, nil
}

func (v jwtVerifier) MiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http_error.ErrUnauthorized.
				CauseString("missing 'Authorization' header").
				Write(w, r)
			return
		}
		if !strings.HasPrefix(authHeader, bearerTokenPrefix) {
			http_error.ErrUnauthorized.
				Causef("missing 'Bearer ' prefix in auth header: %s", authHeader).
				Write(w, r)
			return
		}
		authHeader = authHeader[len(bearerTokenPrefix):]
		claims, err := v.validateToken(authHeader)
		if err != nil {
			http_error.ErrBadRequest.Cause(err).Write(w, r)
			return
		}
		log.Println("Got user ID: ", claims.Subject)
		ctx := context.WithValue(r.Context(), ContextKeyUserId, claims.Subject)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newJWTVerifierFromDomain(jwksUrl string, expected jwt.Expected) (*jwtVerifier, error) {
	jwks, err := getKeySetFromURL(jwksUrl)
	if err != nil {
		log.Fatal(err)
	}
	verifier := new(jwtVerifier)
	verifier.Expected = expected
	verifier.JSONWebKeySet = jwks
	return verifier, nil
}

var bearerTokenPrefix = "Bearer "

type ContextKey string

const ContextKeyUserId ContextKey = "userID"

func getUserIDFromRequest(r *http.Request) string {
	return r.Context().Value(ContextKeyUserId).(string)
}

func NewJWTMiddleWareFromURL(jwksUrl string, expected jwt.Expected) func(http.Handler) http.Handler {
	verifier, err := newJWTVerifierFromDomain(jwksUrl, expected)
	if err != nil {
		log.Fatal(err)
	}
	return verifier.MiddleWare
}
