package config

import (
	"errors"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

// Represents expected contents of configuration file.
type Config struct {
	WORKSPACE string `env:"WORKSPACE" env-required:"true"`
	Database
}
type Database struct {
	DBUser
	DBAddress string `env:"DBAddress" env-required:"true"`
	DBPort    string `env:"DBPort"    env-required:"true"`
	DBName    string `env:"DBName"    env-required:"true"`
}
type DBUser struct {
	DBRole     string `env:"DBRole"     env-required:"true"`
	DBPassword string `env:"DBPassword" env-required:"true"`
}

// Get reads from CONFIG_FILE.
// Return config or zero value config and error.
func Get() (Config, error) {
	var c Config

	err := cleanenv.ReadEnv(&c)
	if err != nil {

		// cleanenv errors are fully dynamic and do not suppor errors.Is().
		if strings.Contains(err.Error(), " is required but the value is not provided") {
			err = errors.Join(err, ErrMissingRequiredEnv)
		}

		return Config{}, err
	}

	return c, nil
}

// MustGet reads from CONFIG_FILE.
// Return config or panics, if any error happened.
func MustGet() Config {

	cfg, err := Get()
	if err != nil {
		panic(err)
	}

	return cfg
}
