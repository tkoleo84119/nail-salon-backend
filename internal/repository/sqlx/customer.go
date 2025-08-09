package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// CustomerRepositoryInterface defines the interface for customer repository
type CustomerRepositoryInterface interface {
	GetAllCustomersByFilter(ctx context.Context, params GetAllCustomersByFilterParams) (int, []GetAllCustomersByFilterItem, error)
	UpdateMyCustomer(ctx context.Context, customerID int64, req customerModel.UpdateMyCustomerRequest) (*customerModel.UpdateMyCustomerResponse, error)
	UpdateCustomer(ctx context.Context, customerID int64, params UpdateCustomerParams) (UpdateCustomerResponse, error)
}

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

type GetAllCustomersByFilterParams struct {
	Name          *string
	LineName      *string
	Phone         *string
	Level         *string
	IsBlacklisted *bool
	MinPastDays   *int
	Limit         *int
	Offset        *int
	Sort          *[]string
}

type GetAllCustomersByFilterItem struct {
	ID            int64              `db:"id"`
	Name          string             `db:"name"`
	LineName      pgtype.Text        `db:"line_name"`
	Phone         string             `db:"phone"`
	Birthday      pgtype.Date        `db:"birthday"`
	City          pgtype.Text        `db:"city"`
	Level         pgtype.Text        `db:"level"`
	IsBlacklisted pgtype.Bool        `db:"is_blacklisted"`
	LastVisitAt   pgtype.Timestamptz `db:"last_visit_at"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at"`
}

// GetAllCustomers retrieves all customers with filtering, pagination and sorting
func (r *CustomerRepository) GetAllCustomersByFilter(ctx context.Context, req GetAllCustomersByFilterParams) (int, []GetAllCustomersByFilterItem, error) {
	// set default value for limit and offset
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	// set default value for sort
	sort := utils.HandleSort([]string{"created_at", "updated_at", "is_blacklisted", "last_visit_at"}, "last_visit_at", "DESC", req.Sort)

	whereConditions := []string{}
	args := []interface{}{}

	if req.Name != nil && *req.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*req.Name+"%")
	}

	if req.LineName != nil && *req.LineName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("line_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*req.LineName+"%")
	}

	if req.Phone != nil && *req.Phone != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("phone ILIKE $%d", len(args)+1))
		args = append(args, "%"+*req.Phone+"%")
	}

	if req.Level != nil && *req.Level != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("level = $%d", len(args)+1))
		args = append(args, *req.Level)
	}

	if req.IsBlacklisted != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_blacklisted = $%d", len(args)+1))
		args = append(args, *req.IsBlacklisted)
	}

	if req.MinPastDays != nil && *req.MinPastDays > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("(last_visit_at IS NOT NULL AND last_visit_at < NOW() - INTERVAL '%d days')", *req.MinPastDays))
	}

	// Build WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM customers
		%s`, whereClause)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, err
	}

	dataArgs := append(args, limit, offset)
	limitArgIndex := len(args) + 1
	offsetArgIndex := limitArgIndex + 1

	// Data query
	dataQuery := fmt.Sprintf(`
		SELECT
			id, name, line_name, phone, birthday, city,
			level, is_blacklisted, last_visit_at, updated_at
		FROM customers
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, sort, limitArgIndex, offsetArgIndex)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var results []GetAllCustomersByFilterItem
	for rows.Next() {
		var item GetAllCustomersByFilterItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.LineName,
			&item.Phone,
			&item.Birthday,
			&item.City,
			&item.Level,
			&item.IsBlacklisted,
			&item.LastVisitAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return 0, nil, err
		}
		results = append(results, item)
	}

	return total, results, nil
}

// ------------------------------------------------------------------------------------------------------

