package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewGetAll(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminBookingModel.GetAllParsedRequest, role string, storeIds []int64) (*adminBookingModel.GetAllResponse, error) {
	// Verify store exists
	_, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for staff (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIds); err != nil {
			return nil, err
		}
	}

	// Get booking list from repository
	total, bookings, err := s.repo.Booking.GetAllStoreBookingsByFilter(ctx, storeID, sqlxRepo.GetAllStoreBookingsByFilterParams{
		StylistID: req.StylistID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    req.Status,
		Limit:     &req.Limit,
		Offset:    &req.Offset,
		Sort:      &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get bookings", err)
	}

	bookingIDs := make([]int64, len(bookings))
	for i, booking := range bookings {
		bookingIDs[i] = booking.ID
	}

	if len(bookingIDs) == 0 {
		return &adminBookingModel.GetAllResponse{
			Total: total,
			Items: []adminBookingModel.GetAllItem{},
		}, nil
	}

	batchDetails, err := s.queries.GetBookingDetailsByBookingIDs(ctx, bookingIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get booking details", err)
	}

	detailsByBookingID := make(map[int64][]dbgen.GetBookingDetailsByBookingIDsRow)
	for _, detail := range batchDetails {
		detailsByBookingID[detail.BookingID] = append(detailsByBookingID[detail.BookingID], detail)
	}

	// Build response items by assembling data in service layer
	items := make([]adminBookingModel.GetAllItem, len(bookings))
	for i, booking := range bookings {
		details := detailsByBookingID[booking.ID]

		// Separate main service and sub services based on is_addon flag
		var mainService *adminBookingModel.GetAllMainService
		subServices := []adminBookingModel.GetAllSubService{}

		for _, detail := range details {
			if utils.PgBoolToBool(detail.IsAddon) {
				// This is a sub service (addon)
				subServices = append(subServices, adminBookingModel.GetAllSubService{
					ID:   utils.FormatID(detail.ServiceID),
					Name: detail.ServiceName,
				})
			} else {
				// This is the main service (non-addon)
				mainService = &adminBookingModel.GetAllMainService{
					ID:   utils.FormatID(detail.ServiceID),
					Name: detail.ServiceName,
				}
			}
		}

		// Ensure we have a main service
		if mainService == nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Booking missing main service", nil)
		}

		// Assemble the complete booking item DTO
		item := adminBookingModel.GetAllItem{
			ID: utils.FormatID(booking.ID),
			Customer: adminBookingModel.GetAllCustomer{
				ID:   utils.FormatID(booking.CustomerID),
				Name: booking.CustomerName,
			},
			Stylist: adminBookingModel.GetAllStylist{
				ID:   utils.FormatID(booking.StylistID),
				Name: utils.PgTextToString(booking.StylistName),
			},
			TimeSlot: adminBookingModel.GetAllTimeSlot{
				ID:        utils.FormatID(booking.TimeSlotID),
				WorkDate:  utils.PgDateToDateString(booking.WorkDate),
				StartTime: utils.PgTimeToTimeString(booking.StartTime),
				EndTime:   utils.PgTimeToTimeString(booking.EndTime),
			},
			MainService:    *mainService,
			SubServices:    subServices,
			ActualDuration: utils.PgInt4ToInt32Ptr(booking.ActualDuration),
			Status:         booking.Status,
		}

		items[i] = item
	}

	response := &adminBookingModel.GetAllResponse{
		Total: total,
		Items: items,
	}

	return response, nil
}
