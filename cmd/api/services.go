package main

import (
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/services/economic"
	"golang.org/x/time/rate"
	"time"
)

type Services struct {
	alphaVantageEconomicService AlphaVantageEconomicService
	economicdashservice         EconomicDashboardService
}

func NewAlphaServices(models data.Models, client *alpha.Client, logger *jsonlog.Logger) Services {
	return Services{
		alphaVantageEconomicService: economic.AlphaVantageEconomicService{
			Models: models,
			Client: client,
			Logger: logger,
			Limiter: economic.AlphaVantageLimiter{
				MinuteLimiter: rate.NewLimiter(rate.Every(1*time.Minute), 5),
				DailyLimiter:  rate.NewLimiter(rate.Every(24*time.Hour), 500),
			},
		},
		economicdashservice: economic.DashboardService{Models: models, Logger: logger},
	}
}

type AlphaVantageEconomicService interface {
	GetAll(reportType data.ReportType) (*[]data.Economic, error)
	GetIntervalWithPercentChange(reportType data.ReportType, years int, paging data.Paging) (*[]data.EconomicWithChange, data.Metadata, error)
	StartDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]data.SummaryHeader, error)
}
