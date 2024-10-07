package config

import (
	"time"
)

type Config struct {
	DB Database `yaml:"db" json:"db"`
}

type Database struct {
	Port         string        `yaml:"port" json:"port"`
	Host         string        `yaml:"host" json:"host"`
	User         string        `yaml:"user" json:"user"`
	Password     string        `yaml:"password" json:"password"`
	Database     string        `yaml:"database" json:"database"`
	MaxConns     int32         `yaml:"maxConns" json:"maxConns"`
	MaxIdleConns int32         `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxLifeTime  time.Duration `yaml:"maxLifeTime" json:"maxLifeTime"`
}
