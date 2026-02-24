package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Validator struct {
	cachedSet jwk.CachedSet
	cache     *jwk.Cache
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewValidator(jwksURL string) (*Validator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create JWKS cache: %w", err)
	}

	registerCtx, registerCancel := context.WithTimeout(ctx, 5*time.Second)
	defer registerCancel()

	if err := c.Register(registerCtx, jwksURL); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to register JWKS URL: %w", err)
	}

	if _, err := c.Refresh(ctx, jwksURL); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to fetch initial JWKS: %w", err)
	}

	cachedSet, err := c.CachedSet(jwksURL)

	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to get cached set: %w", err)
	}

	return &Validator{
		cachedSet: cachedSet,
		cache:     c,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (v *Validator) Close() {
	v.cancel()
}

func (v *Validator) ValidateToken(tokenString string) (jwt.Token, error) {
	token, err := jwt.Parse([]byte(tokenString),
		jwt.WithKeySet(v.cachedSet),
		jwt.WithValidate(true),
		jwt.WithAcceptableSkew(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
