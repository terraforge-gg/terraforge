package utils

import "github.com/labstack/echo/v5"

// Extracts the userId from the Echo request context.
// Returns the userId string and ok bool indicating success
func GetSessionUserId(c *echo.Context) (string, bool) {
	val := c.Get("userId")
	if val == nil {
		return "", false
	}

	s, ok := val.(string)
	if !ok || s == "" {
		return "", false
	}

	return s, true
}
