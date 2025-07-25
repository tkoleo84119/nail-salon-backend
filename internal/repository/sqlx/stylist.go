package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StylistRepositoryInterface defines the interface for stylist repository
type StylistRepositoryInterface interface {
	UpdateStylist(ctx context.Context, staffUserID int64, req adminStylistModel.UpdateMyStylistRequest) (*adminStylistModel.UpdateMyStylistResponse, error)
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
