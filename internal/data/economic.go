package data

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

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
)

func (r ReportType) ToTable() string {
	return TableFromReportType(r)
}

type Economic struct {
	Date  time.Time       `db:"time" json:"date"`
	Value decimal.Decimal `db:"value" json:"value"`
}

type EconomicWithChange struct {
	Date   time.Time        `db:"time" json:"date"`
	Value  decimal.Decimal  `db:"value" json:"value"`
	Change *decimal.Decimal `db:"percentage_change" json:"change"`
}

type EconomicModel struct {
	DB *sqlx.DB
}

func (m *EconomicModel) LatestWithPercentChange(ctx context.Context, table string) (*EconomicWithChange, error) {
	res := EconomicWithChange{}
	sql := fmt.Sprintf(`
			SELECT
		    	time,
		    	value,
		    	100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value) AS percentage_change
			FROM %s
			ORDER BY time DESC
			LIMIT 1`, table)

	err := m.DB.GetContext(ctx, &res, sql)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *EconomicModel) GetIntervalWithPercentChange(ctx context.Context, table string, years int, paging Paging) (*[]EconomicWithChange, Metadata, error) {
	res := []EconomicWithChange{}

	yearsParam := fmt.Sprintf("'%d year'", years)

	sql := fmt.Sprintf(`
			SELECT
				count(*) OVER(),
		    	time,
		    	value,
		    	100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value) AS percentage_change
			FROM %s
			WHERE time > current_date - INTERVAL %s
			ORDER BY time DESC
			LIMIT $1 OFFSET $2`, table, yearsParam,
	)
	args := []interface{}{paging.limit(), paging.offset()}

	rows, err := m.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	for rows.Next() {
		var economic EconomicWithChange

		err := rows.Scan(
			&totalRecords,
			&economic.Date,
			&economic.Value,
			&economic.Change,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		res = append(res, economic)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, paging.Page, paging.PageSize)
	return &res, metadata, nil
}

func (m *EconomicModel) GetAll(ctx context.Context, table string) (*[]Economic, error) {
	data := []Economic{}
	err := m.DB.SelectContext(ctx, &data, fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC`, table))
	return &data, err
}

func (m *EconomicModel) Insert(ctx context.Context, table string, data *Economic) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (time, value) VALUES (:time, :value)`, table), *data)
	tx.Commit()
	return err
}

func (m *EconomicModel) InsertMany(ctx context.Context, table string, data *[]Economic) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExec(fmt.Sprintf(`INSERT INTO %s (time, value) VALUES (:time, :value)`, table), *data)
	tx.Commit()
	return err
}
