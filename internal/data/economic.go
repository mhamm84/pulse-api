package data

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

type Report struct {
	Id                      int64     `db:"id" json:"id"`
	Slug                    string    `db:"slug" json:"slug"`
	DisplayName             string    `db:"display_name" json:"displayName"`
	Description             string    `db:"description" json:"description"`
	Image                   string    `db:"image" json:"image"`
	LastPullDate            time.Time `db:"last_data_pull" json:"lastPullDate"`
	InitialSyncDelayMinutes int       `db:"initial_sync_delay_minutes" json:"initialSyncDelayMinutes"`
	Extras                  Extras    `json:"extras"`
}

type Extras map[string]interface{}

func (e Extras) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Extras) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &e)
}

type ReportType int8
type TreasuryMaturity string

const (
	CPI ReportType = iota
	ConsumerSentiment
	RetailSales
	TreasuryYieldThreeMonth
	TreasuryYieldTwoYear
	TreasuryYieldFiveYear
	TreasuryYieldSevenYear
	TreasuryYieldTenYear
	TreasuryYieldThirtyYear
	RealGDP
	RealGdpPerCapita
	Unknown
)

func ReportTypeTreasuryYieldMaturity(maturity string) ReportType {
	switch maturity {
	case "3m":
		return TreasuryYieldThreeMonth
	case "2y":
		return TreasuryYieldTwoYear
	case "5y":
		return TreasuryYieldFiveYear
	case "7y":
		return TreasuryYieldSevenYear
	case "10y":
		return TreasuryYieldTenYear
	case "30y":
		return TreasuryYieldThirtyYear
	default:
		return Unknown
	}
}

func MaturityFromReportType(report ReportType) TreasuryMaturity {
	switch report {
	case TreasuryYieldThreeMonth:
		return "3m"
	case TreasuryYieldTwoYear:
		return "2y"
	case TreasuryYieldFiveYear:
		return "5y"
	case TreasuryYieldSevenYear:
		return "7y"
	case TreasuryYieldTenYear:
		return "10y"
	case TreasuryYieldThirtyYear:
		return "30y"
	default:
		return "Unknown"
	}
}

func (r ReportType) String() string {
	switch r {
	case CPI:
		return "CPI"
	case ConsumerSentiment:
		return "CONSUMER_SENTIMENT"
	case RetailSales:
		return "RETAIL_SALES"
	case TreasuryYieldThreeMonth:
		return "TREASURY_YIELD_THREE_MONTH"
	case TreasuryYieldTwoYear:
		return "TREASURY_YIELD_TWO_YEAR"
	case TreasuryYieldFiveYear:
		return "TREASURY_YIELD_FIVE_YEAR"
	case TreasuryYieldSevenYear:
		return "TREASURY_YIELD_SEVEN_YEAR"
	case TreasuryYieldTenYear:
		return "TREASURY_YIELD_TEN_YEAR"
	case TreasuryYieldThirtyYear:
		return "TREASURY_YIELD_THIRTY_YEAR"
	case RealGDP:
		return "REAL_GDP"
	case RealGdpPerCapita:
		return "REAL_GDP_PER_CAPITA"
	default:
		return "Unknown"
	}
}

type tableName string

func TableFromReportType(reportType ReportType) string {
	switch reportType {
	case CPI:
		return string(cpiTableName)
	case ConsumerSentiment:
		return string(consumerSentimentTableName)
	case RetailSales:
		return string(retailSalesTableName)
	case TreasuryYieldThreeMonth:
		return string(treasuryYieldThreeMonthTableName)
	case TreasuryYieldTwoYear:
		return string(treasuryYieldTwoYearTableName)
	case TreasuryYieldFiveYear:
		return string(treasuryYieldFiveYearTableName)
	case TreasuryYieldSevenYear:
		return string(treasuryYieldSevenYearTableName)
	case TreasuryYieldTenYear:
		return string(treasuryYieldTenYearTableName)
	case TreasuryYieldThirtyYear:
		return string(treasuryYieldThirtyYearTableName)
	case RealGDP:
		return string(realGdpTableName)
	case RealGdpPerCapita:
		return string(realGdpPerCapitaTableName)
	default:
		return "unknown"
	}
}

const (
	cpiTableName                     tableName = "cpi"
	consumerSentimentTableName       tableName = "consumer_sentiment"
	retailSalesTableName             tableName = "retail_sales"
	treasuryYieldThreeMonthTableName tableName = "treasury_yield_three_month"
	treasuryYieldTwoYearTableName    tableName = "treasury_yield_two_year"
	treasuryYieldFiveYearTableName   tableName = "treasury_yield_five_year"
	treasuryYieldSevenYearTableName  tableName = "treasury_yield_seven_year"
	treasuryYieldTenYearTableName    tableName = "treasury_yield_ten_year"
	treasuryYieldThirtyYearTableName tableName = "treasury_yield_thirty_year"
	realGdpTableName                 tableName = "real_gdp"
	realGdpPerCapitaTableName        tableName = "real_gdp_per_capita"
)

func (r ReportType) ToTable() string {
	return TableFromReportType(r)
}

type Economic struct {
	Date  time.Time       `db:"time" json:"date"`
	Value decimal.Decimal `db:"value" json:"value"`
}

type EconomicWithChange struct {
	Date   time.Time       `db:"time" json:"date"`
	Value  decimal.Decimal `db:"value" json:"value"`
	Change decimal.Decimal `db:"percentage_change" json:"change"`
}

type EconomicWithChangeResult struct {
	Data *[]EconomicWithChange
	Meta *Metadata
}

type EconomicRepository interface {
	LatestWithPercentChange(ctx context.Context, table string) (*EconomicWithChange, error)
	GetIntervalWithPercentChange(ctx context.Context, table string, years int, paging Paging) (*EconomicWithChangeResult, error)
	GetAll(ctx context.Context, table string) (*[]Economic, error)
	Insert(ctx context.Context, table string, data *Economic) error
	InsertMany(ctx context.Context, table string, data *[]Economic) error
}

type ReportRepository interface {
	GetAllReports(ctx context.Context) ([]*Report, error)
	UpdateReportLastPullDate(ctx context.Context, slug string) error
	GetReportBySlug(ctx context.Context, slug string) (*Report, error)
	GetReports(ctx context.Context) (*[]Report, error)
}
