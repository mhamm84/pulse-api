package economic

import (
	"context"
	"fmt"
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

func (s DashboardService) GetDashboardSummary() (*[]economic.SummaryHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]economic.SummaryHeader, 0, 10)

	if cpiSummary := s.createDashSummary(ctx, cpiTableName, "Monthly CPI"); cpiSummary != nil {
		dashData = append(dashData, economic.SummaryHeader{
			HeaderName: "Monthly CPI",
			Summaries:  []economic.Summary{*cpiSummary},
		})
	}
	if consumerSummary := s.createDashSummary(ctx, consumerSentimentTableName, "Monthly Consumer Sentiment"); consumerSummary != nil {
		dashData = append(dashData, economic.SummaryHeader{
			HeaderName: "Monthly Consumer Sentiment",
			Summaries:  []economic.Summary{*consumerSummary},
		})
	}

	treasurySummaries := []economic.Summary{}
	s.add(ctx, &treasurySummaries, treasuryYieldThreeMonthTableName, "Treasury Yield - 3 Months")
	s.add(ctx, &treasurySummaries, treasuryYieldTwoYearTableName, "Treasury Yield - 2 Years")
	s.add(ctx, &treasurySummaries, treasuryYieldThreeMonthTableName, "Treasury Yield - 5 Years")
	s.add(ctx, &treasurySummaries, treasuryYieldFiveYearTableName, "Treasury Yield - 7 Years")
	s.add(ctx, &treasurySummaries, treasuryYieldTenYearTableName, "Treasury Yield - 10 Years")
	s.add(ctx, &treasurySummaries, treasuryYieldThirtyYearTableName, "Treasury Yield - 30 Years")

	fmt.Println(treasurySummaries)

	dashData = append(dashData, economic.SummaryHeader{
		HeaderName: "Treasury Yields",
		Summaries:  treasurySummaries,
	})

	return &dashData, nil
}

func (s DashboardService) add(ctx context.Context, summaries *[]economic.Summary, tableName, dashHeader string) {
	if summary := s.createDashSummary(ctx, tableName, dashHeader); summary != nil {
		*summaries = append(*summaries, *summary)
	}
}

func (s DashboardService) createDashSummary(ctx context.Context, tableName, dashHeader string) *economic.Summary {
	latestWithChange, err := s.Models.EconomicModel.LatestWithPercentChange(ctx, tableName)
	if err != nil {
		s.Logger.PrintInfo("error getting LatestWithPercentChange data for dashboard summary", map[string]interface{}{
			"dataType":   tableName,
			"dashHeader": dashHeader,
			"error":      err.Error(),
		})
		return nil
	}
	return &economic.Summary{
		Name:       dashHeader,
		LastUpdate: latestWithChange.Date,
		Value:      latestWithChange.Value,
		Change:     latestWithChange.Change,
	}
}
