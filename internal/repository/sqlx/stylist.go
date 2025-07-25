package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StylistRepositoryInterface defines the interface for stylist repository
type StylistRepositoryInterface interface {
	UpdateStylist(ctx context.Context, staffUserID int64, req adminStylistModel.UpdateMyStylistRequest) (*adminStylistModel.UpdateMyStylistResponse, error)
	GetStoreStylists(ctx context.Context, storeID int64, limit, offset int) ([]storeModel.GetStoreStylistsItemModel, int, error)
}

type StylistRepository struct {
	db *sqlx.DB
}

func NewStylistRepository(db *sqlx.DB) *StylistRepository {
	return &StylistRepository{db: db}
}

// UpdateStylist updates stylist with dynamic fields
func (r *StylistRepository) UpdateStylist(ctx context.Context, staffUserID int64, req adminStylistModel.UpdateMyStylistRequest) (*adminStylistModel.UpdateMyStylistResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"staff_user_id": staffUserID,
	}

	if req.StylistName != nil {
		setParts = append(setParts, "name = :name")
		args["name"] = *req.StylistName
	}

	if req.GoodAtShapes != nil {
		setParts = append(setParts, "good_at_shapes = :good_at_shapes")
		args["good_at_shapes"] = *req.GoodAtShapes
	}

	if req.GoodAtColors != nil {
		setParts = append(setParts, "good_at_colors = :good_at_colors")
		args["good_at_colors"] = *req.GoodAtColors
	}

	if req.GoodAtStyles != nil {
		setParts = append(setParts, "good_at_styles = :good_at_styles")
		args["good_at_styles"] = *req.GoodAtStyles
	}

	if req.IsIntrovert != nil {
		setParts = append(setParts, "is_introvert = :is_introvert")
		args["is_introvert"] = *req.IsIntrovert
	}

	query := fmt.Sprintf(`
		UPDATE stylists
		SET %s
		WHERE staff_user_id = :staff_user_id
		RETURNING
			id,
			staff_user_id,
			name,
			good_at_shapes,
			good_at_colors,
			good_at_styles,
			is_introvert,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "))

	var result dbgen.Stylist
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

	response := &adminStylistModel.UpdateMyStylistResponse{
		ID:           utils.FormatID(result.ID),
		StaffUserID:  utils.FormatID(staffUserID),
		StylistName:  result.Name.String,
		GoodAtShapes: result.GoodAtShapes,
		GoodAtColors: result.GoodAtColors,
		GoodAtStyles: result.GoodAtStyles,
		IsIntrovert:  result.IsIntrovert.Bool,
	}

	return response, nil
}

// GetStoreStylistsModel represents the database model for store stylists
type GetStoreStylistsModel struct {
	ID            int64    `db:"id"`
	Name          string   `db:"name"`
	GoodAtShapes  []string `db:"good_at_shapes"`
	GoodAtColors  []string `db:"good_at_colors"`
	GoodAtStyles  []string `db:"good_at_styles"`
	IsIntrovert   bool     `db:"is_introvert"`
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
			ID:            utils.FormatID(stylist.ID),
			Name:          stylist.Name,
			GoodAtShapes:  stylist.GoodAtShapes,
			GoodAtColors:  stylist.GoodAtColors,
			GoodAtStyles:  stylist.GoodAtStyles,
			IsIntrovert:   stylist.IsIntrovert,
		}
	}

	return items, total, nil
}
