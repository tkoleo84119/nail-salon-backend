package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StylistRepositoryInterface defines the interface for stylist repository
type StylistRepositoryInterface interface {
	CreateStylistTx(ctx context.Context, tx *sqlx.Tx, params CreateStylistTxParams) (int64, error)
	GetStylistByStaffUserID(ctx context.Context, staffUserID int64) (*GetStylistByStaffUserIDResponse, error)
	UpdateStylist(ctx context.Context, staffUserID int64, params UpdateStylistParams) (UpdateStylistResponse, error)
	GetStoreStylists(ctx context.Context, storeID int64, limit, offset int) ([]storeModel.GetStoreStylistsItemModel, int, error)
	GetStoreStylistList(ctx context.Context, storeID int64, params GetStoreStylistListParams) ([]GetStoreStylistListModel, int, error)
}

type StylistRepository struct {
	db *sqlx.DB
}

func NewStylistRepository(db *sqlx.DB) *StylistRepository {
	return &StylistRepository{db: db}
}

type CreateStylistTxParams struct {
	ID          int64 `db:"id"`
	StaffUserID int64 `db:"staff_user_id"`
}

func (r *StylistRepository) CreateStylistTx(ctx context.Context, tx *sqlx.Tx, params CreateStylistTxParams) (int64, error) {
	query := `
		INSERT INTO stylists
		VALUES (:id, :staff_user_id)
		RETURNING id
	`

	var id int64
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to create store: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, params).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create store: %w", err)
	}

	return id, nil
}

type GetStylistByStaffUserIDResponse struct {
	ID           int64              `db:"id"`
	Name         pgtype.Text        `db:"name"`
	GoodAtShapes []string           `db:"good_at_shapes"`
	GoodAtColors []string           `db:"good_at_colors"`
	GoodAtStyles []string           `db:"good_at_styles"`
	IsIntrovert  pgtype.Bool        `db:"is_introvert"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}

func (r *StylistRepository) GetStylistByStaffUserID(ctx context.Context, staffUserID int64) (*GetStylistByStaffUserIDResponse, error) {
	query := `
		SELECT id,
		name,
		COALESCE(good_at_shapes, '{}'::text[]) AS good_at_shapes,
		COALESCE(good_at_colors, '{}'::text[]) AS good_at_colors,
		COALESCE(good_at_styles, '{}'::text[]) AS good_at_styles,
		is_introvert,
		created_at,
		updated_at
		FROM stylists
		WHERE staff_user_id = $1
	`

	m := pgtype.NewMap()
	row := r.db.QueryRowContext(ctx, query, staffUserID)

	var result GetStylistByStaffUserIDResponse
	err := row.Scan(
		&result.ID,
		&result.Name,
		m.SQLScanner(&result.GoodAtShapes),
		m.SQLScanner(&result.GoodAtColors),
		m.SQLScanner(&result.GoodAtStyles),
		&result.IsIntrovert,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type UpdateStylistParams struct {
	Name         *string   `db:"name"`
	GoodAtShapes *[]string `db:"good_at_shapes"`
	GoodAtColors *[]string `db:"good_at_colors"`
	GoodAtStyles *[]string `db:"good_at_styles"`
	IsIntrovert  *bool     `db:"is_introvert"`
}

type UpdateStylistResponse struct {
	ID           int64              `db:"id"`
	StaffUserID  int64              `db:"staff_user_id"`
	Name         pgtype.Text        `db:"name"`
	GoodAtShapes []string           `db:"good_at_shapes"`
	GoodAtColors []string           `db:"good_at_colors"`
	GoodAtStyles []string           `db:"good_at_styles"`
	IsIntrovert  pgtype.Bool        `db:"is_introvert"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}

