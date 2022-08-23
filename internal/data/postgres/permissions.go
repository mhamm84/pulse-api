package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/internal/data"
	"time"

	"github.com/lib/pq" // New import
)

type pgpermissions struct {
	db *sqlx.DB
}

func NewPermissionsRepository(db *sqlx.DB) data.PermissionsRepository {
	return &pgpermissions{db: db}
}

func (p *pgpermissions) AddForUser(userId int64, codes ...string) error {
	insert := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx, insert, userId, pq.Array(codes))
	return err
}

func (p *pgpermissions) GetAllForUser(userId int64) (data.Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		INNER JOIN users ON users.id = users_permissions.user_id
		WHERE users.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var permissions data.Permissions

	rows, err := p.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
