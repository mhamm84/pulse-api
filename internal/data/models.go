package data

import (
	"github.com/jmoiron/sqlx"
)

type Models struct {
	CpiModel CpiModel
}

func NewModels(db *sqlx.DB) Models {

	return Models{
		CpiModel: CpiModel{db: db},
	}
}
