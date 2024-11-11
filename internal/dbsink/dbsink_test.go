package dbsink

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"walletapp/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (

	// In-Code skip hook, should be replaced with proper IE test.
	_IgnoreDBIntegrationTest = true
)

type TS struct {
	TMM time.Time `json:"TMM"`
	STR string    `json:"STR"`
}

// SelfUnmarshal should be used when you have to compare
// basicaly two exactly the same objects, but one got processed by
// Marshal=>Unmarshal, and its [time.Time] field(s) got corrupted,
// resulting in False Negative comparation via reflect.DeepEqual or assert.Equal.
//
//nolint:all // Insane function by itself.
func SelfUnmarshal[V any](t *testing.T, v *V) {
	t.Helper()
	bts, err := json.Marshal(v)
	require.NoError(t, err)
	err = json.Unmarshal(bts, v)
	require.NoError(t, err)
}

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
