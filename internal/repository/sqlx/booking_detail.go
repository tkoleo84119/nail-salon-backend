package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type BookingDetailRepository struct {
	db *sqlx.DB
}

func NewBookingDetailRepository(db *sqlx.DB) *BookingDetailRepository {
	return &BookingDetailRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type BulkCreateBookingDetailsParams struct {
	ID        int64          `db:"id"`
	BookingID int64          `db:"booking_id"`
	ServiceID int64          `db:"service_id"`
	Price     pgtype.Numeric `db:"price"`
}

// BulkCreateBookingDetailsTx bulk creates booking details
func (r *BookingDetailRepository) BulkCreateBookingDetailsTx(ctx context.Context, tx *sqlx.Tx, params []BulkCreateBookingDetailsParams) error {
	const batchSize = 1000

	stmt1000 := buildInsertSQL(batchSize)

	for i := 0; i < len(params); i += batchSize {
		end := i + batchSize
		if end > len(params) {
			end = len(params)
		}
		batch := params[i:end]

		// if not full batch, use smaller sql
		sql := stmt1000
		if len(batch) != batchSize {
			sql = buildInsertSQL(len(batch))
		}

		// prepare args
		args := make([]interface{}, 0, len(batch)*4)
		for _, v := range batch {
			args = append(args, v.ID, v.BookingID, v.ServiceID, v.Price)
		}

		if _, err := tx.ExecContext(ctx, sql, args...); err != nil {
			return fmt.Errorf("batch insert failed at batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteBookingDetailsByBookingIDTx deletes booking details by booking id
func (r *BookingDetailRepository) DeleteBookingDetailsByBookingIDTx(ctx context.Context, tx *sqlx.Tx, bookingID int64) error {
	query := `
		DELETE FROM booking_details
		WHERE booking_id = $1
	`

	if _, err := tx.ExecContext(ctx, query, bookingID); err != nil {
		return fmt.Errorf("delete booking details failed: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// buildInsertSQL builds sql string that inserts into booking_details table in batch
func buildInsertSQL(batchSize int) string {
	var sb strings.Builder
	sb.WriteString("INSERT INTO booking_details (id, booking_id, service_id, price) VALUES ")
	param := 1
	for i := 0; i < batchSize; i++ {
		sb.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d)", param, param+1, param+2, param+3))
		if i < batchSize-1 {
			sb.WriteByte(',')
		}
		param += 4
	}
	return sb.String()
}
