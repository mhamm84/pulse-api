package main

import (
	"fmt"
	"github.com/mhamm84/pulse-api/cmd/config"
	"github.com/mhamm84/pulse-api/cmd/pulse/api"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

const (
	version            = "1.0.0"
	logLevel           = "log-level"
	port               = "port"
	env                = "env"
	dbDsn              = "db-dsn"
	dbMaxOpenConns     = "db-max-open-conns"
	dbMaxIdleConns     = "db-max-idle-conns"
	dbMaxIdleTime      = "db-max-idle-time"
	alphaVantageUrl    = "alpha-vantage-base-url"
	alphaVantageToken  = "alpha-vantage-api-token"
	rateLimiterRPS     = "limiter-rps"
	rateLimiterBurst   = "limiter-burst"
	rateLimiterEnabled = "limiter-enabled"
	dataSyncEnable     = "data-sync-enable"

	dev        = "dev"
	staging    = "stg"
	uat        = "uat"
	production = "prod"
	cors       = "cors-trusted-origins"

	defaultPort           = 9091
	defaultMaxOpenConns   = 25
	defaultMaxIdleConns   = 25
	defaultMaxIdleTime    = "15m"
	defaultRatePerSeconds = 2
	defaultRateBurst      = 4
	defaultCors           = "http://localhost:9090"
)

var cfg config.ApiConfig

func RunApiCmd() *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "run-api",
		Short: "Launches the Pulse HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("RUNNING THE API")
			api.StartApi(&cfg)
		},
	} // End Run CMD

	// Parse arguments passed in on startup

	runCmd.Flags().StringVar(&cfg.LogLevel, logLevel, "INFO", "logging level [DEBUG,INFO,WARNING,ERROR,FATAL]")
	runCmd.Flags().IntVar(&cfg.Port, port, defaultPort, "Pulse API port number")
	runCmd.Flags().StringVar(&cfg.Env, env, dev, fmt.Sprintf("%s|%s|%s|%s", dev, staging, uat, production))
	// POSTGRESQL
	runCmd.Flags().StringVar(&cfg.DB.Dsn, dbDsn, os.Getenv("PULSE_POSTGRES_DSN"), "Postgres DSN")
	runCmd.Flags().IntVar(&cfg.DB.MaxOpenConns, dbMaxOpenConns, defaultMaxOpenConns, "PostgreSQL max open connections")
	runCmd.Flags().IntVar(&cfg.DB.MaxIdleConns, dbMaxIdleConns, defaultMaxIdleConns, "PostgreSQL max open connections")
	runCmd.Flags().StringVar(&cfg.DB.MaxIdleTime, dbMaxIdleTime, defaultMaxIdleTime, "PostgreSQL max connection idle time")
	// Alpha Vantage
	runCmd.Flags().StringVar(&cfg.AlphaVantage.BaseUrl, alphaVantageUrl, os.Getenv("ALPHA_VANTAGE_BASE_URL"), "Base Url for Alpha Vantage API - https://www.alphavantage.co/")
	runCmd.Flags().StringVar(&cfg.AlphaVantage.Token, alphaVantageToken, os.Getenv("ALPHA_VANTAGE_API_TOKEN"), "Auth Token for Alpha Vantage API - https://www.alphavantage.co/")
	// API Rate Limiter
	runCmd.Flags().Float64Var(&cfg.Limiter.RPS, rateLimiterRPS, defaultRatePerSeconds, "Rate limiter maximum requests per second")
	runCmd.Flags().IntVar(&cfg.Limiter.Burst, rateLimiterBurst, defaultRateBurst, "Rate limiter maximum burst")
	runCmd.Flags().BoolVar(&cfg.Limiter.Enabled, rateLimiterEnabled, true, "Enable rate limiter")
	// CORS
	runCmd.Flags().StringSliceVar(&cfg.Cors.TrustedOrigins, cors, []string{defaultCors}, "all the Cors trusted origin URLS, usage: --cors-trusted-origin=url1,url2")
	cfg.Cors.TrustedOrigins = strings.Fields(os.Getenv("PULSE_CORS_TRUSTED_ORIGIN"))
	// DATA SYNC TASKS
	dataSyncEnabled, _ := strconv.ParseBool(os.Getenv("PULSE_DATA_SYNC_ENABLE"))
	runCmd.Flags().BoolVar(&cfg.DataSync, dataSyncEnable, dataSyncEnabled, "enable/disable the data sync updates to data providers")

	return runCmd
} // End CMD
