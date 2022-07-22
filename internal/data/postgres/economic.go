package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
)

type economicPG struct {
	db *sqlx.DB
}

func NewEconomicRepository(db *sqlx.DB) data.EconomicRepository {
	return &economicPG{db: db}
}

func (p *economicPG) LatestWithPercentChange(ctx context.Context, table string) (*data.EconomicWithChange, error) {
	res := data.EconomicWithChange{}
	sql := fmt.Sprintf(`
			SELECT
		    	time,
		    	value,
		    	100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value) AS percentage_change
			FROM %s
			ORDER BY time DESC
			LIMIT 1`, table)

	err := p.db.GetContext(ctx, &res, sql)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *economicPG) GetIntervalWithPercentChange(ctx context.Context, table string, years int, paging data.Paging) (*data.EconomicWithChangeResult, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	res := []data.EconomicWithChange{}

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
	args := []interface{}{paging.Limit(), paging.Offset()}

	rows, err := p.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalRecords := 0
	for rows.Next() {
		var economic data.EconomicWithChange
		err := rows.Scan(
			&totalRecords,
			&economic.Date,
			&economic.Value,
			&economic.Change,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, economic)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	metadata := data.CalculateMetadata(totalRecords, paging.Page, paging.PageSize)

	return &data.EconomicWithChangeResult{Data: &res, Meta: &metadata}, nil
}

func (p *economicPG) GetAll(ctx context.Context, table string) (*[]data.Economic, error) {
	data := []data.Economic{}
	err := p.db.SelectContext(ctx, &data, fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC`, table))
	return &data, err
}

func (p *economicPG) Insert(ctx context.Context, table string, data *data.Economic) error {
	tx := p.db.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx, fmt.Sprintf(`INSERT INTO %s (time, value) VALUES (:time, :value)`, table), *data)
	tx.Commit()
	return err
}

func (p *economicPG) InsertMany(ctx context.Context, table string, data *[]data.Economic) error {
	tx := p.db.MustBeginTx(ctx, nil)
	_, err := tx.NamedExec(fmt.Sprintf(`INSERT INTO %s (time, value) VALUES (:time, :value)`, table), *data)
	tx.Commit()
	return err
}
