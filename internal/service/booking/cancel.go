package booking

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Cancel struct {
	queries       *dbgen.Queries
	db            *pgxpool.Pool
	lineMessenger *utils.LineMessageClient
	activityLog   cache.ActivityLogCacheInterface
}

func NewCancel(queries *dbgen.Queries, db *pgxpool.Pool, lineMessenger *utils.LineMessageClient, activityLog cache.ActivityLogCacheInterface) CancelInterface {
	return &Cancel{
		queries:       queries,
		db:            db,
		lineMessenger: lineMessenger,
		activityLog:   activityLog,
	}
}

func (s *Cancel) Cancel(ctx context.Context, bookingID int64, req bookingModel.CancelRequest, customerID int64) (*bookingModel.CancelResponse, error) {
	// Verify booking exists and belongs to customer
	bookingInfo, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Check if booking belongs to the customer
	if bookingInfo.CustomerID != customerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if booking is in a cancelable state (only SCHEDULED bookings can be canceled)
	if bookingInfo.Status != common.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToCancel)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Cancel booking with optional cancel reason
	_, err = qtx.CancelBooking(ctx, dbgen.CancelBookingParams{
		ID:           bookingID,
		Status:       common.BookingStatusCancelled,
		CancelReason: utils.StringPtrToPgText(req.CancelReason, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to cancel booking", err)
	}

	// update time slot status to available
	isAvailable := true
	_, err = qtx.UpdateTimeSlotIsAvailable(ctx, dbgen.UpdateTimeSlotIsAvailableParams{
		ID:          bookingInfo.TimeSlotID,
		IsAvailable: utils.BoolPtrToPgBool(&isAvailable),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot status", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction commit failed", err)
	}

	newBooking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking details", err)
	}

	subServiceNames := []string{}
	for _, detail := range bookingDetails {
		if detail.IsAddon.Bool {
			subServiceNames = append(subServiceNames, detail.ServiceName)
		}
	}

	// Build response
	response := &bookingModel.CancelResponse{
		ID:              utils.FormatID(newBooking.ID),
		StoreId:         utils.FormatID(newBooking.StoreID),
		StoreName:       newBooking.StoreName,
		StylistId:       utils.FormatID(newBooking.StylistID),
		StylistName:     utils.PgTextToString(newBooking.StylistName),
		CustomerName:    newBooking.CustomerName,
		CustomerPhone:   newBooking.CustomerPhone,
		Date:            utils.PgDateToDateString(newBooking.WorkDate),
		TimeSlotId:      utils.FormatID(newBooking.TimeSlotID),
		StartTime:       utils.PgTimeToTimeString(newBooking.StartTime),
		EndTime:         utils.PgTimeToTimeString(newBooking.EndTime),
		MainServiceName: bookingDetails[0].ServiceName,
		SubServiceNames: subServiceNames,
		Status:          newBooking.Status,
		CreatedAt:       utils.PgTimestamptzToTimeString(newBooking.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(newBooking.UpdatedAt),
	}

	// if customer no chat permission (this mean customer not give permission to liff app, so can't send message in liff app) send line message, but not return error
	if req.HasChatPermission != nil && !*req.HasChatPermission {
		err = s.lineMessenger.SendBookingNotification(newBooking.CustomerLineUid, common.BookingActionCancelled, &utils.BookingData{
			StoreName:       response.StoreName,
			Date:            response.Date,
			StartTime:       response.StartTime,
			EndTime:         response.EndTime,
			CustomerName:    &response.CustomerName,
			CustomerPhone:   &response.CustomerPhone,
			StylistName:     response.StylistName,
			MainServiceName: bookingDetails[0].ServiceName,
			SubServiceNames: subServiceNames,
		})
		if err != nil {
			log.Printf("failed to send line message: %v", err)
		}
	}

	// Log activity
	go func() {
		logCtx := context.Background()
		if err := s.activityLog.LogCustomerBookingCancel(logCtx, bookingInfo.CustomerName); err != nil {
			log.Printf("failed to log customer booking cancel activity: %v", err)
		}
	}()

	return response, nil
}
