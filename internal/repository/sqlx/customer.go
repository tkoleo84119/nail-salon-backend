package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// CustomerRepositoryInterface defines the interface for customer repository
type CustomerRepositoryInterface interface {
	UpdateMyCustomer(ctx context.Context, customerID int64, req customer.UpdateMyCustomerRequest) (*customer.UpdateMyCustomerResponse, error)
}

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// UpdateMyCustomer updates customer's own profile information
func (r *CustomerRepository) UpdateMyCustomer(ctx context.Context, customerID int64, req customer.UpdateMyCustomerRequest) (*customer.UpdateMyCustomerResponse, error) {
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

	return &customer.UpdateMyCustomerResponse{
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
