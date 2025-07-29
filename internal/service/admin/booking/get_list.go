package adminBooking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetBookingListServiceInterface defines the interface for getting booking list
type GetBookingListServiceInterface interface {
	GetBookingList(ctx context.Context, storeID string, req adminBookingModel.GetBookingListRequest, staffContext common.StaffContext) (*adminBookingModel.GetBookingListResponse, error)
}

type GetBookingListService struct {
	queries           *dbgen.Queries
	bookingRepository *sqlx.BookingRepository
}

func NewGetBookingListService(queries *dbgen.Queries, bookingRepository *sqlx.BookingRepository) *GetBookingListService {
	return &GetBookingListService{
		queries:           queries,
		bookingRepository: bookingRepository,
	}
}

func (s *GetBookingListService) GetBookingList(ctx context.Context, storeID string, req adminBookingModel.GetBookingListRequest, staffContext common.StaffContext) (*adminBookingModel.GetBookingListResponse, error) {
	// Parse and validate store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	// Verify store exists
	_, err = s.queries.GetStoreByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for staff (except SUPER_ADMIN)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin {
		hasAccess, err := utils.CheckOneStoreAccess(storeIDInt, staffContext)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Parse and validate stylist ID if provided
	var stylistIDInt *int64
	if req.StylistID != nil && *req.StylistID != "" {
		parsed, err := utils.ParseID(*req.StylistID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid stylist ID", err)
		}
		stylistIDInt = &parsed
	}

	// Parse date filters if provided
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		parsedTime, err := utils.DateStringToTime(*req.StartDate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid start date format, expected YYYY-MM-DD", err)
		}
		startDate = &parsedTime
	}

	if req.EndDate != nil && *req.EndDate != "" {
		parsedTime, err := utils.DateStringToTime(*req.EndDate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid end date format, expected YYYY-MM-DD", err)
		}
		endDate = &parsedTime
	}

	// Prepare repository parameters
	repoParams := sqlx.GetStoreBookingListParams{
		StylistID: stylistIDInt,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	// Get booking list from repository
	bookings, total, err := s.bookingRepository.GetStoreBookingList(ctx, storeIDInt, repoParams)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get bookings", err)
	}

	// Build response items by assembling data in service layer
	items := make([]adminBookingModel.BookingListItemDTO, 0, len(bookings))
	for _, booking := range bookings {
		// Get booking details (services) for this booking
		details, err := s.queries.GetBookingDetailsByBookingID(ctx, booking.ID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get booking details", err)
		}

		// Separate main service and sub services based on is_addon flag
		var mainService *adminBookingModel.BookingMainServiceDTO
		subServices := []adminBookingModel.BookingSubServiceDTO{}

		for _, detail := range details {
			if utils.PgBoolToBool(detail.IsAddon) {
				// This is a sub service (addon)
				subServices = append(subServices, adminBookingModel.BookingSubServiceDTO{
					ID:   utils.FormatID(detail.ServiceID),
					Name: detail.ServiceName,
				})
			} else {
				// This is the main service (non-addon)
				mainService = &adminBookingModel.BookingMainServiceDTO{
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
		item := adminBookingModel.BookingListItemDTO{
			ID: utils.FormatID(booking.ID),
			Customer: adminBookingModel.BookingCustomerDTO{
				ID:   utils.FormatID(booking.CustomerID),
				Name: booking.CustomerName,
			},
			Stylist: adminBookingModel.BookingStylistDTO{
				ID:   utils.FormatID(booking.StylistID),
				Name: utils.PgTextToString(booking.StylistName),
			},
			TimeSlot: adminBookingModel.BookingTimeSlotDTO{
				ID:        utils.FormatID(booking.TimeSlotID),
				WorkDate:  utils.PgDateToDateString(booking.WorkDate),
				StartTime: utils.PgTimeToTimeString(booking.StartTime),
				EndTime:   utils.PgTimeToTimeString(booking.EndTime),
			},
			MainService: *mainService,
			SubServices: subServices,
			Status:      booking.Status,
		}

		items = append(items, item)
	}

	response := &adminBookingModel.GetBookingListResponse{
		Total: total,
		Items: items,
	}

	return response, nil
}