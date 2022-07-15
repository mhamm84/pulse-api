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
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	version = "1.0.0"

	dev        = "dev"
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
		baseUrl string
		token   string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
	logLevel string
}

type application struct {
	cfg      config
	logger   *jsonlog.Logger
	services Services
}

func main() {
	// Parse arguments passed in on startup
	var cfg config

	flag.StringVar(&cfg.logLevel, "log-level", "INFO", "logging level")
	flag.IntVar(&cfg.port, "port", 9091, "Pulse API port number")
	flag.StringVar(&cfg.env, "env", dev, fmt.Sprintf("%s|%s|%s|%s", dev, staging, uat, production))
	// DB jdbc:postgresql://localhost:5432/pulse
	flag.StringVar(&cfg.db.dsn, "db-dsn", "db-dsn", "Postgres DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// Alpha Vantage
	flag.StringVar(&cfg.alphaVantage.baseUrl, "alpha-vantage-base-url", os.Getenv("ALPHA_VANTAGE_BASE_URL"), "Base Url for Alpha Vantage API - https://www.alphavantage.co/")
	flag.StringVar(&cfg.alphaVantage.token, "alpha-vantage-api-token", os.Getenv("ALPHA_VANTAGE_API_TOKEN"), "Auth Token for Alpha Vantage API - https://www.alphavantage.co/")
	// API Rate Limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	// CORS
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()

	// Setup logging
	logger := jsonlog.New(os.Stdout, jsonlog.GetLevel(cfg.logLevel))

	// Fancy ascii splash when starting the app
	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	// Connect to the database
	db, err := openDB(cfg, *logger)
	if err != nil {
		panic(err)
	}

	// Create Alpha Client
	alphaClient := alpha.NewClient(cfg.alphaVantage.baseUrl, cfg.alphaVantage.token)
	// Create the app
	app := application{
		cfg:      cfg,
		logger:   logger,
		services: NewAlphaServices(data.NewModels(db), alphaClient, logger),
	}
	//TODO - Uncomment
	// Start the data sync tasks to keep data from the API up to date in the DB
	//app.startDataSyncs()

	// Serve the API
	err = app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

/*
 * Connect to DB
 */
func openDB(cfg config, logger jsonlog.Logger) (*sqlx.DB, error) {
	logger.PrintInfo("connecting and pinging postgres", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sqlx.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
