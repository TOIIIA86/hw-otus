package main

import "fmt"

type Config struct {
	PostgreSQL PostgreSQL
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
