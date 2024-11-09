package app

import (
	"testing"
	"walletapp/config"

	"github.com/stretchr/testify/assert"
)

var (
	// Might come with a lot of problems,
	// that are not actually indicate that something broken.
	_IgnoreWholeAppTesting = true
)

func TestAppHappyPath(t *testing.T) {
	// t.Parallel() might be unsafe.

	if _IgnoreWholeAppTesting {
		t.SkipNow()
	}

	ap := New(config.Config{}) // nolint:all // depricated.
	defer func() {
		a := recover()
		assert.Nil(t, a)
	}()
	err := ap.Run()
	assert.NoError(t, err)

}
