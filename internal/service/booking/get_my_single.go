package booking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyBookingService struct {
	queries dbgen.Querier
}

func NewGetMyBookingService(queries dbgen.Querier) GetMyBookingServiceInterface {
	return &GetMyBookingService{
		queries: queries,
	}
}

func (s *GetMyBookingService) GetMyBooking(ctx context.Context, bookingIDStr string, customerContext common.CustomerContext) (*bookingModel.GetMyBookingResponse, error) {
	// Parse booking ID
	bookingID, err := utils.ParseID(bookingIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid booking ID", err)
	}

	// Get booking basic information with ownership validation
	booking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return 404 for both non-existent and unauthorized access
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Validate ownership - customer can only access their own bookings
	if booking.CustomerID != customerContext.CustomerID {
		// Return 404 instead of 403 to prevent data leakag
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
	}

	// Get booking services/details
	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking details", err)
	}

	// Build services list
	services := make([]bookingModel.GetMyBookingServiceModel, len(bookingDetails))
	for i, detail := range bookingDetails {
		services[i] = bookingModel.GetMyBookingServiceModel{
			ID:   utils.FormatID(detail.ServiceID),
			Name: detail.ServiceName,
		}
	}

	// Build response
	var note *string
	if booking.Note.Valid && booking.Note.String != "" {
		note = &booking.Note.String
	}

	response := &bookingModel.GetMyBookingResponse{
		ID:          utils.FormatID(booking.ID),
		StoreId:     utils.FormatID(booking.StoreID),
		StoreName:   booking.StoreName,
		StylistId:   utils.FormatID(booking.StylistID),
		StylistName: utils.PgTextToString(booking.StylistName),
		Date:        utils.PgDateToDateString(booking.WorkDate),
		TimeSlot: bookingModel.GetMyBookingTimeSlotModel{
			ID:        utils.FormatID(booking.TimeSlotID),
			StartTime: utils.PgTimeToTimeString(booking.StartTime),
			EndTime:   utils.PgTimeToTimeString(booking.EndTime),
		},
		Services: services,
		Note:     note,
		Status:   booking.Status,
	}

	return response, nil
}
