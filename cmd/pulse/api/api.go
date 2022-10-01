package api

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/cmd/config"
	"github.com/mhamm84/pulse-api/cmd/pulse/helper"
	"github.com/mhamm84/pulse-api/internal/mailer"
	"github.com/mhamm84/pulse-api/internal/repo"
	"github.com/mhamm84/pulse-api/internal/services"
	"github.com/mhamm84/pulse-api/internal/utils"
	"go.uber.org/zap"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	cfg      config.ApiConfig
	services services.ServicesModel
	mailer   *mailer.Mailer
	wg       *sync.WaitGroup
}

func StartApi(cfg *config.ApiConfig) {
	// Fancy ascii splash when starting the app
	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	db, err := helper.OpenDB(&cfg.DB, 5, time.Second*2)
	if err != nil {
		panic(err)
	}

	// Create Alpha Client
	alphaClient := alpha.NewClient(cfg.AlphaVantage.BaseUrl, cfg.AlphaVantage.Token)

	// SMTP mailer
	mailer := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)

	// Create the app
	app := application{
		cfg:      *cfg,
		services: services.NewServicesModel(repo.NewModels(db), alphaClient, mailer),
		mailer:   mailer,
	}

	// Start the data sync tasks to keep data from the API up to date in the DB
	if cfg.DataSync {
		utils.Logger.Info("Starting startEconomicReportDataSync")
		app.startEconomicReportDataSync()
	}

	logConfig(cfg)

	// Serve the API
	err = app.serve()
	if err != nil {
		utils.Logger.Fatal("fatal error while serving the api", zap.Error(err))
	}
}

func logConfig(cfg *config.ApiConfig) {
	utils.Logger.Info("API",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("env", cfg.Env),
		zap.Strings("cors", cfg.Cors.TrustedOrigins),
		zap.String("logLevel", cfg.LogLevel),
	)
	utils.Logger.Info("SMTP Server Config",
		zap.String("host", cfg.SMTP.Host),
		zap.Int("port", cfg.SMTP.Port),
		zap.String("username", cfg.SMTP.Username),
		zap.String("password", cfg.SMTP.Password),
		zap.String("sender", cfg.SMTP.Sender),
	)
	utils.Logger.Info("Rate Limiter",
		zap.Bool("enabled", cfg.Limiter.Enabled),
		zap.Float64("rps", cfg.Limiter.RPS),
		zap.Int("username", cfg.Limiter.Burst),
	)
	utils.Logger.Info("DB",
		zap.String("dsn", cfg.DB.Dsn),
		zap.Int("port", cfg.DB.MaxOpenConns),
		zap.Int("env", cfg.DB.MaxIdleConns),
		zap.String("cors", cfg.DB.MaxIdleTime),
	)
	utils.Logger.Info("AlphaVantage Config",
		zap.String("baseUrl", cfg.AlphaVantage.BaseUrl),
		zap.String("token", cfg.AlphaVantage.Token),
	)
}
