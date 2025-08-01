package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type StaffUserStoreAccessRepositoryInterface interface {
	CreateStaffUserStoreAccessTx(ctx context.Context, tx *sqlx.Tx, req CreateStaffUserStoreAccessTxParams) error
	GetByStaffId(ctx context.Context, staffId int64, isActive *bool) ([]GetByStaffIdItem, error)
}

type StaffUserStoreAccessRepository struct {
	db *sqlx.DB
}

func NewStaffUserStoreAccessRepository(db *sqlx.DB) *StaffUserStoreAccessRepository {
	return &StaffUserStoreAccessRepository{
		db: db,
	}
}

type CreateStaffUserStoreAccessTxParams struct {
	StoreID     int64 `db:"store_id"`
	StaffUserID int64 `db:"staff_user_id"`
}

func (r *StaffUserStoreAccessRepository) CreateStaffUserStoreAccessTx(ctx context.Context, tx *sqlx.Tx, req CreateStaffUserStoreAccessTxParams) (int64, error) {
	query := `
		INSERT INTO staff_user_store_access (store_id, staff_user_id)
		VALUES (:store_id, :staff_user_id)
		RETURNING id
	`

	var id int64
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to create staff user store access: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, req).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create staff user store access: %w", err)
	}

	return id, nil
}

type GetByStaffIdItem struct {
	StoreID  int64       `db:"store_id"`
	Name     string      `db:"name"`
	Address  pgtype.Text `db:"address"`
	Phone    pgtype.Text `db:"phone"`
	IsActive pgtype.Bool `db:"is_active"`
}

func (r *StaffUserStoreAccessRepository) GetByStaffId(ctx context.Context, staffId int64, isActive *bool) ([]GetByStaffIdItem, error) {
	whereParts := []string{"sc.staff_user_id = $1"}
	args := []interface{}{staffId}

	if isActive != nil {
		whereParts = append(whereParts, "s.is_active = $2")
		args = append(args, *isActive)
	}

	query := fmt.Sprintf(`
		SELECT sc.store_id, s.name, s.address, s.phone, s.is_active
		FROM staff_user_store_access sc
		JOIN stores s ON s.id = sc.store_id
		WHERE %s
	`, strings.Join(whereParts, " AND "))

	var results []GetByStaffIdItem
	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return results, nil
}
