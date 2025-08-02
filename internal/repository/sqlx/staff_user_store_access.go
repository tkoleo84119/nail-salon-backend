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
	CreateStaffUserStoreAccess(ctx context.Context, req CreateStaffUserStoreAccessTxParams) error
	BatchCreateStaffUserStoreAccessTx(ctx context.Context, tx *sqlx.Tx, params []CreateStaffUserStoreAccessTxParams) error
	GetStaffUserStoreAccessByStaffId(ctx context.Context, staffId int64, isActive *bool) ([]GetStaffUserStoreAccessByStaffIdItem, error)
	CheckStoreAccessExists(ctx context.Context, staffUserID int64, storeID int64) (bool, error)
	BatchDeleteStaffUserStoreAccess(ctx context.Context, staffUserID int64, storeIDs []int64) error
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

func (r *StaffUserStoreAccessRepository) CreateStaffUserStoreAccessTx(ctx context.Context, tx *sqlx.Tx, req CreateStaffUserStoreAccessTxParams) error {
	query := `
		INSERT INTO staff_user_store_access (store_id, staff_user_id)
		VALUES (:store_id, :staff_user_id)
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create staff user store access: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create staff user store access: %w", err)
	}

	return nil
}

func (r *StaffUserStoreAccessRepository) CreateStaffUserStoreAccess(ctx context.Context, req CreateStaffUserStoreAccessTxParams) error {
	query := `
		INSERT INTO staff_user_store_access (store_id, staff_user_id)
		VALUES (:store_id, :staff_user_id)
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create staff user store access: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	return nil
}

func (r *StaffUserStoreAccessRepository) BatchCreateStaffUserStoreAccessTx(ctx context.Context, tx *sqlx.Tx, params []CreateStaffUserStoreAccessTxParams) error {
	const batchSize = 1000

	var (
		sb   strings.Builder
		args []interface{}
	)

	for i := 0; i < len(params); i += batchSize {
		end := i + batchSize
		if end > len(params) {
			end = len(params)
		}

		sb.Reset()
		args = args[:0]

		sb.WriteString(
			"INSERT INTO staff_user_store_access (store_id, staff_user_id) VALUES ",
		)

		param := 1
		for j, v := range params[i:end] {
			sb.WriteString(fmt.Sprintf("($%d,$%d)", param, param+1))
			if j < end-i-1 {
				sb.WriteByte(',')
			}
			args = append(args, v.StoreID, v.StaffUserID)
			param += 2
		}

		if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}
	return nil
}

type GetStaffUserStoreAccessByStaffIdItem struct {
	StoreID  int64       `db:"store_id"`
	Name     string      `db:"name"`
	Address  pgtype.Text `db:"address"`
	Phone    pgtype.Text `db:"phone"`
	IsActive pgtype.Bool `db:"is_active"`
}

func (r *StaffUserStoreAccessRepository) GetStaffUserStoreAccessByStaffId(ctx context.Context, staffId int64, isActive *bool) ([]GetStaffUserStoreAccessByStaffIdItem, error) {
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

	var results []GetStaffUserStoreAccessByStaffIdItem
	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return results, nil
}

func (r *StaffUserStoreAccessRepository) CheckStoreAccessExists(ctx context.Context, staffUserID int64, storeID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM staff_user_store_access WHERE staff_user_id = $1 AND store_id = $2
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, staffUserID, storeID)
	if err != nil {
		return false, fmt.Errorf("failed to check store access existence: %w", err)
	}

	return exists, nil
}

func (r *StaffUserStoreAccessRepository) BatchDeleteStaffUserStoreAccess(ctx context.Context, staffUserID int64, storeIDs []int64) error {
	query := `
		DELETE FROM staff_user_store_access WHERE staff_user_id = $1 AND store_id = ANY($2)
	`

	_, err := r.db.ExecContext(ctx, query, staffUserID, storeIDs)
	if err != nil {
		return fmt.Errorf("failed to delete store access: %w", err)
	}

	return nil
}
