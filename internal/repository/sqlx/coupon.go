package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CouponRepository struct {
	db *sqlx.DB
}

func NewCouponRepository(db *sqlx.DB) *CouponRepository {
	return &CouponRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllCouponsByFilterParams struct {
	Name     *string
	Code     *string
	IsActive *bool
	Limit    *int
	Offset   *int
	Sort     *[]string
}

type GetAllCouponsByFilterItem struct {
	ID             int64              `db:"id"`
	Name           string             `db:"name"`
	DisplayName    string             `db:"display_name"`
	Code           string             `db:"code"`
	DiscountRate   pgtype.Numeric     `db:"discount_rate"`
	DiscountAmount pgtype.Numeric     `db:"discount_amount"`
	IsActive       pgtype.Bool        `db:"is_active"`
	Note           pgtype.Text        `db:"note"`
	CreatedAt      pgtype.Timestamptz `db:"created_at"`
	UpdatedAt      pgtype.Timestamptz `db:"updated_at"`
}

// GetAllCouponsByFilter retrieves coupons with filtering, pagination and sorting
func (r *CouponRepository) GetAllCouponsByFilter(ctx context.Context, params GetAllCouponsByFilterParams) (int, []GetAllCouponsByFilterItem, error) {
	whereConditions := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.Code != nil && *params.Code != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("code ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Code+"%")
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM coupons
		%s`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllCouponsByFilterItem{}, nil
	}

	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"isActive":  "is_active",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			display_name,
			code,
			discount_rate,
			discount_amount,
			is_active,
			COALESCE(note, '') as note,
			created_at,
			updated_at
		FROM coupons
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
    `, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllCouponsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, results, nil
}
