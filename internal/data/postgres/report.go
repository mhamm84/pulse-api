package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
)

type reportPG struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) data.ReportRepository {
	return &reportPG{db: db}
}

func (p *reportPG) GetAllReports(ctx context.Context) ([]*data.Report, error) {
	reports := []*data.Report{}
	query := `
		SELECT
		    slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras
		FROM economic_report`

	err := p.db.SelectContext(ctx, &reports, query)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (p *reportPG) UpdateReportLastPullDate(ctx context.Context, slug string) error {
	query := `UPDATE economic_report SET last_data_pull = NOW() WHERE slug = $1`

	_, err := p.db.ExecContext(ctx, query, slug)
	if err != nil {
		return err
	}
	return nil
}

func (p *reportPG) GetReportBySlug(ctx context.Context, slug string) (*data.Report, error) {
	report := data.Report{}
	query := `
		SELECT
			slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras
		FROM economic_report
		WHERE slug = $1`

	err := p.db.GetContext(ctx, &report, query, slug)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (p *reportPG) GetReports(ctx context.Context) (*[]data.Report, error) {
	reports := []data.Report{}
	query := `
		SELECT
			slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras
		FROM economic_report`

	err := p.db.SelectContext(ctx, &reports, query)
	if err != nil {
		return nil, err
	}
	return &reports, nil
}
