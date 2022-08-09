package api

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/repo"
	"github.com/mhamm84/pulse-api/internal/services/economic"
	"github.com/mhamm84/pulse-api/internal/services/economic/alpha"
	"golang.org/x/time/rate"
	"time"
)

type Services struct {
	alphaVantageEconomicService EconomicService
	economicdashservice         EconomicDashboardService
}

func NewAlphaServices(models repo.Models, client alpha.ClientInterface, logger *jsonlog.Logger) Services {
	return Services{
		alphaVantageEconomicService: alpha.AlphaVantageEconomicService{
			EconomicRepository: models.EconomicRepository,
			ReportRepository:   models.ReportRepository,
			Client:             client,
			Logger:             logger,
			Limiter: alpha.AlphaVantageLimiter{
				MinuteLimiter: rate.NewLimiter(rate.Every(1*time.Minute), 5),
				DailyLimiter:  rate.NewLimiter(rate.Every(24*time.Hour), 500),
			},
		},
		economicdashservice: economic.DashboardService{EconomicRepository: models.EconomicRepository, Logger: logger},
	}
}

type EconomicService interface {
	GetAll(reportType data.ReportType) (*[]data.Economic, error)
	GetIntervalWithPercentChange(ctx context.Context, dataChan chan data.EconomicWithChangeResult, errChan chan error, reportType data.ReportType, years int, paging data.Paging)
	StartDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]data.Summary, error)
}