// UpdateMyCustomer updates customer's own profile information
func (r *CustomerRepository) UpdateMyCustomer(ctx context.Context, customerID int64, req customerModel.UpdateMyCustomerRequest) (*customerModel.UpdateMyCustomerResponse, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	// Build dynamic SET clause
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}
	if req.Birthday != nil {
		setParts = append(setParts, fmt.Sprintf("birthday = $%d", argIndex))
		args = append(args, *req.Birthday)
		argIndex++
	}
	if req.City != nil {
		setParts = append(setParts, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *req.City)
		argIndex++
	}
	if req.FavoriteShapes != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_shapes = $%d", argIndex))
		args = append(args, *req.FavoriteShapes)
		argIndex++
	}
	if req.FavoriteColors != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_colors = $%d", argIndex))
		args = append(args, *req.FavoriteColors)
		argIndex++
	}
	if req.FavoriteStyles != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_styles = $%d", argIndex))
		args = append(args, *req.FavoriteStyles)
		argIndex++
	}
	if req.IsIntrovert != nil {
		setParts = append(setParts, fmt.Sprintf("is_introvert = $%d", argIndex))
		args = append(args, *req.IsIntrovert)
		argIndex++
	}
	if req.CustomerNote != nil {
		setParts = append(setParts, fmt.Sprintf("customer_note = $%d", argIndex))
		args = append(args, *req.CustomerNote)
		argIndex++
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add customer ID for WHERE clause
	args = append(args, customerID)

	query := fmt.Sprintf(`
		UPDATE customers
		SET %s
		WHERE id = $%d
		RETURNING id, name, phone, birthday, city, favorite_shapes, favorite_colors,
				  favorite_styles, is_introvert, referral_source, referrer, customer_note`,
		strings.Join(setParts, ", "), argIndex)

	var result struct {
		ID             int64     `db:"id"`
		Name           string    `db:"name"`
		Phone          string    `db:"phone"`
		Birthday       *string   `db:"birthday"`
		City           *string   `db:"city"`
		FavoriteShapes *[]string `db:"favorite_shapes"`
		FavoriteColors *[]string `db:"favorite_colors"`
		FavoriteStyles *[]string `db:"favorite_styles"`
		IsIntrovert    *bool     `db:"is_introvert"`
		ReferralSource *[]string `db:"referral_source"`
		Referrer       *string   `db:"referrer"`
		CustomerNote   *string   `db:"customer_note"`
	}

	err := r.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return &customerModel.UpdateMyCustomerResponse{
		ID:             utils.FormatID(result.ID),
		Name:           result.Name,
		Phone:          result.Phone,
		Birthday:       result.Birthday,
		City:           result.City,
		FavoriteShapes: result.FavoriteShapes,
		FavoriteColors: result.FavoriteColors,
		FavoriteStyles: result.FavoriteStyles,
		IsIntrovert:    result.IsIntrovert,
		ReferralSource: result.ReferralSource,
		Referrer:       result.Referrer,
		CustomerNote:   result.CustomerNote,
	}, nil
}

// ------------------------------------------------------------------------------------------------------

type UpdateCustomerParams struct {
	StoreNote     *string
	Level         *string
	IsBlacklisted *bool
}

type UpdateCustomerResponse struct {
	ID            int64              `db:"id"`
	Name          string             `db:"name"`
	Phone         string             `db:"phone"`
	Birthday      pgtype.Date        `db:"birthday"`
	City          pgtype.Text        `db:"city"`
	Level         pgtype.Text        `db:"level"`
	IsBlacklisted pgtype.Bool        `db:"is_blacklisted"`
	LastVisitAt   pgtype.Timestamptz `db:"last_visit_at"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at"`
}

func (r *CustomerRepository) UpdateCustomer(ctx context.Context, customerID int64, params UpdateCustomerParams) (UpdateCustomerResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.StoreNote != nil {
		setParts = append(setParts, fmt.Sprintf("store_note = $%d", len(args)+1))
		args = append(args, *params.StoreNote)
	}

	if params.Level != nil {
		setParts = append(setParts, fmt.Sprintf("level = $%d", len(args)+1))
		args = append(args, *params.Level)
	}

	if params.IsBlacklisted != nil {
		setParts = append(setParts, fmt.Sprintf("is_blacklisted = $%d", len(args)+1))
		args = append(args, *params.IsBlacklisted)
	}

	args = append(args, customerID)

	query := fmt.Sprintf(`
		UPDATE customers
		SET %s
		WHERE id = $%d
		RETURNING id, name, phone, birthday, city, level, is_blacklisted, last_visit_at`,
		strings.Join(setParts, ", "), len(args))

	var result UpdateCustomerResponse
	err := r.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return UpdateCustomerResponse{}, err
	}

	return result, nil
}
