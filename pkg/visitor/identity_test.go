package visitor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeriveKeyUsesStableAnonymousIDWithoutLeakingIt(t *testing.T) {
	const anonymousID = "550e8400-e29b-41d4-a716-446655440000"

	first, err := DeriveKey("secret", anonymousID, "203.0.113.8")
	require.NoError(t, err)
	second, err := DeriveKey("secret", anonymousID, "198.51.100.9")
	require.NoError(t, err)

	require.Equal(t, first, second)
	require.Len(t, first, 64)
	require.False(t, strings.Contains(first, anonymousID))
}

func TestDeriveKeyFallsBackToIP(t *testing.T) {
	first, err := DeriveKey("secret", "invalid", "203.0.113.8")
	require.NoError(t, err)
	second, err := DeriveKey("secret", "", "203.0.113.8")
	require.NoError(t, err)

	require.Equal(t, first, second)
	require.NotContains(t, first, "203.0.113.8")
}

func TestDeriveKeyRejectsMissingSecretOrIdentity(t *testing.T) {
	_, err := DeriveKey("", "550e8400-e29b-41d4-a716-446655440000", "")
	require.Error(t, err)

	_, err = DeriveKey("secret", "", "")
	require.Error(t, err)
}
