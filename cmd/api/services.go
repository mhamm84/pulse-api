package main

import (
	"context"
	"github.com/mhamm84/gofinance-alpha/alpha"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/services/economic"
)

type Services struct {
	cpiService          CpiService
	economicdashservice EconomicDashboardService
}

func NewAlphaServices(models data.Models, client *alpha.AlphaClient, logger *jsonlog.Logger) Services {
	return Services{
		cpiService:          economic.CpiAlphaService{Models: models, Client: client, Logger: logger},
		economicdashservice: economic.DashboardService{Models: models, Logger: logger},
	}
}

type CpiService interface {
	CpiGetAll(ctx context.Context) (*[]data.Cpi, error)
	StartCpiDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]data.EconomicSummary, error)
}
