package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStaffService struct {
	queries *dbgen.Queries
}

func NewGetStaffService(queries *dbgen.Queries) *GetStaffService {
	return &GetStaffService{
		queries: queries,
	}
}

func (s *GetStaffService) GetStaff(ctx context.Context, staffID string) (*adminStaffModel.GetStaffResponse, error) {
	// Parse staff ID
	staffUserID, err := utils.ParseID(staffID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid staff ID", err)
	}

	// Get staff user information
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get staff user", err)
	}

	// Prepare response with staff information
	response := &adminStaffModel.GetStaffResponse{
		ID:        utils.FormatID(staffUser.ID),
		Username:  staffUser.Username,
		Email:     staffUser.Email,
		Role:      staffUser.Role,
		IsActive:  staffUser.IsActive.Bool,
		CreatedAt: staffUser.CreatedAt.Time,
		Stylist:   nil, // Default to nil
	}

	// Try to get stylist information if exists
	stylist, err := s.queries.GetStylistByStaffUserID(ctx, staffUserID)
	if err != nil {
		// If no stylist found, that's okay - staff member might not be a stylist
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get stylist information", err)
		}
	} else {
		// Convert stylist information
		response.Stylist = &adminStaffModel.StaffStylistInfo{
			ID:           utils.FormatID(stylist.ID),
			Name:         stylist.Name.String,
			GoodAtShapes: stylist.GoodAtShapes,
			GoodAtColors: stylist.GoodAtColors,
			GoodAtStyles: stylist.GoodAtStyles,
			IsIntrovert:  stylist.IsIntrovert.Bool,
		}
	}

	return response, nil
}
