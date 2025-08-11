package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StoreRepositoryInterface interface {
	GetAllStoreByFilter(ctx context.Context, params GetAllStoreByFilterParams) (int, []GetAllStoreByFilterItem, error)
	UpdateStore(ctx context.Context, storeID int64, req UpdateStoreParams) (*UpdateStoreResponse, error)
}

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

type GetAllStoreByFilterParams struct {
	Name     *string
	IsActive *bool
	Limit    *int
	Offset   *int
	Sort     *[]string
}

type GetAllStoreByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Address   pgtype.Text        `db:"address"`
	Phone     pgtype.Text        `db:"phone"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

// GetAllByFilter retrieves all stores, can filter by name and is_active
func (r *StoreRepository) GetAllStoreByFilter(ctx context.Context, params GetAllStoreByFilterParams) (int, []GetAllStoreByFilterItem, error) {
	// Set default pagination values
	limit := 20
	offset := 0
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}

	// Set default sort values
	sort := utils.HandleSort([]string{"created_at", "updated_at", "is_active"}, "created_at", "ASC", params.Sort)

	whereParts := []string{}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, "name ILIKE :name")
		args["name"] = "%" + *params.Name + "%"
	}

	if params.IsActive != nil {
		whereParts = append(whereParts, "is_active = :is_active")
		args["is_active"] = *params.IsActive
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count query for total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM stores
		%s
	`, whereClause)

	var total int
	rows, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, nil, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	query := fmt.Sprintf(`
		SELECT id, name, address, phone, is_active, created_at, updated_at
		FROM stores
		%s
		ORDER BY %s
		LIMIT :limit OFFSET :offset
	`, whereClause, sort)

	var results []GetAllStoreByFilterItem
	rows, err = r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetAllStoreByFilterItem
		if err := rows.StructScan(&item); err != nil {
			return 0, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return total, results, nil
}

// ------------------------------------------------------------------------------------------------

type UpdateStoreParams struct {
	Name     *string
	Address  *string
	Phone    *string
	IsActive *bool
}

type UpdateStoreResponse struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Address   pgtype.Text        `db:"address"`
	Phone     pgtype.Text        `db:"phone"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func (r *StoreRepository) UpdateStore(ctx context.Context, storeID int64, req UpdateStoreParams) (*UpdateStoreResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": storeID,
	}

	if req.Name != nil && *req.Name != "" {
		setParts = append(setParts, "name = :name")
		args["name"] = *req.Name
	}

	if req.Address != nil {
		setParts = append(setParts, "address = :address")
		args["address"] = utils.StringPtrToPgText(req.Address, false)
	}

	if req.Phone != nil {
		setParts = append(setParts, "phone = :phone")
		args["phone"] = utils.StringPtrToPgText(req.Phone, false)
	}

	if req.IsActive != nil {
		setParts = append(setParts, "is_active = :is_active")
		args["is_active"] = utils.BoolPtrToPgBool(req.IsActive)
	}

	query := fmt.Sprintf(`
		UPDATE stores
		SET %s
		WHERE id = :id
		RETURNING
			id,
			name,
			address,
			phone,
			is_active,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "))

	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}
	defer rows.Close()

	var result UpdateStoreResponse
	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned from update")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	return &result, nil
}
