package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraforge-gg/terraforge/internal/auth"
)

func generateTestToken(t *testing.T, testAuth *auth.TestAuth, userId string, username string, email string) string {
	t.Helper()
	token, err := testAuth.GenerateToken(userId, username, email)
	require.NoError(t, err)
	return token
}
