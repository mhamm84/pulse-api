package economic

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

type Economic struct {
	Date  time.Time       `db:"time"`
	Value decimal.Decimal `db:"value"`
}

type EconomicWithChange struct {
	Date   time.Time       `db:"time"`
	Value  decimal.Decimal `db:"value"`
	Change decimal.Decimal `db:"percentage_change"`
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
		    	ROUND(100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value),2) AS percentage_change
			FROM %s
			ORDER BY time DESC
			LIMIT 1`, table)

	err := m.DB.GetContext(ctx, &res, sql)
	if err != nil {
		return nil, err
	}
	return &res, nil
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
