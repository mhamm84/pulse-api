package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const version = "1.0.0"
const (
	dev        = "dev"
	devCloud   = "dev-cloud"
	staging    = "stg"
	uat        = "uat"
	production = "prod"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	alphaVantage struct {
		token string
	}
}

type application struct {
	cfg      config
	logger   *jsonlog.Logger
	services Services
}

func main() {
	var cfg config

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	flag.IntVar(&cfg.port, "port", 9091, "Pulse API port number")
	flag.StringVar(&cfg.env, "env", devCloud, fmt.Sprintf("%s|%s|%s|%s|%s", dev, devCloud, staging, uat, production))
	// DB jdbc:postgresql://localhost:5432/pulse
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("PULSE_POSTGRES_DSN"), "Postgres DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.alphaVantage.token, "alpha-vantage-api-token", os.Getenv("ALPHA_VANTAGE_API_TOKEN"), "https://www.alphavantage.co/")

	flag.Parse()

	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	db, err := openDB(cfg, *logger)
	if err != nil {
		panic(err)
	}

	// Create Alpha Client
	alphaClient := alpha.NewClient(cfg.alphaVantage.token)

	app := application{
		cfg:      cfg,
		logger:   logger,
		services: NewAlphaServices(data.NewModels(db), alphaClient, logger),
	}

	app.startDataSyncs()

	err = app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

}

/*
 * Connect to MongoDB
 */
func openDB(cfg config, logger jsonlog.Logger) (*sqlx.DB, error) {
	logger.PrintInfo("connecting and pinging postgres", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sqlx.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that // passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value // less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string // to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
