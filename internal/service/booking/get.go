package booking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) GetInterface {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, bookingID int64, customerID int64) (*bookingModel.GetResponse, error) {
	booking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	if booking.CustomerID != customerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
	}

	// Get booking services/details
	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking details", err)
	}

	var mainService bookingModel.GetServiceItem
	subServices := make([]bookingModel.GetServiceItem, 0)
	for _, detail := range bookingDetails {
		if detail.IsAddon.Bool {
			subServices = append(subServices, bookingModel.GetServiceItem{
				ID:   utils.FormatID(detail.ServiceID),
				Name: detail.ServiceName,
			})
		} else {
			mainService = bookingModel.GetServiceItem{
				ID:   utils.FormatID(detail.ServiceID),
				Name: detail.ServiceName,
			}
		}
	}

	response := &bookingModel.GetResponse{
		ID:            utils.FormatID(booking.ID),
		StoreId:       utils.FormatID(booking.StoreID),
		StoreName:     booking.StoreName,
		StylistId:     utils.FormatID(booking.StylistID),
		StylistName:   utils.PgTextToString(booking.StylistName),
		Date:          utils.PgDateToDateString(booking.WorkDate),
		TimeSlotId:    utils.FormatID(booking.TimeSlotID),
		StartTime:     utils.PgTimeToTimeString(booking.StartTime),
		EndTime:       utils.PgTimeToTimeString(booking.EndTime),
		MainService:   mainService,
		SubServices:   subServices,
		IsChatEnabled: utils.PgBoolToBool(booking.IsChatEnabled),
		Note:          utils.PgTextToString(booking.Note),
		Status:        booking.Status,
		CreatedAt:     utils.PgTimestamptzToTimeString(booking.CreatedAt),
		UpdatedAt:     utils.PgTimestamptzToTimeString(booking.UpdatedAt),
	}

	return response, nil
}
