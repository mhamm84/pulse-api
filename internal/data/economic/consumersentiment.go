package economic

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

// CONSUMER_SENTIMENT

type ConsumerSentiment struct {
	Date  time.Time       `db:"time"`
	Value decimal.Decimal `db:"value"`
}

type ConsumerSentimentWithChange struct {
	Date   time.Time       `db:"time"`
	Value  decimal.Decimal `db:"value"`
	Change decimal.Decimal `db:"percentage_change"`
}

type ConsumerSentimentModel struct {
	DB *sqlx.DB
}

func (m *ConsumerSentimentModel) LatestConsumerSentimentWithPercentChange(ctx context.Context) (*ConsumerSentimentWithChange, error) {
	res := ConsumerSentimentWithChange{}
	sql := `SELECT
		    	time,
		    	value,
		    	ROUND(100.0 * (1 - LEAD(value) OVER (ORDER BY time desc) / value),2) AS percentage_change
			FROM consumer_sentiment
			ORDER BY time DESC
			LIMIT 1`

	err := m.DB.GetContext(ctx, &res, sql)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *ConsumerSentimentModel) GetAll(ctx context.Context) (*[]ConsumerSentiment, error) {
	csData := []ConsumerSentiment{}
	sql := `SELECT * FROM consumer_sentiment`
	err := m.DB.SelectContext(ctx, &csData, sql)
	return &csData, err
}

func (m *ConsumerSentimentModel) Insert(ctx context.Context, data *ConsumerSentiment) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx, "", *data)
	tx.Commit()
	return err
}

func (m *ConsumerSentimentModel) InsertMany(ctx context.Context, data *[]ConsumerSentiment) error {
	tx := m.DB.MustBeginTx(ctx, nil)
	_, err := tx.NamedExec(`INSERT INTO consumer_sentiment (time, value) VALUES (:time, :value)`, *data)
	tx.Commit()
	return err
}
