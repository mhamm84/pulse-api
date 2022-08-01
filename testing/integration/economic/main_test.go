//go:build integration

package economic

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mhamm84/pulse-api/testing/helper"
	"github.com/mhamm84/pulse-api/testing/integration/config"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// add -> // +build integration
	// SETUP

	// postgres://username:password@url.com:5432/dbName
	time.Sleep(5 * time.Second)
	db, err := helper.OpenDB(helper.TestDns)
	_, err = db.Exec("CREATE EXTENSION citext")
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE")
	if err != nil {
		panic(err)
	}

	mig, err := runMigration(db, "../../../migrations/")
	config.TestingConfig.Migration = mig
	config.TestingConfig.DB = db

	os.Exit(m.Run())
}

func runMigration(dbConn *sqlx.DB, migrationsFolderLocation string) (*migrate.Migrate, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := postgres.WithInstance(dbConn.DB, &postgres.Config{})

	fmt.Println("running migration files from", pathToMigrate)
	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "postgres", driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}
