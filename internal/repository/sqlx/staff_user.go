package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffUserRepository struct {
	db *sqlx.DB
}

func NewStaffUserRepository(db *sqlx.DB) *StaffUserRepository {
	return &StaffUserRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStaffByFilterParams struct {
	Username *string
	Email    *string
	Role     *string
	IsActive *bool
	Limit    *int
	Offset   *int
	Sort     *[]string
}

type GetAllStaffByFilterResponse struct {
	ID        int64              `db:"id"`
	Username  string             `db:"username"`
	Email     string             `db:"email"`
	Role      string             `db:"role"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

// GetAllStaffByFilter retrieves staff list with dynamic filtering and pagination
func (r *StaffUserRepository) GetAllStaffByFilter(ctx context.Context, params GetAllStaffByFilterParams) (int, []GetAllStaffByFilterResponse, error) {
	// where conditions
	whereConditions := []string{}
	args := []interface{}{}

	if params.Username != nil && *params.Username != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("username ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Username+"%")
	}

	if params.Email != nil && *params.Email != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("email ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Email+"%")
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	if params.Role != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("role = $%d", len(args)+1))
		args = append(args, *params.Role)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM staff_users
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllStaffByFilterResponse{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"created_at ASC"}
	sort := utils.HandleSortByMap(map[string]string{
		"isActive":  "is_active",
		"role":      "role",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT
			id,
			username,
			email,
			role,
			is_active,
			created_at,
			updated_at
		FROM staff_users
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var result []GetAllStaffByFilterResponse
	if err := r.db.SelectContext(ctx, &result, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, result, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateStaffUserParams struct {
	Email    *string
	Role     *string
	IsActive *bool
}

type UpdateStaffUserResponse struct {
	ID        int64              `db:"id"`
	Username  string             `db:"username"`
	Email     string             `db:"email"`
	Role      string             `db:"role"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

// UpdateStaffUser updates staff user with dynamic fields
func (r *StaffUserRepository) UpdateStaffUser(ctx context.Context, id int64, params UpdateStaffUserParams) (UpdateStaffUserResponse, error) {
	// set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Email != nil && *params.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, *params.Email)
	}

	if params.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role = $%d", len(args)+1))
		args = append(args, *params.Role)
	}

	if params.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateStaffUserResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE staff_users
		SET %s
		WHERE id = $%d
		RETURNING id, username, email, role, is_active, created_at, updated_at
	`, strings.Join(setParts, ", "), len(args))

	var result UpdateStaffUserResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateStaffUserResponse{}, fmt.Errorf("failed to execute update query: %w", err)
	}

	return result, nil
}
