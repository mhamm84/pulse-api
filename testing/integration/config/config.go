package config

import "github.com/golang-migrate/migrate"

type Config struct {
	Migration *migrate.Migrate
	Dsn       string
}

var TestingConfig = Config{
	Migration: nil,
	Dsn:       "postgres://pulse_user:password@localhost:5433/pulse_testing?sslmode=disable",
}
