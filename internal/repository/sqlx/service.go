package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
)

type ServiceRepositoryInterface interface {
	UpdateService(ctx context.Context, serviceID int64, req service.UpdateServiceRequest) (*service.UpdateServiceResponse, error)
}

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (r *ServiceRepository) UpdateService(ctx context.Context, serviceID int64, req service.UpdateServiceRequest) (*service.UpdateServiceResponse, error) {
	setParts := []string{}
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

	// Always update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))

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

	return &service.UpdateServiceResponse{
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