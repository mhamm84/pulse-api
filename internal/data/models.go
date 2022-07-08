package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data/economic"
)

type Models struct {
	CpiModel               economic.CpiModel
	ConsumerSentimentModel economic.ConsumerSentimentModel
}

func NewModels(db *sqlx.DB) Models {

	return Models{
		CpiModel:               economic.CpiModel{DB: db},
		ConsumerSentimentModel: economic.ConsumerSentimentModel{DB: db},
	}
}
