package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StoreRepositoryInterface interface {
	CreateStoreTx(ctx context.Context, tx *sqlx.Tx, req CreateStoreTxParams) (int64, error)
	GetAllStoreByFilter(ctx context.Context, params GetAllStoreByFilterParams) (int, []GetAllStoreByFilterItem, error)
	GetAllStore(ctx context.Context, isActive *bool) ([]GetAllStoreItem, error)
	UpdateStore(ctx context.Context, storeID int64, req adminStoreModel.UpdateStoreRequest) (*adminStoreModel.UpdateStoreResponse, error)
	GetStores(ctx context.Context, limit, offset int) ([]storeModel.GetStoresItemModel, int, error)
	GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error)
	CheckStoreNameExists(ctx context.Context, name string) (bool, error)
}

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

type CreateStoreTxParams struct {
	ID      int64       `db:"id"`
	Name    string      `db:"name"`
	Address pgtype.Text `db:"address"`
	Phone   pgtype.Text `db:"phone"`
}

// CreateStoreTx creates a new store in a transaction
func (r *StoreRepository) CreateStoreTx(ctx context.Context, tx *sqlx.Tx, req CreateStoreTxParams) (int64, error) {
	query := `
		INSERT INTO stores (id, name, address, phone)
		VALUES (:id, :name, :address, :phone)
		RETURNING id
	`

	var id int64
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to create store: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, req).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create store: %w", err)
	}

	return id, nil
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
	sort := utils.HandleSort([]string{"created_at", "is_active"}, "created_at", "ASC", params.Sort)

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

type GetAllStoreItem struct {
	ID       int64       `db:"id"`
	Name     string      `db:"name"`
	Address  pgtype.Text `db:"address"`
	Phone    pgtype.Text `db:"phone"`
	IsActive pgtype.Bool `db:"is_active"`
}

// GetAll retrieves all stores, can filter by is_active
func (r *StoreRepository) GetAllStore(ctx context.Context, isActive *bool) ([]GetAllStoreItem, error) {
	whereParts := []string{}
	args := []interface{}{}

	if isActive != nil {
		whereParts = append(whereParts, "is_active = $1")
		args = append(args, isActive)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, address, phone, is_active
		FROM stores
		%s
	`, whereClause)

	var results []GetAllStoreItem
	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return results, nil
}

func (r *StoreRepository) UpdateStore(ctx context.Context, storeID int64, req adminStoreModel.UpdateStoreRequest) (*adminStoreModel.UpdateStoreResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": storeID,
	}

	if req.Name != nil {
		setParts = append(setParts, "name = :name")
		args["name"] = *req.Name
	}

	if req.Address != nil {
		setParts = append(setParts, "address = :address")
		args["address"] = utils.StringPtrToPgText(req.Address, true)
	}

	if req.Phone != nil {
		setParts = append(setParts, "phone = :phone")
		args["phone"] = utils.StringPtrToPgText(req.Phone, true)
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
			is_active
	`, strings.Join(setParts, ", "))

	var result struct {
		ID       int64       `db:"id"`
		Name     string      `db:"name"`
		Address  pgtype.Text `db:"address"`
		Phone    pgtype.Text `db:"phone"`
		IsActive pgtype.Bool `db:"is_active"`
	}

	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned from update")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	return &adminStoreModel.UpdateStoreResponse{
		ID:       utils.FormatID(result.ID),
		Name:     result.Name,
		Address:  result.Address.String,
		Phone:    result.Phone.String,
		IsActive: result.IsActive.Bool,
	}, nil
}

// GetStoresModel represents the database model for store queries
type GetStoresModel struct {
	ID      int64       `db:"id"`
	Name    string      `db:"name"`
	Address pgtype.Text `db:"address"`
	Phone   pgtype.Text `db:"phone"`
}

// GetStores retrieves stores with pagination, filtering by is_active=true
func (r *StoreRepository) GetStores(ctx context.Context, limit, offset int) ([]storeModel.GetStoresItemModel, int, error) {
	// Count query for total records
	countQuery := `
		SELECT COUNT(*)
		FROM stores
		WHERE is_active = true
	`

	var total int
	if err := r.db.Get(&total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Main query with pagination
	query := `
		SELECT id, name,
		       COALESCE(address, '') as address,
		       COALESCE(phone, '') as phone
		FROM stores
		WHERE is_active = true
		ORDER BY name
		LIMIT :limit OFFSET :offset
	`

	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	var results []GetStoresModel
	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetStoresModel
		if err := rows.StructScan(&item); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	// Build response models
	items := make([]storeModel.GetStoresItemModel, len(results))
	for i, result := range results {
		items[i] = storeModel.GetStoresItemModel{
			ID:      utils.FormatID(result.ID),
			Name:    result.Name,
			Address: utils.PgTextToString(result.Address),
			Phone:   utils.PgTextToString(result.Phone),
		}
	}

	return items, total, nil
}

// GetStoreListModel represents the database model for admin store list queries
type GetStoreListModel struct {
	ID       int64       `db:"id"`
	Name     string      `db:"name"`
	Address  pgtype.Text `db:"address"`
	Phone    pgtype.Text `db:"phone"`
	IsActive pgtype.Bool `db:"is_active"`
}

func (r *StoreRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM stores
			WHERE name = $1
		)
	`

	var exists bool
	if err := r.db.Get(&exists, query, name); err != nil {
		return false, fmt.Errorf("failed to check name existence: %w", err)
	}

	return exists, nil
}
