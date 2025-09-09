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
			COALESCE(discount_rate, 0) AS discount_rate,
			COALESCE(discount_amount, 0) AS discount_amount,
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

// ---------------------------------------------------------------------------------------------------------------------

type UpdateCouponParams struct {
	Name     *string
	IsActive *bool
	Note     *string
}

// UpdateCoupon updates coupon fields that are provided
func (r *CouponRepository) UpdateCoupon(ctx context.Context, couponID int64, params UpdateCouponParams) error {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, utils.BoolPtrToPgBool(params.IsActive))
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, utils.StringPtrToPgText(params.Note, false))
	}

	if len(setParts) == 1 {
		return fmt.Errorf("no fields to update")
	}

	args = append(args, couponID)
	query := fmt.Sprintf(`
		UPDATE coupons
		SET %s
		WHERE id = $%d
	`, strings.Join(setParts, ", "), len(args))

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute coupon update: %w", err)
	}

	return nil
}
