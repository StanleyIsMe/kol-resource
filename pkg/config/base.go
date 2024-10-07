package config

import "time"

type Config[CustomConfigT any] struct {
	Env       string `env:"APP_ENV" env-default:"local"`
	Name      string `yaml:"name" json:"name"`
	PrettyLog bool   `yaml:"prettylog" json:"prettylog"`
	LogLevel  string `yaml:"logLevel" json:"logLevel"`
	Debug     bool   `yaml:"debug" json:"debug"`
	HTTP      struct {
		Port     int `yaml:"port" json:"port"`
		Timeouts struct {
			ReadTimeout       time.Duration `yaml:"readTimeout" json:"readTimeout"`
			ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" json:"readHeaderTimeout"`
			WriteTimeout      time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
			IdleTimeout       time.Duration `yaml:"idleTimeout" json:"idleTimeout"`
		} `yaml:"timeouts" json:"timeouts"`
	} `yaml:"http" json:"http"`

	CustomConfig CustomConfigT `yaml:"custom" json:"custom"`
}
