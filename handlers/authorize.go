package handlers

import (
	"github.com/coffeemakr/wedo"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"time"
)

var UsedTokenIssuer *JwtTokenIssuer

type JwtTokenIssuer struct {
	PrivateKey *jose.JSONWebKey
}

type JwtTokenVerifier struct {
	KeySet *jose.JSONWebKeySet
}

func (v *JwtTokenVerifier) verifyToken(rawToken string) (*wedo.DecodedToken, error) {
	var claims jwt.Claims
	var result wedo.DecodedToken
	token, err := jwt.ParseSigned(rawToken)
	if err != nil {
		return nil, err
	}

	err = token.Claims(v.KeySet, &claims)
	if err != nil {
		return nil, err
	}
	result.UserName = claims.Subject
	return &result, nil
}

func (i JwtTokenIssuer) IssueToken(decodedToken *wedo.DecodedToken) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: i.PrivateKey}, nil)
	if err != nil {
		return "", nil
	}

	// TODO: set expiry lower && refresh
	issuedAt := time.Now()
	expiry := issuedAt.AddDate(0, 1, 0)
	claims := jwt.Claims{
		ID:       RandStringRunes(16),
		Subject:  decodedToken.UserName,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(expiry),
	}

	raw, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		return "", err
	}
	return raw, err
}
