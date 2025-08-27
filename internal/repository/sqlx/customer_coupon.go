package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CustomerCouponRepository struct {
	db *sqlx.DB
}

func NewCustomerCouponRepository(db *sqlx.DB) *CustomerCouponRepository {
	return &CustomerCouponRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllCustomerCouponsByFilterParams struct {
	CustomerID *int64
	CouponID   *int64
	IsUsed     *bool
	Limit      *int
	Offset     *int
	Sort       *[]string
}

type GetAllCustomerCouponsByFilterItem struct {
	ID         int64              `db:"id"`
	CustomerID int64              `db:"customer_id"`
	CouponID   int64              `db:"coupon_id"`
	ValidFrom  pgtype.Timestamptz `db:"valid_from"`
	ValidTo    pgtype.Timestamptz `db:"valid_to"`
	IsUsed     pgtype.Bool        `db:"is_used"`
	UsedAt     pgtype.Timestamptz `db:"used_at"`
	CreatedAt  pgtype.Timestamptz `db:"created_at"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at"`
}

// GetAllCustomerCouponsByFilter retrieves customer_coupons with dynamic filtering and pagination
func (r *CustomerCouponRepository) GetAllCustomerCouponsByFilter(ctx context.Context, params GetAllCustomerCouponsByFilterParams) (int, []GetAllCustomerCouponsByFilterItem, error) {
	whereConditions := []string{}
	args := []interface{}{}

	if params.CustomerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("customer_id = $%d", len(args)+1))
		args = append(args, *params.CustomerID)
	}

	if params.CouponID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("coupon_id = $%d", len(args)+1))
		args = append(args, *params.CouponID)
	}

	if params.IsUsed != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_used = $%d", len(args)+1))
		args = append(args, *params.IsUsed)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM customer_coupons
		%s`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllCustomerCouponsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"is_used, created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"isUsed":    "is_used",
		"validTo":   "valid_to",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query with joins to customers and coupons
	query := fmt.Sprintf(`
		SELECT
			id,
			customer_id,
			coupon_id,
			valid_from,
			valid_to,
			is_used,
			used_at,
			created_at,
			updated_at
		FROM customer_coupons
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
    `, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllCustomerCouponsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, results, nil
}
