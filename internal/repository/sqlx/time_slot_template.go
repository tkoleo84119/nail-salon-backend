package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TimeSlotTemplateRepositoryInterface defines the interface for time slot template repository
type TimeSlotTemplateRepositoryInterface interface {
	UpdateTimeSlotTemplate(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error)
	GetTimeSlotTemplateList(ctx context.Context, params GetTimeSlotTemplateListParams) ([]TimeSlotTemplateModel, int, error)
}

type TimeSlotTemplateRepository struct {
	db *sqlx.DB
}

func NewTimeSlotTemplateRepository(db *sqlx.DB) *TimeSlotTemplateRepository {
	return &TimeSlotTemplateRepository{db: db}
}

// UpdateTimeSlotTemplate updates time slot template with dynamic fields
func (r *TimeSlotTemplateRepository) UpdateTimeSlotTemplate(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": templateID,
	}

	if req.Name != nil {
		setParts = append(setParts, "name = :name")
		args["name"] = *req.Name
	}

	if req.Note != nil {
		setParts = append(setParts, "note = :note")
		args["note"] = *req.Note
	}

	query := fmt.Sprintf(`
		UPDATE time_slot_templates
		SET %s
		WHERE id = :id
		RETURNING id, name, note, updater, created_at, updated_at`,
		strings.Join(setParts, ", "))

	var result struct {
		ID        int64  `db:"id"`
		Name      string `db:"name"`
		Note      string `db:"note"`
		Updater   int64  `db:"updater"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update time slot template: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("time slot template not found")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	return &adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse{
		ID:   utils.FormatID(result.ID),
		Name: result.Name,
		Note: result.Note,
	}, nil
}

type TimeSlotTemplateModel struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Note      string `db:"note"`
	Updater   int64  `db:"updater"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type GetTimeSlotTemplateListParams struct {
	Name   *string
	Limit  *int
	Offset *int
}

// GetTimeSlotTemplateList retrieves time slot templates with pagination and name filtering
func (r *TimeSlotTemplateRepository) GetTimeSlotTemplateList(ctx context.Context, params GetTimeSlotTemplateListParams) ([]TimeSlotTemplateModel, int, error) {
	// Build WHERE clause
	whereParts := []string{}
	args := map[string]interface{}{}

	// Add name filter if provided
	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, "name ILIKE :name")
		args["name"] = "%" + *params.Name + "%"
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Set pagination defaults
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	args["limit"] = limit
	args["offset"] = offset

	// Query for total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM time_slot_templates %s`, whereClause)
	var total int
	rows, err := r.db.NamedQuery(countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count time slot templates: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	// Query for templates
	templatesQuery := fmt.Sprintf(`
		SELECT id, name, note, updater, created_at, updated_at
		FROM time_slot_templates
		%s
		ORDER BY created_at DESC, id DESC
		LIMIT :limit OFFSET :offset`,
		whereClause)

	var results []TimeSlotTemplateModel

	rows, err = r.db.NamedQuery(templatesQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query time slot templates: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result TimeSlotTemplateModel
		if err := rows.StructScan(&result); err != nil {
			return nil, 0, fmt.Errorf("failed to scan time slot template: %w", err)
		}
		results = append(results, result)
	}

	return results, total, nil
}
