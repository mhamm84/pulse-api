package economic

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"time"
)

const (
	dashboardTimeout = 10
)

type DashboardService struct {
	Models data.Models
	Logger *jsonlog.Logger
}

func (s DashboardService) GetDashboardSummary() (*[]data.EconomicSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]data.EconomicSummary, 0, 10)

	latestCpi, err := s.Models.CpiModel.LatestCpiWithPercentChange(ctx)
	if err != nil {
		return nil, err
	}

	dashData = append(dashData, data.EconomicSummary{
		Name:       "CPI",
		LastUpdate: latestCpi.Date,
		Value:      latestCpi.Value,
		Change:     latestCpi.Change,
	})
	return &dashData, nil
}
