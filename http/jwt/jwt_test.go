package jwt

import (
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/http/session"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	session := session.NewSession("1234", "test@test.com", enum.AccountRoleUser, nil)
	token, err := GenerateJwt(session, 10*time.Second, false)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	t.Log(token)
}

// nolint
func TestParseJWT(t *testing.T) {
	ses, refresh, err := ParseJwt("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X3JvbGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiZXhwIjoxNzQxNTQxMTU4LCJpc19yZWZyZXNoIjpmYWxzZSwidXNlcl9pZCI6IjEyMzQifQ.wRKMJSraYhvZUzOBCHr7A5brst0q-AZq3kS6mrJk4Jo")
	require.Nil(t, ses)
	require.False(t, refresh)
	require.Error(t, err)
	t.Log(err)
}
