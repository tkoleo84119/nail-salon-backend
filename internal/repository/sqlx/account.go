package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllAccountsByFilterParams struct {
	Name     *string
	IsActive *bool
	Limit    *int
	Offset   *int
	Sort     *[]string
}

type GetAllAccountsByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Note      pgtype.Text        `db:"note"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func (r *AccountRepository) GetAllAccountsByFilter(ctx context.Context, params GetAllAccountsByFilterParams) (int, []GetAllAccountsByFilterItem, error) {
	// where conditions
	whereConditions := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM accounts
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllAccountsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"created_at ASC"}
	sort := utils.HandleSortByMap(map[string]string{
		"isActive":  "is_active",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT id, name, note, is_active, created_at, updated_at
		FROM accounts
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllAccountsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}
