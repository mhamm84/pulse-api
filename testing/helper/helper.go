package helper

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const TestDns = "postgres://pulse_user:password@localhost:5433/pulse_testing?sslmode=disable"

func OpenDB(dsn string) (*sqlx.DB, error) {
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
