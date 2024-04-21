package security_test

import (
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/security"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClaims(t *testing.T) {
	t.Parallel()

	gid := uuid.New()
	pid := uuid.New()

	token, err := security.NewToken(gid, pid)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	claims, err := security.ParseToken(token)
	assert.NoError(t, err)

	assert.Equal(t, gid, claims.GameId)
	assert.Equal(t, pid, claims.PlayerId)
	assert.NotEmpty(t, claims.ServerId)

	assert.True(t, claims.IssuedAt.Before(claims.ExpiresAt))
	assert.True(t, claims.IssuedAt.Before(time.Now()))
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestALotOfClaims(t *testing.T) {
	t.Parallel()

	gid := uuid.New()
	pid := uuid.New()

	for i := 0; i < 1000; i++ {
		token, err := security.NewToken(gid, pid)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		claims, err := security.ParseToken(token)
		assert.NoError(t, err)

		assert.Equal(t, gid, claims.GameId)
		assert.Equal(t, pid, claims.PlayerId)
		assert.NotEmpty(t, claims.ServerId)

		assert.True(t, claims.IssuedAt.Before(claims.ExpiresAt))
		assert.True(t, claims.IssuedAt.Before(time.Now()))
		assert.True(t, claims.ExpiresAt.After(time.Now()))
	}
}

func BenchmarkClaims(b *testing.B) {
	gid := uuid.New()
	pid := uuid.New()

	for i := 0; i < b.N; i++ {
		token, _ := security.NewToken(gid, pid)
		security.ParseToken(token)
	}
}
