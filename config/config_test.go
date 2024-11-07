package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint:all // Test must contain bad practicies such as dynamic errors.
func TestGet(t *testing.T) {
	t.Parallel()

	os.Unsetenv("WORKSPACE") // TODO: May it harm other tests?
	_, err := Get()
	assert.ErrorIs(t, err, ErrMissingRequiredEnv)

	func() {
		defer func() {
			a := recover()
			if assert.IsType(t, err, a) {
				err, _ = a.(error)
				assert.ErrorIs(t, err, ErrMissingRequiredEnv)
			}
		}()
		_ = MustGet()

	}()

	os.Setenv("WORKSPACE", "dummy")

	cfg, err := Get()
	assert.NotZero(t, cfg)
	assert.NoError(t, err)

}
