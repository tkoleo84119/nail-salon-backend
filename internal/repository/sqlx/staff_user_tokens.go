package sqlx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type StaffUserTokensRepositoryInterface interface {
	Create(ctx context.Context, params CreateParams) (int64, error)
	CheckValid(ctx context.Context, refreshToken string) (bool, error)
	GetValid(ctx context.Context, refreshToken string) (*GetValidResponse, error)
	Revoke(ctx context.Context, refreshToken string) error
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

func (r *StaffUserTokensRepository) CheckValid(ctx context.Context, refreshToken string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM staff_user_tokens WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false
		)
	`

	var valid bool
	err := r.db.GetContext(ctx, &valid, query, refreshToken)
	if err != nil {
		return false, fmt.Errorf("check valid failed: %w", err)
	}
	return valid, nil
}

type GetValidResponse struct {
	ID          int64              `db:"id"`
	StaffUserID int64              `db:"staff_user_id"`
	ExpiredAt   pgtype.Timestamptz `db:"expired_at"`
	IsRevoked   pgtype.Bool        `db:"is_revoked"`
}

func (r *StaffUserTokensRepository) GetValid(ctx context.Context, refreshToken string) (*GetValidResponse, error) {
	query := `
		SELECT id, staff_user_id, expired_at, is_revoked
		FROM staff_user_tokens
		WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false
	`

	var result GetValidResponse
	err := r.db.GetContext(ctx, &result, query, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get valid failed: %w", err)
	}
	return &result, nil
}

// Revoke revokes a refresh token by setting is_revoked = true
func (r *StaffUserTokensRepository) Revoke(ctx context.Context, refreshToken string) error {
	query := `
		UPDATE staff_user_tokens
		SET is_revoked = true, updated_at = NOW()
		WHERE refresh_token = $1
	`

	_, err := r.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("revoke token failed: %w", err)
	}

	return nil
}
