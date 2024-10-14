package config

import (
	"time"
)

type Config struct {
	DB    Database `yaml:"db" json:"db"`
	Email Email    `yaml:"email" json:"email"`
	Auth  Auth     `yaml:"auth" json:"auth"`
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

type Email struct {
	AdminEmail string `yaml:"adminEmail" json:"adminEmail"`
	AdminPass  string `yaml:"adminPass" json:"adminPass"`
	ServerHost string `yaml:"serverHost" json:"serverHost"`
	ServerPort int    `yaml:"serverPort" json:"serverPort"`
}

type Auth struct {
	JWTKey string        `yaml:"jwtKey" json:"jwtKey"`
	JWTExp time.Duration `yaml:"jwtExp" json:"jwtExp"`
}
