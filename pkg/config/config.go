package config

import (
	"context"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// LoadWithEnv loads the configuration from the given path.
func LoadWithEnv[CustomConfigT any](_ context.Context, configPath string) (*Config[CustomConfigT], error) {
	var cfg Config[CustomConfigT]

	currEnv := "local"
	if e := os.Getenv("APP_ENV"); e != "" {
		currEnv = e
	}

	if err := cleanenv.ReadConfig(configPath+"/base.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("read base config failed: %w", &MissingBaseConfigError{Err: err})
	}

	cfg.Env = currEnv

	if err := cleanenv.ReadConfig(fmt.Sprintf("%s/%s.yaml", configPath, cfg.Env), &cfg); err != nil {
		return nil, fmt.Errorf("read %s config failed: %w", cfg.Env, &MissingEnvConfigError{Env: cfg.Env, Err: err})
	}

	return &cfg, nil
}
