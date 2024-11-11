package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (

	// Any change to the config requires test of config for it to pass validation,
	// by itself the fact that config even fails IS GOOD, cuz SOME KIND of validation ->
	// -> is happening, this flag should be set to FALSE only in pre-Release version,
	// where it's okay to hardcode correct config once, when its structure should stay,
	// more-or-less the same till release.
	_IgnoreInconvenientConfigTest = true
)

// nolint:all // Test must contain bad practicies such as dynamic errors.
func TestGet(t *testing.T) {
	if _IgnoreInconvenientConfigTest {
		t.Skip()
	}

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
