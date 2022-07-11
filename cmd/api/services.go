package main

import (
	"context"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/internal/data"
	economic2 "github.com/mhamm84/pulse-api/internal/data/economic"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/services/economic"
)

type Services struct {
	alphaVantageEconomicService AlphaVantageEconomicService
	economicdashservice         EconomicDashboardService
}

func NewAlphaServices(models data.Models, client *alpha.Client, logger *jsonlog.Logger) Services {
	return Services{
		alphaVantageEconomicService: economic.AlphaVantageEconomicService{Models: models, Client: client, Logger: logger},
		economicdashservice:         economic.DashboardService{Models: models, Logger: logger},
	}
}

type AlphaVantageEconomicService interface {
	GetAll(ctx context.Context, tableName string) (*[]economic2.Economic, error)
	StartDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]economic2.SummaryHeader, error)
}
