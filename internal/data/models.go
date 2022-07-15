package data

import (
	"github.com/jmoiron/sqlx"
)

type Models struct {
	EconomicModel EconomicModel
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		EconomicModel: EconomicModel{DB: db},
	}
}
