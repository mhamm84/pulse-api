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
	// Give the DB and Pulse API time to work out a connection
	time.Sleep(20 * time.Second)

	db, err := helper.OpenDB(helper.TestDns)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	mig, err := runMigration(db, "../../../migrations/")
	config.TestingConfig.Migration = mig
	config.TestingConfig.DB = db

	// Run the tests
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
