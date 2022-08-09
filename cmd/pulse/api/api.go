package api

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/cmd/config"
	"github.com/mhamm84/pulse-api/cmd/pulse/helper"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/repo"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	cfg      config.ApiConfig
	logger   *jsonlog.Logger
	services Services
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
		services: NewAlphaServices(repo.NewModels(db), alphaClient, logger),
	}

	// Start the data sync tasks to keep data from the API up to date in the DB
	if cfg.DataSync {
		logger.PrintInfo("Starting startEconomicReportDataSync", nil)
		app.startEconomicReportDataSync()
	}

	// Serve the API
	err = app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
