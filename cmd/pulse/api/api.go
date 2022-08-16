package api

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/cmd/config"
	"github.com/mhamm84/pulse-api/cmd/pulse/helper"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/mailer"
	"github.com/mhamm84/pulse-api/internal/repo"
	"github.com/mhamm84/pulse-api/internal/services"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	cfg      config.ApiConfig
	logger   *jsonlog.Logger
	services services.ServicesModel
	mailer   mailer.Mailer
}

func StartApi(cfg *config.ApiConfig) {
	logger := jsonlog.New(os.Stdout, jsonlog.GetLevel(cfg.LogLevel))

	// Fancy ascii splash when starting the app
	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	db, err := helper.OpenDB(&cfg.DB, 5, time.Second*2)
	if err != nil {
		panic(err)
	}

	// Create Alpha Client
	alphaClient := alpha.NewClient(cfg.AlphaVantage.BaseUrl, cfg.AlphaVantage.Token)

	// Create the app
	app := application{
		cfg:      *cfg,
		logger:   logger,
		services: services.NewServicesModel(repo.NewModels(db), alphaClient, logger),
		mailer:   mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	}

	// Start the data sync tasks to keep data from the API up to date in the DB
	if cfg.DataSync {
		logger.PrintInfo("Starting startEconomicReportDataSync", nil)
		app.startEconomicReportDataSync()
	}

	logConfig(logger, cfg)

	// Serve the API
	err = app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

func logConfig(logger *jsonlog.Logger, cfg *config.ApiConfig) {
	logger.PrintInfo("SMTP Server Config", map[string]interface{}{
		"host":     cfg.SMTP.Host,
		"port":     cfg.SMTP.Port,
		"username": cfg.SMTP.Username,
		"password": cfg.SMTP.Password,
		"sender":   cfg.SMTP.Sender,
	})
	logger.PrintInfo("AlphaVantage Config", map[string]interface{}{
		"baseUrl": cfg.AlphaVantage.BaseUrl,
		"token":   cfg.AlphaVantage.Token,
	})
	logger.PrintInfo("Rate Limiter", map[string]interface{}{
		"enabled": cfg.Limiter.Enabled,
		"rps":     cfg.Limiter.RPS,
		"burst":   cfg.Limiter.Burst,
	})
	logger.PrintInfo("API", map[string]interface{}{
		"port":     cfg.Port,
		"env":      cfg.Env,
		"cors":     cfg.Cors.TrustedOrigins,
		"dataSync": cfg.DataSync,
		"logLevel": cfg.LogLevel,
	})
	logger.PrintInfo("DB", map[string]interface{}{
		"dsn":          cfg.DB.Dsn,
		"maxOpenConns": cfg.DB.MaxOpenConns,
		"maxIdleConns": cfg.DB.MaxIdleConns,
		"maxIdleTime":  cfg.DB.MaxIdleTime,
	})
}
