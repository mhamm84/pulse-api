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
	EconomicRepository data.EconomicRepository
	Logger             *jsonlog.Logger
}

func (s DashboardService) GetDashboardSummary() (*[]data.Summary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout*time.Second)
	defer cancel()

	dashData := make([]data.Summary, 0, 10)

	if cpiSummary := createDashSummary(ctx, s.Logger, s.EconomicRepository, data.CPI.ToTable(), "CPI", nil); cpiSummary != nil {
		dashData = append(dashData, *cpiSummary)
	}
	if consumerSummary := createDashSummary(ctx, s.Logger, s.EconomicRepository, data.ConsumerSentiment.ToTable(), "Consumer Sentiment", nil); consumerSummary != nil {
		dashData = append(dashData, *consumerSummary)
	}

	if retailSalesSummary := createDashSummary(ctx, s.Logger, s.EconomicRepository, data.RetailSales.ToTable(), "Retail Sales", nil); retailSalesSummary != nil {
		dashData = append(dashData, *retailSalesSummary)
	}

	treasurySummaries := []data.Summary{}
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldThreeMonth.ToTable(), "3M Treasury Yield", addTreasuryExtras(data.TreasuryYieldThreeMonth))
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldTwoYear.ToTable(), "2Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldTwoYear))
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldFiveYear.ToTable(), "5Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldFiveYear))
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldSevenYear.ToTable(), "7Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldSevenYear))
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldTenYear.ToTable(), "10Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldTenYear))
	add(ctx, s.Logger, s.EconomicRepository, &treasurySummaries, data.TreasuryYieldThirtyYear.ToTable(), "30Y Treasury Yield", addTreasuryExtras(data.TreasuryYieldThirtyYear))

	dashData = append(dashData, treasurySummaries...)
	return &dashData, nil
}

func addTreasuryExtras(reportType data.ReportType) map[string]interface{} {
	return map[string]interface{}{"maturity": data.MaturityFromReportType(reportType)}
}

func add(ctx context.Context, logger *jsonlog.Logger, economyRepo data.EconomicRepository, summaries *[]data.Summary, tableName, dashHeader string, extras map[string]interface{}) {
	if summary := createDashSummary(ctx, logger, economyRepo, tableName, dashHeader, extras); summary != nil {
		*summaries = append(*summaries, *summary)
	}
}

func createDashSummary(ctx context.Context, logger *jsonlog.Logger, economyRepo data.EconomicRepository, tableName, dashHeader string, extras map[string]interface{}) *data.Summary {
	latestWithChange, err := economyRepo.LatestWithPercentChange(ctx, tableName)
	if err != nil {
		logger.PrintInfo("error getting LatestWithPercentChange data for dashboard summary", map[string]interface{}{
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
		Change:     latestWithChange.Change,
		Slug:       tableName,
		Extras:     extras,
	}
}
