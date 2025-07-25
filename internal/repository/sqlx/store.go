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
	UpdateStore(ctx context.Context, storeID int64, req adminStoreModel.UpdateStoreRequest) (*adminStoreModel.UpdateStoreResponse, error)
	GetStores(ctx context.Context, limit, offset int) ([]storeModel.GetStoresItemModel, int, error)
}

type StoreRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
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
