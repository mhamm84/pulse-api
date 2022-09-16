package services

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/mailer"
	"github.com/mhamm84/pulse-api/internal/repo"
	"github.com/mhamm84/pulse-api/internal/services/economic"
	"github.com/mhamm84/pulse-api/internal/services/economic/alpha"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type ServicesModel struct {
	AlphaVantageEconomicService EconomicService
	Economicdashservice         EconomicDashboardService
	UserService                 UserService
	PermissionsService          PermissionsService
	TokenService                TokenService
}

func NewServicesModel(models repo.Models, client alpha.ClientInterface, mailer *mailer.Mailer, logger *jsonlog.Logger) ServicesModel {
	newTokenService := NewTokenService(models.TokenRepository, logger)
	newUserService := NewUserService(models.UserRepository, models.PermissionsRepository, newTokenService, mailer, logger)

	return ServicesModel{
		AlphaVantageEconomicService: alpha.AlphaVantageEconomicService{
			EconomicRepository: models.EconomicRepository,
			ReportRepository:   models.ReportRepository,
			Client:             client,
			Logger:             logger,
			Limiter: alpha.AlphaVantageLimiter{
				MinuteLimiter: rate.NewLimiter(rate.Every(1*time.Minute), 5),
				DailyLimiter:  rate.NewLimiter(rate.Every(24*time.Hour), 500),
			},
		},
		Economicdashservice: economic.DashboardService{EconomicRepository: models.EconomicRepository, Logger: logger},
		TokenService:        newTokenService,
		UserService:         newUserService,
		PermissionsService:  NewPermissionsService(models.PermissionsRepository, logger),
	}
}

type EconomicService interface {
	GetAll(reportType data.ReportType) (*[]data.Economic, error)
	GetIntervalWithPercentChange(ctx context.Context, wg *sync.WaitGroup, dataChan chan data.EconomicWithChangeResult, errChan chan error, reportType data.ReportType, years int, paging data.Paging)
	GetStats(ctx context.Context, wg *sync.WaitGroup, dataChan chan data.EconomicStatsResult, errChan chan error, reportType data.ReportType, years int, timeBucket int, paging data.Paging)
	StartDataSyncTask()
}

type EconomicDashboardService interface {
	GetDashboardSummary() (*[]data.Summary, error)
}

type UserService interface {
	RegisterUser(user *data.User) error
	ActivateUser(token string) (*data.User, error)
	GetByEmail(email string) (*data.User, error)
	GetFromToken(tokenScope, tokenplaintext string) (*data.User, error)
}

type PermissionsService interface {
	GetAllForUser(userId int64) (data.Permissions, error)
}

type TokenService interface {
	New(userID int64, ttl time.Duration, scope string) (*data.Token, error)
	DeleteAllForUser(userId int64, scope string) error
}
