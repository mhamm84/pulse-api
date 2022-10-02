package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
)

type tokenpg struct {
	DB *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) data.TokenRepository {
	return &tokenpg{DB: db}
}

func (t *tokenpg) Insert(ctx context.Context, token *data.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4) 
	`

	args := []interface{}{token.Hash, token.UserId, token.Expiry, token.Scope}
	_, err := t.DB.ExecContext(ctx, query, args...)

	return err
}

func (t *tokenpg) DeleteAllForUser(ctx context.Context, userId int64, scope string) error {
	query := `
		DELETE FROM tokens 
		WHERE user_id = $1 
		AND scope = $2
	`

	args := []interface{}{userId, scope}
	_, err := t.DB.ExecContext(ctx, query, args...)

	return err
}
