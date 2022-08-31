package main

import (
	"fmt"
	"time"
)

type Config struct {
	App        App
	Server     Server
	GRPC       GRPC
	Logger     Logger
	PostgreSQL PostgreSQL
}

type App struct {
	Name        string `envconfig:"APP_NAME" default:"calendar" required:"true"`
	StorageType string `envconfig:"APP_STORAGE_TYPE" required:"true"`
	RouterType  string `envconfig:"APP_ROUTER_TYPE" required:"true"`
}

type Server struct {
	Host         string        `envconfig:"SERVER_HOST" default:"127.0.0.1" required:"true"`
	Port         string        `envconfig:"SERVER_PORT" default:"80" required:"true"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"30s" required:"true"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s" required:"true"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"30s" required:"true"`
}

type GRPC struct {
	Host string `envconfig:"GRPC_HOST" default:"127.0.0.1" required:"true"`
	Port string `envconfig:"GRPC_PORT" default:"50051" required:"true"`
}

type Logger struct {
	Level       string `envconfig:"LOGGER_LEVEL" default:"info" required:"true"`
	Development bool   `envconfig:"LOGGER_DEVELOPMENT" default:"false" required:"true"`
}

type PostgreSQL struct {
	Host     string `envconfig:"PG_HOST" required:"true"`
	Port     int    `envconfig:"PG_PORT" required:"true"`
	Username string `envconfig:"PG_USERNAME" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Database string `envconfig:"PG_DATABASE" required:"true"`
}

func (p *PostgreSQL) BuildDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database)
}

func NewConfig() Config {
	return Config{}
}
