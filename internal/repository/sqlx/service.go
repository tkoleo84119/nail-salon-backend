package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ServiceRepositoryInterface interface {
	GetAllServiceByFilter(ctx context.Context, params GetAllServiceByFilterParams) (int, []GetAllServiceByFilterItem, error)
	UpdateService(ctx context.Context, serviceID int64, req adminServiceModel.UpdateServiceRequest) (*adminServiceModel.UpdateServiceResponse, error)
	GetStoreServices(ctx context.Context, storeID int64, isAddon *bool, limit, offset int) ([]storeModel.GetStoreServicesItemModel, int, error)
}

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// GetStoreServiceListModel represents the database model for admin service list queries
type GetAllServiceByFilterItem struct {
	ID              int64              `db:"id"`
	Name            string             `db:"name"`
	Price           pgtype.Numeric     `db:"price"`
	DurationMinutes int32              `db:"duration_minutes"`
	IsAddon         pgtype.Bool        `db:"is_addon"`
	IsActive        pgtype.Bool        `db:"is_active"`
	IsVisible       pgtype.Bool        `db:"is_visible"`
	Note            pgtype.Text        `db:"note"`
	CreatedAt       pgtype.Timestamptz `db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `db:"updated_at"`
}

type GetAllServiceByFilterParams struct {
	Name      *string
	IsAddon   *bool
	IsActive  *bool
	IsVisible *bool
	Limit     *int
	Offset    *int
	Sort      *[]string
}

// GetStoreServiceList retrieves services for a specific store with admin filtering and pagination
func (r *ServiceRepository) GetAllServiceByFilter(ctx context.Context, params GetAllServiceByFilterParams) (int, []GetAllServiceByFilterItem, error) {
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
	sort := utils.HandleSort([]string{"created_at", "updated_at", "is_active", "is_visible", "is_addon"}, "created_at", "ASC", params.Sort)

	// Build WHERE clause parts
	whereParts := []string{}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, "name ILIKE :name")
		args["name"] = "%" + *params.Name + "%"
	}

	if params.IsAddon != nil {
		whereParts = append(whereParts, "is_addon = :is_addon")
		args["is_addon"] = *params.IsAddon
	}

	if params.IsActive != nil {
		whereParts = append(whereParts, "is_active = :is_active")
		args["is_active"] = *params.IsActive
	}

	if params.IsVisible != nil {
		whereParts = append(whereParts, "is_visible = :is_visible")
		args["is_visible"] = *params.IsVisible
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(whereParts, " AND "))
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM services
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

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			price,
			duration_minutes,
			is_addon,
			is_active,
			is_visible,
			COALESCE(note, '') as note,
			created_at,
			updated_at
		FROM services
		%s
		ORDER BY %s
		LIMIT :limit OFFSET :offset
	`, whereClause, sort)

	var results []GetAllServiceByFilterItem
	rows, err = r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetAllServiceByFilterItem
		if err := rows.StructScan(&item); err != nil {
			return 0, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return total, results, nil
}

func (r *ServiceRepository) UpdateService(ctx context.Context, serviceID int64, req adminServiceModel.UpdateServiceRequest) (*adminServiceModel.UpdateServiceResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	if req.DurationMinutes != nil {
		setParts = append(setParts, fmt.Sprintf("duration_minutes = $%d", argIndex))
		args = append(args, *req.DurationMinutes)
		argIndex++
	}

	if req.IsAddon != nil {
		setParts = append(setParts, fmt.Sprintf("is_addon = $%d", argIndex))
		args = append(args, *req.IsAddon)
		argIndex++
	}

	if req.IsVisible != nil {
		setParts = append(setParts, fmt.Sprintf("is_visible = $%d", argIndex))
		args = append(args, *req.IsVisible)
		argIndex++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if req.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", argIndex))
		args = append(args, *req.Note)
		argIndex++
	}

	// Add WHERE clause
	args = append(args, serviceID)
	whereClause := fmt.Sprintf("WHERE id = $%d", argIndex)

	query := fmt.Sprintf(`
		UPDATE services
		SET %s
		%s
		RETURNING id, name, price, duration_minutes, is_addon, is_visible, is_active, note
	`, strings.Join(setParts, ", "), whereClause)

	var result struct {
		ID              int64  `db:"id"`
		Name            string `db:"name"`
		Price           int64  `db:"price"`
		DurationMinutes int32  `db:"duration_minutes"`
		IsAddon         bool   `db:"is_addon"`
		IsVisible       bool   `db:"is_visible"`
		IsActive        bool   `db:"is_active"`
		Note            string `db:"note"`
	}

	err := r.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return &adminServiceModel.UpdateServiceResponse{
		ID:              fmt.Sprintf("%d", result.ID),
		Name:            result.Name,
		Price:           result.Price,
		DurationMinutes: result.DurationMinutes,
		IsAddon:         result.IsAddon,
		IsVisible:       result.IsVisible,
		IsActive:        result.IsActive,
		Note:            result.Note,
	}, nil
}

// GetStoreServicesModel represents the database model for store services
type GetStoreServicesModel struct {
	ID              int64  `db:"id"`
	Name            string `db:"name"`
	DurationMinutes int32  `db:"duration_minutes"`
	IsAddon         bool   `db:"is_addon"`
	Note            string `db:"note"`
}

// GetStoreServices retrieves services for a specific store with flexible filtering
func (r *ServiceRepository) GetStoreServices(ctx context.Context, storeID int64, isAddon *bool, limit, offset int) ([]storeModel.GetStoreServicesItemModel, int, error) {
	// Build WHERE conditions - services are always filtered by visibility and active status
	whereParts := []string{
		"is_visible = true",
		"is_active = true",
	}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// Add isAddon filtering if provided
	if isAddon != nil {
		whereParts = append(whereParts, "is_addon = :is_addon")
		args["is_addon"] = *isAddon
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Query for services - Note: This API gets ALL services, not store-specific ones
	// Based on the API spec, it seems to be getting services available for a store
	// but the database schema doesn't show a direct store-service relationship
	// Following the spec literally - getting all services with filters
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			duration_minutes,
			is_addon,
			COALESCE(note, '') as note
		FROM services
		WHERE %s
		ORDER BY name ASC
		LIMIT :limit OFFSET :offset
	`, whereClause)

	var services []GetStoreServicesModel
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("query services failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service GetStoreServicesModel
		if err := rows.StructScan(&service); err != nil {
			return nil, 0, fmt.Errorf("scan service failed: %w", err)
		}
		services = append(services, service)
	}

	// Count total records with same conditions
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM services
		WHERE %s
	`, whereClause)

	var total int
	countRow, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("count services failed: %w", err)
	}
	defer countRow.Close()

	if countRow.Next() {
		if err := countRow.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("scan count failed: %w", err)
		}
	}

	// Convert to response models
	items := make([]storeModel.GetStoreServicesItemModel, len(services))
	for i, service := range services {
		var note *string
		if service.Note != "" {
			note = &service.Note
		}

		items[i] = storeModel.GetStoreServicesItemModel{
			ID:              fmt.Sprintf("%d", service.ID),
			Name:            service.Name,
			DurationMinutes: int(service.DurationMinutes),
			IsAddon:         service.IsAddon,
			Note:            note,
		}
	}

	return items, total, nil
}
