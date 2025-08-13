package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type BookingDetailRepositoryInterface interface {
	DeleteBookingDetailsByBookingIDTx(ctx context.Context, tx *sqlx.Tx, bookingID int64) error
	BulkCreateBookingDetailsTx(ctx context.Context, tx *sqlx.Tx, params []BulkCreateBookingDetailsParams) error
}

type BookingDetailRepository struct {
	db *sqlx.DB
}

func NewBookingDetailRepository(db *sqlx.DB) *BookingDetailRepository {
	return &BookingDetailRepository{
		db: db,
	}
}

func (r *BookingDetailRepository) DeleteBookingDetailsByBookingIDTx(ctx context.Context, tx *sqlx.Tx, bookingID int64) error {
	query := `
		DELETE FROM booking_details
		WHERE booking_id = $1
	`
	_, err := tx.ExecContext(ctx, query, bookingID)
	return err
}

// ---------------------------------------------------------------------------------------------------------------------

type BulkCreateBookingDetailsParams struct {
	ID        int64          `db:"id"`
	BookingID int64          `db:"booking_id"`
	ServiceID int64          `db:"service_id"`
	Price     pgtype.Numeric `db:"price"`
}

func (r *BookingDetailRepository) BulkCreateBookingDetailsTx(ctx context.Context, tx *sqlx.Tx, params []BulkCreateBookingDetailsParams) error {
	const batchSize = 1000

	var (
		sb   strings.Builder
		args []interface{}
	)

	for i := 0; i < len(params); i += batchSize {
		end := i + batchSize
		if end > len(params) {
			end = len(params)
		}

		sb.Reset()
		args = args[:0]

		sb.WriteString("INSERT INTO booking_details (id, booking_id, service_id, price) VALUES ")

		param := 1
		for j, v := range params[i:end] {
			sb.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d)", param, param+1, param+2, param+3))
			if j < end-i-1 {
				sb.WriteByte(',')
			}
			args = append(args, v.ID, v.BookingID, v.ServiceID, v.Price)
			param += 4
		}

		if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	return nil
}
