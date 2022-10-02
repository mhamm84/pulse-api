package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
	"time"
)

var AnonymousUser = &data.User{}

type userpg struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) data.UserRepository {
	return &userpg{db: db}
}

func (p *userpg) GetFromToken(ctx context.Context, tokenScope, tokenplaintext string) (*data.User, error) {

	tokenHash := sha256.Sum256([]byte(tokenplaintext))

	query := `
		SELECT u.id, u.created_at, u.name, u.email, u.password_hash, u.activated, u.version
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1
		AND t.scope = $2
		AND t.expiry > $3
	`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user data.User
	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (p *userpg) Insert(ctx context.Context, user *data.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return data.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (p *userpg) GetByEmail(ctx context.Context, email string) (*data.User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version FROM users
		WHERE email = $1`

	var user data.User

	err := p.db.QueryRowContext(ctx, query, email).Scan(&user.ID,
		&user.CreatedAt, &user.Name, &user.Email, &user.Password.Hash, &user.Activated, &user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (p *userpg) Update(ctx context.Context, user *data.User) error {
	query := ` 
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1 
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{user.Name,
		user.Email, user.Password.Hash, user.Activated, user.ID, user.Version,
	}

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return data.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return data.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
