package economic

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/data/economic"
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

func (s DashboardService) GetDashboardSummary() (*[]economic.Summary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]economic.Summary, 0, 10)

	// CPI
	latestCpi, err := s.Models.CpiModel.LatestCpiWithPercentChange(ctx)
	if err != nil {
		return nil, err
	}

	dashData = append(dashData, economic.Summary{
		Name:       "CPI",
		LastUpdate: latestCpi.Date,
		Value:      latestCpi.Value,
		Change:     latestCpi.Change,
	})

	// CONSUMER SENTIMENT
	latestConsumerSentiment, err := s.Models.ConsumerSentimentModel.LatestConsumerSentimentWithPercentChange(ctx)
	if err != nil {
		return nil, err
	}

	dashData = append(dashData, economic.Summary{
		Name:       "Consumer Sentiment",
		LastUpdate: latestConsumerSentiment.Date,
		Value:      latestConsumerSentiment.Value,
		Change:     latestConsumerSentiment.Change,
	})

	return &dashData, nil
}
