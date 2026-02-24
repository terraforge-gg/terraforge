package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/dto"
)

func JWTMiddleware(v *auth.Validator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
					Title:  "Unauthorized",
					Status: http.StatusUnauthorized,
					Detail: "Missing authorization header.",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
					Title:  "Unauthorized",
					Status: http.StatusUnauthorized,
					Detail: "Invalid authorization header format.",
				})
			}

			tokenString := parts[1]

			token, err := v.ValidateToken(tokenString)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
					Title:  "Unauthorized",
					Status: http.StatusUnauthorized,
					Detail: "Invalid token.",
				})
			}

			var id string
			err = token.Get("id", &id)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
					Title:  "Unauthorized",
					Status: http.StatusUnauthorized,
					Detail: "Invalid token.",
				})
			}

			c.Set("userId", id)

			return next(c)
		}
	}
}
