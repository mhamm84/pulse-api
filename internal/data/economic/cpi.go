package economic

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

// "data":{"Name":"Consumer Price Index for all Urban Consumers","Interval":"monthly","Unit":"index 1982-1984=100",
// "Data":[{"Date":"2022-05-01","Value":"292.296"}, ...]}

type Cpi struct {
	Date  time.Time       `db:"time"`
	Value decimal.Decimal `db:"value"`
}

type CpiWithChange struct {
	Date   time.Time       `db:"time"`
	Value  decimal.Decimal `db:"value"`
	Change decimal.Decimal `db:"percentage_change"`
}

type CpiModel struct {
	DB *sqlx.DB
}

func (m *CpiModel) LatestCpiWithPercentChange(ctx context.Context) (*CpiWithChange, error) {
	res := CpiWithChange{}
	sql := `SELECT
		    	time,
		    	value,
		    	ROUND(100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value),2) AS percentage_change
			FROM cpi
			ORDER BY time DESC
			LIMIT 1`

	err := m.DB.GetContext(ctx, &res, sql)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *CpiModel) GetAll(ctx context.Context) (*[]Cpi, error) {
	cpiData := []Cpi{}
	sql := `SELECT * FROM cpi`
	err := m.DB.SelectContext(ctx, &cpiData, sql)
	return &cpiData, err
}

func (m *CpiModel) Insert(ctx context.Context, data *Cpi) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx, "", *data)
	tx.Commit()
	return err
}

func (m *CpiModel) InsertMany(ctx context.Context, data *[]Cpi) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExec(`INSERT INTO cpi (time, value) VALUES (:time, :value)`, *data)
	tx.Commit()
	return err
}
