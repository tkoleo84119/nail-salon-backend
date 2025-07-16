package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StaffUserRepositoryInterface defines the interface for staff user repository
type StaffUserRepositoryInterface interface {
	UpdateStaffUser(ctx context.Context, id int64, req staff.UpdateStaffRequest) (*staff.UpdateStaffResponse, error)
	UpdateStaffMe(ctx context.Context, id int64, req staff.UpdateStaffMeRequest) (*staff.UpdateStaffMeResponse, error)
}

type StaffUserRepository struct {
	db *sqlx.DB
}

func NewStaffUserRepository(db *sqlx.DB) *StaffUserRepository {
	return &StaffUserRepository{db: db}
}

// UpdateStaffUser updates staff user with dynamic fields
func (r *StaffUserRepository) UpdateStaffUser(ctx context.Context, id int64, req staff.UpdateStaffRequest) (*staff.UpdateStaffResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": id,
	}

	if req.Role != nil {
		setParts = append(setParts, "role = :role")
		args["role"] = *req.Role
	}

	if req.IsActive != nil {
		setParts = append(setParts, "is_active = :is_active")
		args["is_active"] = *req.IsActive
	}

	query := fmt.Sprintf(`
		UPDATE staff_users
		SET %s
		WHERE id = :id
		RETURNING
			id,
			username,
			email,
			role,
			is_active,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "))

	var result dbgen.StaffUser
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

	response := &staff.UpdateStaffResponse{
		ID:       utils.FormatID(result.ID),
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
		IsActive: result.IsActive.Bool,
	}

	return response, nil
}

// UpdateStaffMe updates current staff user's information with dynamic fields
func (r *StaffUserRepository) UpdateStaffMe(ctx context.Context, id int64, req staff.UpdateStaffMeRequest) (*staff.UpdateStaffMeResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": id,
	}

	if req.Email != nil {
		setParts = append(setParts, "email = :email")
		args["email"] = *req.Email
	}

	query := fmt.Sprintf(`
		UPDATE staff_users
		SET %s
		WHERE id = :id
		RETURNING
			id,
			username,
			email,
			role
	`, strings.Join(setParts, ", "))

	var result dbgen.StaffUser
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

	response := &staff.UpdateStaffMeResponse{
		ID:       utils.FormatID(result.ID),
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
	}

	return response, nil
}

