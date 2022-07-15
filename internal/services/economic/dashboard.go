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

func (s DashboardService) GetDashboardSummary() (*[]data.Summary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]data.Summary, 0, 10)

	if cpiSummary := s.createDashSummary(ctx, data.CPI.ToTable(), "CPI", nil); cpiSummary != nil {
		dashData = append(dashData, *cpiSummary)
	}
	if consumerSummary := s.createDashSummary(ctx, data.ConsumerSentiment.ToTable(), "Consumer Sentiment", nil); consumerSummary != nil {
		dashData = append(dashData, *consumerSummary)
	}

	if retailSalesSummary := s.createDashSummary(ctx, data.RetailSales.ToTable(), "Retail Sales", nil); retailSalesSummary != nil {
		dashData = append(dashData, *retailSalesSummary)
	}

	treasurySummaries := []data.Summary{}
	s.add(ctx, &treasurySummaries, data.TreasuryYieldThreeMonth.ToTable(), "3M Treasury Yield", addTreasuryExtras(data.TreasuryYieldThreeMonth))
	s.add(ctx, &treasurySummaries, data.TreasuryYieldTwoYear.ToTable(), "2Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldTwoYear))
	s.add(ctx, &treasurySummaries, data.TreasuryYieldFiveYear.ToTable(), "5Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldFiveYear))
	s.add(ctx, &treasurySummaries, data.TreasuryYieldSevenYear.ToTable(), "7Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldSevenYear))
	s.add(ctx, &treasurySummaries, data.TreasuryYieldTenYear.ToTable(), "10Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldTenYear))
	s.add(ctx, &treasurySummaries, data.TreasuryYieldThirtyYear.ToTable(), "30Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldThirtyYear))

	dashData = append(dashData, treasurySummaries...)
	return &dashData, nil
}

func addTreasuryExtras(reportType data.ReportType) map[string]interface{} {
	return map[string]interface{}{"maturity": data.MaturityFromReportType(reportType)}
}

func (s DashboardService) add(ctx context.Context, summaries *[]data.Summary, tableName, dashHeader string, extras map[string]interface{}) {
	if summary := s.createDashSummary(ctx, tableNa:qme, dashHeader, extras); summary != nil {
		*summaries = append(*summaries, *summary)
	}
}

func (s DashboardService) createDashSummary(ctx context.Context, tableName, dashHeader string, extras map[string]interface{}) *data.Summary {
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
		Slug:       tableName,
		Extras:     extras,
	}
}
