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
	cpiService               CpiService
	consumerSentimentService ConsumerSentimentService
	economicdashservice      EconomicDashboardService
}

func NewAlphaServices(models data.Models, client *alpha.Client, logger *jsonlog.Logger) Services {
	return Services{
		cpiService:               economic.AlphaVantageCpiService{Models: models, Client: client, Logger: logger},
		consumerSentimentService: economic.AlphaVantageConsumerSentimentService{Models: models, Client: client, Logger: logger},
		economicdashservice:      economic.DashboardService{Models: models, Logger: logger},
	}
}

type CpiService interface {
	CpiGetAll(ctx context.Context) (*[]economic2.Cpi, error)
	StartDataSyncTask()
}

type ConsumerSentimentService interface {
	ConsumerSentimentGetAll(ctx context.Context) (*[]economic2.ConsumerSentiment, error)
	StartDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]economic2.Summary, error)
}
