package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data/economic"
)

type Models struct {
	EconomicModel economic.EconomicModel
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		EconomicModel: economic.EconomicModel{DB: db},
	}
}
