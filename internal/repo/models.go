package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/data/postgres"
)

type Models struct {
	EconomicRepository data.EconomicRepository
	ReportRepository   data.ReportRepository
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		EconomicRepository: postgres.NewEconomicRepository(db),
		ReportRepository:   postgres.NewReportRepository(db),
	}
}
