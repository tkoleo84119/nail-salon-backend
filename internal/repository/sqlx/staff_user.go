package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StaffUserRepositoryInterface defines the interface for staff user repository
type StaffUserRepositoryInterface interface {
	GetAllStaffByFilter(ctx context.Context, params GetAllStaffByFilterParams) (int, []GetAllStaffByFilterResponse, error)
	UpdateStaffUser(ctx context.Context, id int64, params UpdateStaffUserParams) (*UpdateStaffUserResponse, error)
}

type StaffUserRepository struct {
	db *sqlx.DB
}

func NewStaffUserRepository(db *sqlx.DB) *StaffUserRepository {
	return &StaffUserRepository{
		db: db,
	}
}

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

// GetStaffList retrieves staff list with dynamic filtering and pagination
func (r *StaffUserRepository) GetAllStaffByFilter(ctx context.Context, params GetAllStaffByFilterParams) (int, []GetAllStaffByFilterResponse, error) {
	// Set default values
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)

	// Set default sort values
	sort := utils.HandleSort([]string{"created_at", "updated_at", "is_active", "role"}, "created_at", "ASC", params.Sort)

	whereParts := []string{}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if params.Username != nil && *params.Username != "" {
		whereParts = append(whereParts, "username ILIKE :username")
		args["username"] = "%" + *params.Username + "%"
	}

	if params.Email != nil && *params.Email != "" {
		whereParts = append(whereParts, "email ILIKE :email")
		args["email"] = "%" + *params.Email + "%"
	}

	if params.IsActive != nil {
		whereParts = append(whereParts, "is_active = :is_active")
		args["is_active"] = *params.IsActive
	}

	if params.Role != nil {
		whereParts = append(whereParts, "role = :role")
		args["role"] = *params.Role
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count query for total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM staff_users
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
			return 0, nil, fmt.Errorf("failed to scan count result: %w", err)
		}
	}

	// Get staff list
	listQuery := fmt.Sprintf(`
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
		LIMIT :limit OFFSET :offset
	`, whereClause, sort)

	var result []GetAllStaffByFilterResponse
	rows, err = r.db.NamedQueryContext(ctx, listQuery, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute list query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetAllStaffByFilterResponse
		if err := rows.StructScan(&item); err != nil {
			return 0, nil, fmt.Errorf("failed to scan result: %w", err)
		}
		result = append(result, item)
	}

	return total, result, nil
}

// ------------------------------------------------------------------------------------------------

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
func (r *StaffUserRepository) UpdateStaffUser(ctx context.Context, id int64, params UpdateStaffUserParams) (*UpdateStaffUserResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": id,
	}

	if params.Email != nil && *params.Email != "" {
		setParts = append(setParts, "email = :email")
		args["email"] = *params.Email
	}

	if params.Role != nil {
		setParts = append(setParts, "role = :role")
		args["role"] = *params.Role
	}

	if params.IsActive != nil {
		setParts = append(setParts, "is_active = :is_active")
		args["is_active"] = *params.IsActive
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

	var result UpdateStaffUserResponse
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

	return &result, nil
}
