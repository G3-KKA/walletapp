package app

import (
	"testing"
	"walletapp/config"

	"github.com/stretchr/testify/assert"
)

var (
	_IgnoreWholeAppTesting = false
)

func TestAppHappyPath(t *testing.T) {
	// t.Parallel() might be unsafe.

	if _IgnoreWholeAppTesting {
		t.SkipNow()
	}

	ap := New(config.Config{})
	defer func() {
		a := recover()
		assert.Nil(t, a)
	}()
	err := ap.Run()
	assert.NoError(t, err)

}
