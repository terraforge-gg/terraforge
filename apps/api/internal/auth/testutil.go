package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type TestAuth struct {
	Key  *ecdsa.PrivateKey
	JWKS jwk.Set
}

func NewTestAuth() (*TestAuth, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	pubKey, err := jwk.Import(key.Public())
	if err != nil {
		return nil, err
	}

	_ = pubKey.Set("kid", "test-key-id")
	_ = pubKey.Set("alg", "ES256")

	jwks := jwk.NewSet()
	_ = jwks.AddKey(pubKey)

	return &TestAuth{Key: key, JWKS: jwks}, nil
}

func (ta *TestAuth) GenerateToken(userId string, username string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":       userId,
		"name":     username,
		"username": username,
		"email":    email,
		"sub":      userId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iss":      "test-issuer",
	})
	token.Header["kid"] = "test-key-id"
	return token.SignedString(ta.Key)
}
