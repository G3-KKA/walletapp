package config

import "errors"

var (
	ErrMissingRequiredEnv = errors.New("required env not defined")
)
