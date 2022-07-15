package economic

import (
	"context"
	"fmt"
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

func (s DashboardService) GetDashboardSummary() (*[]data.SummaryHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]data.SummaryHeader, 0, 10)

	if cpiSummary := s.createDashSummary(ctx, data.CPI.ToTable(), "CPI"); cpiSummary != nil {
		dashData = append(dashData, data.SummaryHeader{
			HeaderName: "Monthly CPI",
			Summaries:  []data.Summary{*cpiSummary},
		})
	}
	if consumerSummary := s.createDashSummary(ctx, data.ConsumerSentiment.ToTable(), "Consumer Sentiment"); consumerSummary != nil {
		dashData = append(dashData, data.SummaryHeader{
			HeaderName: "Monthly Consumer Sentiment",
			Summaries:  []data.Summary{*consumerSummary},
		})
	}

	treasurySummaries := []data.Summary{}
	s.add(ctx, &treasurySummaries, data.TreasuryYieldThreeMonth.ToTable(), "Treasury Yield - 3 Months")
	s.add(ctx, &treasurySummaries, data.TreasuryYieldTwoYear.ToTable(), "Treasury Yield - 2 Years")
	s.add(ctx, &treasurySummaries, data.TreasuryYieldFiveYear.ToTable(), "Treasury Yield - 5 Years")
	s.add(ctx, &treasurySummaries, data.TreasuryYieldSevenYear.ToTable(), "Treasury Yield - 7 Years")
	s.add(ctx, &treasurySummaries, data.TreasuryYieldTenYear.ToTable(), "Treasury Yield - 10 Years")
	s.add(ctx, &treasurySummaries, data.TreasuryYieldThirtyYear.ToTable(), "Treasury Yield - 30 Years")

	if retailSalesSummary := s.createDashSummary(ctx, data.RetailSales.ToTable(), "Retail Sales"); retailSalesSummary != nil {
		dashData = append(dashData, data.SummaryHeader{
			HeaderName: "Monthly Retail Sales",
			Summaries:  []data.Summary{*retailSalesSummary},
		})
	}

	fmt.Println(treasurySummaries)

	dashData = append(dashData, data.SummaryHeader{
		HeaderName: "Treasury Yields",
		Summaries:  treasurySummaries,
	})

	return &dashData, nil
}

func (s DashboardService) add(ctx context.Context, summaries *[]data.Summary, tableName, dashHeader string) {
	if summary := s.createDashSummary(ctx, tableName, dashHeader); summary != nil {
		*summaries = append(*summaries, *summary)
	}
}

func (s DashboardService) createDashSummary(ctx context.Context, tableName, dashHeader string) *data.Summary {
	latestWithChange, err := s.Models.EconomicModel.LatestWithPercentChange(ctx, tableName)
	if err != nil {
		s.Logger.PrintInfo("error getting LatestWithPercentChange data for dashboard summary", map[string]interface{}{
			"dataType":   tableName,
			"dashHeader": dashHeader,
			"error":      err.Error(),
		})
		return nil
	}
	return &data.Summary{
		Name:       dashHeader,
		LastUpdate: latestWithChange.Date,
		Value:      latestWithChange.Value,
		Change:     *latestWithChange.Change,
	}
}