// UpdateStylist updates stylist with dynamic fields
func (r *StylistRepository) UpdateStylist(ctx context.Context, staffUserID int64, params UpdateStylistParams) (UpdateStylistResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{staffUserID}
	argIndex := 2

	if params.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *params.Name)
		argIndex++
	}

	if params.GoodAtShapes != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_shapes = $%d", argIndex))
		args = append(args, *params.GoodAtShapes)
		argIndex++
	}

	if params.GoodAtColors != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_colors = $%d", argIndex))
		args = append(args, *params.GoodAtColors)
		argIndex++
	}

	if params.GoodAtStyles != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_styles = $%d", argIndex))
		args = append(args, *params.GoodAtStyles)
		argIndex++
	}

	if params.IsIntrovert != nil {
		setParts = append(setParts, fmt.Sprintf("is_introvert = $%d", argIndex))
		args = append(args, *params.IsIntrovert)
		argIndex++
	}

	query := fmt.Sprintf(`
		UPDATE stylists
		SET %s
		WHERE staff_user_id = $1
		RETURNING
			id,
			staff_user_id,
			name,
			COALESCE(good_at_shapes, '{}'::text[]) AS good_at_shapes,
			COALESCE(good_at_colors, '{}'::text[]) AS good_at_colors,
			COALESCE(good_at_styles, '{}'::text[]) AS good_at_styles,
			is_introvert,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "))

	row := r.db.QueryRowxContext(ctx, query, args...)
	m := pgtype.NewMap()

	var result UpdateStylistResponse

	err := row.Scan(
		&result.ID,
		&result.StaffUserID,
		&result.Name,
		m.SQLScanner(&result.GoodAtShapes),
		m.SQLScanner(&result.GoodAtColors),
		m.SQLScanner(&result.GoodAtStyles),
		&result.IsIntrovert,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return UpdateStylistResponse{}, fmt.Errorf("scan result failed: %w", err)
	}

	return result, nil
}

// GetStoreStylistsModel represents the database model for store stylists
type GetStoreStylistsModel struct {
	ID           int64    `db:"id"`
	Name         string   `db:"name"`
	GoodAtShapes []string `db:"good_at_shapes"`
	GoodAtColors []string `db:"good_at_colors"`
	GoodAtStyles []string `db:"good_at_styles"`
	IsIntrovert  bool     `db:"is_introvert"`
}

// GetStoreStylists retrieves stylists for a specific store with flexible filtering
func (r *StylistRepository) GetStoreStylists(ctx context.Context, storeID int64, limit, offset int) ([]storeModel.GetStoreStylistsItemModel, int, error) {
	args := map[string]interface{}{
		"store_id": storeID,
		"limit":    limit,
		"offset":   offset,
	}

	// Query for stylists working at the specified store
	// Join stylists -> staff_users -> staff_user_store_access to find stylists with store access
	// Only return active staff users as specified
	query := `
		SELECT
			s.id,
			s.name,
			COALESCE(s.good_at_shapes, '{}') as good_at_shapes,
			COALESCE(s.good_at_colors, '{}') as good_at_colors,
			COALESCE(s.good_at_styles, '{}') as good_at_styles,
			COALESCE(s.is_introvert, false) as is_introvert
		FROM stylists s
		INNER JOIN staff_users su ON s.staff_user_id = su.id
		INNER JOIN staff_user_store_access susa ON su.id = susa.staff_user_id
		WHERE susa.store_id = :store_id
		  AND su.is_active = true
		ORDER BY s.name ASC
		LIMIT :limit OFFSET :offset
	`

	var stylists []GetStoreStylistsModel
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("query stylists failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stylist GetStoreStylistsModel
		if err := rows.StructScan(&stylist); err != nil {
			return nil, 0, fmt.Errorf("scan stylist failed: %w", err)
		}
		stylists = append(stylists, stylist)
	}

	// Count total records with same conditions
	countQuery := `
		SELECT COUNT(DISTINCT s.id)
		FROM stylists s
		INNER JOIN staff_users su ON s.staff_user_id = su.id
		INNER JOIN staff_user_store_access susa ON su.id = susa.staff_user_id
		WHERE susa.store_id = :store_id
		  AND su.is_active = true
	`

	var total int
	countRow, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("count stylists failed: %w", err)
	}
	defer countRow.Close()

	if countRow.Next() {
		if err := countRow.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("scan count failed: %w", err)
		}
	}

	// Convert to response models
	items := make([]storeModel.GetStoreStylistsItemModel, len(stylists))
	for i, stylist := range stylists {
		items[i] = storeModel.GetStoreStylistsItemModel{
			ID:           utils.FormatID(stylist.ID),
			Name:         stylist.Name,
			GoodAtShapes: stylist.GoodAtShapes,
			GoodAtColors: stylist.GoodAtColors,
			GoodAtStyles: stylist.GoodAtStyles,
			IsIntrovert:  stylist.IsIntrovert,
		}
	}

	return items, total, nil
}

type GetStoreStylistListModel struct {
	ID           int64    `db:"id"`
	StaffUserID  int64    `db:"staff_user_id"`
	Name         string   `db:"name"`
	GoodAtShapes []string `db:"good_at_shapes"`
	GoodAtColors []string `db:"good_at_colors"`
	GoodAtStyles []string `db:"good_at_styles"`
	IsIntrovert  bool     `db:"is_introvert"`
}

type GetStoreStylistListParams struct {
	Name        *string
	IsIntrovert *bool
	Limit       *int
	Offset      *int
}

// GetStoreStylistList retrieves stylists for a specific store with dynamic filtering
func (r *StylistRepository) GetStoreStylistList(ctx context.Context, storeID int64, params GetStoreStylistListParams) ([]GetStoreStylistListModel, int, error) {
	// Set default values
	limit := 20
	offset := 0
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}

	// Build WHERE conditions
	whereParts := []string{"susa.store_id = :store_id", "su.is_active = true"}
	args := map[string]interface{}{
		"store_id": storeID,
		"limit":    limit,
		"offset":   offset,
	}

	// Add name filter (case-insensitive partial match)
	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, "s.name ILIKE :name")
		args["name"] = "%" + *params.Name + "%"
	}

	// Add is_introvert filter
	if params.IsIntrovert != nil {
		whereParts = append(whereParts, "s.is_introvert = :is_introvert")
		args["is_introvert"] = *params.IsIntrovert
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Query for stylists with filtering
	query := fmt.Sprintf(`
		SELECT
			s.id,
			s.staff_user_id,
			s.name,
			COALESCE(s.good_at_shapes, '{}') as good_at_shapes,
			COALESCE(s.good_at_colors, '{}') as good_at_colors,
			COALESCE(s.good_at_styles, '{}') as good_at_styles,
			COALESCE(s.is_introvert, false) as is_introvert
		FROM stylists s
		INNER JOIN staff_users su ON s.staff_user_id = su.id
		INNER JOIN staff_user_store_access susa ON su.id = susa.staff_user_id
		WHERE %s
		ORDER BY s.name ASC
		LIMIT :limit OFFSET :offset
	`, whereClause)

	var stylists []GetStoreStylistListModel
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("query stylists failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stylist GetStoreStylistListModel
		if err := rows.StructScan(&stylist); err != nil {
			return nil, 0, fmt.Errorf("scan stylist failed: %w", err)
		}
		stylists = append(stylists, stylist)
	}

	// Count total records with same filtering conditions
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT s.id)
		FROM stylists s
		INNER JOIN staff_users su ON s.staff_user_id = su.id
		INNER JOIN staff_user_store_access susa ON su.id = susa.staff_user_id
		WHERE %s
	`, whereClause)

	var total int
	countRow, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("count stylists failed: %w", err)
	}
	defer countRow.Close()

	if countRow.Next() {
		if err := countRow.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("scan count failed: %w", err)
		}
	}

	return stylists, total, nil
}
