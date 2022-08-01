package config

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
)

type Config struct {
	Migration *migrate.Migrate
	Dsn       string
	DB        *sqlx.DB
}

var TestingConfig = Config{
	Migration: nil,
	Dsn:       "postgres://pulse_user:password@localhost:5433/pulse_testing?sslmode=disable",
}

func (c *Config) InsertTestData() error {
	files, err := ioutil.ReadDir("../sql")
	if err != nil {
		return err
	}

	for _, f := range files {
		fb, err := os.ReadFile("../sql/" + f.Name())
		if err != nil {
			return err
		}

		line := string(fb)
		fmt.Println(line)
		_, err = c.DB.Exec(line)
		if err != nil {
			return err
		}
	}
	return nil
}
