package dbsink

import (
	"context"
	"testing"
	"walletapp/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (

	// In-Code skip hook, should be replaced with proper IE test.
	_IgnoreDBIntegrationTest = false
)

func TestNewNotConnect(t *testing.T) {
	if _IgnoreDBIntegrationTest {
		t.SkipNow()
	}

	// Linter is angry, and he is right on point, test invalid.
	cfg := config.Database{}

	l := zerolog.New(zerolog.NewTestWriter(t))

	db := New(&l, cfg)

	err := db.Ping(context.TODO())

	require.NoError(t, err)

	err = db.Shutdown(context.TODO())

	// Useless, no error will be returned anyway.
	assert.NoError(t, err)
}
