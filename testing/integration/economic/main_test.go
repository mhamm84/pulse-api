//go:build integration

package economic

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	db, err := openDB("postgres://pulse_user:password@localhost:5433/pulse_testing?sslmode=disable")
	_, err = db.Exec("CREATE EXTENSION citext")
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE")
	if err != nil {
		panic(err)
	}

	mig, err := runMigration(db, "../../../migrations/")
	config.TestingConfig.Migration = mig

	os.Exit(m.Run())
}

func openDB(dsn string) (*sqlx.DB, error) {
	fmt.Println("connecting and pinging postgres")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	//duration, err := time.ParseDuration("15m")
	//if err != nil {
	//	return nil, err
	//}
	//db.SetConnMaxIdleTime(duration)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
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
