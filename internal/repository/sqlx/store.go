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
	GetAll(ctx context.Context, isActive *bool) ([]GetAllItem, error)
	UpdateStore(ctx context.Context, storeID int64, req adminStoreModel.UpdateStoreRequest) (*adminStoreModel.UpdateStoreResponse, error)
	GetStores(ctx context.Context, limit, offset int) ([]storeModel.GetStoresItemModel, int, error)
	GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error)
}

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

type GetAllItem struct {
	ID       int64       `db:"id"`
	Name     string      `db:"name"`
	Address  pgtype.Text `db:"address"`
	Phone    pgtype.Text `db:"phone"`
	IsActive pgtype.Bool `db:"is_active"`
}

// GetAll retrieves all stores, can filter by is_active
func (r *StoreRepository) GetAll(ctx context.Context, isActive *bool) ([]GetAllItem, error) {
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

	var results []GetAllItem
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

// GetStoreList retrieves stores with filtering and pagination for admin
func (r *StoreRepository) GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error) {
	// Set default pagination values
	limit := 20
	offset := 0

	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	// Build WHERE clause parts
	whereParts := []string{}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// Add keyword filter for name and address
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		whereParts = append(whereParts, "(name ILIKE :keyword OR address ILIKE :keyword)")
		args["keyword"] = keyword
	}

	// Add isActive filter
	if req.IsActive != nil {
		whereParts = append(whereParts, "is_active = :is_active")
		args["is_active"] = *req.IsActive
	}

	// Build WHERE clause
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
	rows, err := r.db.NamedQuery(countQuery, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, name, address, phone, is_active
		FROM stores
		%s
		ORDER BY name
		LIMIT :limit OFFSET :offset
	`, whereClause)

	var results []GetStoreListModel
	rows, err = r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetStoreListModel
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Build response models
	items := make([]adminStoreModel.StoreListItemDTO, len(results))
	for i, result := range results {
		items[i] = adminStoreModel.StoreListItemDTO{
			ID:       utils.FormatID(result.ID),
			Name:     result.Name,
			Address:  result.Address.String,
			Phone:    result.Phone.String,
			IsActive: result.IsActive.Bool,
		}
	}

	return &adminStoreModel.GetStoreListResponse{
		Total: total,
		Items: items,
	}, nil
}
