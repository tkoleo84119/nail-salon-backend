package sqlx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type StaffUserTokensRepositoryInterface interface {
	Create(ctx context.Context, params CreateParams) (int64, error)
}

type StaffUserTokensRepository struct {
	db *sqlx.DB
}

func NewStaffUserTokensRepository(db *sqlx.DB) *StaffUserTokensRepository {
	return &StaffUserTokensRepository{
		db: db,
	}
}

type CreateParams struct {
	ID           int64              `db:"id"`
	StaffUserID  int64              `db:"staff_user_id"`
	RefreshToken string             `db:"refresh_token"`
	UserAgent    pgtype.Text        `db:"user_agent"`
	IpAddress    *string            `db:"ip_address"`
	ExpiredAt    pgtype.Timestamptz `db:"expired_at"`
}

func (r *StaffUserTokensRepository) Create(ctx context.Context, params CreateParams) (int64, error) {
	query := `
		INSERT INTO staff_user_tokens (id, staff_user_id, refresh_token, user_agent, ip_address, expired_at)
		VALUES (:id, :staff_user_id, :refresh_token, :user_agent, :ip_address, :expired_at)
		RETURNING id
	`

	var id int64
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, params).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}
	return id, nil
}
